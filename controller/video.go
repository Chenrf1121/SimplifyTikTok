package controller

import (
	"SimpliftTikTok/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//有关视频的功能

// 刷视频
func Feed(c *gin.Context) {
	videoService := GetVideo()
	videos, latestTime, err := videoService.Feed()
	if err != nil {
		log.Printf("视频读取错误，err=%v", err)
		c.JSON(http.StatusOK,
			Response{
				http.StatusOK,
				"视频读取错误",
			})
		return
	}
	c.JSON(
		http.StatusOK,
		VideoResponse{
			latestTime.String(),
			videos,
		},
	)
}

// 发布视频
func Publish(c *gin.Context) {
	//	c.HTML(http.StatusOK, "/web/publish.html", nil)
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", c.GetString("userId")), 10, 64)
	//获取封面，如果封面为空，就默认用视频的第一针作为封面
	userId = 1
	picture, err := c.FormFile("picture")
	if err != nil {
		//没用封面
		c.JSON(http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  "没有封面",
			})
	}
	data, err := c.FormFile("data") //获取视频
	title := c.PostForm("title")    //获取标题
	if err != nil {
		log.Printf("上传视频流失败:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 确定视频文件的存储路径和文件名
	// 这里假设使用当前时间戳作为文件名
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + title + "_" + data.Filename
	storageVideoPath := "./source/videos/" + filename

	if title == "" {
		filepoint, err := data.Open() //打开文件
		defer filepoint.Close()
		if err != nil {
			log.Println("打开视频文件失败")
			c.HTML(200, "打开视频文件失败", nil)
			return
		}
		imagesPath := "./source/images/" + filename + ".jpg"
		ffmpegArgs := []string{
			"-i", "pipe:0",
			"-ss", "0", // 提取第一帧
			"-vframes", "1",
			"-vf", "scale=320:-1", // 设置缩略图尺寸
			imagesPath,
		}
		cmd := exec.Command("ffmpeg", ffmpegArgs...)
		// 将视频文件作为输入管道
		cmd.Stdin = filepoint

		// 将命令输出连接到标准输出和标准错误输出
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// 执行命令
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		defer filepoint.Close()

	}
	// 保存视频文件到本地
	if err := c.SaveUploadedFile(data, storageVideoPath); err != nil {
		log.Printf("视频保存失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "视频保存失败"})
		return
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
	return
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

// 用户视频发布列表
func PublishList(c *gin.Context) {
	userId := fmt.Sprintf("%v", c.GetString("userId"))
	//获取要观看哪个用户的视频列表
	targetId := c.DefaultQuery("userId", userId)
	//得到id后需要去数据库查看
	targetuserId, _ := strconv.ParseInt(targetId, 10, 64)
	videoservice := GetVideo()
	videolist, err := videoservice.GetVideoPublishList(targetuserId)
	if err != nil {
		log.Printf("获取用户视频列表失败，err = %v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(
		http.StatusOK,
		VideoResponse{VideoList: videolist},
	)
}

// GetVideo 拼装videoService
func GetVideo() service.VideoServiceImpl {
	var videoService service.VideoServiceImpl

	return videoService
}
