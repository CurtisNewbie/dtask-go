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
	StartTime *dto.TTime `json:"startTime"`

	/** end time */
	EndTime *dto.TTime `json:"endTime"`

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
	StartTime *dto.TTime `json:"startTime"`

	/** end time */
	EndTime *dto.TTime `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	Paging *dto.Paging `json:"pagingVo"`
}

type RecordTaskHistoryReq struct {

	/** task's id */
	TaskId int `json:"taskId"`

	/** start time */
	StartTime *dto.TTime `json:"startTime"`

	/** end time */
	EndTime *dto.TTime `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	/** result of last execution */
	RunResult *string `json:"runResult"`
}

func RecordTaskHistory(req *RecordTaskHistoryReq) error {

	db := config.GetDB().Table("task_history")
	m := make(map[string]any)

	m["task_id"] = req.TaskId
	m["start_time"] = time.Time(*req.StartTime)
	m["end_time"] = time.Time(*req.EndTime)
	m["run_by"] = req.RunBy
	m["run_result"] = req.RunResult
	m["create_time"] = time.Now()

	if e := db.Create(m).Error; e != nil {
		return e
	}

	return nil
}

// List tasks
func ListTaskHistoryByPage(user *util.User, req *ListTaskHistoryByPageReq) (*ListTaskHistoryByPageResp, error) {

	util.RequireRole(user, util.ADMIN)

	if req.Paging == nil {
		req.Paging = &dto.Paging{Limit: 30, Page: 1}
	}

	var histories []TaskHistoryWebVo
	selectq := config.GetDB().
		Table("task_history th").
		Select("th.id, t.job_name, th.task_id, th.start_time, th.end_time, th.run_by, th.run_result").
		Joins("LEFT JOIN task t ON th.task_id = t.id").
		Offset(dto.CalcOffset(req.Paging)).
		Limit(req.Paging.Limit).
		Order("th.id DESC")

	// dynamic where conditions
	_addWhereForListTaskHistoryByPage(req, selectq)

	tx := selectq.Scan(&histories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if histories == nil {
		histories = []TaskHistoryWebVo{}
	}

	countq := config.GetDB().
		Table("task_history th").
		Select("count(th.id)").
		Joins("LEFT JOIN task t ON th.task_id = t.id")

	// dynamic where conditions
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
		*query = *query.Where("th.task_id = ?", *req.TaskId)
	}
	if !util.IsEmpty(req.JobName) {
		*query = *query.Where("t.job_name like ?", "%"+*req.JobName+"%")
	}
	if req.StartTime != nil {
		t := time.Time(*req.StartTime)
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		*query = *query.Where("th.start_time >= ?", t)
	}
	if req.EndTime != nil {
		t := time.Time(*req.EndTime)
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		*query = *query.Where("th.end_time <= ?", t)
	}
	return query
}
