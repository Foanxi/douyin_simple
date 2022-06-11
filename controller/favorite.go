package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	//用户喜欢的视频
	var FavouriteVideos = []Video{}
	//获取用户的id编号
	token := c.Query("token")
	m := dbm.GerAllUser()
	userId := m[token].Id
	videosList := dbm.FavouriteByUserId(userId)
	for _, video := range videosList {
		FavouriteVideos = append(FavouriteVideos, video)
	}
	if FavouriteVideos == nil {
		fmt.Println("用户没有喜欢的视频")
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: nil,
		})
	} else {
		fmt.Println("用户有喜欢的视频")
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: FavouriteVideos,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
