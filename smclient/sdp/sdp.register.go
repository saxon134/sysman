package sdp

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

// Register 注册服务
// app、host、port: 服务信息
func (m *Client) Register(app string, host string, port int) {
	if host == "" || port <= 0 {
		saLog.Err("RPC register error: leak params")
		return
	}

	var params = map[string]interface{}{"app": app, "host": saUrl.QueryEncode(host), "port": port}
	var headers = map[string]string{}
	if m.secret != "" {
		var timestamp = saData.I64tos(time.Now().Unix())
		headers["timestamp"] = timestamp
		headers["sign"] = saData.Md5(m.secret+timestamp, true)
	}

	var pingParams map[string]interface{}
	pingParams = params
	pingParams["address"] = saUrl.ConnectUri(m.sysmainUrl, "sdp/ping")
	pingParams["secret"] = m.secret

	_, err := saHttp.PostRequest(saUrl.ConnectUri(m.sysmainUrl, "sdp/register"), params, headers)
	if err != nil {
		saLog.Err(saError.NewError(err))
		if m.redis != nil {
			//通过接口注册失败时，如果有Redis，则直接通过ping注册一下，避免要延迟一段时间服务才可用
			go m.ping(pingParams)
		}
	} else {
		saLog.Log("[SDP] " + app + " Register OK.")
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(m.pingSecond))
			go m.ping(pingParams)
		}
	}()
}

func (m *Client) ping(pingParams map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			saLog.Err("Sdp ping panic:", err)
		}
	}()

	if m.redis != nil {
		var key = fmt.Sprintf(RedisAppKey, pingParams["app"])
		var sdpAry = make([]*Instance, 0, 10)
		_ = m.redis.GetObj(key, &sdpAry)

		var existed *Instance
		var now = time.Now().UnixMilli()
		for _, c := range sdpAry {
			if c.Host == pingParams["host"] && c.Port == pingParams["port"] {
				existed = c
				break
			}
		}

		if existed == nil {
			sdpAry = append(sdpAry, &Instance{
				Host:   saData.String(pingParams["host"]),
				Port:   saData.Int(pingParams["port"]),
				Weight: 100,
				Time:   now,
			})
		} else {
			existed.Time = now
			existed.Weight = 100
		}
		_ = m.redis.SetObj(key, sdpAry, time.Second*time.Duration(m.pingSecond+2))
	} else {
		var params = map[string]string{"app": saData.String(pingParams["app"]), "host": saData.String(pingParams["host"]), "port": saData.String(pingParams["port"])}
		var secret = pingParams["secret"]
		if secret != "" {
			var timestamp = saData.I64tos(time.Now().Unix())
			params["timestamp"] = timestamp
			params["sign"] = saData.Md5(saData.String(secret)+timestamp, true)
		}

		_, err := saHttp.Get(saData.String(pingParams["address"]), params)
		if err != nil {
			saLog.Err("Sdp ping error:", err)
			return
		}
	}
}
