package controller

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/saxon134/go-utils/saData"
//	"github.com/saxon134/go-utils/saData/saError"
//	"github.com/saxon134/sysman/pkg/api"
//	"github.com/saxon134/sysman/pkg/db"
//	"github.com/saxon134/sysman/pkg/sdp"
//	"net/http"
//)
//
//func SdpRegisterHandler(w http.ResponseWriter, r *http.Request) {
//	var query = r.URL.Query()
//	if checkSign(r) == false {
//		api.ResError(w, saError.ErrUnauthorized.Error())
//		return
//	}
//
//	var in = sdp.Request{}
//	decoder := json.NewDecoder(r.Body)
//	_ = decoder.Decode(&in)
//
//	//兼容GET方法
//	if in.App == "" {
//		in.App = query.Get("app")
//	}
//	if in.Host == "" {
//		in.Host = query.Get("host")
//	}
//	if in.Port <= 0 {
//		in.Port = saData.Int(query.Get("port"))
//	}
//	if in.Cpu <= 0 {
//		in.Cpu, _ = saData.ToFloat32(query.Get("cpu"))
//	}
//	if in.Memo <= 0 {
//		in.Memo, _ = saData.ToFloat32(query.Get("memo"))
//	}
//	if in.App == "" || in.Host == "" || in.Port <= 0 {
//		api.ResErr(w, saError.Stack(saError.ErrParams).Error())
//		return
//	}
//
//	//注册app服务
//	sdp.Chan <- in
//	api.Res(w, nil)
//}
//
//func SdpPingHandler(w http.ResponseWriter, r *http.Request) {
//	SdpRegisterHandler(w, r)
//}
//
//func SdpDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
//	var query = r.URL.Query()
//	if checkSign(r) == false {
//		api.ResError(w, saError.ErrUnauthorized.Error())
//		return
//	}
//
//	var app = query.Get("app")
//	if app == "" {
//		api.ResError(w, saError.Stack(saError.ErrParams).Error())
//		return
//	}
//
//	var sdpAry = make([]*sdp.Instance, 0, 10)
//	var key = fmt.Sprintf(sdp.RedisAppKey, app)
//	_ = db.Redis.GetObj(key, &sdpAry)
//	api.Res(w, sdpAry)
//}
