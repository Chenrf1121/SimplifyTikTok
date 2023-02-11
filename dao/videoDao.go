package dao

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/middleware/ftp"
	"io"
	"log"
	"time"
)

type Video struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	Title       string    `json:"title"`
	PlayUrl     string    `json:"playUrl"`
	CoverUrl    string    `json:"coverUrl"`
	PublishTime time.Time `json:"publishTime"`
}

// VideoFTP
// 通过ftp将视频传入服务器
func VideoFTP(file io.Reader, videoName string) error {
	//进入视频地址
	err := ftp.MyFTP.Cwd("video")
	if err != nil {
		ftp.MyFTP.Mkd("video")
	}
	err = ftp.MyFTP.Stor(videoName+".mp4", file)
	if err != nil {
		log.Printf("上传至服务器失败,err = %v", err)
	}
	return nil
}
func ImageFTP(file io.Reader, imageName string) error {
	//进入视频地址
	err := ftp.MyFTP.Cwd("image")
	if err != nil {
		ftp.MyFTP.Mkd("image")
	}
	err = ftp.MyFTP.Stor(imageName+".jpg", file)
	if err != nil {
		log.Printf("封面上传至服务器失败,err = %v", err)
	}
	return nil
}

func Save(authorId int64, videoName, imageName, title string) error {
	tmpvideo := Video{
		AuthorId:    authorId,
		Title:       title,
		PlayUrl:     config.PlayUrl + videoName + ".mp4",
		CoverUrl:    config.CoverUrl + imageName + ".png",
		PublishTime: time.Now(),
	}
	result := Db.Save(&tmpvideo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func FindVideosPublishLatest(n int) ([]Video, error) {
	videoList := make([]Video, n)
	Db.Debug().Order("publish_time").Limit(n).Find(&videoList)
	return videoList, nil
}
func FindVideoListbyUserId(userId int64) ([]Video, error) {
	videoList := make([]Video, config.MaxCacheVideo)
	result := Db.Debug().Where("author_id = ?", userId).Find(&videoList).Limit(config.MaxCacheVideo)
	if result.Error != nil {
		log.Printf("搜索视频数据库失败")
		return nil, result.Error
	}
	log.Printf("数据库的视频长度%v", len(videoList))
	return videoList, nil
}
