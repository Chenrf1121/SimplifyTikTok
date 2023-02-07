package service

import "mime/multipart"

type VideoServiceImpl struct {
	UserService
	LikeService
	CommentService
}

// 视频上传
func (v VideoServiceImpl) Publish(userId int64, data *multipart.FileHeader, title string) error {
	return nil
}
