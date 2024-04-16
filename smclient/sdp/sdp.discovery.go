package sdp

import (
	"fmt"
	"github.com/saxon134/go-utils/saCache"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"math/rand"
	"strings"
	"time"
)

var _last map[string]string

// Discovery 查找可用服务
func (m *Client) Discovery(app string) (host string, port int) {
	if app == "" {
		return "", 0
	}

	var last = ""
	if _last == nil {
		_last = map[string]string{}
	}
	last = _last[app]

	var sdpAry = m.getSdpAry(app)
	var lastIsOk = false //尽量选择一个跟上次的不一样的实例
	if sdpAry != nil && len(sdpAry) > 0 {
		var weight = 0
		for _, v := range sdpAry {
			var link = fmt.Sprintf("%s:%d", v.Host, v.Port)
			if last != link {
				weight += v.Weight
			} else {
				lastIsOk = true
			}
		}

		if weight > 0 {
			var r = rand.Intn(weight)
			for _, v := range sdpAry {
				var link = fmt.Sprintf("%s:%d", v.Host, v.Port)
				if last != link {
					if weight >= r {
						_last[app] = fmt.Sprintf("%s:%d", v.Host, v.Port)
						return saUrl.QueryDecode(v.Host), v.Port
					}
					weight += v.Weight
				}
			}
		}

		if lastIsOk {
			var ary = strings.Split(last, ":")
			if len(ary) == 2 {
				return saUrl.QueryDecode(ary[0]), saData.Int(ary[1])
			}
		}
	}
	return "", 0
}

func (m *Client) getSdpAry(app string) []*Instance {
	var cacheTime = time.Minute
	if m.redis != nil {
		cacheTime = time.Second * time.Duration((m.pingSecond + 3))
	}
	value, _ := saCache.MSetWithFunc("discoveryApp:"+app, cacheTime, func() (interface{}, error) {
		if m.redis != nil {
			var sdpAry = make([]*Instance, 0, 10)
			var key = fmt.Sprintf(RedisAppKey, app)
			err := m.redis.GetObj(key, &sdpAry)
			if m.redis.IsError(err) {
				return nil, err
			}
			return sdpAry, nil
		} else {
			var params = map[string]string{"app": app}
			if m.secret != "" {
				var timestamp = saData.I64tos(time.Now().Unix())
				params["timestamp"] = timestamp
				params["sign"] = saData.Md5(m.secret+timestamp, true)
			}

			res, err := saHttp.Get(saUrl.ConnectUri(m.sysmainUrl, "sdp/discovery"), params)
			if err != nil {
				saLog.Err("discovery error:", err)
				return nil, err
			}

			var resObj = struct {
				Result []*Instance
			}{Result: make([]*Instance, 0, 5)}
			_ = saData.StrToModel(res, &resObj)
			if len(resObj.Result) == 0 {
				return nil, saError.Stack(saError.ErrNotExisted)
			}
			return resObj.Result, err
		}
	})
	if v, ok := value.([]*Instance); ok {
		return v
	}
	return nil
}
