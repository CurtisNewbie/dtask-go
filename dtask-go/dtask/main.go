package main

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/gocommon/server"
)

const (
	MNG_TASK_CODE = "manage-tasks"
	MNG_TASK_NAME = "Manage Tasks"
)

func main() {

	server.PostServerBootstrapped(func(rail common.Rail) error {
		if goauth.IsEnabled() {
			if e := goauth.AddResourceAsync(rail, goauth.AddResourceReq{Code: MNG_TASK_CODE, Name: MNG_TASK_NAME}); e != nil {
				return fmt.Errorf("gclient.AddResource, %v", e)
			}
		}
		return nil
	})

	// open-api routes
	server.IPost("/open/api/task/list",
		web.ListTaskByPageEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List tasks", Code: MNG_TASK_CODE}))

	server.IPost("/open/api/task/history",
		web.ListTaskHistoryByPageEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List task execution history", Code: MNG_TASK_CODE}))

	server.IPost("/open/api/task/update",
		web.UpdateTaskEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Update task", Code: MNG_TASK_CODE}))

	server.IPost("/open/api/task/trigger",
		web.TriggerTaskEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Trigger task", Code: MNG_TASK_CODE}))

	// internal endpoints (these are protected by the gateway)
	server.Get("/remote/task/all", web.ListAllTaskRpc)
	server.Get("/remote/task/valid", web.ValidTaskRpc)
	server.IPost("/remote/task/lastRunInfo/update", web.UpdateTaskLastRunInfoRpc)
	server.IPost("/remote/task/disable", web.DisableTaskRpc)
	server.IPost("/remote/task/history", web.RecordTaskHistoryRpc)
	server.IPost("/remote/task/declare", web.DeclareTaskRpc)

	server.BootstrapServer(os.Args)
}
