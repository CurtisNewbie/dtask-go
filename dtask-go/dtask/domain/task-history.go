package domain

import (
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
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
	StartTime *common.TTime `json:"startTime"`

	/** end time */
	EndTime *common.TTime `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	/** result of last execution */
	RunResult *string `json:"runResult"`
}

type ListTaskHistoryByPageResp struct {
	Histories *[]TaskHistoryWebVo `json:"list"`
	Paging    common.Paging       `json:"pagingVo"`
}

type ListTaskHistoryByPageReq struct {

	/** task's id */
	TaskId *int `json:"taskId"`

	/** task' name */
	JobName *string `json:"jobName"`

	/** start time */
	StartTime *common.TTime `json:"startTime"`

	/** end time */
	EndTime *common.TTime `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	Paging *common.Paging `json:"pagingVo"`
}

type RecordTaskHistoryReq struct {

	/** task's id */
	TaskId int `json:"taskId"`

	/** start time */
	StartTime *common.TTime `json:"startTime"`

	/** end time */
	EndTime *common.TTime `json:"endTime"`

	/** task triggered by */
	RunBy *string `json:"runBy"`

	/** result of last execution */
	RunResult *string `json:"runResult"`
}

type DeclareTaskReq struct {

	/** job's name */
	JobName *string `json:"jobName" validation:"notNil"`

	/** name of bean that will be executed */
	TargetBean *string `json:"targetBean" validation:"notNil"`

	/** cron expression */
	CronExpr *string `json:"cronExpr" validation:"notNil"`

	/** app group that runs this task */
	AppGroup *string `json:"appGroup" validation:"notNil"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled *int `json:"enabled" validation:"notNil"`

	/** whether the task can be executed concurrently: 0-disabled, 1-enabled */
	ConcurrentEnabled *int `json:"concurrentEnabled" validation:"notNil"`

	/** Whether this declaration overrides existing configuration */
	Overridden *bool `json:"overridden" validation:"notNil"`
}

func RecordTaskHistory(ec common.ExecContext, req RecordTaskHistoryReq) error {

	db := mysql.GetMySql().Table("task_history")
	m := make(map[string]any)

	st := time.Time(*req.StartTime)
	et := time.Time(*req.EndTime)

	m["task_id"] = req.TaskId
	m["start_time"] = st
	m["end_time"] = et
	m["run_by"] = req.RunBy
	m["run_result"] = req.RunResult
	m["create_time"] = time.Now()

	if e := db.Create(m).Error; e != nil {
		return e
	}

	return nil
}

// List tasks
func ListTaskHistoryByPage(ec common.ExecContext, req ListTaskHistoryByPageReq) (*ListTaskHistoryByPageResp, error) {

	if req.Paging == nil {
		req.Paging = &common.Paging{Limit: 30, Page: 1}
	}

	var histories []TaskHistoryWebVo
	selectq := mysql.GetMySql().
		Table("task_history th").
		Select("th.id, t.job_name, th.task_id, th.start_time, th.end_time, th.run_by, th.run_result").
		Joins("LEFT JOIN task t ON th.task_id = t.id").
		Offset(req.Paging.GetOffset()).
		Limit(req.Paging.Limit).
		Order("th.id DESC")

	// dynamic where conditions
	_addWhereForListTaskHistoryByPage(&req, selectq)

	tx := selectq.Scan(&histories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if histories == nil {
		histories = []TaskHistoryWebVo{}
	}

	countq := mysql.GetMySql().
		Table("task_history th").
		Select("count(th.id)").
		Joins("LEFT JOIN task t ON th.task_id = t.id")

	// dynamic where conditions
	_addWhereForListTaskHistoryByPage(&req, countq)

	var total int
	tx = countq.Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListTaskHistoryByPageResp{Histories: &histories, Paging: req.Paging.ToRespPage(total)}, nil
}

func _addWhereForListTaskHistoryByPage(req *ListTaskHistoryByPageReq, query *gorm.DB) *gorm.DB {
	if req.TaskId != nil {
		*query = *query.Where("th.task_id = ?", *req.TaskId)
	}
	if !common.IsEmpty(req.JobName) {
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

// Declare task
func DeclareTask(ec common.ExecContext, req DeclareTaskReq) error {
	appGroup := req.AppGroup
	_, e := redis.RLockRun(ec, "task:declare:dtaskgo:"+*appGroup, func() (any, error) {

		slt := "select id from task where app_group = ? and target_bean = ? limit 1"
		var id int
		tx := mysql.GetMySql().Raw(slt, *appGroup, *req.TargetBean).Scan(&id)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if tx.RowsAffected < 1 {
			ist := "insert into task (job_name, cron_expr, enabled, concurrent_enabled, target_bean, app_group, update_by, update_date) values (?, ?, ?, ?, ?, ?, ?, ?)"
			tx := mysql.GetMySql().Exec(ist, *req.JobName, *req.CronExpr, *req.Enabled, *req.ConcurrentEnabled, *req.TargetBean, *req.AppGroup, "JobDeclaration", time.Now())
			return nil, tx.Error
		}

		if !*req.Overridden {
			return nil, nil
		}

		udt := "update task set cron_expr = ?, concurrent_enabled = ?, enabled = ?, update_by = ?, update_date = ? where id = ?"
		return nil, mysql.GetMySql().Exec(udt, *req.CronExpr, *req.ConcurrentEnabled, *req.Enabled, "JobDeclaration", time.Now(), id).Error
	})
	return e
}
