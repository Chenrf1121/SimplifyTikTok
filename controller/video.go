package controller

import (
	"SimpliftTikTok/common"
	"SimpliftTikTok/config"
	"SimpliftTikTok/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//有关视频的功能

// 刷视频
func Feed(c *gin.Context) {
	c.Header("Pragma", "no-cache")
	video_id_str := c.DefaultQuery("video_id", "-1")
	video_id, err := strconv.ParseInt(video_id_str, 10, 64)
	videoService := GetVideo()
	videos, _, err := videoService.Feed()

	if err != nil {
		log.Printf("视频读取错误，err=%v", err)
		c.JSON(http.StatusOK,
			Response{
				http.StatusOK,
				"视频读取错误",
			})
		return
	}
	if video_id != -1 {
		//查视频
		videoService.GetVideo(video_id)
	}
	c.HTML(200, "m3u8.html", gin.H{
		"Videos": videos,
	})

}

// 发布视频
func Publish(c *gin.Context) {
	//	c.HTML(http.StatusOK, "/web/publish.html", nil)
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", c.GetString("userId")), 10, 64)
	//获取封面，如果封面为空，就默认用视频的第一针作为封面
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
	filename = common.HashString(filename)
	storageVideoPath := config.Pwd + "/source/videos/" + filename
	storageVideoPath = strings.Split(storageVideoPath, ".")[0]
	log.Println("storageVideoPath = ", storageVideoPath)
	imagesPath := config.Pwd + "/source/images/" + filename + ".jpg"
	if picture == nil {
		err = getFrame(data, imagesPath)
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Thumbnail extracted successfully",
				"image":   imagesPath,
			})
		}
	} else {
		if err = c.SaveUploadedFile(picture, imagesPath); err != nil {
			log.Printf("封面保存失败：%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "封面保存失败"})
			return
		}
	}
	//将视频转为m3u8类型
	go convertToM3U8(data, storageVideoPath, "playlist.m3u8")

	//封装视频接口
	videoServiceImpl := GetVideo()
	log.Printf("%v publish video", userId)
	err = videoServiceImpl.Publish(userId, filename+"/playlist.m3u8", filename+".jpg", title)
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

// 获取视频某一帧
func getFrame(data *multipart.FileHeader, imagesPath string) error {
	videoFile, err := data.Open()
	if err != nil {
		return errors.New("Failed to open video file")
	}
	defer videoFile.Close()

	// 创建一个临时文件来保存视频文件
	tmpVideoFile, err := os.CreateTemp("", "video.*")
	if err != nil {
		return errors.New("Failed to create temporary video file")
	}
	defer os.Remove(tmpVideoFile.Name())
	defer tmpVideoFile.Close()

	// 将上传的视频文件保存到临时文件中
	_, err = io.Copy(tmpVideoFile, videoFile)
	if err != nil {
		return errors.New("Failed to save video file")
	}

	// 定义 FFmpeg 命令行参数
	ffmpegArgs := []string{
		"-i", tmpVideoFile.Name(),
		"-ss", "0", // 提取第一帧
		"-vframes", "1",
		"-vf", "scale=320:-1", // 设置缩略图尺寸
		imagesPath,
	}

	// 创建一个 cmd 对象，指定 FFmpeg 命令及参数
	cmd := exec.Command("ffmpeg", ffmpegArgs...)

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return errors.New("Failed to extract thumbnail from video")
	}
	return nil
}

// 将视频转为m3u8格式
func convertToM3U8(data *multipart.FileHeader, outputDir, outputFile string) error {

	// 创建文件夹
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		return errors.New("outputdir already exist!")
	}
	videoFile, err := data.Open()
	if err != nil {
		return errors.New("Failed to open video file")
	}
	defer videoFile.Close()

	// 创建一个临时文件来保存视频文件
	tmpVideoFile, err := os.CreateTemp("", "video.*")
	if err != nil {
		return errors.New("Failed to create temporary video file")
	}
	defer os.Remove(tmpVideoFile.Name())
	defer tmpVideoFile.Close()

	// 将上传的视频文件保存到临时文件中
	_, err = io.Copy(tmpVideoFile, videoFile)
	if err != nil {
		return errors.New("Failed to save video file")
	}

	cmd := exec.Command("ffmpeg", "-i", tmpVideoFile.Name(), "-c:v", "libx264", "-c:a", "aac", "-hls_time", "10", "-hls_list_size", "0", "-hls_segment_filename", outputDir+"/%03d.ts", outputDir+"/"+outputFile)

	// 设置输出到标准输出，便于查看转换过程
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
