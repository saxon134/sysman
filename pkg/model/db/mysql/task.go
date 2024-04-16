package mysql

import (
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNTask = "sm_task"

type TblTask struct {
	Id       int64       `json:"id" form:"id"`
	SysId    int64       `json:"sysId" form:"sysId"`
	Name     string      `json:"name" form:"name"`     //任务名称，后台维护
	Key      string      `json:"key" form:"key"`       //任务的唯一标识，不可重复
	Status   sm.Status   `json:"status" form:"status"` //0-信息不全 1-处理中 2-成功 -1-失败
	Spec     string      `json:"spec" form:"spec"`     //执行周期
	Params   string      `json:"params" form:"params"` //执行任务时参数
	Remark   string      `json:"remark" form:"remark"`
	PauseAt  *saOrm.Time `json:"pauseAt" form:"pauseAt"`
	DeleteAt *saOrm.Time `json:"deleteAt" form:"deleteAt"`
}

func (m *TblTask) TableName() string {
	return TBNTask
}
