package router

import (
	"SimpliftTikTok/controller"
	"SimpliftTikTok/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

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
