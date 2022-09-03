package web

import (
	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
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

func ListTaskByPageEndpoint(c *gin.Context) {
	user, e := util.ExtractUser(c)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var req domain.ListTaskByPageReqWebVo
	if e := c.ShouldBindJSON(&req); e != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if _, e := domain.ListTaskByPage(user, &req); e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	util.DispatchOk(c)
}

func listTaskHistoryByPageEndpoint(c *gin.Context) {
	// todo
}

func UpdateTaskEndpoint(c *gin.Context) {
	// todo
}

func TriggerTaskEndpoint(c *gin.Context) {
	// todo
}
