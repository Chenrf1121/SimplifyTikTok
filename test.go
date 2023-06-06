package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		fileHeader, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		videoFile, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to open video file",
			})
			return
		}
		defer videoFile.Close()

		// 创建一个临时文件来保存视频文件
		tmpVideoFile, err := os.CreateTemp("", "video.*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create temporary video file",
			})
			return
		}
		defer os.Remove(tmpVideoFile.Name())
		defer tmpVideoFile.Close()

		// 将上传的视频文件保存到临时文件中
		_, err = io.Copy(tmpVideoFile, videoFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save video file",
			})
			return
		}

		// 定义 FFmpeg 命令行参数
		outputFile := "./output.jpg"
		ffmpegArgs := []string{
			"-i", tmpVideoFile.Name(),
			"-ss", "0", // 提取第一帧
			"-vframes", "1",
			"-vf", "scale=320:-1", // 设置缩略图尺寸
			outputFile,
		}

		// 创建一个 cmd 对象，指定 FFmpeg 命令及参数
		cmd := exec.Command("ffmpeg", ffmpegArgs...)

		// 执行命令
		err = cmd.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to extract thumbnail from video",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Thumbnail extracted successfully",
			"image":   outputFile,
		})
	})

	r.Run(":8081")
}
