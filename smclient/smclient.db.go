package smclient

import (
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/go-utils/saRedis"
)

var redis *saRedis.Redis
var mySql *saOrm.DB
