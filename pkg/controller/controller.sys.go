package controller

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/sysman/pkg/api"
	"github.com/saxon134/sysman/pkg/model/db/mysql"
	"github.com/saxon134/sysman/pkg/sm"
)

func SysList(ctx *api.Context) (res *api.Response, err error) {
	var sysAry = make([]*mysql.TblSys, 0, 10)
	err = saError.Stack(sm.MySql.Raw("select * from sm_sys where status = 2 order by seq asc").Scan(&sysAry).Error)
	if err != nil {
		return nil, err
	}
	return &api.Response{Result: sysAry}, nil
}

func SysSave(ctx *api.Context) (res *api.Response, err error) {
	var obj = new(mysql.TblSys)
	ctx.BindByApi(obj)
	if obj.Name == "" || obj.Code == "" {
		return nil, saError.Stack(saError.ErrParams)
	}

	err = saError.Stack(sm.MySql.Save(obj).Error)
	if err != nil {
		return nil, err
	}
	return &api.Response{Result: nil}, nil
}
