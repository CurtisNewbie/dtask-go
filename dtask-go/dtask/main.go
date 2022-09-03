package main

import (
	"fmt"
	"log"
	"os"

	"github.com/curtisnewbie/dtask/web"
	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

func main() {
	profile := config.ParseProfile(os.Args[1:])
	log.Printf("Using profile: %v", profile)

	conf, err := config.ParseJsonConfig(fmt.Sprintf("app-conf-%v.json", profile))
	if err != nil {
		panic(err)
	}
	config.SetGlobalConfig(conf)

	// init handle for database
	if err := config.InitDBFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}
	// init handle for redis
	config.InitRedisFromConfig(&conf.RedisConf)

	// bootstrap web server
	err = server.BootstrapServer(&conf.ServerConf, config.IsProd(profile), func(router *gin.Engine) {
		web.RegisterTaskRoutes(router)
	})
	if err != nil {
		panic(err)
	}
}
