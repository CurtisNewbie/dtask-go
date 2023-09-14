package main

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/miso"
)

const (
	MNG_TASK_CODE = "manage-tasks"
	MNG_TASK_NAME = "Manage Tasks"
)

func main() {

	miso.PostServerBootstrapped(func(rail miso.Rail) error {
		if goauth.IsEnabled() {
			if e := goauth.AddResourceAsync(rail, goauth.AddResourceReq{Code: MNG_TASK_CODE, Name: MNG_TASK_NAME}); e != nil {
				return fmt.Errorf("gclient.AddResource, %v", e)
			}
		}
		return nil
	})

	// open-api routes
	miso.IPost("/open/api/task/list",
		web.ListTaskByPageEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List tasks", Code: MNG_TASK_CODE}))

	miso.IPost("/open/api/task/history",
		web.ListTaskHistoryByPageEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List task execution history", Code: MNG_TASK_CODE}))

	miso.IPost("/open/api/task/update",
		web.UpdateTaskEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Update task", Code: MNG_TASK_CODE}))

	miso.IPost("/open/api/task/trigger",
		web.TriggerTaskEndpoint,
		goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Trigger task", Code: MNG_TASK_CODE}))

	// internal endpoints (these are protected by the gateway)
	miso.Get("/remote/task/all", web.ListAllTaskRpc)
	miso.Get("/remote/task/valid", web.ValidTaskRpc)
	miso.IPost("/remote/task/lastRunInfo/update", web.UpdateTaskLastRunInfoRpc)
	miso.IPost("/remote/task/disable", web.DisableTaskRpc)
	miso.IPost("/remote/task/history", web.RecordTaskHistoryRpc)
	miso.IPost("/remote/task/declare", web.DeclareTaskRpc)

	miso.BootstrapServer(os.Args)
}
