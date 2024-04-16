package smclient

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saLog"
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/go-utils/saRedis"
	"github.com/saxon134/sysman/smclient/rbac"
	"github.com/saxon134/sysman/smclient/sdp"
	"github.com/saxon134/sysman/smclient/task"
)

// RoleManager 角色权限管理
var RoleManager *rbac.Client

// SDiscovery 服务发现
var SDiscovery *sdp.Client

// Task 定时任务
var Task *task.Client

func Init(config Config) (err error) {

	// Redis、MySql必须配置
	if config.Redis.Host == "" || config.MySql.Host == "" {
		err = saError.New(saError.ErrParams)

		return
	}

	//初始化Redis数据库
	{
		redis, err = saRedis.Init(config.Redis.Host, config.Redis.Pass, config.Redis.Db)
		if err != nil {
			panic("Redis初始化出错" + err.Error())
		}

		_, err = redis.Pool.Dial()
		if err != nil {
			panic("Redis初始化出错" + err.Error())
		}

		saLog.Log("smclient Redis init ok")
	}

	//初始化MySql数据库
	{
		var dns = config.MySql.User + ":" + config.MySql.Pass + "@tcp(" + config.MySql.Host + ")/" + config.MySql.Db + "?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local"
		mySql = saOrm.Open(dns, saOrm.Conf{})
		saLog.Log("smclient MySql init ok")
	}

	//初始化 rbac
	RoleManager = rbac.New(redis)

	//初始化task
	Task = task.Init(config.Sysman.Root, config.Sysman.Secret)

	//初始化 sdp

	return nil
}
