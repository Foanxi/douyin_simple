package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	//用户的点赞,1代表点赞,2代表取消点赞
	actionType := c.Query("action_type")
	//获取视频的id
	videoId := c.Query("video_id")
	isActionService := service.FavoriteActionService(token, videoId, actionType)
	if isActionService {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//获取用户的token检测是否合法
	token := c.Query("token")
	videosList := service.FavoriteListService(token)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		VideoList: videosList,
	})
}
