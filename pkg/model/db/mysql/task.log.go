package mysql

import (
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNTaskLog = "sm_task_log"

type TblTaskLog struct {
	Id        int64       `json:"id" form:"id"`
	Pid       int64       `json:"pid" form:"pid"`
	Specified bool        `json:"specified" form:"specified"` //是否是指定步骤执行的
	TaskId    int64       `json:"taskId" form:"taskId"`
	StepId    int64       `json:"stepId" form:"stepId"`
	Status    sm.Status   `json:"status" form:"status"`
	Params    string      `json:"params" form:"params"` //执行任务时参数
	Result    string      `json:"result" form:"result"`
	Err       string      `json:"err" form:"err"`
	CreateAt  *saOrm.Time `json:"createAt" form:"createAt"`
	DoneAt    *saOrm.Time `json:"doneAt" form:"doneAt"`
	Ms        int64       `json:"ms" form:"ms"` //执行消耗时间(毫秒)
}

func (m *TblTaskLog) TableName() string {
	return TBNTaskLog
}
