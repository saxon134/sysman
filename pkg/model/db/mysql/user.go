package mysql

import (
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/sysman/pkg/sm"
)

const TBNUser = "sm_user"

type TblUser struct {
	Id          int64       `json:"id" form:"id"`
	Name        string      `json:"name" form:"name"`
	Account     string      `json:"account" form:"account"`
	Password    string      `json:"password" form:"password"`
	Status      sm.Status   `json:"status" form:"status"`
	LastLoginAt *saOrm.Time `json:"lastLoginAt" form:"lastLoginAt"`
}

func (m *TblUser) TableName() string {
	return TBNUser
}
