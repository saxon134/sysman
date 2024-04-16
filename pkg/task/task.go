package task

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saData/saTime"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/go-utils/saTask"
	"github.com/saxon134/sysman/pkg/model/db/mysql"
	"github.com/saxon134/sysman/pkg/sm"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Init() {
	var taskAry = make([]*mysql.TblTask, 0, 100)
	sm.MySql.Raw(
		fmt.Sprintf("select * from %s where pause_at is null and delete_at is null;", mysql.TBNTask),
	).Scan(&taskAry)
	if len(taskAry) == 0 {
		return
	}

	var caseAry = make([]saTask.Case, 0, len(taskAry))
	for _, v := range taskAry {
		if saTask.CheckSpec(v.Spec) {
			caseAry = append(caseAry, saTask.Case{Key: v.Key, Spec: v.Spec, Handler: Handle, Params: v.Params})
		}
	}
	saTask.Init(caseAry...)
}

// Handle
// @Description: 任务执行方法，注册任务时用
func Handle(key string, params string) error {
	var taskObj = new(mysql.TblTask)
	sm.MySql.Raw(
		fmt.Sprintf("select * from %s where `key` = '%s' and pause_at is null and delete_at is null;", mysql.TBNTask, key),
	).First(taskObj)
	if taskObj.Id <= 0 {
		return nil
	}

	var err = RunStep(taskObj, 0, nil, 0, params)
	sm.MySql.Exec(fmt.Sprintf(
		"update %s set status = %d where id = %d;",
		mysql.TBNTask, saHit.Int(err != nil, -1, 2), taskObj.Id,
	))
	return err
}

