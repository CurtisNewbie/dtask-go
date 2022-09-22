package main

import (
	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

func main() {
	profile, conf := config.DefaultParseProfConf()

	// init handle for database
	if err := config.InitDBFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}
	// init handle for redis
	config.InitRedisFromConfig(&conf.RedisConf)

	// bootstrap web server
	err := server.BootstrapServer(&conf.ServerConf, config.IsProd(profile), func(router *gin.Engine) {
		web.RegisterTaskOpenRoutes(router)
		web.RegisterTaskInternalRoutes(router)
	})

	if err != nil {
		panic(err)
	}
}
