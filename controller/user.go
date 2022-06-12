package controller

import (
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
)

// UsersLoginInfo use map to store user info, and key is username+password for demo
var UsersLoginInfo = Dbm.GerAllUser()

var userIdSequence = Dbm.GetLastId()

type UserLoginResponse struct {
	_type.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	_type.Response
	User _type.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if _, exist := UsersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: _type.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := _type.User{
			Id:            userIdSequence,
			Password:      password,
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		//usersLoginInfo[token] = newUser
		Dbm.AddUser(newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: _type.Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	UsersLoginInfo = Dbm.GerAllUser()

	if user, exist := UsersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: _type.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	UsersLoginInfo = Dbm.GerAllUser()

	if user, exist := UsersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: _type.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
