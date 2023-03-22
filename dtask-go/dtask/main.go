package main

import (
	"log"
	"os"

	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/goauth/client/goauth-client-go/gclient"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
)

const (
	MNG_TASK_CODE = "manage-tasks"
	MNG_TASK_NAME = "Manage Tasks"
)

func main() {

	ec := common.EmptyExecContext()
	server.OnServerBootstrapped(func() {
		if e := gclient.AddResource(ec.Ctx, gclient.AddResourceReq{Code: MNG_TASK_CODE, Name: MNG_TASK_NAME}); e != nil {
			log.Fatalf("gclient.AddResource, %v", e)
		}
	})

	// open-api routes
	server.PostJ(server.OpenApiPath("/task/list"), web.ListTaskByPageEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/task/list"), Type: gclient.PT_PROTECTED, Desc: "List tasks", Method: "POST"})

	server.PostJ(server.OpenApiPath("/task/history"), web.ListTaskHistoryByPageEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/task/histroy"), Type: gclient.PT_PROTECTED, Desc: "List task execution history",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/task/update"), web.UpdateTaskEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/task/update"), Type: gclient.PT_PROTECTED, Desc: "Update task", Method: "POST"})

	server.PostJ(server.OpenApiPath("/task/trigger"), web.TriggerTaskEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/task/trigger"), Type: gclient.PT_PROTECTED, Desc: "Trigger task", Method: "POST"})

	// Internal RPC Calls (these should be protected by the gateway)
	server.Get(server.InternalApiPath("/task/all"), web.ListAllTaskRpc)
	server.Get(server.InternalApiPath("/task/valid"), web.ValidTaskRpc)
	server.PostJ(server.InternalApiPath("/task/lastRunInfo/update"), web.UpdateTaskLastRunInfoRpc)
	server.PostJ(server.InternalApiPath("/task/disable"), web.DisableTaskRpc)
	server.PostJ(server.InternalApiPath("/task/history"), web.RecordTaskHistoryRpc)
	server.PostJ(server.InternalApiPath("/task/declare"), web.DeclareTaskRpc)

	server.DefaultBootstrapServer(os.Args)
}

func reportPath(ec common.ExecContext, r gclient.CreatePathReq) {
	server.OnServerBootstrapped(func() {
		r.Url = "/dtaskgo" + r.Url
		r.Group = "dtaskgo"
		r.ResCode = MNG_TASK_CODE
		if e := gclient.AddPath(ec.Ctx, r); e != nil {
			log.Fatalf("gclient.AddPath, req: %+v, %v", r, e)
		}
	})
}
