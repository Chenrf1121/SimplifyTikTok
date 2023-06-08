package service

import "SimpliftTikTok/dao"

// 视频结构体
type Video struct {
	MetaVideo    dao.MetaVideo
	DynamicVideo dao.DynamicVideo
	Comments     []Comment `json:"comments"` //评论
}

// VideoService 定义视频接口以及各种方法
type VideoService interface {
	//用户发布视频
	Publish(userId int64)
	//获取视频
	Feed()
	//获取用户发布视频列表
	GetVideoPublishList(userId int64) ([]Video, error)
}
