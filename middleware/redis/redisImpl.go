package redis

import (
	"SimpliftTikTok/config"
	"strconv"
)

// 把关注列表放入redis
func AddFollowingToRedis(userId int, ids []int64) {
	RdbFollowing.SAdd(Ctx, strconv.Itoa(userId), -1)
	for _, id := range ids {
		RdbFollowing.SAdd(Ctx, strconv.Itoa(userId), id)
	}
	// 更新following的过期时间
	RdbFollowing.Expire(Ctx, strconv.Itoa(userId), config.ExpireTime)
}

// 把粉丝列表放入redis
func AddFollowerToRedis(userId int, ids []int64) {
	RdbFollowers.SAdd(Ctx, strconv.Itoa(userId), -1)
	for _, id := range ids {
		RdbFollowers.SAdd(Ctx, strconv.Itoa(userId), id)
	}
	// 更新followers的过期时间。
	RdbFollowers.Expire(Ctx, strconv.Itoa(userId), config.ExpireTime)
}
