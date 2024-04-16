package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saLog"
	"github.com/saxon134/sysman/pkg/api"
	"github.com/saxon134/sysman/pkg/api/middleware"
	"github.com/saxon134/sysman/pkg/controller"
	"github.com/saxon134/sysman/pkg/sm"
	net "net/http"
	"time"
)

var net_server net.Server

// Init 阻塞
func Init() {
	gin.SetMode(gin.ReleaseMode)

	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(middleware.CrossDomain())

	//初始化路由
	{
		api.Gett(g.Group(""), "", controller.Index)
		initRoutes(g.Group(sm.Conf.Http.Root))
	}

	var port = saData.String(sm.Conf.Http.Port)
	saLog.Info("[HTTP] Listening on : " + port)

	net_server = net.Server{
		Addr:    ":" + port,
		Handler: g,
	}
	err := net_server.ListenAndServe()
	if err != nil && errors.Is(err, net.ErrServerClosed) {
		panic("[HTTP] error:" + err.Error())
	}
}

func Shutdown() {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_ = net_server.Shutdown(ctx)
}
