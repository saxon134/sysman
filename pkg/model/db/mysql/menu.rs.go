package mysql

import (
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNMenuRs = "sm_menu_rs"

type TblMenuRs struct {
	Id     int64     `json:"id" form:"id"`
	SysId  int64     `json:"sysId" form:"sysId"`
	Status sm.Status `json:"status" form:"status"`
	RoleId int64     `json:"roleId" form:"roleId"`
	MenuId int64     `json:"menuId" form:"menuId"`
}

func (m *TblMenuRs) TableName() string {
	return TBNMenuRs
}
