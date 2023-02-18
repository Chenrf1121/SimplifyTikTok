package router

import (
	"SimpliftTikTok/controller"
	"SimpliftTikTok/middleware/jwt"
	"SimpliftTikTok/middleware/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter(r *gin.Engine) {
	gp := prometheus.New(r)

	// public directory is used to serve static resources
	//	registry := prometheus.NewRegistry()
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(
		prometheus.Registry, promhttp.HandlerOpts{
			Registry: prometheus.Registry,
		})))
	r.Static("/static", "./public")
	r.Use(gp.PromeMiddleware())
	apiRouter := r.Group("/tiktok")

	// 用户注册，登录，用户信息，上传视频，刷视频，显示发布视频
	apiRouter.GET("/feed/", controller.Feed) //视频观看
	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/", jwt.Auth(), controller.Publish)
	apiRouter.GET("/publishlist/", jwt.AuthWithoutLogin(), controller.PublishList)

	//用户关注，给视频点赞，评论，收藏
	//用户关注
	apiRouter.POST("/follow/", controller.FollowAction)
	//apiRouter.GET("/favorite/list/", controller.FavoriteList)
	//apiRouter.POST("/comment/action/", controller.CommentAction)
	//apiRouter.GET("/comment/list/", controller.CommentList)
	//
	//// extra apis - II
	//apiRouter.POST("/relation/action/", controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	//apiRouter.GET("/relation/friend/list/", controller.FriendList)
	//apiRouter.GET("/message/chat/", controller.MessageChat)
	//apiRouter.POST("/message/action/", controller.MessageAction)
}
