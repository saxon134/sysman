package task

import (
	"github.com/saxon134/go-utils/saData"
	"time"
)

func checkSign(sign string, timestamp string) bool {
	if _client.smSecret == "" {
		return true
	}

	if sign == "" || timestamp == "" {
		return false
	}

	sign2 := saData.Md5(_client.smSecret+timestamp+_client.smSecret, true)
	return sign == sign2
}

func genSign() (sign string, timestamp string) {
	if _client.smSecret == "" {
		return
	}

	timestamp = saData.String(time.Now().Unix())
	sign = saData.Md5(_client.smSecret+timestamp+_client.smSecret, true)
	return sign, timestamp
}
