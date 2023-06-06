package main

import (
	"log"
	"os/exec"
)

func main() {
	videoPath := "./source/videos/1686040774529576000-video1_sample.mp4"
	// 使用 ffmpeg 将视频流推送到流媒体服务器
	cmdArgs := []string{
		"-i", videoPath,
		"-c:v", "copy",
		"-f", "flv",
		"rtmp://127.0.0.1/videos/test.mp4",
	}
	cmd := exec.Command("ffmpeg", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to push video stream: %s\n", string(output))
		//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push video stream"})
		return
	}

}