// RunStep
// @Description: 执行步骤，可以指定步骤ID，不指定则循环任务下所有步骤
func RunStep(taskObj *mysql.TblTask, stepId int64, preStep *mysql.TblTaskStep, logId int64, params string) (err error) {
	if taskObj == nil {
		return saError.Stack(saError.ErrNotExisted)
	}

	//任务步骤信息
	var stepObj = new(mysql.TblTaskStep)

	//指定任务步骤
	if stepId > 0 {
		sm.MySql.Raw(fmt.Sprintf(
			"select * from %s where id = %d and task_id = %d;", mysql.TBNTaskStep, stepId, taskObj.Id,
		)).Scan(stepObj)
		if stepObj.Id <= 0 {
			return saError.Stack(saError.ErrNotExisted)
		}
	} else
	//查询下一个步骤执行
	{
		var stepAry = make([]*mysql.TblTaskStep, 0, 20)
		sm.MySql.Raw(fmt.Sprintf(
			"select * from %s where task_id = %d and pause_at is null and delete_at is null order by seq asc;", mysql.TBNTaskStep, taskObj.Id,
		)).Scan(&stepAry)
		if len(stepAry) == 0 {
			return nil
		}

		//查找下一个可执行步骤
		var isOk = false
		if preStep == nil || preStep.Id <= 0 {
			isOk = true //无preStep，则第一个可执行步骤即可
		}
		for _, v := range stepAry {
			if isOk {
				if v.RelyPreSuccess == 1 {
					if preStep == nil || preStep.Id <= 0 || preStep.Status == sm.StatusSuccess || preStep.Status == sm.StatusPartiallySuccess {
						stepObj = v
						break
					}
				} else {
					stepObj = v
					break
				}
			}

			if preStep != nil && preStep.Id == v.Id {
				isOk = true
			}
		}
	}

	//新增task执行日志
	if logId <= 0 {
		var taskLog = &mysql.TblTaskLog{
			Pid:       0,
			Specified: stepId > 0,
			TaskId:    taskObj.Id,
			CreateAt:  saOrm.Now(),
			DoneAt:    nil,
			Ms:        0,
			Status:    sm.StatusProcessing,
			Params:    params,
		}
		err = sm.MySql.Save(taskLog).Error
		if err != nil {
			return saError.Stack(err)
		}
		logId = taskLog.Id
	}

	//有可执行的步骤，则开始执行
	if stepObj.Id > 0 {
		//步骤日志
		var stepLog = &mysql.TblTaskLog{
			Pid:       logId,
			Specified: stepId > 0,
			TaskId:    taskObj.Id,
			StepId:    stepObj.Id,
			Ms:        0,
			CreateAt:  saOrm.Now(),
			Status:    sm.StatusProcessing,
			Params:    params,
		}

		//缺少配置
		if strings.HasPrefix(stepObj.Url, "http") == false {
			saLog.Err(saError.Stack(saError.ErrData))
			stepObj.Status = sm.StatusFail
			stepLog.Status = sm.StatusFail
			stepLog.Err = saError.ErrData
		} else {
			//发起任务
			var sign, timestamp = genSign()
			var resp = new(RunResponse)
			err = saHttp.Do(saHttp.Params{
				Method: "POST", Url: stepObj.Url,
				Header: map[string]interface{}{"sign": sign, "timestamp": timestamp},
				Body:   map[string]interface{}{"key": stepObj.Key, "logId": logId, "async": stepObj.Type == 2, "params": saHit.OrStr(params, stepObj.Params)},
			}, resp)
			if err != nil || resp.Code != 0 {
				stepObj.Status = sm.StatusFail
				stepLog.Status = sm.StatusFail
				if err != nil {
					stepLog.Err = err.Error()
				} else {
					stepLog.Err = resp.Msg
				}
			} else {
				//步骤执行状态
				if stepObj.Type == 1 {
					stepObj.Status = sm.StatusSuccess
					stepLog.Result = resp.Result.Result
				} else {
					stepObj.Status = sm.StatusProcessing
				}
				stepLog.Status = stepObj.Status
			}
			stepLog.DoneAt = saOrm.Now()
			stepLog.Ms = stepLog.DoneAt.T().UnixMilli() - stepLog.CreateAt.T().UnixMilli()
		}

		//修改步骤状态
		err = sm.MySql.Exec(fmt.Sprintf(
			"update %s set status = %d where id = %d;",
			mysql.TBNTaskStep, stepObj.Status.Int(), stepObj.Id,
		)).Error
		if err != nil {
			saLog.Err(saError.Stack(err))
		}

		//保存步骤执行日志
		err = sm.MySql.Save(stepLog).Error
		if err != nil {
			saLog.Err(saError.Stack(err))
		}
	}

	//判断task状态
	if stepId > 0 {
		taskObj.Status = stepObj.Status

		//保存任务&日志状态
		err = sm.MySql.Exec(fmt.Sprintf(
			"update %s set status = %d where id = %d; \n update %s set status = %d where id = %d;",
			mysql.TBNTaskLog, taskObj.Status.Int(), logId,
			mysql.TBNTask, taskObj.Status.Int(), taskObj.Id,
		)).Error
		if err != nil {
			saLog.Err(saError.Stack(err))
		}
	} else {
		//无可执行步骤，或者说所有步骤都执行完了
		if stepObj.Id <= 0 {
			var failCnt = 0
			sm.MySql.Raw(fmt.Sprintf(
				"select count(*) from %s where task_id = %d and pause_at is null and delete_at is null and status = -1;",
				mysql.TBNTaskStep, taskObj.Id,
			)).Scan(&failCnt)
			if failCnt > 0 {
				taskObj.Status = sm.StatusFail
			} else {
				taskObj.Status = sm.StatusSuccess
			}

			//执行时间
			var ms int64
			var now = time.Now()
			{
				var createAt string
				sm.MySql.Raw(`select date_format(create_at, '%Y-%m-%d %H:%i:%s') from sm_task_log where id = ` + saData.String(logId)).Scan(&createAt)
				if createAt != "" {
					ms = now.UnixMilli() - saTime.TimeFromStr(createAt, time.DateTime).UnixMilli()
				}
			}

			//保存任务&日志状态
			err = sm.MySql.Exec(fmt.Sprintf(
				"update %s set status = %d, done_at = '%s', ms = %d where id = %d; \n update %s set status = %d where id = %d;",
				mysql.TBNTaskLog, taskObj.Status.Int(), time.Now().Format(time.DateTime), ms, logId,
				mysql.TBNTask, taskObj.Status.Int(), taskObj.Id,
			)).Error
			if err != nil {
				saLog.Err(saError.Stack(err))
			}
		} else {
			_ = RunStep(taskObj, 0, stepObj, logId, params)
		}
	}

	return nil
}

