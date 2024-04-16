package api

import (
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saHit"
	"strings"
)

type Handle func(c *Context) (res *Response, err error)

func Get(group *gin.RouterGroup, path string, handler Handle, auth ...AuthType) {
	group.GET(path, func(c *gin.Context) {
		get(c, handler, auth...)
	})
}

// Gett 仅为了对齐好看点
func Gett(group *gin.RouterGroup, path string, handler Handle, auth ...AuthType) {
	Get(group, path, handler, auth...)
}

func Post(group *gin.RouterGroup, path string, handler Handle, auth ...AuthType) {
	group.POST(path, func(c *gin.Context) {
		post(c, handler, auth...)
	})
}

func GetAndPost(group *gin.RouterGroup, path string, handler Handle, auth ...AuthType) {
	group.POST(path, func(c *gin.Context) {
		post(c, handler, auth...)
	})

	group.GET(path, func(c *gin.Context) {
		get(c, handler, auth...)
	})
}

func get(gc *gin.Context, handle Handle, auth ...AuthType) {
	var err error
	var c = &Context{
		Context: gc,
	}

	//权限校验
	if err = AuthCheck(c, auth...); err != nil {
		resErr(c, err)
		return
	}

	//分页
	c.Paging.Limit = 20
	c.Paging.Offset = 0
	c.Paging.Valid = false

	var limit = 0
	{
		if s, ok := gc.GetQuery("pageSize"); ok {
			if v, _ := saData.ToInt(s); v > 0 {
				limit = v
			}
		}
		if limit == 0 {
			if s, ok := gc.GetQuery("size"); ok {
				if v, _ := saData.ToInt(s); v > 0 {
					limit = v
				}
			}
		}
		if limit == 0 {
			if s, ok := gc.GetQuery("limit"); ok {
				if v, _ := saData.ToInt(s); v > 0 {
					limit = v
				}
			}
		}

		if limit > 0 {
			c.Paging.Valid = true
			c.Paging.Limit = limit
		} else {
			c.Paging.Valid = false
			c.Paging.Limit = 20
		}
	}

	var offset, page = 0, 0
	if c.Paging.Valid {
		s, _ := gc.GetQuery("pageNumber")
		if page, _ = saData.ToInt(s); page > 0 {
			offset = c.Paging.Limit * (page - 1)
		}

		if offset <= 0 {
			s, _ = gc.GetQuery("page")
			if page, _ = saData.ToInt(s); page > 0 {
				offset = c.Paging.Limit * (page - 1)
			}
		}

		if offset <= 0 {
			s, _ = gc.GetQuery("current")
			if page, _ = saData.ToInt(s); page > 0 {
				offset = c.Paging.Limit * (page - 1)
			}
		}

		if offset <= 0 {
			s, _ = gc.GetQuery("offset")
			offset, _ = saData.ToInt(s)
		}

		c.Paging.Offset = saHit.Int(offset >= 0, offset, 0)

		// 总数，一般第一次后端count，之后前端传
		s, _ = gc.GetQuery("total")
		c.Paging.Total, _ = saData.ToInt64(s)
	}

	//排序
	if s, _ := gc.GetQuery("sort"); len(s) > 0 {
		ary := strings.Split(s, "__")
		c.Sort.Key = saData.SnakeStr(ary[0])
		c.Sort.Desc = len(ary) == 2 && ary[1] == "desc"
	}

	//执行handle函数
	var res *Response
	res, err = handle(c)
	if err != nil {
		resErr(c, err)
		return
	}
	resSuccess(c, res)
}

func post(gc *gin.Context, handle Handle, auth ...AuthType) {
	var err error

	c := &Context{Context: gc}

	//权限校验
	if err = AuthCheck(c, auth...); err != nil {
		resErr(c, err)
		return
	}

	//执行handle函数
	var res *Response
	res, err = handle(c)
	if err != nil {
		resErr(c, err)
		return
	}
	resSuccess(c, res)
}
