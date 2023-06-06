package config

import "time"

// Secret 密钥
var Secret = "TikTok"

// OneDayOfHours 时间
var OneDayOfHours = 60 * 60 * 24
var OneMinute = 60 * 1
var OneMonth = 60 * 60 * 24 * 30
var OneYear = 365 * 60 * 60 * 24

// IP
// mysqlIp
var IP = "10.252.45.17"
var MysqlPort = ":3306"

// redis
var RedisPort = ":6379"
var RedisPassword = "123456"
var ExpireTime = time.Hour * 48 //过期时间

// ftpconfig
var FtpConfig = IP + ":21"
var FtpUser = "admin"
var FtpPwd = "123456"
var HeartbeatTime = 60

// video
var PlayUrl = "http://" + IP + "/videos/"
var CoverUrl = "http://" + IP + "/images/"
var MaxCacheVideo = 4
