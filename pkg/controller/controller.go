package controller

import (
	"github.com/saxon134/sysman/pkg/api"
)

// Index handler，避免每个项目都搞一个
func Index(c *api.Context) (res *api.Response, err error) {
	return &api.Response{Result: "hi"}, nil
}

//func checkSign(r *http.Request) bool {
//	if conf.Conf.Http.Secret == "" {
//		return true
//	}
//
//	var sign = r.Header.Get("sign")
//	var timestamp = r.Header.Get("timestamp")
//	if sign == "" || timestamp == "" {
//		return false
//	}
//
//	sign2 := saData.Md5(conf.Conf.Http.Secret+timestamp, true)
//	return sign == sign2
//}
//
//func genSign() (sign string, timestamp string) {
//	if conf.Conf.Http.Secret == "" {
//		return
//	}
//
//	timestamp = saData.String(time.Now().Unix())
//	sign = saData.Md5(conf.Conf.Http.Secret+timestamp, true)
//	return sign, timestamp
//}
