package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()
	router.POST("/upload", handleUpload)
	router.GET("/videos/:filename", handleStream)
	router.GET("/videos1/:filename", handleStream1)
	router.Run(":8082")
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to upload file"})
		return
	}

	filename := filepath.Base(file.Filename)
	dst := filepath.Join("videos", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}

func handleStream(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("videos", filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	c.Writer.Header().Set("Content-Type", "video/mp4")
	c.Writer.Header().Set("Content-Length", string(stat.Size()))
	c.Writer.Header().Set("Accept-Ranges", "bytes")

	io.Copy(c.Writer, file)
}

func handleStream1(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("videos", filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	buffer := make([]byte, stat.Size())
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	c.Writer.Header().Set("Content-Type", "video/mp4")
	c.Writer.Header().Set("Content-Length", string(stat.Size()))
	c.Writer.Header().Set("Accept-Ranges", "bytes")

	c.Writer.Write(buffer)
}
