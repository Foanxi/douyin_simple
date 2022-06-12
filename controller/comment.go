package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/Dao"
	_type "github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	_type.Response
	CommentList []_type.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	_type.Response
	Comment _type.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	token := c.Query("token")
	//获取用户的id，用于鉴别是否有权限看评论
	if token == "" {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	videoId := c.Query("video_id")
	videoIdInt, err := strconv.ParseInt(videoId, 10, 8)
	if err != nil {
	}
	actionType := c.Query("action_type")
	if user, exist := Dao.Udi.GerAllUser()[token]; exist {
		//等于1时说明用户要发布评论
		var comment _type.Comment
		if actionType == "1" {
			commentText := c.Query("comment_text")
			//首先先把评论添加进数据库中
			comment = Dao.Cdi.AddComment(user.Id, videoIdInt, commentText)
			c.JSON(http.StatusOK, CommentActionResponse{Response: _type.Response{StatusCode: 0},
				Comment: comment})
		} else if actionType == "2" {
			//删除评论
			//获取该评论的id
			commentId := c.Query("comment_id")
			commentIdInt, err := strconv.ParseInt(commentId, 10, 20)
			if err != nil {
				fmt.Println("在转换评论id时出错，err = ", err)
			}
			//进行删除操作
			isDelete := Dao.Cdi.DeleteCommentById(commentIdInt, videoIdInt)
			fmt.Println("是否成功删除该评论：", isDelete)
			//返回全部列表
			CommentList := Dao.Cdi.GetAllComment(videoIdInt)
			c.JSON(http.StatusOK, CommentListResponse{
				Response:    _type.Response{StatusCode: 0},
				CommentList: CommentList,
			})
		}
	} else {
		c.JSON(http.StatusOK, _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 20)
	if err != nil {
		fmt.Println("在videoId转换时出错,err = ", err)
	}
	//用户有权限则可以进入到评论列表查看该视频的所有评论,查询所有评论
	CommentList := Dao.Cdi.GetAllComment(videoId)
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    _type.Response{StatusCode: 0},
		CommentList: CommentList,
	})
}
