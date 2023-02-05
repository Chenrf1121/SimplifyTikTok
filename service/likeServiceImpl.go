package service

type LikeServiceImpl struct {
}

// IsFavorite 根据当前视频id判断是否点赞了该视频。
func (t LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	return false, nil
}

// FavouriteCount 根据当前视频id获取当前视频点赞数量。
func (t LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	return 0, nil
}

// TotalFavourite 根据userId获取这个用户总共被点赞数量
func (t LikeServiceImpl) TotalFavourite(userId int64) (int64, error) {
	return 0, nil
}

// FavouriteVideoCount 根据userId获取这个用户点赞视频数量
func (t LikeServiceImpl) FavouriteVideoCount(userId int64) (int64, error) {
	return 0, nil
}

/*
   2.request需要实现的功能
*/
//当前用户对视频的点赞操作 ,并把这个行为更新到like表中。
//当前操作行为，1点赞，2取消点赞。
func (t LikeServiceImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	return nil
}

// GetFavouriteList 获取当前用户的所有点赞视频，调用videoService的方法/
func (t LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	return nil, nil
}
