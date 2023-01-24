package web

import (
	"strconv"

	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/gin-gonic/gin"
)

func ListAllTaskRpc(c *gin.Context, ec common.ExecContext) (any, error) {
	appGroup := c.Query("appGroup")
	return domain.ListAllTasks(ec, appGroup)
}

func UpdateTaskLastRunInfoRpc(c *gin.Context, ec common.ExecContext, req domain.UpdateLastRunInfoReq) (any, error) {
	return nil, domain.UpdateTaskLastRunInfo(ec, req)
}

func ValidTaskRpc(c *gin.Context, ec common.ExecContext) (any, error) {
	taskId := c.Query("taskId")
	cvtd, e := strconv.Atoi(taskId)
	if e != nil {
		return nil, e
	}

	return nil, domain.IsEnabledTask(ec, cvtd)
}

func DisableTaskRpc(c *gin.Context, ec common.ExecContext, req domain.DisableTaskReqVo) (any, error) {
	return nil, domain.DisableTask(ec, req)
}

func RecordTaskHistoryRpc(c *gin.Context, ec common.ExecContext, req domain.RecordTaskHistoryReq) (any, error) {
	return nil, domain.RecordTaskHistory(ec, req)
}

func DeclareTaskRpc(c *gin.Context, ec common.ExecContext, req domain.DeclareTaskReq) (any, error) {
	return nil, domain.DeclareTask(ec, req)
}
