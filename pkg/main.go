package main

import (
	"github.com/saxon134/sysman/pkg/router"
	"github.com/saxon134/sysman/pkg/sdp"
	"github.com/saxon134/sysman/pkg/sm"
	"github.com/saxon134/sysman/pkg/task"
)

func main() {
	//初始化
	sm.Init()

	//定时任务
	task.Init()

	//初始化http服务
	go router.Init()

	//初始化SDP
	sdp.Init()

	//防止应用退出
	<-make(chan bool)
}
