package domain

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
	"github.com/curtisnewbie/gocommon/weberr"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskEnabled int
type TaskConcurrentEnabled int

const (
	TASK_ENABLED  TaskEnabled = 1
	TASK_DISABLED TaskEnabled = 0

	TASK_CONCURRENT_ENABLED  TaskConcurrentEnabled = 1
	TASK_CONCURRENT_DISABLED TaskConcurrentEnabled = 0
)

type UpdateTaskReq struct {

	/** id */
	Id int `json:"id"`

	/** job's name */
	JobName *string `json:"jobName"`

	/** name of bean that will be executed */
	TargetBean *string `json:"targetBean"`

	/** cron expression */
	CronExpr *string `json:"cronExpr"`

	/** app group that runs this task */
	AppGroup *string `json:"appGroup"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled *TaskEnabled `json:"enabled"`

	/** whether the task can be executed concurrently: 0-disabled, 1-enabled */
	ConcurrentEnabled *TaskConcurrentEnabled `json:"concurrentEnabled"`
}

type TaskWebVo struct {

	/** id */
	Id int `json:"id"`

	/** job's name */
	JobName string `json:"jobName"`

	/** cron expression */
	CronExpr string `json:"cronExpr"`

	/** app group that runs this task */
	AppGroup string `json:"appGroup"`

	/** the last time this task was executed */
	LastRunStartTime *dto.WTime `json:"lastRunStartTime"`

	/** the last time this task was finished */
	LastRunEndTime *dto.WTime `json:"lastRunEndTime"`

	/** app that previously ran this task */
	LastRunBy string `json:"lastRunBy"`

	/** result of last execution */
	LastRunResult string `json:"lastRunResult"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled TaskEnabled `json:"enabled"`

	/** whether the task can be executed concurrently: 0-disabled, 1-enabled */
	ConcurrentEnabled TaskConcurrentEnabled `json:"concurrentEnabled"`

	/** update date */
	UpdateDate *dto.WTime `json:"updateDate"`

	/** updated by */
	UpdateBy string `json:"updateBy"`
}

type ListTaskByPageRespWebVo struct {
	Tasks  *[]TaskWebVo `json:"list"`
	Paging *dto.Paging  `json:"pagingVo"`
}

