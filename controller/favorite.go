package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	userId := usersLoginInfo[token].Id

	//用户的点赞,1代表点赞,2代表取消点赞
	actionType := c.Query("action_type")
	//获取视频的id
	viedoId := c.Query("video_id")
	fmt.Println(viedoId)

	//如果点赞不为1，说明需要取消点赞
	dbm.UpdateUserFavorite(userId, viedoId, actionType)

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {

	//获取用户的token检测是否合法
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请先登录后再操作"})
	}

	//调用FavouriteByUserId查出该用户喜爱的视频列表
	userid := c.Query("user_id")
	user_id, _ := strconv.ParseInt(userid, 10, 8)
	videosList := dbm.FavouriteByUserId(user_id)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videosList,
	})

}
