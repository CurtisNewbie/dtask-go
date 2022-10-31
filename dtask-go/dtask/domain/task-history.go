package domain

import (
	"time"

	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
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

type DeclareTaskReq struct {

	/** job's name */
	JobName *string `json:"jobName"`

	/** name of bean that will be executed */
	TargetBean *string `json:"targetBean"`

	/** cron expression */
	CronExpr *string `json:"cronExpr"`

	/** app group that runs this task */
	AppGroup *string `json:"appGroup"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled *int `json:"enabled"`

	/** whether the task can be executed concurrently: 0-disabled, 1-enabled */
	ConcurrentEnabled *int `json:"concurrentEnabled"`

	/** Whether this declaration overrides existing configuration */
	Overridden *bool `json:"overridden"`
}

func RecordTaskHistory(req *RecordTaskHistoryReq) error {

	db := mysql.GetDB().Table("task_history")
	m := make(map[string]any)

	st := time.Time(*req.StartTime)
	et := time.Time(*req.EndTime)
	// log.Infof("recordTaskHistory, start: %s, end: %s", dto.TimePrettyPrint(&st), dto.TimePrettyPrint(&et))

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
func ListTaskHistoryByPage(user *util.User, req *ListTaskHistoryByPageReq) (*ListTaskHistoryByPageResp, error) {

	util.RequireRole(user, util.ADMIN)

	if req.Paging == nil {
		req.Paging = &dto.Paging{Limit: 30, Page: 1}
	}

	var histories []TaskHistoryWebVo
	selectq := mysql.GetDB().
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

	countq := mysql.GetDB().
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

// Declare task
func DeclareTask(req *DeclareTaskReq) error {
	util.NonNil(req, "req is nil")
	util.NonNil(req.JobName, "jobName is nil")
	util.NonNil(req.CronExpr, "cronExpr is nil")
	util.NonNil(req.Enabled, "enabled is nil")
	util.NonNil(req.ConcurrentEnabled, "concurrentEnabled is nil")
	util.NonNil(req.Overridden, "overriden is nil")
	util.NonNil(req.TargetBean, "targetBean is nil")

	appGroup := util.NonNil(req.AppGroup, "appGroup is nil")
	_, e := redis.LockRun("task:declare:dtaskgo:"+*appGroup, func() any {

		slt := "select id from task where app_group = ? and target_bean = ? limit 1"
		var id int
		tx := mysql.GetDB().Raw(slt, *appGroup, *req.TargetBean).Scan(&id)
		if tx.Error != nil {
			return tx.Error
		}

		if tx.RowsAffected < 1 {
			ist := "insert into task (job_name, cron_expr, enabled, concurrent_enabled, target_bean, app_group, update_by, update_date) values (?, ?, ?, ?, ?, ?, ?, ?)"
			tx := mysql.GetDB().Exec(ist, *req.JobName, *req.CronExpr, *req.Enabled, *req.ConcurrentEnabled, *req.TargetBean, *req.AppGroup, "JobDeclaration", time.Now())
			return tx.Error
		}

		if !*req.Overridden {
			return nil
		}

		udt := "update task set cron_expr = ?, concurrent_enabled = ?, enabled = ?, update_by = ?, update_date = ? where id = ?"
		return mysql.GetDB().Exec(udt, *req.CronExpr, *req.ConcurrentEnabled, *req.Enabled, "JobDeclaration", time.Now(), id).Error
	})
	return e
}