type ListTaskByPageReqWebVo struct {
	Paging *dto.Paging `json:"pagingVo"`

	/** job's name */
	JobName *string `json:"jobName"`

	/** app group that runs this task */
	AppGroup *string `json:"appGroup"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled *TaskEnabled `json:"enabled"`
}

type TriggerTaskReqVo struct {
	Id *int `json:"id"`
}

// JobKey for manually triggered jobs
type TriggeredJobKey struct {
	Name      string
	Group     string
	TriggerBy string
}

type TaskIdAppGroup struct {
	Id       *int
	AppGroup *string
}

type UpdateLastRunInfoReq struct {

	/** id */
	Id int `json:"id"`

	/** the last time this task was executed */
	LastRunStartTime dto.TTime `json:"lastRunStartTime"`

	/** the last time this task was finished */
	LastRunEndTime dto.TTime `json:"lastRunEndTime"`

	/** app that previously ran this task */
	LastRunBy string `json:"lastRunBy"`

	/** result of last execution */
	LastRunResult string `json:"lastRunResult"`
}

type DisableTaskReqVo struct {

	/** id */
	Id int `json:"id"`

	/** result of last execution */
	LastRunResult string `json:"lastRunResult"`

	/** update date */
	UpdateDate dto.TTime `json:"updateDate"`

	/** updated by */
	UpdateBy string `json:"updateBy"`
}

func DisableTask(req *DisableTaskReqVo) error {
	qry := config.GetDB()
	qry = qry.Debug().Table("task").Where("id = ?", req.Id)

	umap := make(map[string]any)
	umap["last_run_result"] = req.LastRunResult
	umap["update_by"] = req.UpdateBy
	umap["update_date"] = time.Time(req.UpdateDate)

	tx := qry.Updates(umap)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func IsEnabledTask(taskId *int) error {
	var id int
	tx := config.GetDB().Raw("select id from task where id = ? and enabled = 1", *taskId).Scan(&id)

	if tx.Error != nil {
		return tx.Error
	}

	if id < 1 {
		return weberr.NewWebErr("Task not found or disabled")
	}
	return nil
}

func UpdateTaskLastRunInfo(req *UpdateLastRunInfoReq) error {

	qry := config.GetDB()
	qry = qry.Table("task").Where("id = ?", req.Id)

	umap := make(map[string]any)
	umap["last_run_start_time"] = time.Time(req.LastRunStartTime)
	umap["last_run_end_time"] = time.Time(req.LastRunEndTime)
	umap["last_run_by"] = req.LastRunBy
	umap["last_run_result"] = req.LastRunResult

	tx := qry.Updates(umap)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Trigger a task
func TriggerTask(user *util.User, req *TriggerTaskReqVo) error {

	util.RequireRole(user, util.ADMIN)

	ta, e := FindTaskAppGroup(*req.Id)
	if e != nil {
		return e
	}

	// push the TriggeredJobKey into redis list, let the master poll and execute it
	tjk := TriggeredJobKey{Name: strconv.Itoa(*ta.Id), Group: *ta.AppGroup, TriggerBy: user.Username}
	key := _buildTriggeredJobListKey(*ta.AppGroup)
	log.Infof("Triggering task, key: %v, TriggeredJobKey: %+v", key, tjk)

	val, e := json.Marshal(tjk)
	if e != nil {
		return e
	}
	cmd := config.GetRedis().LPush(key, string(val))
	if e := cmd.Err(); e != nil {
		return e
	}

	return nil
}

func FindTaskAppGroup(id int) (*TaskIdAppGroup, error) {
	var ta TaskIdAppGroup
	tx := config.GetDB().Raw("select id, app_group from task where id = ?", id).Scan(&ta)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ta, nil
}

// Update task
func UpdateTask(user *util.User, req *UpdateTaskReq) error {

	util.RequireRole(user, util.ADMIN)

	qry := config.GetDB()
	qry = qry.Where("id = ?", req.Id)

	umap := make(map[string]any)

	if util.IsEmpty(req.JobName) {
		umap["job_name"] = *req.JobName
	}
	if util.IsEmpty(req.TargetBean) {
		umap["target_bean"] = *req.TargetBean
	}
	if util.IsEmpty(req.CronExpr) {
		umap["cron_expr"] = *req.CronExpr
	}
	if util.IsEmpty(req.AppGroup) {
		umap["app_group"] = *req.AppGroup
	}
	if req.Enabled != nil {
		umap["enabled"] = *req.Enabled
	}
	if req.ConcurrentEnabled != nil {
		umap["concurrent_enabled"] = *req.ConcurrentEnabled
	}

	tx := qry.Updates(umap)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// List all tasks for the appGroup
func ListAllTasks(appGroup *string) (*[]TaskWebVo, error) {

	var tasks []TaskWebVo
	selectq := config.GetDB().Table("task").Where("app_group = ?", *appGroup)

	tx := selectq.Scan(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tasks == nil {
		tasks = []TaskWebVo{}
	}

	return &tasks, nil
}

// List tasks
func ListTaskByPage(user *util.User, req *ListTaskByPageReqWebVo) (*ListTaskByPageRespWebVo, error) {

	util.RequireRole(user, util.ADMIN)

	if req.Paging == nil {
		req.Paging = &dto.Paging{Limit: 30, Page: 1}
	}

	var tasks []TaskWebVo
	selectq := config.GetDB().Table("task").Limit(req.Paging.Limit).Offset(dto.CalcOffset(req.Paging))
	_addWhereForListTaskByPage(req, selectq)

	tx := selectq.Scan(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tasks == nil {
		tasks = []TaskWebVo{}
	}

	countq := config.GetDB().Table("task").Select("COUNT(*)")
	_addWhereForListTaskByPage(req, countq)
	var total int
	tx = countq.Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListTaskByPageRespWebVo{Tasks: &tasks, Paging: dto.BuildResPage(req.Paging, total)}, nil
}

func _addWhereForListTaskByPage(req *ListTaskByPageReqWebVo, query *gorm.DB) *gorm.DB {
	if !util.IsEmpty(req.JobName) {
		*query = *query.Where("job_name like ?", "%"+*req.JobName+"%")
	}
	if !util.IsEmpty(req.AppGroup) {
		*query = *query.Where("app_group = ?", *req.AppGroup)
	}
	if req.Enabled != nil {
		*query = *query.Where("enabled = ?", *req.Enabled)
	}
	return query
}

// Build Redis's key for list of manually triggered job
func _buildTriggeredJobListKey(appGroup string) string {
	return "task:trigger:group:" + appGroup
}
