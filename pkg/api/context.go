package api

import (
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saLog"
	"strings"
)

/************************************/
/********** Context Model ********/
/************************************/

type Context struct {
	*gin.Context

	// 请求ID
	RequestId string `json:"requestId"`

	// header获取的数据
	Headers Headers `json:"headers"`

	// query中获取的分页、排序数据
	Paging Paging `json:"paging"`
	Sort   Sort   `json:"sort"`

	// 请求原始参数会保存在这里，bind只是获取想要的数据，可能会丢失一部分数据
	// 其他服务调用context方法需要设置参数也可以通过对其赋值进行传参
	// Bind的时候，优先使用该参数，无数据才会从请求中获取
	RawParam []byte `json:"-"`
}

type Headers struct {
	// token解析出的数据
	User User
}

type Paging struct {
	Limit  int //默认值为10，即便Valid为false，Limit也不会空
	Offset int
	Total  int64 //默认为0，前端如果传了后端就不需要再count了
	Valid  bool  //有些场景，不传分页参数，表示需要获取所有数据。具体业务代码控制
}

type Sort struct {
	Key  string `json:"key"`
	Desc bool   `json:"desc"`
}

type User struct {
	Id    int64  `json:"id"`
	Roles string `json:"roles"` //角色code，多个逗号分隔
}

type Response struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
	Ext    interface{} `json:"ext,omitempty"`
	Total  int64       `json:"total"` //分页时的总数
}

/************************************/
/********** Context Handler ********/
/************************************/

// BindByApi bind request params
func (ctx *Context) BindByApi(objPtr interface{}) {
	var err error
	if ctx.RawParam != nil {
		str, _ := saData.ToStr(ctx.RawParam)
		if str != "" {
			params, _ := objPtr.(map[string]interface{})
			if params != nil {
				err = saData.StrToModel(str, &params)
			} else {
				err = saData.StrToModel(str, objPtr)
			}
			if err == nil {
				return
			}
		}
	}

	if ctx.Request.Method == "GET" {
		if params, ok := objPtr.(map[string]interface{}); ok && params != nil {
			var values = ctx.Request.URL.Query()
			for k, v := range values {
				params[k] = v[0]
			}
		} else if params, ok := objPtr.(map[string]string); ok && params != nil {
			var values = ctx.Request.URL.Query()
			for k, v := range values {
				params[k] = v[0]
			}
		} else {
			err = ctx.ShouldBindQuery(objPtr)
			if err != nil {
				saLog.Err("gin bind params error:", err)
			}
		}
	} else if ctx.Request.Method == "POST" {
		var ctype = ctx.Request.Header.Get("Content-Type")
		ctype = saHit.Str(ctype != "", ctype, ctx.Request.Header.Get("content-type"))
		if strings.Contains(ctype, "application/json") {
			// 备份rawData 目的是为了能够多次bind，并且http.ResErr时打印传参
			rawData, _ := ctx.GetRawData()
			if len(rawData) > 0 {
				if params, ok := objPtr.(map[string]interface{}); ok && params != nil {
					err = saData.BytesToModel(rawData, &params)
				} else if params, ok := objPtr.(map[string]string); ok && params != nil {
					err = saData.BytesToModel(rawData, &params)
				} else {
					err = saData.BytesToModel(rawData, objPtr)
				}
				if err != nil {
					saLog.Err("api bind params error:", err)
				}
				ctx.RawParam = rawData
			}
		} else {
			err = ctx.ShouldBind(objPtr)
			if err != nil {
				saLog.Err("gin bind error:", err)
			}
		}
	} else {
		saLog.Err("GET/POST之外不支持")
	}
}

func (m Paging) IsAll() bool {
	return int64(m.Offset) >= m.Total
}
