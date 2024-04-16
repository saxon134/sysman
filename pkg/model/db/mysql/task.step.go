package mysql

import (
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNTaskStep = "sm_task_step"

type TblTaskStep struct {
	Id             int64       `json:"id" form:"id"`
	Name           string      `json:"name" form:"name"` //步骤名称，后台维护
	TaskId         int64       `json:"taskId" form:"taskId"`
	Seq            int         `json:"seq" form:"seq"`       //顺序
	Status         sm.Status   `json:"status" form:"status"` //0-信息不全 1-处理中 2-成功 -1-失败
	Type           int         `json:"type" form:"type"`     //1-同步 2-异步
	Key            string      `json:"key" form:"key"`
	Url            string      `json:"url" form:"url"`                       //https://www.com/link
	Params         string      `json:"params" form:"params"`                 //执行任务时参数
	DelaySecond    int         `json:"delaySecond" form:"delaySecond"`       //延迟执行
	RelyPreSuccess int         `json:"relyPreSuccess" form:"relyPreSuccess"` //是否依赖上一步执行成功
	PauseAt        *saOrm.Time `json:"pauseAt" form:"pauseAt"`
	DeleteAt       *saOrm.Time `json:"deleteAt" form:"deleteAt"`
}

func (m *TblTaskStep) TableName() string {
	return TBNTaskStep
}
