package mysql

import "github.com/saxon134/sysman/pkg/sm"

const TBNSys = "sm_sys"

type TblSys struct {
	Id     int64     `json:"id" form:"id"`
	Code   string    `json:"code" form:"code"`
	Name   string    `json:"name" form:"name"`
	Desc   string    `json:"desc" form:"desc"`
	Tag    string    `json:"tag" form:"tag"` //角标
	Logo   string    `json:"logo" form:"logo"`
	Url    string    `json:"url" form:"url"` //首页网址
	Seq    int       `json:"seq" form:"seq"`
	Status sm.Status `json:"status" form:"status"`
}

func (m *TblSys) TableName() string {
	return TBNSys
}
