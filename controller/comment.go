package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//给视频评论
func AddVideoComment(c *gin.Context) {

}

// LikeVideo handles the route for liking a video
func LikeVideo(c *gin.Context) {
	// Parse request data
	videoID := c.Param("video_id")

	// Update the video's like count in the database
	// ...

	c.JSON(http.StatusOK, gin.H{
		"message": "Video liked successfully" + videoID,
	})
}

// FavoriteVideo handles the route for favoriting a video
func FavoriteVideo(c *gin.Context) {
	// Parse request data
	videoID := c.Param("video_id")

	// Add the video to the user's favorites in the database
	// ...

	c.JSON(http.StatusOK, gin.H{
		"message": "Video favorited successfully" + videoID,
	})
}

// ShareVideo handles the route for sharing a video
func ShareVideo(c *gin.Context) {
	// Parse request data
	videoID := c.Param("video_id")

	// Add the video to the user's shares in the database
	// ...

	c.JSON(http.StatusOK, gin.H{
		"message": "Video shared successfully" + videoID,
	})
}

// ReplyToComment handles the route for replying to a comment
func ReplyToComment(c *gin.Context) {
	// Parse request data
	parentID := c.Param("comment_id")
	content := c.PostForm("content")

	// Create a new comment as a reply
	comment := Comment{
		ParentID: parentID,
		Content:  content,
		// Set other comment fields
	}
	log.Println(comment)
	// Save the reply comment to the database
	// ...

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment replied successfully",
	})
}

// LikeComment handles the route for liking a comment
func LikeComment(c *gin.Context) {
	// Parse request data
	commentID := c.Param("comment_id")
	log.Printf(commentID)
	// Update the comment's like count in the database
	// ...

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment liked successfully",
	})
}

// GetReplies handles the route for retrieving replies to a comment
func GetReplies(c *gin.Context) {
	// Parse request data
	commentID := c.Param("comment_id")
	log.Printf(commentID)

	// Retrieve the replies to the comment from the database
	// ...

	replies := []Comment{
		// Populate the replies data
	}

	c.JSON(http.StatusOK, gin.H{
		"replies": replies,
	})
}

// GetVideoComments handles the route for retrieving comments for a video
func GetVideoComments(c *gin.Context) {
	// Parse request data
	videoID := c.Param("video_id")
	log.Printf(videoID)

	// Retrieve the comments for the video from the database
	// ...

	comments := []Comment{
		// Populate the comments data
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
