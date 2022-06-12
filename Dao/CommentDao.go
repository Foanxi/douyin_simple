package Dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
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
	resultComment, err := global.Db.Query("select * from comment where video_id = ?", videoId)
	if err != nil {
		fmt.Println("查询该视频的全部评论时出错")
	}
	if resultComment == nil {
		fmt.Println("该视频并没有评论")
		return nil
	}
	//获取评论的用户信息
	resultUser, err := global.Db.Query("select user.Id,user.Name,user.Password,user.Followcount,user.FollowerCount from comment,user where video_id = ? and comment.user_id = user.Id", videoId)
	defer resultUser.Close()

	if err != nil {
		fmt.Println("在拿取视频评论用户信息时出错")
	}
	//将用户信息存储在一个map中
	userList := make(map[int64]_type.User)
	//var userNum int64
	for resultUser.Next() {
		var user _type.User
		resultUser.Scan(&user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)
		userList[user.Id] = user
		//userNum++
	}
	var commentNum int64
	for resultComment.Next() {
		var commentUserId int64
		var comment _type.Comment
		var videoId int64
		resultComment.Scan(&comment.Id, &commentUserId, &videoId, &comment.Content, &comment.CreateDate)
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
	var result int64
	global.Db.QueryRow("SELECT comment_id FROM COMMENT ORDER BY comment_id DESC LIMIT 1;").Scan(&result)
	if result == 0 {
		return 1
	}
	return result
}

// AddComment 将用户评论添加进数据库
func (Cdi *commentDaoImp) AddComment(useId int64, videoId int64, commentText string) _type.Comment {
	//获取当前日期
	timeUnix := time.Now().Unix()
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	fmt.Println(formatTimeStr)
	rs := []rune(formatTimeStr)
	//格式为mm-dd
	formatTimeStr = string(rs[5:10])
	//获取最后一个id
	lastId := Cdi.GetLastCommentId() + 1
	//获取指定id的用户
	user := Udi.GetUserById(useId)
	//插入评论
	_, err := global.Db.Exec("insert into comment value(?,?,?,?,?)", lastId, useId, videoId, commentText, formatTimeStr)
	if err != nil {
		fmt.Println("在插入评论时出错")
	}
	//修改评论数
	_, err = global.Db.Exec("update video set comment_count = comment_count + 1 where id = ?", videoId)
	if err != nil {
		fmt.Println("在修改评论数时出错 ，err = ", err)
	}
	return _type.Comment{Id: lastId, User: user, Content: commentText, CreateDate: formatTimeStr}
}

func (Cdi *commentDaoImp) DeleteCommentById(commentIdInt int64, videoIdInt int64) bool {
	//已删除评论
	_, err := global.Db.Exec("delete from comment where comment_id = ?", commentIdInt)
	if err != nil {
		fmt.Println("删除评论时出错，err = ", err)
		return false
	}
	//将视频的评论数减一
	_, err = global.Db.Exec("update video set comment_count = comment_count-1 where id = ?", videoIdInt)
	if err != nil {
		fmt.Println("将视频的评论数减一时出错，err = ", err)
		return false
	}
	return true
}
