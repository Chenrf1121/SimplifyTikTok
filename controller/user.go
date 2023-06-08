package controller

import (
	"SimpliftTikTok/dao"
	"SimpliftTikTok/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}
type UserBaseInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register POST tiktok/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println("user,pass = ", username, password)
	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response{StatusCode: 200, StatusMsg: "输入用户名和密码"},
			-1, "",
		})
		return
	}

	usi := service.UserServiceImpl{}

	u := usi.GetTableUserByUsername(username)
	if username == u.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := dao.TableUser{
			Name:     username,
			Password: service.EnCoder(password),
		}
		if usi.InsertTableUser(&newUser) != true {
			log.Println("Insert Data Fail")
		}
		u := usi.GetTableUserByUsername(username)
		token := service.GenerateToken(username)
		log.Println("注册返回的id: ", u.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

var (
	maxLoginAttempts = 4
	lockoutDuration  = time.Minute
	failedAttempts   = make(map[string]int)
	lockoutTimes     = make(map[string]time.Time)
	lockedUsers      = make(map[string]bool)
	failedIPs        = make(map[string]int)
	lockedIPs        = make(map[string]time.Time)
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	usi := service.UserServiceImpl{}
	u := usi.GetTableUserByUsername(username)
	if u.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "用户名或密码错误，请重新登录"},
		})
		return
	}
	ip := c.ClientIP()
	if lockedIPs[ip].After(time.Now()) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "IP地址已被锁定，请稍后再试"},
		})
		return
	}
	if failedIPs[ip] >= maxLoginAttempts {
		lockedIPs[ip] = time.Now().Add(24 * time.Hour)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "IP地址已被锁定，请稍后再试"},
		})
		return
	}
	if lockedUsers[username] {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "账号已被锁定，请稍后再试"},
		})
		return
	}
	if service.EnCoder(password) != u.Password {
		failedAttempts[username]++
		failedIPs[ip]++
		if failedAttempts[username] >= maxLoginAttempts {
			lockedUsers[username] = true
			lockoutTimes[username] = time.Now()
		}
		if failedIPs[ip] >= maxLoginAttempts {
			lockedIPs[ip] = time.Now().Add(24 * time.Hour)
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "用户名或密码错误，请重新登录"},
		})
		return
	}
	failedAttempts[username] = 0
	failedIPs[ip] = 0
	lockedUsers[username] = false
	lockedIPs[ip] = time.Time{}
	token := service.GenerateToken(username)
	log.Println("登录返回的id: ", u.Id)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   u.Id,
		Token:    token,
	})
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			for username, t := range lockoutTimes {
				if time.Since(t) > lockoutDuration {
					failedAttempts[username] = 0
					lockedUsers[username] = false
				}
			}
			for ip, t := range lockedIPs {
				if t.Before(time.Now()) {
					delete(lockedIPs, ip)
				}
			}
		}
	}()
}

func UserInfo(c *gin.Context) {
	//按照抖音"我的信息"界面，需要显示抖音号，抖音名称，获赞，朋友，关注，粉丝等信息
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", c.GetString("userId")), 10, 64)
	usi := service.UserServiceImpl{
		service.FollowServiceImpl{},
		service.LikeServiceImpl{}}
	u := usi.GetTableUserById(userId)
	if u.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 3, StatusMsg: "用户登录过期，请重新登录"},
		})
	}
	userInfo, _ := usi.GetUserById(u.Id)
	userInfo.Name = u.Name
	c.JSON(http.StatusOK, UserResponse{
		Response{StatusCode: 0},
		userInfo,
	})
}

// 用户关注和不关注
func FollowAction(c *gin.Context) {
	//目标id从前端读取，自身id通过token读取
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", c.GetString("userId")), 10, 64)
	targetId, _ := strconv.ParseInt(c.Query("targetId"), 10, 64)
	usi := service.UserServiceImpl{}
	//首先判断他是不是在关注列表里
	tag, err := usi.IsFollowing(userId, targetId)
	if err != nil {
		log.Printf("查看关注用户失败，err=%v", err)
		c.JSON(http.StatusOK,
			"查看关注用户失败")
		return
	}
	//进行操作
	if tag == false {
		_, err := usi.AddFollowRelation(userId, targetId)
		if err != nil {
			log.Printf("用户%v，关注目标用户%v 失败", userId, targetId)
			c.JSON(http.StatusOK,
				"网络延迟，操作失败")
		}
		c.JSON(http.StatusOK, "关注成功")
	} else {
		_, err := usi.DeleteFollowRelation(userId, targetId)
		if err != nil {
			log.Printf("用户%v，取消关注目标用户%v 失败", userId, targetId)
			c.JSON(http.StatusOK,
				"网络延迟，操作失败")
		}
		c.JSON(http.StatusOK, "取消关注成功")
	}
	return
}
