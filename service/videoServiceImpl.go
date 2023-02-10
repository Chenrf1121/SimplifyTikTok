package service

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/dao"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
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
		tmpvideo.Video = j
		videos[i] = tmpvideo
	}
	return videos, time.Now(), nil
}

// 视频上传
func (v VideoServiceImpl) Publish(userId int64, data *multipart.FileHeader, title string, picture *multipart.FileHeader) error {
	videofile, err := data.Open()
	imagefile, imageerr := picture.Open()
	defer videofile.Close()
	defer imagefile.Close()
	if err != nil {
		log.Printf("打开视频错误，data.open err,%v", err)
		return err
	}
	if imageerr != nil {
		log.Printf("打开封面错误，data.open err,%v", err)
		return err
	}

	//生成一个视频名称
	videoName := uuid.New().String()
	err = dao.VideoFTP(videofile, videoName)
	if err != nil {
		log.Printf("视频上传失败,dao.VideoFTP方法失败，err = %v", err)
		return err
	}
	imageName := uuid.New().String()
	err = dao.ImageFTP(imagefile, imageName)
	err = dao.Save(userId, videoName, imageName, title)
	if err != nil {
		log.Printf("视频保存至数据库错误，err=%v", err)
		return err
	}
	//上传成功
	return nil
}
