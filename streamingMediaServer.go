package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	// 创建HTTP处理函数
	http.HandleFunc("/stream", handleStream)

	// 启动HTTP服务器
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	// 处理流媒体请求
	// 在此处调用FFmpeg处理视频流，并将数据写入ResponseWriter
	// 根据需要设置适当的HTTP头部
	// 调用FFmpeg处理视频流
	cmd := exec.Command("ffmpeg", "-i", "input.mp4", "-c", "copy", "-f", "mpegts", "pipe:1")
	output, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("FFmpeg stdout pipe error: ", err)
		return
	}

	// 启动FFmpeg进程
	err = cmd.Start()
	if err != nil {
		log.Println("FFmpeg start error: ", err)
		return
	}

	// 将FFmpeg输出写入ResponseWriter
	_, err = io.Copy(w, output)
	if err != nil {
		log.Println("Copying FFmpeg output error: ", err)
		return
	}

	// 等待FFmpeg进程退出
	err = cmd.Wait()
	if err != nil {
		log.Println("FFmpeg exited with error: ", err)
	}

}
