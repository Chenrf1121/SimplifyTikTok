package service

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/dao"
	"log"
	"time"
)

type VideoServiceImpl struct {
	UserService
	LikeService
	CommentService
}

// 刷视频
func (v VideoServiceImpl) Feed() ([]Video, time.Time, error) {
	//从文件服务器读视频
	videos := make([]Video, config.MaxCacheVideo)
	//从数据库得到视频信息
	tableVideoes, err := dao.FindVideosPublishLatest(config.MaxCacheVideo)
	if err != nil {
		log.Println("从数据库读取视频失败")
		return nil, time.Now(), err
	}
	for i, j := range tableVideoes {
		//通过视频表中playurl和coveurl获得视频和封面在ftp服务器的位置
		tmpvideo := Video{}
		tmpvideo.TableVideo = j
		videos[i] = tmpvideo
	}
	return videos, time.Now(), nil
}

// 视频上传
func (v VideoServiceImpl) Publish(userId int64, videoname, imagename, title string) error {
	log.Printf("%v publish video", userId)
	err := dao.Save(userId, videoname, imagename, title)
	if err != nil {
		log.Printf("视频保存至数据库错误，err=%v", err)
		return err
	}
	//上传成功
	return nil
}

// 获取用户发布视频列表
func (v VideoServiceImpl) GetVideoPublishList(userId int64) ([]Video, error) {
	tablevideo, err := dao.FindVideoListbyUserId(userId)
	if err != nil {
		log.Printf("从数据库读数据失败，err = %v", err)
		return nil, err
	}
	video := make([]Video, config.MaxCacheVideo)
	log.Printf("video len == %v", len(video))
	for i, j := range tablevideo {
		video[i].TableVideo = j
	}
	return video, nil
}
