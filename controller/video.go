package controller

import (
	"SimpliftTikTok/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//有关视频的功能

// 刷视频
func Feed(c *gin.Context) {
	videoService := GetVideo()
	videoService.Feed()
	c.JSON(
		http.StatusOK,
		VideoResponse{},
	)
}

// 发布视频
func Publish(c *gin.Context) {
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", c.GetString("userId")), 10, 64)

	data, err := c.FormFile("data") //获取视频
	if err != nil {
		log.Printf("上传视频流失败:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//获取封面，如果封面为空，就默认用视频的第一针作为封面
	picture, err := c.FormFile("picture")
	if err != nil {
		//没用封面
	}

	if data.Size > 1024*1024*100 {
		//如果视频超过100M，就无法上传
		log.Println("视频过大，上传失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title") //获取标题
	//封装视频接口
	videoServiceImpl := GetVideo()
	//发布视频到ftp
	err = videoServiceImpl.Publish(userId, data, title, picture)
	if err != nil {
		log.Printf("视频上传至文件服务器失败，%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Printf("视频发布成功，发布用户,%v", userId)
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "视频发布成功！",
	})
	return
}

// GetVideo 拼装videoService
func GetVideo() service.VideoServiceImpl {
	var videoService service.VideoServiceImpl

	return videoService
}
