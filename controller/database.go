package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/initalize"
	"github.com/RaymondCode/simple-demo/type"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync/atomic"
	"time"
)

var db = global.Db

type Manager interface {
	SearchUser(userid int64) _type.User
	FindAllUser() []_type.User
	AddUser(user _type.User) bool
	GetLastId() int64
	GerAllUser() map[string]_type.User
	GetLastVideoId() int64
	InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool
	FavouriteByUserId(id int64) []_type.Video
	GetVideoList(token string) []_type.Video
	UpdateUserFavorite(userId int64, videoId string, favourite string)
	UserFavoriteUser(userId string, favouriteUserId string, actionType string) bool
	GetUserPublish(userId string) []_type.Video
	GetAllComment(videoId int64) []_type.Comment
	GetLastCommentId() int64
	AddComment(useId int64, videoId int64, commentText string) _type.Comment
}

type manager struct {
}

var Dbm Manager = &manager{}

func InitDB() {
	initalize.LoadConfig()
	initalize.Mysql()
	db = global.Db
}

func (mgr *manager) SearchUser(userid int64) _type.User {
	row := db.QueryRow("select * from user where id=?", &userid)
	if row == nil {
		fmt.Print("查询失败")
		return _type.User{}
	}
	var u _type.User
	err := row.Scan(&u.Id, &u.Name, &u.Password, &u.FollowCount, &u.FollowerCount)
	if err != nil {
		fmt.Print("添加至结构体失败")
	}
	return u
}

func (mgr *manager) FindAllUser() []_type.User {
	u := _type.User{}
	users := make([]_type.User, 0)
	rows, _ := db.Query("select * from user")
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Password)
		if err != nil {
			return nil
		}
		users = append(users, u)
	}
	return users
}

func (mgr *manager) AddUser(user _type.User) bool {
	_, err := db.Exec("INSERT INTO user(Id,Name,Password,FollowCount,FollowerCount)VALUES (?,?,?,?,?)", &user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)
	if err != nil {
		return false
	}
	return true
}

func (mgr *manager) GetLastId() int64 {
	var id int64
	err := db.QueryRow("select id from user order by id desc limit 1").Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func (mgr *manager) GerAllUser() map[string]_type.User {
	if db == nil {
		InitDB()
	}
	m := make(map[string]_type.User)
	rows, _ := db.Query("select id, name, password, followcount, followercount from user")
	var u _type.User
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Password, &u.FollowCount, &u.FollowerCount)
		if err != nil {
			return nil
		}
		token := u.Name + u.Password
		m[token] = u
	}

	return m
}

func (mgr *manager) GetLastVideoId() int64 {
	var id int64
	err := db.QueryRow("select id from video order by id desc limit 1").Scan(&id)
	if err != nil {
		fmt.Println("获取视频最后id的时候失败")
	}
	if id == 0 {
		return 1
	}
	return id
}

// InsertVideo 插入新的视频信息
func (mgr *manager) InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool {
	var id = Dbm.GetLastVideoId()
	atomic.AddInt64(&id, 1)
	_, err := db.Exec("INSERT INTO video(id,author_id,play_url,cover_url,favourite_count,comment_count,create_time) value (?,?,?,?,?,?,?)", &id, &authorId, &playUrl, &coverUrl, 0, 0, time)
	if err != nil {
		return false
	}
	return true
}

func (mgr *manager) FavouriteByUserId(id int64) []_type.Video {
	var count int64
	db.QueryRow("select count(*) from video v1,favouriter_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id).Scan(&count)
	if count > 30 {
		count = 30
	}
	//找到该用户所喜欢的视频
	rows, _ := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video v1,favouriter_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id)
	m := make([]_type.Video, count)

	//获取视频的信息
	numCount := 0
	for rows.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u _type.VideoDB
		err := rows.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video _type.Video

		//调用SearchUser方法获取到User对象
		author := Dbm.SearchUser(u.AuthorId)

		row, _ := db.Query("select * from favouriter_video where user_id = ? and video_id =? and favourite = 1", id, u.Id)

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = row.Next()
		if err != nil {
			fmt.Println("err = ", err)
			break
		}

		//为第numCount个视频赋值
		m[numCount] = video
		numCount++
	}
	return m
}

func (mgr *manager) GetVideoList(token string) []_type.Video {
	var count int64
	db.QueryRow("select count(*) from video").Scan(&count)
	if count > 30 {
		count = 30
	}
	videoList, err := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video order by create_time desc")
	if err != nil {
		return nil
	}
	m := make([]_type.Video, count)
	numCount := 0
	//获取视频的信息
	for videoList.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u _type.VideoDB
		err := videoList.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video _type.Video

		//调用SearchUser方法获取到User对象
		author := Dbm.SearchUser(u.AuthorId)

		user := UsersLoginInfo[token]
		row, _ := db.Query("select * from favouriter_video where user_id = ? and video_id =? and favourite = 1", user.Id, u.Id)

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = row.Next()
		if err != nil {
			fmt.Println("err = ", err)
			break
		}

		//为第numCount个视频赋值
		m[numCount] = video
		numCount++
	}
	return m
}

