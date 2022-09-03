package domain

import (
	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
	"gorm.io/gorm"
)

type TaskEnabled int

const (
	TASK_ENABLED  TaskEnabled = 1
	TASK_DISABLED TaskEnabled = 0
)

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
	LastRunStartTime dto.WTime `json:"lastRunStartTime"`

	/** the last time this task was finished */
	LastRunEndTime dto.WTime `json:"lastRunEndTime"`

	/** app that previously ran this task */
	LastRunBy string `json:"lastRunBy"`

	/** result of last execution */
	LastRunResult string `json:"lastRunResult"`

	/** whether the task is enabled: 0-disabled, 1-enabled */
	Enabled TaskEnabled `json:"enabled"`

	/** whether the task can be executed concurrently: 0-disabled, 1-enabled */
	ConcurrentEnabled int `json:"concurrentEnabled"`

	/** update date */
	UpdateDate dto.WTime `json:"updateDate"`

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

// List tasks
func ListTaskByPage(user *util.User, req *ListTaskByPageReqWebVo) (*ListTaskByPageRespWebVo, error) {

	util.RequireRole(user, util.ADMIN)

	var tasks []TaskWebVo
	selectq := config.GetDB().Limit(req.Paging.Limit).Offset(dto.CalcOffset(req.Paging))
	_addWhereForListTaskByPage(req, selectq)

	tx := selectq.Scan(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tasks == nil {
		tasks = []TaskWebVo{}
	}

	countq := config.GetDB().Select("COUNT(*)")
	_addWhereForListTaskByPage(req, countq)
	var total int
	tx = countq.Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListTaskByPageRespWebVo{Tasks: &tasks, Paging: dto.BuildResPage(req.Paging, total)}, nil
}

func _addWhereForListTaskByPage(req *ListTaskByPageReqWebVo, query *gorm.DB) *gorm.DB {
	if req.JobName != nil {
		query.Where("job_name like ?", "%"+*req.JobName+"%")
	}
	if req.AppGroup != nil {
		query.Where("app_group = ?", *req.AppGroup)
	}
	if req.Enabled != nil {
		query.Where("enabled = ?", *req.Enabled)
	}
	return query
}
