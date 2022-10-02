package web

import (
	"strconv"

	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

// Register internal routes
func RegisterTaskInternalRoutes(router *gin.Engine) {

	// Internal RPC Calls (these should be protected by the gateway)
	router.GET(server.ResolvePath("/task/all", false), util.BuildRouteHandler(ListAllTaskRpc))
	router.POST(server.ResolvePath("/task/lastRunInfo/update", false), util.BuildJRouteHandler(UpdateTaskLastRunInfoRpc))
	router.GET(server.ResolvePath("/task/valid", false), util.BuildRouteHandler(ValidTaskRpc))
	router.POST(server.ResolvePath("/task/disable", false), util.BuildJRouteHandler(DisableTaskRpc))
	router.POST(server.ResolvePath("/task/history", false), util.BuildJRouteHandler(RecordTaskHistoryRpc))
}

/*
	curl "http://localhost:8083/remote/task/all?appGroup=file-service"
*/
func ListAllTaskRpc(c *gin.Context) (any, error) {
	appGroup := c.Query("appGroup")
	return domain.ListAllTasks(&appGroup)
}

/*
	curl -X POST http://localhost:8083/remote/task/lastRunInfo/update -d ' { "id": 1, "lastRunStartTime" : "2022-09-10 15:04:05", "lastRunEndTime" : "2022-09-10 15:04:10", "lastRunBy" : "Yongj Zhuang", "lastRunResult" : "Looks good to me" } '
*/
func UpdateTaskLastRunInfoRpc(c *gin.Context, req *domain.UpdateLastRunInfoReq) (any, error) {
	e := domain.UpdateTaskLastRunInfo(req)
	return nil, e
}

/*
	curl "http://localhost:8083/remote/task/valid?taskId=1"
*/
func ValidTaskRpc(c *gin.Context) (any, error) {
	taskId := c.Query("taskId")
	cvtd, e := strconv.Atoi(taskId)
	if e != nil {
		panic(e)
	}

	e = domain.IsEnabledTask(&cvtd)
	return nil, e
}

/*
	curl -X POST http://localhost:8083/remote/task/disable -d ' { "id": 1, "lastRunResult" : "Something is wrong", "updateBy" : "scheduler", "updateDate" : "2022-09-10 17:04:10" }'
*/
func DisableTaskRpc(c *gin.Context, req *domain.DisableTaskReqVo) (any, error) {
	e := domain.DisableTask(req)
	return nil, e
}

/*
	curl -X POST http://localhost:8083/remote/task/history -d ' { "taskId": 1, "runResult" : "Very good", "runBy" : "scheduler", "startTime" : "2022-09-10 17:04:10", "endTime" : "2022-09-10 17:05:10" }'
*/
func RecordTaskHistoryRpc(c *gin.Context, req *domain.RecordTaskHistoryReq) (any, error) {
	e := domain.RecordTaskHistory(req)
	return nil, e
}
