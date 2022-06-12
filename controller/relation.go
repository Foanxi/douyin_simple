package controller

import (
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
	//token := c.Query("token")
	token := c.Query("token")
	favouriteUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	statue := Dbm.UserFavoriteUser(token, favouriteUserId, actionType)
	if statue == true {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 1, StatusMsg: "Can't focus on yourself"})
	}
}

// FollowList all users have same followed list
func FollowList(c *gin.Context) {
	var userList []_type.User
	//获取use_id
	useId := c.Query("user_id")
	userList = Dbm.GetAuthorById(useId)
	c.JSON(http.StatusOK, UserListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		UserList: []_type.User{DemoUser},
	})
}
