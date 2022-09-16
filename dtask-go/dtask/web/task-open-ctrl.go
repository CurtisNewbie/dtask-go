package web

import (
	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

// Register OPEN API routes
func RegisterTaskOpenRoutes(router *gin.Engine) {
	router.POST(server.ResolvePath("/task/list", true), ListTaskByPageEndpoint)
	router.POST(server.ResolvePath("/task/history", true), ListTaskHistoryByPageEndpoint)
	router.POST(server.ResolvePath("/task/update", true), UpdateTaskEndpoint)
	router.POST(server.ResolvePath("/task/trigger", true), TriggerTaskEndpoint)
}

// List tasks
func ListTaskByPageEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.ListTaskByPageReqWebVo
	util.MustBindJson(c, &req)

	r, e := domain.ListTaskByPage(user, &req)
	if e != nil {
		panic(e)
	}
	util.DispatchOkWData(c, r)
}

// List task histories
func ListTaskHistoryByPageEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.ListTaskHistoryByPageReq
	util.MustBindJson(c, &req)

	r, e := domain.ListTaskHistoryByPage(user, &req)
	if e != nil {
		panic(e)
	}
	util.DispatchOkWData(c, r)
}

// Update task info
func UpdateTaskEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.UpdateTaskReq
	util.MustBindJson(c, &req)

	e := domain.UpdateTask(user, &req)
	if e != nil {
		panic(e)
	}
	util.DispatchOk(c)
}

// Trigger a task
func TriggerTaskEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.TriggerTaskReqVo
	util.MustBindJson(c, &req)

	e := domain.TriggerTask(user, &req)
	if e != nil {
		panic(e)
	}
	util.DispatchOk(c)
}
