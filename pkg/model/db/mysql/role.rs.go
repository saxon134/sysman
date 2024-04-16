package mysql

import "github.com/saxon134/sysman/pkg/sm"

const TBNRoleRs = "sm_role_rs"

type TblRoleRs struct {
	Id     int64     `json:"id" form:"id"`
	SysId  int64     `json:"sysId" form:"sysId"`
	Status sm.Status `json:"status" form:"status"`
	UserId int64     `json:"userId" form:"userId"`
	RoleId int64     `json:"roleId" form:"roleId"`
}

func (m *TblRoleRs) TableName() string {
	return TBNRoleRs
}
