package dto

import "github.com/saxon134/sysman/pkg/model/db/mysql"

type Menu struct {
	*mysql.TblMenu
	Buttons []*mysql.TblMenu `json:"buttons" form:"buttons"`
	Subs    []*Menu          `json:"subs" form:"subs"`
}
