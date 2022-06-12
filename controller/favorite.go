package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/Dao"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	userId := Dao.Udi.GerAllUser()[token].Id

	//用户的点赞,1代表点赞,2代表取消点赞
	actionType := c.Query("action_type")
	//获取视频的id
	videoId := c.Query("video_id")
	fmt.Println(videoId)

	//如果点赞不为1，说明需要取消点赞
	Dao.Fdi.UpdateUserFavorite(userId, videoId, actionType)

	if _, exist := Dao.Udi.GerAllUser()[token]; exist {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {

	//获取用户的token检测是否合法
	token := c.Query("token")

	//调用FavouriteByUserId查出该用户喜爱的视频列表
	userId, _ := Dao.Udi.GerAllUser()[token]
	videosList := Dao.Vdi.FavouriteByUserId(userId.Id)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		VideoList: videosList,
	})

}
