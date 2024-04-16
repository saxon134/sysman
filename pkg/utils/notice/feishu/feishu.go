package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

type FeiShu struct {
	webhookUrl string
	secret     string
}

func New(webhookUrl string, secret string) *FeiShu {
	return &FeiShu{webhookUrl: webhookUrl, secret: secret}
}

func (m *FeiShu) SendTxt(txt string) {
	var args = map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": txt},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}

	var response = new(Response)
	var err = saHttp.Do(saHttp.Params{Method: "POST", Url: m.webhookUrl, Body: args, Header: map[string]interface{}{"Content-Type": "application/json"}}, response)
	if err != nil || response.Code != 0 {
		saLog.Err("飞书通知发送失败：", saHit.If(err != nil, err, response.Msg))
	}
}

func (m *FeiShu) SendTxtWithTitle(title string, txt string) {
	var contents = make([]*PostItem, 0, 1)
	contents = append(contents, &PostItem{Tag: "text", Text: txt})
	var args = map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   title,
					"content": [][]*PostItem{contents},
				},
			},
		},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}

	var response = new(Response)
	var err = saHttp.Do(saHttp.Params{Method: "POST", Url: m.webhookUrl, Body: args, Header: map[string]interface{}{"Content-Type": "application/json"}}, response)
	if err != nil || response.Code != 0 {
		saLog.Err("飞书通知发送失败：", saHit.If(err != nil, err, response.Msg))
	}
}

// SendPostTxt
// @Description: 发送富文本
func (m *FeiShu) SendPostTxt(title string, contents []*PostItem) {
	var args = map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   title,
					"content": [][]*PostItem{contents},
				},
			},
		},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}

	var response = new(Response)
	var err = saHttp.Do(saHttp.Params{Method: "POST", Url: m.webhookUrl, Body: args, Header: map[string]interface{}{"Content-Type": "application/json"}}, response)
	if err != nil || response.Code != 0 {
		saLog.Err("飞书通知发送失败：", saHit.If(err != nil, err, response.Msg))
	}
}

func sign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
