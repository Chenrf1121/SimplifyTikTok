package main

import (
	"SimpliftTikTok/dao"
	"SimpliftTikTok/middleware/ftp"
	"SimpliftTikTok/middleware/redis"
	"SimpliftTikTok/router"
	"SimpliftTikTok/service"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()
	initDevp()

	r := gin.Default()

	router.InitRouter(r)
	pprof.Register(r)
	r.Run(":8090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func initDevp() {
	dao.Init()
	redis.InitRedis()
	ftp.InitFTP()
}
