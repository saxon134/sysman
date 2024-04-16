package rbac

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saRedis"
)

type Client struct {
	redis *saRedis.Redis
}

func New(redis *saRedis.Redis) *Client {
	return &Client{redis: redis}
}

// TokenUser 通过token读取Redis中用户信息
// 通过该方法查询当前接口使用用户
func (c *Client) TokenUser(token string) (info *User, err error) {
	if token == "" {
		err = saError.ErrNotExisted
		return
	}

	if c.redis == nil {
		err = saError.ErrNotSupport
		return
	}

	info = new(User)
	err = c.redis.GetObj("sysman:tokens:"+token, info)
	if err != nil {
		return
	}

	if info.Id <= 0 {
		err = saError.ErrNotExisted
	}

	return info, nil
}
