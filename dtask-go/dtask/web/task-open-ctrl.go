package web

import (
	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

// Register OPEN API routes
func RegisterTaskOpenRoutes(router *gin.Engine) {
	router.POST(server.ResolvePath("/task/list", true), util.BuildAuthJRouteHandler(ListTaskByPageEndpoint))
	router.POST(server.ResolvePath("/task/history", true), util.BuildAuthJRouteHandler(ListTaskHistoryByPageEndpoint))
	router.POST(server.ResolvePath("/task/update", true), util.BuildAuthJRouteHandler(UpdateTaskEndpoint))
	router.POST(server.ResolvePath("/task/trigger", true), util.BuildAuthJRouteHandler(TriggerTaskEndpoint))
}

// List tasks
func ListTaskByPageEndpoint(c *gin.Context, user *util.User, req *domain.ListTaskByPageReqWebVo) (any, error) {
	return domain.ListTaskByPage(user, req)
}

// List task histories
func ListTaskHistoryByPageEndpoint(c *gin.Context, user *util.User, req *domain.ListTaskHistoryByPageReq) (any, error) {
	return domain.ListTaskHistoryByPage(user, req)
}

// Update task info
func UpdateTaskEndpoint(c *gin.Context, user *util.User, req *domain.UpdateTaskReq) (any, error) {
	e := domain.UpdateTask(user, req)
	return nil, e
}

// Trigger a task
func TriggerTaskEndpoint(c *gin.Context, user *util.User, req *domain.TriggerTaskReqVo) (any, error) {
	e := domain.TriggerTask(user, req)
	return nil, e
}
