package dto

import "github.com/saxon134/sysman/pkg/sm"

type UserListRequest struct {
	Status  sm.Status `form:"status"`
	Account string    `form:"account"`
	Name    string    `form:"name"`
	RoleId  int64     `form:"roleId"`
}
