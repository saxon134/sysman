package router

import (
	"github.com/gin-gonic/gin"
	"github.com/saxon134/sysman/pkg/api"
	"github.com/saxon134/sysman/pkg/controller"
)

func initRoutes(group *gin.RouterGroup) {
	api.GetAndPost(group, "", controller.Index)

	api.Gett(group, "task/list", controller.TaskList, api.AuthToken)
	api.Gett(group, "task", controller.TaskDetail, api.AuthToken)
	api.Post(group, "task/save", controller.TaskSave, api.AuthToken)
	api.Post(group, "task/once", controller.TaskOnce, api.AuthToken)
	api.Gett(group, "task/steps", controller.TaskSteps, api.AuthToken)
	api.Post(group, "task/step/save", controller.TaskStepSave, api.AuthToken)
	api.Post(group, "task/step/run", controller.TaskStepRun, api.AuthToken)
	api.Gett(group, "task/log", controller.TaskLogList, api.AuthToken)
	api.Post(group, "task/callback", controller.TaskStepCallback, api.AuthSign)

	api.Post(group, "user/login", controller.UserLogin, api.AuthSign)
	api.Gett(group, "user/menus", controller.UserMenus, api.AuthToken)
	api.Post(group, "user/save", controller.UserSave, api.AuthToken)
	api.Gett(group, "user/list", controller.UserList, api.AuthToken)

	api.Gett(group, "sys/list", controller.SysList, api.AuthToken)
	api.Post(group, "sys/save", controller.SysSave, api.AuthToken)

	api.Gett(group, "menu/list", controller.MenuList, api.AuthToken)
	api.Post(group, "menu/save", controller.MenuSave, api.AuthToken)
	api.Post(group, "menu/topSeq.save", controller.MenuTopSeqSave, api.AuthToken)

}
