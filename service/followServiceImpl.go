package service

import (
	"SimpliftTikTok/config"
	"SimpliftTikTok/dao"
	"SimpliftTikTok/middleware/redis"
	"log"
	"strconv"
)

type FollowServiceImpl struct {
}

func (f FollowServiceImpl) IsFollowing(userId int64, targetId int64) (bool, error) {
	return false, nil
}

// GetFollowerCnt 根据用户id来查询用户被多少其他用户关注
func (f FollowServiceImpl) GetFollowerCnt(userId int64) (int64, error) {
	if cnt, err := redis.RdbFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}
	ids, err := dao.NewFollowDaoInstance().GetFollowerCnt(userId)
	if err != nil {
		log.Println("查询数据库失败，err = ", err)
		return 0, err
	}
	//更新redis
	go addFollowerToRedis(int(userId), ids)
	return int64(len(ids)), nil
}

// GetFollowingCnt 根据用户id来查询用户关注了多少其它用户
func (f FollowServiceImpl) GetFollowingCnt(userId int64) (int64, error) {
	//SCard获取集合的长度
	if cnt, err := redis.RdbFollowing.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}
	//查不到需要查询mysql
	ids, err := dao.NewFollowDaoInstance().GetFollowingCnt(userId)
	if err != nil {
		log.Println("查询数据库失败，err = ", err)
		return 0, err
	}
	//更新redis
	go addFollowingToRedis(int(userId), ids)
	return int64(len(ids)), nil
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

// 把关注列表放入redis
func addFollowingToRedis(userId int, ids []int64) {
	redis.RdbFollowing.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	for _, id := range ids {
		redis.RdbFollowing.SAdd(redis.Ctx, strconv.Itoa(userId), id)
	}
	// 更新following的过期时间
	redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(userId), config.ExpireTime)
}

// 把粉丝列表放入redis
func addFollowerToRedis(userId int, ids []int64) {
	redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	for _, id := range ids {
		redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), id)
	}
	// 更新followers的过期时间。
	redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(userId), config.ExpireTime)
}
