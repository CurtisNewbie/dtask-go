package domain

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
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

	/** name of bean that will be executed */
	TargetBean string `json:"targetBean"`

	/** cron expression */
	CronExpr string `json:"cronExpr"`

	/** app group that runs this task */
	AppGroup string `json:"appGroup"`

	/** the last time this task was executed */
	LastRunStartTime *common.TTime `json:"lastRunStartTime"`

	/** the last time this task was finished */
	LastRunEndTime *common.TTime `json:"lastRunEndTime"`

	/** app that previously ran this task */
	LastRunBy string `json:"lastRunBy"`

	/** result of last execution */
	LastRunResult string `json:"lastRunResult"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled TaskEnabled `json:"enabled"`

	/** whether the task can be executed concurrently: 0-disabled, 1-enabled */
	ConcurrentEnabled TaskConcurrentEnabled `json:"concurrentEnabled"`

	/** update date */
	UpdateDate *common.TTime `json:"updateDate"`

	/** updated by */
	UpdateBy string `json:"updateBy"`
}

type ListTaskByPageRespWebVo struct {
	Tasks  *[]TaskWebVo  `json:"list"`
	Paging common.Paging `json:"pagingVo"`
}

type ListTaskByPageReqWebVo struct {
	Paging *common.Paging `json:"pagingVo"`

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
	Name      string `json:"name"`
	Group     string `json:"group"`
	TriggerBy string `json:"triggerBy"`
}

type TaskIdAppGroup struct {
	Id       *int
	AppGroup *string
}

type UpdateLastRunInfoReq struct {

	/** id */
	Id *int `json:"id"`

	/** the last time this task was executed */
	LastRunStartTime *common.TTime `json:"lastRunStartTime"`

	/** the last time this task was finished */
	LastRunEndTime *common.TTime `json:"lastRunEndTime"`

	/** app that previously ran this task */
	LastRunBy *string `json:"lastRunBy"`

	/** result of last execution */
	LastRunResult *string `json:"lastRunResult"`
}

type DisableTaskReqVo struct {

	/** id */
	Id int `json:"id"`

	/** result of last execution */
	LastRunResult string `json:"lastRunResult"`

	/** update date */
	UpdateDate common.TTime `json:"updateDate"`

	/** updated by */
	UpdateBy string `json:"updateBy"`
}

func DisableTask(ec common.Rail, req DisableTaskReqVo) error {
	qry := mysql.GetMySql()
	qry = qry.Table("task").Where("id = ?", req.Id)

	umap := make(map[string]any)
	umap["enabled"] = 0
	umap["last_run_result"] = req.LastRunResult
	umap["update_by"] = req.UpdateBy
	umap["update_date"] = time.Time(req.UpdateDate)

	tx := qry.Updates(umap)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func IsEnabledTask(ec common.Rail, taskId int) error {
	var id int
	if tx := mysql.GetMySql().Raw("select id from task where id = ? and enabled = 1", taskId).Scan(&id); tx.Error != nil {
		return tx.Error
	}

	if id < 1 {
		return common.NewWebErr("Task not found or disabled")
	}
	return nil
}

func UpdateTaskLastRunInfo(ec common.Rail, req UpdateLastRunInfoReq) error {
	ec.Infof("Received: %+v", req)
	if req.Id == nil {
		panic("id is required")
	}
	if req.LastRunBy == nil {
		panic("lastRunBy is required")
	}
	if req.LastRunStartTime == nil {
		panic("lastRunStartTime is required")
	}
	if req.LastRunEndTime == nil {
		panic("lastRunEndTime is required")
	}

	qry := mysql.GetMySql()
	qry = qry.Table("task").Where("id = ?", req.Id)

	st := time.Time(*req.LastRunStartTime)
	et := time.Time(*req.LastRunEndTime)

	umap := make(map[string]any)
	umap["last_run_start_time"] = st
	umap["last_run_end_time"] = et
	umap["last_run_by"] = req.LastRunBy

	end := 255
	curr := len(*req.LastRunResult)
	if curr < 255 {
		end = curr
	}
	umap["last_run_result"] = (*req.LastRunResult)[:end]

	tx := qry.Updates(umap)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Trigger a task
func TriggerTask(ec common.Rail, req TriggerTaskReqVo, user common.User) error {
	ta, e := FindTaskAppGroup(*req.Id)
	if e != nil {
		return e
	}

	// push the TriggeredJobKey into redis list, let the master poll and execute it
	tjk := TriggeredJobKey{Name: strconv.Itoa(*ta.Id), Group: *ta.AppGroup, TriggerBy: user.Username}
	key := _buildTriggeredJobListKey(*ta.AppGroup)
	val, e := json.Marshal(tjk)
	if e != nil {
		return e
	}
	json := string(val)
	ec.Infof("Triggering task, key: %v, TriggeredJobKey: %+v, json: %s", key, tjk, json)

	cmd := redis.GetRedis().LPush(key, json)
	if e := cmd.Err(); e != nil {
		return e
	}

	return nil
}

func FindTaskAppGroup(id int) (*TaskIdAppGroup, error) {
	var ta TaskIdAppGroup
	tx := mysql.GetMySql().Raw("select id, app_group from task where id = ?", id).Scan(&ta)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ta, nil
}

// Update task
func UpdateTask(ec common.Rail, req UpdateTaskReq, user common.User) error {

	qry := mysql.GetMySql()
	qry = qry.Table("task").Where("id = ?", req.Id)

	umap := make(map[string]any)

	if req.JobName != nil && !common.IsBlankStr(*req.JobName) {
		umap["job_name"] = *req.JobName
	}
	if req.TargetBean != nil && !common.IsBlankStr(*req.TargetBean) {
		umap["target_bean"] = *req.TargetBean
	}
	if req.CronExpr != nil && !common.IsBlankStr(*req.CronExpr) {
		umap["cron_expr"] = *req.CronExpr
	}
	if req.AppGroup != nil && !common.IsBlankStr(*req.AppGroup) {
		umap["app_group"] = *req.AppGroup
	}
	if req.Enabled != nil {
		umap["enabled"] = *req.Enabled
	}
	if req.ConcurrentEnabled != nil {
		umap["concurrent_enabled"] = int(*req.ConcurrentEnabled)
	}
	umap["update_by"] = user.Username
	umap["update_date"] = time.Now()

	tx := qry.Updates(umap)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// List all tasks for the appGroup
func ListAllTasks(ec common.Rail, appGroup string) (*[]TaskWebVo, error) {

	var tasks []TaskWebVo
	selectq := mysql.GetMySql().Table("task").Where("app_group = ?", appGroup)

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
func ListTaskByPage(ec common.Rail, req ListTaskByPageReqWebVo) (*ListTaskByPageRespWebVo, error) {
	if req.Paging == nil {
		req.Paging = &common.Paging{Limit: 30, Page: 1}
	}

	var tasks []TaskWebVo
	selectq := mysql.GetMySql().
		Table("task").
		Limit(req.Paging.Limit).
		Offset(req.Paging.GetOffset()).
		Order("app_group, id desc")

	_addWhereForListTaskByPage(&req, selectq)

	tx := selectq.Scan(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tasks == nil {
		tasks = []TaskWebVo{}
	}

	countq := mysql.GetMySql().
		Table("task").
		Select("COUNT(*)")

	_addWhereForListTaskByPage(&req, countq)

	var total int
	tx = countq.Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListTaskByPageRespWebVo{Tasks: &tasks, Paging: req.Paging.ToRespPage(total)}, nil
}

func _addWhereForListTaskByPage(req *ListTaskByPageReqWebVo, query *gorm.DB) *gorm.DB {
	if req.JobName != nil && !common.IsBlankStr(*req.JobName) {
		*query = *query.Where("job_name like ?", "%"+*req.JobName+"%")
	}
	if req.AppGroup != nil && !common.IsBlankStr(*req.AppGroup) {
		*query = *query.Where("app_group like ?", "%"+*req.AppGroup+"%")
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
