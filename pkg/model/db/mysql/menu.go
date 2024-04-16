package mysql

import (
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNMenu = "sm_menu"

type TblMenu struct {
	Id        int64       `json:"id" form:"id"`
	SysId     int64       `json:"sysId" form:"sysId"`
	Type      sm.MenuType `json:"type" form:"type"`
	Status    sm.Status   `json:"status" form:"status"`
	Pid       int64       `json:"pid" form:"pid"` //上级菜单ID
	Seq       int         `json:"seq" form:"seq"` //序号，增序
	Title     string      `json:"title" form:"title"`
	Icon      string      `json:"icon" form:"icon"`
	Path      string      `json:"path" form:"path"`
	Hidden    bool        `json:"hidden" form:"hidden"`       //是否在菜单栏隐藏
	Component string      `json:"component" form:"component"` //页面文件路径
}

func (m *TblMenu) TableName() string {
	return TBNMenu
}
