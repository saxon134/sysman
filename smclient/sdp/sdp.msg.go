package sdp

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"time"
)

// SendMsg
// app: 必选，指定应用
// host/port: 可选，指定服务器；空则向所有实例的应用发送消息
func (m *Client) SendMsg(targetApp string, targetHost string, targetPort int, msg string) (err error) {
	if targetApp == "" {
		return saError.NewError("必须指定应用")
	}

	if msg == "" {
		return saError.NewError("消息不能空")
	}

	var sdpAry = m.getSdpAry(targetApp)
	if sdpAry == nil || len(sdpAry) == 0 {
		return saError.NewError("无可用应用")
	}

	var headers = map[string]string{}
	if m.secret != "" {
		var timestamp = saData.String(time.Now().Unix())
		var sign = saData.Md5(m.secret+timestamp, true)
		headers["sign"] = sign
		headers["timestamp"] = timestamp
	}

	var sendCnt = 0
	if targetHost != "" {
		for _, v := range sdpAry {
			if v.Host == targetHost && v.Port == targetPort {
				var url = saUrl.ConnectUri("http://"+v.Host+":"+saData.String(v.Port), m.clientRoot, "sdp/msg")
				_, err = saHttp.PostRequest(url, map[string]string{"msg": msg}, headers)
				sendCnt++
				return err
			}
		}
	} else {
		for _, v := range sdpAry {
			if v.Host != "" {
				var url = saUrl.ConnectUri("http://"+v.Host+":"+saData.String(v.Port), m.clientRoot, "sdp/msg")
				_, err = saHttp.PostRequest(url, map[string]string{"msg": msg}, headers)
				sendCnt++
				if err != nil {
					return err
				}
			}
		}
	}

	if sendCnt == 0 {
		return saError.NewError("未找到可发送消息的应用")
	}
	return nil
}
