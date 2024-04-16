package mysql

import (
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNRole = "sm_role"

type TblRole struct {
	Id     int64     `json:"id" form:"id"`
	SysId  int64     `json:"sysId" form:"sysId"`
	Status sm.Status `json:"status" form:"status"`
	Code   string    `json:"code" form:"code"`
	Name   string    `json:"name" form:"name"`
}

func (m *TblRole) TableName() string {
	return TBNRole
}
