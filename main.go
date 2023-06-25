package main

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/dao"
	"SimpliftTikTok/middleware/ftp"
	"SimpliftTikTok/middleware/redis"
	"SimpliftTikTok/router"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	//	go service.RunMessageServer()
	initDevp()

	r := gin.Default()

	router.InitRouter(r)
	pprof.Register(r)
	err := r.Run("0.0.0.0:8090")
	log.Println("err ======= ", err) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func initDevp() {
	config.Pwd, _ = os.Getwd()
	log.Println("PWD = ", config.Pwd)
	dao.Init()
	redis.InitRedis()
	ftp.InitFTP()
}
