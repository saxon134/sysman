package task

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/sysman/pkg/sm"
	"time"
)

func checkSign(sign, timestamp string) bool {
	if sm.Conf.Http.Secret == "" {
		return true
	}

	if sign == "" || timestamp == "" {
		return false
	}

	sign2 := saData.Md5(sm.Conf.Http.Secret+timestamp+sm.Conf.Http.Secret, true)
	return sign == sign2
}

func genSign() (sign string, timestamp string) {
	if sm.Conf.Http.Secret == "" {
		return
	}

	timestamp = saData.String(time.Now().Unix())
	sign = saData.Md5(sm.Conf.Http.Secret+timestamp+sm.Conf.Http.Secret, true)
	return sign, timestamp
}
