package controller

import (
	"github.com/RaymondCode/simple-demo/Dao"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync/atomic"
)

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

	if _, exist := Dao.Udi.GerAllUser()[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: _type.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		var userIdSequence = Dao.Udi.GetLastId()
		atomic.AddInt64(&userIdSequence, 1)
		newUser := _type.User{
			Id:            userIdSequence,
			Password:      password,
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		Dao.Udi.AddUser(newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: _type.Response{StatusCode: 0},
			UserId:   Dao.Udi.GetLastId(),
			Token:    username + password,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	if user, exist := Dao.Udi.GerAllUser()[token]; exist {
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
	authorId := c.Query("user_id")
	user := Dao.Udi.GerAllUser()[token]
	id, _ := strconv.ParseInt(authorId, 10, 8)
	user.IsFollow = Dao.Udi.GetUserRelation(id, user.Id)

	c.JSON(http.StatusOK, UserResponse{
		Response: _type.Response{StatusCode: 0},
		User:     user,
	})
}