// UpdateUserFavorite 更新用户的点赞操作
func (mgr *manager) UpdateUserFavorite(userId int64, videoId string, favourite string) {
	video_Id, _ := strconv.ParseInt(videoId, 10, 8)
	favourInt, _ := strconv.ParseInt(favourite, 10, 8)

	//首先查询该用户是否与该视频有关系，如果没有则添加新用户，并将favourite赋值为1
	row := db.QueryRow("select * from favouriter_video where user_id = ? and video_id = ?", userId, video_Id)

	var result _type.Favourite
	err := row.Scan(&result.UserId, &result.VideoId, &result.Favourite)

	//已经建立关系
	if err == nil {
		if favourInt == 1 {
			_, err := db.Exec("update favouriter_video set favourite = ? where user_id = ? and video_id = ?", 1, userId, videoId)
			if err != nil {
				fmt.Println("1在修改用户点赞操作时出现错误，err = ", err)
			}
			_, err = db.Exec("update video set favourite_count = favourite_count+1 where id=?", video_Id)
		} else if favourInt == 2 {
			_, err2 := db.Exec("delete from favouriter_video where user_id=? and video_id=?", userId, videoId)
			if err2 != nil {
				fmt.Println("在删除关系的时候出错，err = ", err2)
			}
			_, err = db.Exec("update video set favourite_count = favourite_count-1 where id=?", video_Id)
		}
	} else {
		//添加用户和视频的关系
		_, err := db.Exec("insert into favouriter_video values(?,?,1)", userId, videoId)
		if err != nil {
			fmt.Println("添加用户与视频的关系出错")
		}
		_, err = db.Exec("update video set favourite_count = favourite_count+1 where id=?", video_Id)
	}
}

func (mgr *manager) UserFavoriteUser(userId string, favouriteUserId string, actionType string) bool {
	favouriteType, _ := strconv.ParseInt(actionType, 10, 8)
	id := UsersLoginInfo[userId].Id
	row := db.QueryRow("select * from user where Id = ?", id)
	if row != nil {
		if strconv.FormatInt(id, 10) != favouriteUserId {
			//关注
			if favouriteType == 1 {
				_, err := db.Exec("insert into author_fans(authod_id,favourite_id) values(?,?)", id, favouriteUserId)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return true
			} else {
				// 取消关注
				_, err := db.Exec("delete from author_fans where authod_id = ? and  favourite_id = ? ", id, favouriteUserId)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return true
			}
		}
	}
	return false
}

func (mgr *manager) GetUserPublish(userId string) []_type.Video {
	var count int64

	id, _ := strconv.ParseInt(userId, 10, 8)
	db.QueryRow("select count(*) from video where author_id=?", id).Scan(&count)
	if count > 30 {
		count = 30
	}
	videoList, err := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video where author_id = ?", id)
	if err != nil {
		fmt.Println("查询失败")
	}
	m := make([]_type.Video, count)
	numCount := 0
	//获取视频的信息
	for videoList.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u _type.VideoDB
		err := videoList.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video _type.Video

		//调用SearchUser方法获取到User对象
		author := Dbm.SearchUser(u.AuthorId)

		row, _ := db.Query("select * from favouriter_video where user_id = ? and video_id =? and favourite = 1", id, u.Id)

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = row.Next()
		if err != nil {
			fmt.Println("err = ", err)
			break
		}

		//为第numCount个视频赋值
		m[numCount] = video
		numCount++
	}
	return m
}

//将用户评论添加进数据库
func (mgr *manager) AddComment(useId int64, videoId int64, commentText string) _type.Comment {
	//获取当前日期
	timeUnix := time.Now().Unix()
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	fmt.Println(formatTimeStr)
	rs := []rune(formatTimeStr)
	//格式为mm-dd
	formatTimeStr = string(rs[5:10])
	//获取最后一个id
	lastId := mgr.GetLastCommentId() + 1
	//获取指定id的用户
	user := mgr.GetUserById(useId)
	//插入评论
	_, err := db.Exec("insert into comment value(?,?,?,?,?)", lastId, useId, videoId, commentText, formatTimeStr)
	if err != nil {
		fmt.Println("在插入评论时出错")
	}
	//修改评论数
	_, err = db.Exec("update video set comment_count = comment_count + 1 where id = ?", videoId)
	if err != nil {
		fmt.Println("在修改评论数时出错 ，err = ", err)
	}
	return _type.Comment{lastId, user, commentText, formatTimeStr}
}

//获取评论编号的最后一个id
func (mgr *manager) GetLastCommentId() int64 {
	var result int64
	db.QueryRow("SELECT comment_id FROM COMMENT ORDER BY comment_id DESC LIMIT 1;").Scan(&result)
	if result == 0 {
		return 1
	}
	return result
}

//返回指定id的用户
func (mgr *manager) GetUserById(id int64) _type.User {
	var user _type.User
	err := db.QueryRow("select * from user where id = ?", id).Scan(&user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)
	if err != nil {
		fmt.Println("返回指定id的用户时出错")
	}
	return user
}

//查询该视频的所有的评论，并返回
func (mgr *manager) GetAllComment(videoId int64) []_type.Comment {
	commentList := make([]_type.Comment, 20)
	//获取评论的信息
	resultComment, err := db.Query("select * from comment where video_id = ?", videoId)
	if err != nil {
		fmt.Println("查询该视频的全部评论时出错")
	}
	if resultComment == nil {
		fmt.Println("该视频并没有评论")
		return nil
	}
	//获取评论的用户信息
	resultUser, err := db.Query("select user.Id,user.Name,user.Password,user.Followcount,user.FollowerCount from comment,user where video_id = ? and comment.user_id = user.Id", videoId)
	if err != nil {
		fmt.Println("在拿取视频评论用户信息时出错")
	}
	//将用户信息存储在一个map中
	userList := make(map[int64]_type.User, 20)
	for resultUser.Next() {
		var user _type.User
		resultUser.Scan(&user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)
		userList[user.Id] = user
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
