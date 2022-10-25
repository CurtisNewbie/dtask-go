package main

import (
	"os"

	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

func main() {
	_, conf := config.DefaultParseProfConf(os.Args)

	server.BootstrapServer(conf, func(router *gin.Engine) {
		web.RegisterTaskOpenRoutes(router)
		web.RegisterTaskInternalRoutes(router)
	})
}
