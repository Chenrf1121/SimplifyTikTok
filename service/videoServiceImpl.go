package service

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/dao"
	"github.com/google/uuid"
	"log"
	"time"
)

type VideoServiceImpl struct {
	UserService
	LikeService
	CommentService
}

// 获取指定videoid的视频
func (v VideoServiceImpl) GetVideo(video_id int64) (Video, error) {
	targetVideo := Video{}
	daoVideo, err := dao.FindVideobyVideoID(video_id)
	if err != nil {
		return targetVideo, err
	}
	targetVideo.MetaVideo = daoVideo
	return targetVideo, nil
}

// 刷视频
func (v VideoServiceImpl) Feed() ([]Video, time.Time, error) {

	//从数据库得到视频信息
	tableVideoes, err := dao.FindVideosPublishLatest(config.MaxCacheVideo)
	//从文件服务器读视频
	videos := make([]Video, len(tableVideoes))
	if err != nil {
		log.Println("从数据库读取视频失败")
		return nil, time.Now(), err
	}
	for i, j := range tableVideoes {
		//通过视频表中playurl和coveurl获得视频和封面在nginx流媒体服务器的位置
		tmpvideo := Video{}
		tmpvideo.MetaVideo = j
		videos[i] = tmpvideo
	}
	return videos, time.Now(), nil
}

// 视频上传
func (v VideoServiceImpl) Publish(userId int64, videoname, imagename, title string) error {
	videoUuid := uuid.New()
	videoId := int64(videoUuid.ID())
	err := dao.Save(userId, videoId, videoname, imagename, title)
	if err != nil {
		log.Printf("视频元数据保存至数据库错误，err=%v", err)
		return err
	}
	err = dao.SaveDynamicVideo(videoId)
	if err != nil {
		log.Printf("视频动态数据保存至数据库错误，err=%v", err)
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
	for i := range tablevideo {
		video[i].MetaVideo = tablevideo[i]
	}

	return video, nil
}
