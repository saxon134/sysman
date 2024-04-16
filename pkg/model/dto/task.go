package dto

import "github.com/saxon134/sysman/pkg/model/db/mysql"

type TaskListRequest struct {
	SysId int64  `json:"sysId" form:"sysId"`
	Name  string `json:"name" form:"name"`
}

type TaskItem struct {
	*mysql.TblTask
	Paused bool                 `json:"paused" form:"paused"`
	Steps  []*mysql.TblTaskStep `json:"steps" form:"steps"`
}

type TaskLogRequest struct {
	TaskId int64 `json:"taskId" form:"taskId"`
	LogId  int64 `json:"logId" form:"logId"`
}

type TaskLogItem struct {
	*mysql.TblTaskLog
	Step *mysql.TblTaskStep `json:"step" form:"step"`
}

type TaskDetail struct {
	*mysql.TblTask
	Sys      string `json:"sys"`
	Paused   bool   `json:"paused"`
	NextTime string `json:"nextTime"`
	PreTime  string `json:"preTime"`
}

type TaskEventRequest struct {
	TaskId int64
	Params string
}

type TaskStepRunRequest struct {
	TaskId int64
	StepId int64
	Params string
}

type TaskStepsSaveRequest struct {
	TaskId int64   `json:"taskId" form:"taskId"`
	Steps  []*Step `json:"steps" form:"steps"`
}

type Step struct {
	*mysql.TblTaskStep
	Paused bool `json:"paused" form:"paused"`
}
