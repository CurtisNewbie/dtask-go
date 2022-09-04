package domain

import (
	"time"

	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
	"gorm.io/gorm"
)

type TaskHistoryWebVo struct {

	/** id */
	Id int

	/** job name */
	JobName *string `json:"jobName"`

	/** task's id */
	TaskId *int `json:"taskId"`

	/** start time */
	StartTime *dto.WTime `json:"startTime"`

	/** end time */
	EndTime *dto.WTime `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	/** result of last execution */
	RunResult *string `json:"runResult"`
}

type ListTaskHistoryByPageResp struct {
	Histories *[]TaskHistoryWebVo `json:"list"`
	Paging    *dto.Paging         `json:"pagingVo"`
}

type ListTaskHistoryByPageReq struct {

	/** task's id */
	TaskId *int `json:"taskId"`

	/** task' name */
	JobName *string `json:"jobName"`

	/** start time */
	StartTime *time.Time `json:"startTime"`

	/** end time */
	EndTime *time.Time `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	Paging *dto.Paging `json:"pagingVo"`
}

// List tasks
func ListTaskHistoryByPage(user *util.User, req *ListTaskHistoryByPageReq) (*ListTaskHistoryByPageResp, error) {

	util.RequireRole(user, util.ADMIN)

	if req.Paging == nil {
		req.Paging = &dto.Paging{Limit: 30, Page: 1}
	}

	var histories []TaskHistoryWebVo
	selectq := config.GetDB().Table("task_history").Limit(req.Paging.Limit).Offset(dto.CalcOffset(req.Paging))
	_addWhereForListTaskHistoryByPage(req, selectq)

	tx := selectq.Scan(&histories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if histories == nil {
		histories = []TaskHistoryWebVo{}
	}

	countq := config.GetDB().Table("task_history").Select("COUNT(*)")
	_addWhereForListTaskHistoryByPage(req, countq)
	var total int
	tx = countq.Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListTaskHistoryByPageResp{Histories: &histories, Paging: dto.BuildResPage(req.Paging, total)}, nil
}

func _addWhereForListTaskHistoryByPage(req *ListTaskHistoryByPageReq, query *gorm.DB) *gorm.DB {
	if req.TaskId != nil {
		query.Where("task_id = ?", *req.TaskId)
	}
	if !util.IsEmpty(req.JobName) {
		query.Where("job_name like ?", "%"+*req.JobName+"%")
	}
	if req.StartTime != nil {
		query.Where("start_time >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		query.Where("end_time <= ?", *req.EndTime)
	}
	return query
}
