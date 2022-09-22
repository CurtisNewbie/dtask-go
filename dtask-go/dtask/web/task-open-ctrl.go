package web

import (
	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

// Register OPEN API routes
func RegisterTaskOpenRoutes(router *gin.Engine) {
	router.POST(server.ResolvePath("/task/list", true), util.BuildAuthRouteHandler(ListTaskByPageEndpoint))
	router.POST(server.ResolvePath("/task/history", true), util.BuildAuthRouteHandler(ListTaskHistoryByPageEndpoint))
	router.POST(server.ResolvePath("/task/update", true), util.BuildAuthRouteHandler(UpdateTaskEndpoint))
	router.POST(server.ResolvePath("/task/trigger", true), util.BuildAuthRouteHandler(TriggerTaskEndpoint))
}

// List tasks
func ListTaskByPageEndpoint(c *gin.Context, user *util.User) any {
	var req domain.ListTaskByPageReqWebVo
	util.MustBindJson(c, &req)

	r, e := domain.ListTaskByPage(user, &req)
	if e != nil {
		panic(e)
	}
	return r
}

// List task histories
func ListTaskHistoryByPageEndpoint(c *gin.Context, user *util.User) any {

	var req domain.ListTaskHistoryByPageReq
	util.MustBindJson(c, &req)

	r, e := domain.ListTaskHistoryByPage(user, &req)
	if e != nil {
		panic(e)
	}
	return r
}

// Update task info
func UpdateTaskEndpoint(c *gin.Context, user *util.User) any {

	var req domain.UpdateTaskReq
	util.MustBindJson(c, &req)

	e := domain.UpdateTask(user, &req)
	if e != nil {
		panic(e)
	}
	return nil
}

// Trigger a task
func TriggerTaskEndpoint(c *gin.Context, user *util.User) any {

	var req domain.TriggerTaskReqVo
	util.MustBindJson(c, &req)

	e := domain.TriggerTask(user, &req)
	if e != nil {
		panic(e)
	}
	return nil
}
