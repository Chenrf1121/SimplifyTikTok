package redis

import (
	"SimpliftTikTok/config"
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var RdbFollowers *redis.Client  //粉丝信息
var RdbFollowing *redis.Client  // 关注信息
var RdbGetlikes *redis.Client   // 获赞数
var RdbUserVideos *redis.Client // 用户发表的视频
var RdbUserInfo *redis.Client   // 用户名和账号
func InitRedis() {
	RdbFollowers = redis.NewClient(&redis.Options{
		Addr:     config.IP + config.RedisPort,
		Password: config.RedisPassword,
		DB:       0, // 粉丝列表信息存入 DB0.
	})
	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     config.IP + config.RedisPort,
		Password: config.RedisPassword,
		DB:       1, // 关注列表信息存入 DB1.
	})
	RdbGetlikes = redis.NewClient(&redis.Options{
		Addr:     config.IP + config.RedisPort,
		Password: config.RedisPassword,
		DB:       2, // 获赞数存入 DB2.
	})
	RdbUserVideos = redis.NewClient(&redis.Options{
		Addr:     config.IP + config.RedisPort,
		Password: config.RedisPassword,
		DB:       3, // 用户视频存入 DB3.
	})
	RdbUserInfo = redis.NewClient(&redis.Options{
		Addr:     config.IP + config.RedisPort,
		Password: config.RedisPassword,
		DB:       4, // 用户抖音号抖音名信息存入 DB4.
	})
}
