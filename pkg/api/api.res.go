package api

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saLog"
)

// 返回正确时的数据组装
func resSuccess(c *Context, v *Response) {
	if c == nil || c.IsAborted() {
		return
	}

	if v == nil {
		v = new(Response)
	}

	var dic = map[string]interface{}{"code": 0}
	if v.Result != nil {
		dic["result"] = v.Result
	}

	if v.Ext != nil {
		dic["ext"] = v.Ext
	}

	if c.Paging.Valid == true {
		v.Total = saHit.Int64(v.Total > 0, v.Total, c.Paging.Total)
		v.Total = saHit.Int64(v.Total > 0, v.Total, int64(c.Paging.Offset+c.Paging.Limit))
		dic["hasNext"] = v.Total > int64(c.Paging.Offset+c.Paging.Limit)
		dic["totalCount"] = v.Total
	}

	c.JSON(200, dic)
	c.Abort()
}

// 返回error时的数据组装
func resErr(c *Context, err interface{}) {
	if c == nil || c.IsAborted() {
		return
	}

	var msg = "接口报错"
	var code = saError.SensitiveErrorCode
	var caller = ""
	var errLog = ""

	if s, ok := err.(string); ok {
		code = saError.NormalErrorCode
		if s != "" {
			msg = s
		}
	} else if e, ok := err.(saError.Error); ok {
		code = e.Code
		msg = e.Msg
		caller = e.Caller
	} else if e, ok := err.(*saError.Error); ok {
		code = e.Code
		msg = e.Msg
		caller = e.Caller
	} else if e, ok := err.(error); ok {
		err_s := e.Error()
		var dic map[string]interface{}
		_ = saData.StrToModel(err_s, &dic)

		//如果是saError，则错误码、信息按saError输出；否则全部按照敏感信息处理
		sa_err := new(saError.Error)
		if saData.StrToModel(err_s, sa_err) == nil && sa_err.Code > 0 {
			code = sa_err.Code
			msg = sa_err.Msg
			caller = sa_err.Caller
		} else {
			code = saError.SensitiveErrorCode
			msg = err_s
		}
	}

	//请求日志
	errLog = fmt.Sprintf("[ERR] %s %s %d", c.Request.Method, c.Request.URL.Path, code)
	errLog += "\n[MSG] " + msg
	{
		if c.Request.URL.RawQuery != "" {
			errLog += "\n[QUE] " + c.Request.URL.RawQuery
		}

		if c.Request.Method == "POST" {
			var raw = saData.String(c.RawParam)
			if raw != "" {
				errLog += "\n[BOD] " + raw
			}
		}
	}
	saLog.Err(saError.Error{
		Msg:    errLog,
		Caller: caller,
	})

	if code != 0 {
		rsp_v := map[string]interface{}{"code": code}
		rsp_v["msg"] = msg

		c.JSON(400, rsp_v)
		c.Abort()
		return
	}

	//异常情况
	c.JSON(500, &map[string]interface{}{"code": saError.NormalErrorCode, "msg": "服务器开了个小差"})
	c.Abort()
}
