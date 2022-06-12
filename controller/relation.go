package controller

import (
	"github.com/RaymondCode/simple-demo/Dao"
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
	statue := Dao.Fdi.UserFavoriteUser(token, favouriteUserId, actionType)
	if statue {
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
	//进入service层
	userList = Dao.Udi.GetAuthorById(useId)
	c.JSON(http.StatusOK, UserListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	var userFanList []_type.User
	useId := c.Query("user_id")
	userFanList = Dao.Udi.GetFanList(useId)
	c.JSON(http.StatusOK, UserListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		UserList: userFanList,
	})
}