// Callback
// @Description: 异步任务步骤，client执行结束后回调结果
func Callback(in *CallbackData) (err error) {
	if in == nil || in.Key == "" || in.LogId <= 0 {
		return saError.Stack(saError.ErrParams)
	}

	//查询日志信息
	var taskLog *mysql.TblTaskLog
	var stepLog *mysql.TblTaskLog
	{
		var logAry = make([]*mysql.TblTaskLog, 0, 2)
		sm.MySql.Raw(
			fmt.Sprintf("select * from %s where id = %d or (pid = %d and `key` = '%s');", mysql.TBNTaskLog, in.LogId, in.LogId, in.Key),
		).Scan(logAry)

		for _, v := range logAry {
			if v.Id == in.LogId {
				taskLog = v
			} else if v.Pid == in.LogId {
				stepLog = v
			}
		}

		if taskLog == nil || stepLog == nil || taskLog.Id <= 0 || stepLog.Id <= 0 {
			return saError.Stack(saError.ErrNotExisted)
		}
	}

	//查询任务信息
	var taskObj = new(mysql.TblTask)
	{
		sm.MySql.Raw(fmt.Sprintf(
			"select * from %s where id = %d;", mysql.TBNTask, taskLog.TaskId,
		)).Scan(taskObj)
		if taskObj.Id <= 0 {
			return saError.Stack(saError.ErrNotExisted)
		}
	}

	//查询步骤信息
	var stepObj = new(mysql.TblTaskStep)
	{
		sm.MySql.Raw(fmt.Sprintf(
			"select * from %s where task_id = %d and `key` = '%s';", mysql.TBNTaskStep, taskLog.TaskId, in.Key,
		)).Scan(stepObj)
		if stepObj.Id <= 0 {
			return saError.Stack(saError.ErrNotExisted)
		}
	}

	//执行结果
	if in.Success {
		stepObj.Status = sm.StatusSuccess
		stepLog.Status = sm.StatusSuccess
	} else {
		stepObj.Status = sm.StatusFail
		stepLog.Status = sm.StatusFail
	}
	stepLog.Result = in.Result
	stepLog.Err = in.Err

	//保存结果数据
	err = sm.MySql.Transaction(func(tx *gorm.DB) error {
		var e error
		e = tx.Save(stepLog).Error
		if e != nil {
			return e
		}

		e = tx.Save(stepObj).Error
		if e != nil {
			return e
		}

		return nil
	})
	if err != nil {
		saLog.Err(saError.Stack(err))
	}

	//指定步骤，则直接修改状态
	if taskLog.Specified {
		taskObj.Status = stepObj.Status
		taskLog.Status = stepObj.Status
		err = sm.MySql.Exec(
			fmt.Sprintf(`
				update %s set status = %d where id = %d;
				update %s set status = %d where id = %d;`,
				mysql.TBNTask, taskObj.Status.Int(), taskObj.Id,
				mysql.TBNTaskLog, taskLog.Status.Int(), taskLog.Id,
			)).Error
		if err != nil {
			saLog.Err(saError.Stack(err))
		}
	} else
	//非指定步骤，查询下一个可执行的步骤并执行
	{
		_ = RunStep(taskObj, 0, stepObj, taskLog.Id, "")
	}

	return nil
}
