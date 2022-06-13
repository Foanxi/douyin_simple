package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoListResponse struct {
	_type.Response
	VideoList []_type.Video `json:"video_list"`
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId := c.Query("user_id")
	videolist := service.PublishListService(userId)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		VideoList: videolist,
	})
}

func Action(c *gin.Context) {
	token := c.PostForm("token")
	//获取用户上传的视频名称
	title := c.PostForm("title")
	isSuccess := service.ActionService(c, token, title)
	if isSuccess {
		c.JSON(http.StatusOK, _type.Response{
			StatusCode: 0,
			StatusMsg:  title + " uploaded successfully",
		})
	}
}
