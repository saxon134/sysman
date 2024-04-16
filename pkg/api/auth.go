package api

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/sysman/pkg/sm"
	"time"
)

type AuthType int8

const (
	AuthNull AuthType = iota
	AuthToken
	AuthSign
)

// AuthCheck  权限校验
func AuthCheck(c *Context, auths ...AuthType) (err error) {
	if auths == nil || len(auths) == 0 {
		return nil
	}

	for _, auth := range auths {
		if auth == AuthSign {
			var t = c.GetHeader("t")
			if t == "" {
				t = c.GetHeader("timestamp")
			}
			var ts = saData.Int64(t)
			var now = time.Now().Unix()
			if now-ts > 5 {
				return saError.Stack(saError.ErrUnauthorized)
			}

			var sign = c.GetHeader("sign")
			if sign == "" || ts <= 1685332856 {
				return saError.Stack(saError.ErrUnauthorized)
			}

			if saData.Md5(sm.Conf.Http.Secret+t, true) != sign {
				return saError.Stack(saError.ErrUnauthorized)
			}
		} else if auth == AuthToken {
			var token = c.GetHeader("auth")
			var uid = c.GetHeader("uid")
			err = sm.Redis.GetObj(fmt.Sprintf("sysman:tokens:%s:%s", uid, token), &c.Headers.User)
			if sm.Redis.IsError(err) {
				return saError.Stack(err)
			}

			if c.Headers.User.Id <= 0 {
				return saError.Stack(saError.ErrLoggedFail)
			}
		}
	}

	return nil
}
