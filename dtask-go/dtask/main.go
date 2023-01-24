package main

import (
	"os"

	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/gocommon/server"
)

func main() {

	// open-api routes
	server.PostJ(server.OpenApiPath("/task/list"), web.ListTaskByPageEndpoint)
	server.PostJ(server.OpenApiPath("/task/history"), web.ListTaskHistoryByPageEndpoint)
	server.PostJ(server.OpenApiPath("/task/update"), web.UpdateTaskEndpoint)
	server.PostJ(server.OpenApiPath("/task/trigger"), web.TriggerTaskEndpoint)

	// Internal RPC Calls (these should be protected by the gateway)
	server.Get(server.InternalApiPath("/task/all"), web.ListAllTaskRpc)
	server.Get(server.InternalApiPath("/task/valid"), web.ValidTaskRpc)
	server.PostJ(server.InternalApiPath("/task/lastRunInfo/update"), web.UpdateTaskLastRunInfoRpc)
	server.PostJ(server.InternalApiPath("/task/disable"), web.DisableTaskRpc)
	server.PostJ(server.InternalApiPath("/task/history"), web.RecordTaskHistoryRpc)
	server.PostJ(server.InternalApiPath("/task/declare"), web.DeclareTaskRpc)

	server.DefaultBootstrapServer(os.Args)
}
