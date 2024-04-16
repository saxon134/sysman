package controller

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saLog"
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/go-utils/saTask"
	"github.com/saxon134/sysman/pkg/api"
	"github.com/saxon134/sysman/pkg/model/db/mysql"
	"github.com/saxon134/sysman/pkg/model/dto"
	"github.com/saxon134/sysman/pkg/sm"
	"github.com/saxon134/sysman/pkg/task"
	"gorm.io/gorm"
	"strings"
)

func TaskList(c *api.Context) (res *api.Response, err error) {
	var in = new(dto.TaskListRequest)
	c.BindByApi(in)

	var query = "from sm_task t1 left join sm_sys t2 on t2.id = t1.sys_id where 1 = 1 "

	if in.SysId > 0 {
		query += "and t1.sys_id = " + saData.String(in.SysId) + " "
	}

	if in.Name != "" {
		query += fmt.Sprintf("and t1.name like '%%%s%%' ", in.Name)
	}

	sm.MySql.Raw("select count(*) " + query).Scan(&c.Paging.Total)
	if c.Paging.IsAll() {
		return &api.Response{Result: nil}, nil
	}

	var resAry = make([]*dto.TaskDetail, 0, c.Paging.Limit)
	err = sm.MySql.Raw("select t1.*, t2.name sys " + query).Scan(&resAry).Error
	if sm.MySql.IsError(err) {
		return nil, saError.Stack(err)
	}

	//执行时间
	for _, v := range resAry {
		var statusDic, _ = saTask.Status(v.TblTask.Key)
		if statusDic != nil {
			v.NextTime = statusDic["nextTime"]
			v.PreTime = statusDic["preTime"]
		}

		//paused状态
		v.Paused = v.PauseAt.IsZero() == false
	}
	return &api.Response{Result: resAry}, nil
}

func TaskDetail(c *api.Context) (res *api.Response, err error) {
	var taskId = saData.Int64(c.Query("taskId"))
	if taskId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	var info = new(dto.TaskDetail)
	info.TblTask = new(mysql.TblTask)
	err = sm.MySql.Raw(fmt.Sprintf("select * from %s where id = %d;", mysql.TBNTask, taskId)).Error
	if sm.MySql.IsError(err) {
		return nil, err
	}

	if info.Id <= 0 {
		return nil, saError.Stack(saError.ErrNotExisted)
	}

	//执行时间
	var statusDic, _ = saTask.Status(info.TblTask.Key)
	if statusDic != nil {
		info.NextTime = statusDic["nextTime"]
		info.PreTime = statusDic["preTime"]
	}

	return &api.Response{Result: info}, nil
}

func TaskSteps(c *api.Context) (res *api.Response, err error) {
	var taskId = saData.Int64(c.Query("taskId"))
	if taskId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	var steps = make([]*dto.Step, 0, 10)
	err = sm.MySql.Raw(fmt.Sprintf(
		"select * from %s where task_id = %d and delete_at is null", mysql.TBNTaskStep, taskId,
	)).Scan(&steps).Error
	if sm.MySql.IsError(err) {
		return nil, saError.Stack(err)
	}

	for _, v := range steps {
		if v.PauseAt.IsZero() == false {
			v.Paused = true
		}
	}
	return &api.Response{Result: steps}, nil
}

