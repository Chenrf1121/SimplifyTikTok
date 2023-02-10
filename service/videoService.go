package service

import "SimpliftTikTok/dao"

// 视频结构体
type Video struct {
	dao.Video
	Id       int64     `json:"id"`       //视频的发布者
	Name     string    `json:"name"`     //视频名称
	Tag      []string  `json:"tag"`      //视频标签
	LikeCnt  int       `json:"likeCnt"`  //点赞数
	Comments []Comment `json:"comments"` //评论
}

// VideoService 定义视频接口以及各种方法
type VideoService interface {
	//用户发布视频
	Publish(userId int64)
	//获取视频
	Feed()
}
