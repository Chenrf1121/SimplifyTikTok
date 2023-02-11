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
	//首先查看redis
	following, err := redis.RdbFollowing.SMembers(redis.Ctx, strconv.Itoa(int(userId))).Result()
	if err != nil {
		log.Printf("访问redis失败")
	}
	redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
	for i := range following {
		if following[i] == strconv.Itoa(int(targetId)) {
			//redis存在
			return true, nil
		}
	}
	//没有则查看mysql
	tag := false
	followingList, err := dao.NewFollowDaoInstance().GetFollowingCnt(userId)
	for i := range followingList {
		if followingList[i] == targetId {
			tag = true
			break
		}
	}
	//更新redis
	go redis.AddFollowingToRedis(int(userId), followingList)
	return tag, nil
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
	go redis.AddFollowerToRedis(int(userId), ids)
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
	go redis.AddFollowingToRedis(int(userId), ids)
	return int64(len(ids)), nil
}

// 添加关注
func (f FollowServiceImpl) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	return false, nil
}

// 取消关注
func (f FollowServiceImpl) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	return false, nil
}

// 获取关注
func (f FollowServiceImpl) GetFollowing(userId int64) ([]User, error) {
	return nil, nil
}

// 获取粉丝
func (f FollowServiceImpl) GetFollowers(userId int64) ([]User, error) {
	return nil, nil
}
