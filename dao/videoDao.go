package dao

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/middleware/ftp"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
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
func aa() {
	filename := "test.mp4"
	width := 640
	height := 360
	cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	if cmd.Run() != nil {
		panic("could not generate frame")
	}
}
