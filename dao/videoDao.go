package dao

import (
	"SimpliftTikTok/config"
	"log"
	"time"
)

type MetaVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	VideoId     int64
	Title       string    `json:"title"`
	PlayUrl     string    `json:"playUrl"`
	CoverUrl    string    `json:"coverUrl"`
	PublishTime time.Time `json:"publishTime"`
}
type DynamicVideo struct {
	Id         int64 `json:"id"`
	VideoId    int64
	LikeNum    int
	FavourNum  int
	CommentNum int
	ShareNum   int
}

func Save(authorId, videoId int64, videoName, imageName, title string) error {
	tmpvideo := MetaVideo{
		AuthorId:    authorId,
		VideoId:     videoId,
		Title:       title,
		PlayUrl:     config.PlayUrl + videoName,
		CoverUrl:    config.CoverUrl + imageName,
		PublishTime: time.Now(),
	}
	result := Db.Save(&tmpvideo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func SaveDynamicVideo(videoId int64) error {
	tmpvideo := DynamicVideo{
		VideoId: videoId,
	}
	result := Db.Save(&tmpvideo)
	if result.Error != nil {
		return result.Error
	}
	return nil

}
func FindVideosPublishLatest(n int) ([]MetaVideo, error) {
	videoList := make([]MetaVideo, n)
	Db.Debug().Order("publish_time").Limit(n).Find(&videoList)
	return videoList, nil
}

// 根据用户id查所有视频
func FindVideoListbyUserId(userId int64) ([]MetaVideo, error) {
	videoList := make([]MetaVideo, config.MaxCacheVideo)
	result := Db.Debug().Where("author_id = ?", userId).Find(&videoList).Limit(config.MaxCacheVideo)
	if result.Error != nil {
		log.Printf("搜索视频数据库失败")
		return nil, result.Error
	}
	log.Printf("数据库的视频长度%v", len(videoList))
	return videoList, nil
}

// 根据videoid查视频
func FindVideobyVideoID(videoId int64) (MetaVideo, error) {
	targetVideo := MetaVideo{}
	result := Db.Debug().Where("VideoID = ?", videoId).Find(&targetVideo)
	if result.Error != nil {
		log.Printf("搜索视频号为：%v 失败", videoId)
		return targetVideo, result.Error
	}
	return targetVideo, nil
}
