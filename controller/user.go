package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserResponse struct {
	_type.Response
	User _type.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	response := service.RegisterService(username, password)
	c.JSON(http.StatusOK, response)
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	response := service.LoginService(username, password)
	c.JSON(http.StatusOK, response)
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	authorId := c.Query("user_id")
	user := service.UserInfoService(token, authorId)
	c.JSON(http.StatusOK, UserResponse{
		Response: _type.Response{StatusCode: 0},
		User:     user,
	})
}
