package web

import (
	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

// Register routes
func RegisterTaskRoutes(router *gin.Engine) {

	router.POST(server.ResolvePath("/dtask/task/list", true), ListTaskByPageEndpoint)
	router.POST(server.ResolvePath("/dtask/task/history", true), listTaskHistoryByPageEndpoint)
	router.POST(server.ResolvePath("/dtask/task/update", true), UpdateTaskEndpoint)
	router.POST(server.ResolvePath("/dtask/task/trigger", true), TriggerTaskEndpoint)
}

// List tasks
func ListTaskByPageEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.ListTaskByPageReqWebVo
	util.MustBindJson(c, &req)

	r, e := domain.ListTaskByPage(user, &req)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}
	util.DispatchJson(c, r)
}

// List task histories
func listTaskHistoryByPageEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.ListTaskHistoryByPageReq
	util.MustBindJson(c, &req)

	r, e := domain.ListTaskHistoryByPage(user, &req)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}
	util.DispatchJson(c, r)
}

// Update task info
func UpdateTaskEndpoint(c *gin.Context) {
	user := util.RequireUser(c)

	var req domain.UpdateTaskReq
	util.MustBindJson(c, &req)

	e := domain.UpdateTask(user, &req)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
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
		util.DispatchErrJson(c, e)
		return
	}
	util.DispatchOk(c)
}
