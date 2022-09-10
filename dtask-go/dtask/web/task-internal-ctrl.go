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
	router.GET(server.ResolvePath(base, "/task/all", false), ListAllTaskRpc)
	router.POST(server.ResolvePath(base, "/task/lastRunInfo/update", false), UpdateTaskLastRunInfoRpc)
	router.GET(server.ResolvePath(base, "/task/valid", false), ValidTaskRpc)
	router.POST(server.ResolvePath(base, "/task/disable", false), DisableTaskRpc)
	router.POST(server.ResolvePath(base, "/task/history", false), RecordTaskHistoryRpc)
}

/*
	curl "http://localhost:8082/dtask/remote/task/all?appGroup=file-service"
*/
func ListAllTaskRpc(c *gin.Context) {
	appGroup := c.Query("appGroup")
	r, e := domain.ListAllTasks(&appGroup)
	if e != nil {
		panic(e)
	}
	util.DispatchOkWData(c, r)
}

/*
	curl -X POST http://localhost:8082/dtask/remote/task/lastRunInfo/update -d ' { "id": 1, "lastRunStartTime" : "2022-09-10 15:04:05", "lastRunEndTime" : "2022-09-10 15:04:10", "lastRunBy" : "Yongj Zhuang", "lastRunResult" : "Looks good to me" } '
*/
func UpdateTaskLastRunInfoRpc(c *gin.Context) {

	var req domain.UpdateLastRunInfoReq
	util.MustBindJson(c, &req)

	e := domain.UpdateTaskLastRunInfo(&req)
	if e != nil {
		panic(e)
	}
	util.DispatchOk(c)
}

/*
	curl "http://localhost:8082/dtask/remote/task/valid?taskId=1"
*/
func ValidTaskRpc(c *gin.Context) {
	taskId := c.Query("taskId")
	cvtd, e := strconv.Atoi(taskId)
	if e != nil {
		panic(e)
	}

	e = domain.IsEnabledTask(&cvtd)
	if e != nil {
		panic(e)
	}
	util.DispatchOk(c)
}

/*
	curl -X POST http://localhost:8082/dtask/remote/task/disable -d ' { "id": 1, "lastRunResult" : "Something is wrong", "updateBy" : "scheduler", "updateDate" : "2022-09-10 17:04:10" }'
*/
func DisableTaskRpc(c *gin.Context) {
	var req domain.DisableTaskReqVo
	util.MustBindJson(c, &req)

	e := domain.DisableTask(&req)
	if e != nil {
		panic(e)
	}
	util.DispatchOk(c)
}

/*
	curl -X POST http://localhost:8082/dtask/remote/task/history -d ' { "taskId": 1, "runResult" : "Very good", "runBy" : "scheduler", "startTime" : "2022-09-10 17:04:10", "endTime" : "2022-09-10 17:05:10" }'
*/
func RecordTaskHistoryRpc(c *gin.Context) {

	var req domain.RecordTaskHistoryReq
	util.MustBindJson(c, &req)

	if e := domain.RecordTaskHistory(&req); e != nil {
		panic(e)
	}

	util.DispatchOk(c)
}
