package controller

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/sysman/pkg/api"
	"github.com/saxon134/sysman/pkg/model/db/mysql"
	"github.com/saxon134/sysman/pkg/model/dto"
	"github.com/saxon134/sysman/pkg/sm"
)

func MenuList(c *api.Context) (res *api.Response, err error) {
	var sysId = saData.Int64(c.Query("sysId"))
	if sysId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	var menus = make([]*mysql.TblMenu, 0, 50)
	err = saError.Stack(sm.MySql.Table(mysql.TBNMenu).Where("sys_id = ?", sysId).Scan(&menus))
	if err != nil {
		return nil, err
	}

	//组装菜单层级
	var menuAry = make([]*dto.Menu, 0, 20)
	{
		for _, v := range menus {
			if v.Pid == 0 && v.Type == sm.Menu {
				var m = &dto.Menu{TblMenu: v}

				//查找子菜单
				m.Subs = make([]*dto.Menu, 0, 5)
				m.Buttons = make([]*mysql.TblMenu, 0, 5)
				for _, v2 := range menus {
					if v2.Pid == v.Id {
						if v2.Type == sm.Menu {
							m.Subs = append(m.Subs, &dto.Menu{TblMenu: v2})
						} else {
							m.Buttons = append(m.Buttons, v2)
						}
					}
				}
				menuAry = append(menuAry, m)
			}
		}
	}
	return &api.Response{Result: menuAry}, nil
}

func MenuSave(c *api.Context) (res *api.Response, err error) {
	var obj = new(dto.Menu)
	c.BindByApi(obj)
	if obj.Type == sm.MenuTypeNull || obj.SysId <= 0 || obj.Title == "" || obj.Path == "" || obj.Component == "" {
		return nil, saError.Stack(saError.ErrParams)
	}

	err = saError.Stack(sm.MySql.Save(obj).Error)
	if err != nil {
		return nil, err
	}
	return &api.Response{Result: nil}, nil
}

func MenuTopSeqSave(c *api.Context) (res *api.Response, err error) {
	var request = map[int64]int{}
	err = saError.Stack(c.Bind(&request))
	if err != nil {
		return nil, err
	}

	var sql = ""
	for menuId, seq := range request {
		sql += fmt.Sprintf("update sm_menu set seq = %d where id = %d and pid = 0 and type = %d;\n", seq, menuId, sm.Menu.Int())
	}

	if sql == "" {
		return &api.Response{Result: nil}, nil
	}

	err = saError.Stack(sm.MySql.Exec(sql).Error)
	if err != nil {
		return nil, err
	}

	return &api.Response{Result: nil}, nil
}