func TaskStepSave(c *api.Context) (res *api.Response, err error) {
	var in = new(dto.TaskStepsSaveRequest)
	c.BindByApi(in)
	if in.TaskId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	//判断task是否存在
	var taskId int64
	sm.MySql.Raw(fmt.Sprintf("select id from %s where id = %d and delete_at is null;", mysql.TBNTask, in.TaskId)).Scan(&taskId)
	if taskId <= 0 {
		return nil, saError.Stack(saError.ErrNotExisted)
	}

	var stepIds = make([]string, 0, 10)
	for _, v := range in.Steps {
		if v.DeleteAt.IsZero() && v.Id > 0 {
			stepIds = saData.AppendStr(stepIds, saData.String(v.Id))
		}
		v.TaskId = taskId
	}

	err = sm.MySql.Transaction(func(tx *gorm.DB) error {
		//删除不存在的
		if len(stepIds) > 0 {
			err = tx.Exec(fmt.Sprintf(
				"update %s set delete_at = '%s' where task_id = %d and delete_at is null and id not in (%s);",
				mysql.TBNTaskStep, saOrm.Now().String(), taskId, saData.ToSQLIds(stepIds),
			)).Error
		} else {
			err = tx.Exec(fmt.Sprintf(
				"update %s set delete_at = '%s' where task_id = %d and delete_at is null;",
				mysql.TBNTaskStep, saOrm.Now().String(), taskId,
			)).Error
		}

		//保存
		var now = saOrm.Now()
		var ary = make([]*mysql.TblTaskStep, 0, len(in.Steps))
		for _, v := range in.Steps {
			if v.Paused == true {
				v.PauseAt = now
			} else {
				v.PauseAt = nil
			}
			ary = append(ary, v.TblTaskStep)
		}
		err = tx.Save(ary).Error
		if err != nil {
			return saError.Stack(err)
		}

		return nil
	})
	if err != nil {
		return nil, saError.Stack(err)
	}

	return &api.Response{Result: nil}, nil
}

func TaskLogList(c *api.Context) (res *api.Response, err error) {
	var in = new(dto.TaskLogRequest)
	c.BindByApi(in)
	if in.TaskId <= 0 && in.LogId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	var resAry = make([]*dto.TaskLogItem, 0, c.Paging.Limit)

	//返回一级任务执行日志
	if in.TaskId > 0 {
		var query = sm.MySql.Table(mysql.TBNTaskLog).Where("task_id = ?", in.TaskId).Where("pid = 0")
		query.Count(&c.Paging.Total)
		if c.Paging.IsAll() {
			return &api.Response{Result: nil}, nil
		}

		query.Limit(c.Paging.Limit).Offset(c.Paging.Offset).Order("id desc").Find(&resAry)
	}

	//任务详情（步骤日志），需返回步骤基本信息，且不分页
	if in.LogId > 0 {
		//查找日志
		var query = sm.MySql.Table(mysql.TBNTaskLog).Where("pid = ?", in.LogId)
		query.Find(&resAry)

		//查找步骤信息
		var steps = make([]*mysql.TblTaskStep, 0, len(resAry))
		var ids = make([]int64, 0, len(resAry))
		for _, v := range resAry {
			ids = saData.AppendId(ids, v.StepId)
		}
		sm.MySql.Table(mysql.TBNTaskStep).Where("id in ?", ids).Find(&steps)

		//补齐信息
		for _, v := range resAry {
			for _, e := range steps {
				if e.Id == v.StepId {
					v.Step = e
					break
				}
			}
		}
	}

	return &api.Response{Result: resAry}, nil
}

