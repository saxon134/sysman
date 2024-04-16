package sdp

import (
	"github.com/saxon134/go-utils/saLog"
	"github.com/saxon134/go-utils/saRedis"
)

type Client struct {
	sysmainUrl string
	clientRoot string
	secret     string
	pingSecond int
	redis      *saRedis.Redis
}

type Conf struct {
	SysmainUrl string //sysmain地址
	ClientRoot string //client跟路由
	Secret     string //接口秘钥
	PingSecond int    //ping间隔
	Redis      struct {
		Host string
		Pass string
	} //redis非空，则ping/discovery就都会使用Redis，无Redis配置则通过接口获取
}

func NewClient(conf Conf) *Client {
	if conf.PingSecond <= 0 {
		conf.PingSecond = 5 //默认5秒
	} else if conf.PingSecond == 1 {
		conf.PingSecond = 2 //最小2秒
	}

	var client = &Client{
		sysmainUrl: conf.SysmainUrl,
		clientRoot: conf.ClientRoot,
		secret:     conf.Secret,
		pingSecond: conf.PingSecond,
	}

	if conf.Redis.Host != "" {
		var err error
		client.redis, err = saRedis.Init(conf.Redis.Host, conf.Redis.Pass, 0)
		if err != nil {
			saLog.Err("SDP Client Redis初始化出错：" + err.Error())
			client.redis = nil
		}

		_, err = client.redis.Pool.Dial()
		if err != nil {
			saLog.Err("SDP Client Redis初始化出错" + err.Error())
			client.redis = nil
		}
	}

	return client
}
