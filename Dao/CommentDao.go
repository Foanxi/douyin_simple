package Dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	_type "github.com/RaymondCode/simple-demo/type"
	"time"
)

type commentDaoImp struct {
}

var Cdi CommentDaoImp = &commentDaoImp{}

type CommentDaoImp interface {
	GetAllComment(videoId int64) []_type.Comment
	GetLastCommentId() int64
	AddComment(useId int64, videoId int64, commentText string) _type.Comment
	DeleteCommentById(commentIdInt int64, videoIdInt int64) bool
}

// GetAllComment 查询该视频的所有的评论，并返回
func (Cdi *commentDaoImp) GetAllComment(videoId int64) []_type.Comment {
	commentList := make([]_type.Comment, 20)
	//获取评论的信息
	var resultComment []model.Comment
	global.Db.Raw("select * from comment where video_id = ?", videoId).Scan(&resultComment)
	if resultComment == nil {
		fmt.Println("该视频并没有评论")
		return nil
	}

	var resultUser []model.User
	//获取评论的用户信息
	global.Db.Raw("select user.id,user.name,user.password,user.follow_count,user.follower_count from comment,user where video_id = ? and comment.user_id = user.id", videoId).Scan(&resultUser)

	if resultUser == nil {
		fmt.Println("在拿取视频评论用户信息时出错")
	}
	//将用户信息存储在一个map中
	userList := make(map[int64]_type.User)
	for i := 0; i < len(resultUser); i++ {
		resultuser := resultUser[i]
		user := _type.User{
			Id:            resultuser.Id,
			Name:          resultuser.Name,
			Password:      resultuser.Password,
			FollowCount:   resultuser.FollowCount,
			FollowerCount: resultuser.FollowCount,
		}
		userList[user.Id] = user
	}

	var commentNum int
	for commentNum = 0; commentNum < len(resultComment); commentNum++ {
		result := resultComment[commentNum]
		var commentUserId int64
		comment := _type.Comment{
			Id: result.CommentId,
			User: _type.User{
				Id: result.UserId,
			},
			Content:    result.CommentText,
			CreateDate: result.CreateTime,
		}
		for i, user := range userList {
			//说明该评论的用户的id为i，则将user赋值给comment
			if i == commentUserId {
				comment.User = user
			}
		}
		commentList[commentNum] = comment
		commentNum++
	}
	commentList = commentList[:commentNum]
	return commentList
}

// GetLastCommentId 获取评论编号的最后一个id
func (Cdi *commentDaoImp) GetLastCommentId() int64 {
	var result model.Comment
	global.Db.Last(&result)
	if result.CommentId == 0 {
		return 1
	}
	return result.CommentId
}

// AddComment 将用户评论添加进数据库
func (Cdi *commentDaoImp) AddComment(useId int64, videoId int64, commentText string) _type.Comment {
	//获取当前日期
	timeUnix := time.Now().Unix()
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	rs := []rune(formatTimeStr)
	//格式为mm-dd
	formatTimeStr = string(rs[5:10])
	//获取最后一个id
	lastId := Cdi.GetLastCommentId() + 1
	//获取指定id的用户
	user := Udi.GetUserById(useId)
	//插入评论
	global.Db.Exec("insert into comment value(?,?,?,?,?)", lastId, useId, videoId, commentText, formatTimeStr)
	//修改评论数
	global.Db.Exec("update video set comment_count = comment_count + 1 where id = ?", videoId)
	return _type.Comment{Id: lastId, User: user, Content: commentText, CreateDate: formatTimeStr}
}

func (Cdi *commentDaoImp) DeleteCommentById(commentIdInt int64, videoIdInt int64) bool {
	//已删除评论
	global.Db.Exec("delete from comment where comment_id = ?", commentIdInt)

	//将视频的评论数减一
	global.Db.Exec("update video set comment_count = comment_count-1 where id = ?", videoIdInt)

	return true
}
