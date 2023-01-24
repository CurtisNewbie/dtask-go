package web

import (
	"github.com/curtisnewbie/dtask/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/gin-gonic/gin"
)

// List tasks
func ListTaskByPageEndpoint(c *gin.Context, ec common.ExecContext, req domain.ListTaskByPageReqWebVo) (any, error) {
	return domain.ListTaskByPage(ec, req)
}

// List task histories
func ListTaskHistoryByPageEndpoint(c *gin.Context, ec common.ExecContext, req domain.ListTaskHistoryByPageReq) (any, error) {
	return domain.ListTaskHistoryByPage(ec, req)
}

// Update task info
func UpdateTaskEndpoint(c *gin.Context, ec common.ExecContext, req domain.UpdateTaskReq) (any, error) {
	return nil, domain.UpdateTask(ec, req)
}

// Trigger a task
func TriggerTaskEndpoint(c *gin.Context, ec common.ExecContext, req domain.TriggerTaskReqVo) (any, error) {
	return nil, domain.TriggerTask(ec, req)
}
