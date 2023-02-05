package service

import "SimpliftTikTok/dao"

type FollowServiceImpl struct {
}

func (f FollowServiceImpl) IsFollowing(userId int64, targetId int64) (bool, error) {
	return false, nil
}

// GetFollowerCnt 根据用户id来查询用户被多少其他用户关注
func (f FollowServiceImpl) GetFollowerCnt(userId int64) (int64, error) {
	ans, err := dao.GetFollowerCnt(userId)
	if err != nil {
		ans = 0
	}
	return ans, nil
}

// GetFollowingCnt 根据用户id来查询用户关注了多少其它用户
func (f FollowServiceImpl) GetFollowingCnt(userId int64) (int64, error) {
	return 0, nil
}

/*
   二、直接request需要的业务方法
*/
// AddFollowRelation 当前用户关注目标用户
func (f FollowServiceImpl) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	return false, nil
}

// DeleteFollowRelation 当前用户取消对目标用户的关注
func (f FollowServiceImpl) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	return false, nil
}

// GetFollowing 获取当前用户的关注列表
func (f FollowServiceImpl) GetFollowing(userId int64) ([]User, error) {
	return nil, nil
}

// GetFollowers 获取当前用户的粉丝列表
func (f FollowServiceImpl) GetFollowers(userId int64) ([]User, error) {
	return nil, nil
}
