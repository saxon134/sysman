package controller

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/sysman/pkg/api"
	"github.com/saxon134/sysman/pkg/model/db/mysql"
	"github.com/saxon134/sysman/pkg/model/dto"
	"github.com/saxon134/sysman/pkg/sm"
	"sort"
	"time"
)

func UserMenus(ctx *api.Context) (res *api.Response, err error) {
	var params map[string]interface{}
	ctx.BindByApi(&params)

	var sysId = saData.Int64(params["sysId"])
	if sysId <= 0 {
		return nil, saError.Stack(saError.ErrParams)
	}

	//获取用户角色
	var roleAry = make([]*mysql.TblRole, 0, 20)
	err = sm.MySql.Raw(fmt.Sprintf(`
		select id from sm_role 
		where id in (
			select role_id from sm_role_rs where user_id = %d and status = 2
		)
		and sys_id = %d and status = 2
	`, ctx.Headers.User.Id, sysId)).Scan(&roleAry).Error
	if sm.MySql.IsError(err) {
		return nil, saError.Stack(err)
	}
	if len(roleAry) == 0 {
		return &api.Response{Result: nil}, nil
	}

	//ID
	var roleIdAry = make([]int64, len(roleAry))
	for _, v := range roleAry {
		roleIdAry = append(roleIdAry, v.Id)
	}

	//获取菜单
	var menus = make([]*mysql.TblMenu, 0, 100)
	err = sm.MySql.Raw(fmt.Sprintf(
		`select * from sm_menu where id in (select menu_id from sm_menu_rs where role_id in(%s)) and sys_id = %d and status = 2`,
		saData.AryToIds(roleIdAry, false), sysId,
	)).Scan(&menus).Error
	if sm.MySql.IsError(err) {
		err = saError.Stack(err)
		return
	}

	//整理菜单层级，最多2层
	var menuAry = make([]*dto.Menu, 0, 20)
	{
		sort.Slice(menus, func(i, j int) bool {
			return menus[i].Seq < menus[j].Seq
		})
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

func UserLogin(ctx *api.Context) (res *api.Response, err error) {
	var params map[string]string
	ctx.BindByApi(&params)
	if params["account"] == "" || params["password"] == "" {
		return nil, saError.Stack(saError.ErrParams)
	}

	var obj = new(mysql.TblUser)
	err = sm.MySql.Table(mysql.TBNUser).Where("account = ?", params["account"]).First(obj).Error
	if sm.MySql.IsError(err) {
		return nil, saError.Stack(err)
	}

	//登录校验
	{
		if obj.Id <= 0 {
			return nil, saError.Stack(saError.ErrNotExisted)
		}

		var password = saData.Md5(params["password"]+"sysman", true)
		if obj.Password != password {
			return nil, saError.Stack(saError.ErrPassword)
		}

		if obj.Status != sm.StatusEnable {
			return nil, saError.Stack(saError.ErrUnauthorized)
		}
	}

	//获取系统&角色
	var roles string
	var sysAry []*mysql.TblSys
	{
		//获取用户角色
		var roleAry = make([]*mysql.TblRole, 0, 20)
		if obj.Account != "admin" {
			err = sm.MySql.Raw(fmt.Sprintf(`
			select id, code, sys_id, name from sm_role 
			where id in (
				select role_id from sm_role_rs where user_id = %d and status = 2
			)
			and status = 2`, obj.Id,
			)).Scan(&roleAry).Error
			if sm.MySql.IsError(err) {
				err = saError.Stack(err)
				return
			}
			if len(roleAry) == 0 {
				return nil, saError.Stack(saError.ErrUnauthorized)
			}
		}

		//sys & role ID
		var sysIdAry = make([]int64, 0, 10)
		var roleIdAry = make([]int64, len(roleAry))
		for _, v := range roleAry {
			sysIdAry = saData.AppendId(sysIdAry, v.SysId)
			roleIdAry = append(roleIdAry, v.Id)
		}
		roles = saData.AryToIds(roleIdAry, false)
		if obj.Account == "admin" {
			roles = "admin"
		}

		//获取系统
		sysAry = make([]*mysql.TblSys, 0, len(sysIdAry))
		if obj.Account == "admin" {
			err = sm.MySql.Raw(`select * from sm_sys where id in ? and status = 2`, sysIdAry).Scan(&sysAry).Error
		} else {
			err = sm.MySql.Raw(`select * from sm_sys`).Scan(&sysAry).Error
		}
		if sm.MySql.IsError(err) {
			err = saError.Stack(err)
			return
		}
	}

	//保存token
	var token = saData.Md5(saData.String(obj.Id)+saData.String(time.Now().UnixMilli()), true)
	var tokenUser = api.User{Id: obj.Id, Roles: roles}
	err = sm.Redis.Set(fmt.Sprintf("sysman:tokens:%d:%s", obj.Id, token), saData.String(tokenUser), time.Hour*24)
	if err != nil {
		return nil, saError.Stack(err)
	}

	return &api.Response{Result: map[string]interface{}{
		"id":      obj.Id,
		"name":    obj.Name,
		"account": obj.Account,
		"token":   token,
		"sysAry":  sysAry,
		"roles":   roles,
	}}, nil
}

func UserSave(ctx *api.Context) (res *api.Response, err error) {
	var obj = new(mysql.TblUser)
	ctx.BindByApi(obj)
	if obj.Account == "" || obj.Name == "" {
		return nil, saError.Stack(saError.ErrParams)
	}

	var existed = new(mysql.TblUser)
	err = sm.MySql.Table(mysql.TBNUser).Where("account = ?", obj.Account).First(existed).Error
	if sm.MySql.IsError(err) {
		return nil, saError.Stack(err)
	}

	//新建
	if obj.Id <= 0 {
		if obj.Password == "" {
			return nil, saError.Stack(saError.ErrParams)
		}

		if existed.Id > 0 {
			return nil, saError.Stack(saError.ErrExisted)
		}

		obj.Status = sm.Status(saHit.Int(obj.Status == sm.StatusEnable, sm.StatusEnable.Int(), sm.StatusDisable.Int()))
		obj.Password = saData.Md5(obj.Password+"sysman", true)
	} else
	//保存信息
	{
		if existed.Id != obj.Id {
			return nil, saError.New("账户不匹配")
		}

		if obj.Password == "" {
			obj.Password = existed.Password
		} else {
			obj.Password = saData.Md5(obj.Password+"sysman", true)
		}
	}

	err = sm.MySql.Save(obj).Error
	if err != nil {
		return nil, saError.Stack(err)
	}
	return &api.Response{Result: nil}, nil
}

func UserList(ctx *api.Context) (res *api.Response, err error) {
	var in = new(dto.UserListRequest)
	ctx.BindByApi(in)

	var sql = `from sm_user t1 left join sm_role_rs t2 on t2.user_id = t1.id where 1= 1 `

	if in.Account != "" {
		sql += `and t1.account like '%` + in.Account + `%' `
	}

	if in.Name != "" {
		sql += `and t1.name like '%` + in.Name + `%' `
	}

	if in.RoleId > 0 {
		sql += fmt.Sprintf("and t2.role_id = %d and t2.status = 2 ", in.RoleId)
	}

	if in.Status != sm.StatusNull {
		sql += fmt.Sprintf("and t1.status = %d ", in.Status)
	}

	err = saError.Stack(sm.MySql.Raw("select count(t1.*) " + sql).Scan(&ctx.Paging.Total).Error)
	if err != nil {
		return nil, err
	}

	var userAry = make([]*mysql.TblUser, 0, ctx.Paging.Limit)
	err = saError.Stack(
		sm.MySql.Raw(fmt.Sprintf("select t1.* %s limit %d %d; ", sql, ctx.Paging.Offset, ctx.Paging.Limit)).Scan(&userAry).Error,
	)
	if err != nil {
		return nil, err
	}

	return &api.Response{Result: userAry, Total: ctx.Paging.Total}, nil
}
