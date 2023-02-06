package controller

import (
	"SimpliftTikTok/dao"
	"SimpliftTikTok/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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

// Register POST tiktok/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

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
			println("Insert Data Fail")
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
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	usi := service.UserServiceImpl{}
	u := usi.GetTableUserByUsername(username)
	if u.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "用户名或密码错误，请重新登录"},
		})
	}
	if service.EnCoder(password) != u.Password {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "用户名或密码错误，请重新登录"},
		})
	}
	token := service.GenerateToken(username)
	log.Println("登录返回的id: ", u.Id)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   u.Id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	//按照抖音"我的信息"界面，需要显示抖音号，抖音名称，获赞，朋友，关注，粉丝等信息
	userId, _ := c.Get("userId")
	userid := fmt.Sprintf("%v", userId)
	id, _ := strconv.ParseInt(userid, 10, 64)

	usi := service.UserServiceImpl{
		service.FollowServiceImpl{},
		service.LikeServiceImpl{}}

	u := usi.GetTableUserById(id)
	if u.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 3, StatusMsg: "用户登录过期，请重新登录"},
		})
	}
	//
	userInfo, _ := usi.GetUserById(u.Id)
	userInfo.Name = u.Name
	log.Println("userInfo ========= ", userInfo)
	c.JSON(http.StatusOK, UserResponse{
		Response{StatusCode: 0},
		userInfo,
	})

}
