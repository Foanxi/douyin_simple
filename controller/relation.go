package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserListResponse struct {
	_type.Response
	UserList []_type.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	favouriteUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	isSuccess := service.RelationActionService(token, favouriteUserId, actionType)
	if isSuccess {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 1, StatusMsg: "Can't focus on yourself"})
	}
}

// FollowList all users have same followed list
func FollowList(c *gin.Context) {
	//获取use_id
	useId := c.Query("user_id")
	userList := service.FollowListService(useId)
	c.JSON(http.StatusOK, UserListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	useId := c.Query("user_id")
	userFanList := service.FollowerListService(useId)
	c.JSON(http.StatusOK, UserListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		UserList: userFanList,
	})
}