func TaskSave(c *api.Context) (res *api.Response, err error) {
	var in = new(dto.TaskItem)
	c.BindByApi(in)

	if saTask.CheckSpec(in.Spec) == false {
		return nil, saError.Stack("执行周期设置有误")
	}

	if in.Key == "" {
		return nil, saError.Stack(saError.ErrParams)
	}

	if in.Status != sm.StatusDisable {
		if in.Spec == "" || in.SysId <= 0 || in.Name == "" {
			return nil, saError.Stack(saError.ErrParams)
		}
	}

	if in.Spec != "" {
		var ary = strings.Split(in.Spec, " ")
		if len(ary) != 5 && len(ary) != 6 {
			return nil, saError.New("Spec格式有误")
		}
	}

	var taskObj = new(mysql.TblTask)
	if in.Id > 0 {
		sm.MySql.Table(mysql.TBNTask).Where("id = ?", in.Id).First(taskObj)
	}

	if taskObj.SysId != in.SysId && taskObj.SysId > 0 {
		return nil, saError.Stack(saError.ErrNotSupport)
	}

	//key必须唯一
	if taskObj.Key != in.Key {
		var exist int64
		sm.MySql.Raw(fmt.Sprintf(
			"select id from %s where `key` = '%s' and id <> %d;", mysql.TBNTask, in.Key, taskObj.Id,
		)).Scan(&exist)
		if exist > 0 {
			return nil, saError.Stack("key已存在")
		}
	}

	//停止任务
	if in.DeleteAt.IsZero() == false || in.Paused {
		if taskObj.Status == sm.StatusProcessing {
			return nil, saError.Stack("任务处理中，不支持该操作")
		}

		in.PauseAt = saOrm.Now()
		if taskObj.DeleteAt.IsZero() && taskObj.PauseAt.IsZero() {
			err = saTask.Event(&saTask.EventRequest{
				Key:   in.Key,
				Event: "stop",
			})
			if err != nil {
				saLog.Err(saError.Stack(err))
				return nil, err
			}
		}
	} else
	//开启任务
	{
		if saTask.IsCaseExist(in.Key) {
			//先暂停
			err = saTask.Event(&saTask.EventRequest{Key: in.Key, Event: "stop"})
			if err != nil {
				return nil, saError.Stack(err)
			}

			//再开启
			err = saTask.Event(&saTask.EventRequest{
				Key:    in.Key,
				Event:  "start",
				Spec:   in.Spec,
				Params: in.Params,
			})
			if err != nil {
				return nil, saError.Stack(err)
			}
		} else {
			err = saTask.AddCase(&saTask.Case{
				Key:     in.Key,
				Spec:    in.Spec,
				Handler: task.Handle,
				Params:  in.Params,
			})
			if err != nil {
				return nil, saError.Stack(err)
			}
		}
	}

	//保存数据
	err = sm.MySql.Save(in).Error
	if err != nil {
		return nil, saError.New(err)
	}
	return &api.Response{Result: nil}, nil
}

func TaskOnce(c *api.Context) (res *api.Response, err error) {
	var in = new(dto.TaskEventRequest)
	c.BindByApi(in)
	if in.TaskId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	var obj = new(mysql.TblTask)
	err = sm.MySql.Table(mysql.TBNTask).Where("id = ?", in.TaskId).First(obj).Error
	if sm.MySql.IsError(err) {
		return nil, saError.Stack(err)
	}

	if obj.Id <= 0 {
		return nil, saError.Stack(saError.ErrNotExisted)
	}

	err = saTask.Event(&saTask.EventRequest{
		Key:    obj.Key,
		Event:  "once",
		Spec:   obj.Spec,
		Params: saHit.OrStr(in.Params, obj.Params),
	})

	if err != nil {
		return nil, saError.Stack(err)
	}
	return &api.Response{Result: nil}, nil
}

func TaskStepRun(c *api.Context) (res *api.Response, err error) {
	var in = new(dto.TaskStepRunRequest)
	c.BindByApi(in)
	if in.TaskId <= 0 || in.StepId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	//任务信息
	var taskObj = new(mysql.TblTask)
	sm.MySql.Raw(fmt.Sprintf("select * from %s where id = %d;", mysql.TBNTask, in.TaskId)).Scan(taskObj)
	if taskObj.Id <= 0 {
		return nil, saError.Stack(saError.ErrNotExisted)
	}

	if taskObj.Status == sm.StatusProcessing {
		return nil, saError.Stack("任务处理中，不允许操作")
	}

	//发起任务
	err = task.RunStep(taskObj, in.StepId, nil, 0, in.Params)
	if err != nil {
		return nil, err
	}
	return &api.Response{Result: nil}, nil
}

func TaskStepCallback(c *api.Context) (res *api.Response, err error) {
	var in = new(task.CallbackData)
	c.BindByApi(in)

	_ = task.Callback(in)
	return &api.Response{Result: nil}, nil
}
