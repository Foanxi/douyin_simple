package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/initalize"
	"github.com/RaymondCode/simple-demo/type"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync/atomic"
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
	GetFavouriteCount(id int64) int64
	FavouriteByUserId(id int64) []_type.Video
	GetVideoList() []_type.Video
	UpdateUserFavorite(userId int64, videoId string, favourite string)
	UserFavoriteUser(userId string, favouriteUserId string, actionType string) bool
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
	//找到该用户所喜欢的视频
	rows, _ := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video v1,favouriter_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id)
	m := make([]_type.Video, 30)

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

		/**
		TODO
		在这里需要通过token去获取到该用户是否有对该视频点赞
		*/

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = false
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

// GetFavouriteCount 获取某id视频的点赞数
func (mgr *manager) GetFavouriteCount(id int64) int64 {
	//获取数据库该id的点赞数
	favouriteCount := db.QueryRow("SELECT COUNT(*) From favouriter_video where video_id= ? AND favouriter_video.favourite = 1", &id)
	var count int64
	err := favouriteCount.Scan(&count)
	if err != nil {
		fmt.Println("获取视频点赞数转换失败")
	}
	return count
}

func (mgr *manager) GetVideoList() []_type.Video {

	videoList, err := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video order by create_time desc")
	if err != nil {
		return nil
	}
	m := make([]_type.Video, 30)
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

		/**
		TODO
		在这里需要通过token去获取到该用户是否有对该视频点赞
		*/

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = false
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
		} else if favourInt == 2 {
			_, err2 := db.Exec("delete from favouriter_video where user_id=? and video_id=?", userId, videoId)
			if err2 != nil {
				fmt.Println("在删除关系的时候出错，err = ", err2)
			}
		}
	} else {
		//添加用户和视频的关系
		_, err := db.Exec("insert into favouriter_video values(?,?,1)", userId, videoId)
		if err != nil {
			fmt.Println("添加用户与视频的关系出错")
		}
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
