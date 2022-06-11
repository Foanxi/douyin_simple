package controller

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync/atomic"
)

var db *sql.DB

type DataBaseManager interface {
	SearchUser(userid int64) User
	FindAllUser() []User
	AddUser(user User) bool
	GetLastId() int64
	GerAllUser() map[string]User
	GetLastVideoId() int64
	InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool
	GetFavouriteCount(id int64) int64
	FavouriteByUserId(id int64) []Video
	GetVideoList() []Video
	UpdateUserFavorite(userId int64, videoId string, favourite string)
}

type manager struct {
}

var dbm DataBaseManager = &manager{}

func (mgr *manager) SearchUser(userid int64) User {
	db, _ = sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	row := db.QueryRow("select * from user where id=?", &userid)
	if row == nil {
		fmt.Print("查询失败")
		return User{}
	}
	var u User
	err := row.Scan(&u.Id, &u.Name, &u.Password, &u.FollowCount, &u.FollowerCount)
	if err != nil {
		fmt.Print("添加至结构体失败")
	}
	return u
}

func (mgr *manager) FindAllUser() []User {
	db, _ = sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	u := User{}
	users := make([]User, 0)
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

func (mgr *manager) AddUser(user User) bool {
	db, err := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	_, err = db.Exec("INSERT INTO user(Id,Name,Password,FollowCount,FollowerCount,IsFollow)VALUES (?,?,?,?,?,?)", &user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount, &user.IsFollow)
	if err != nil {
		return false
	}
	return true
}

func (mgr *manager) GetLastId() int64 {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	var id int64
	err := db.QueryRow("select id from user order by id desc limit 1").Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func (mgr *manager) GerAllUser() map[string]User {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	m := make(map[string]User)
	rows, _ := db.Query("select id, name, password, followcount, followercount from user")
	var u User
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Password, &u.FollowCount, &u.FollowerCount)
		if err != nil {
			return nil
		}
		token := u.Name + u.Password
		m[token] = u
	}
	defer db.Close()
	return m
}

func (mgr *manager) GetLastVideoId() int64 {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	var id int64
	db.QueryRow("select id from video order by id desc limit 1").Scan(&id)
	if id == 0 {
		return 1
	}
	return id
}

// 插入新的视频信息
func (mgr *manager) InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	defer db.Close()
	var id = dbm.GetLastVideoId()
	atomic.AddInt64(&id, 1)
	_, err := db.Exec("INSERT INTO video(id,author_id,play_url,cover_url,favourite_count,comment_count,create_time) value (?,?,?,?,?,?,?)", &id, &authorId, &playUrl, &coverUrl, 0, 0, time)
	if err != nil {
		return false
	}
	return true
}

func (mgr *manager) FavouriteByUserId(id int64) []Video {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	defer db.Close()

	//找到该用户所喜欢的视频
	rows, _ := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video v1,user_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id)
	m := make([]Video, 30)

	//获取视频的信息
	numcount := 0
	for rows.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u VideoDB
		err := rows.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video Video

		//调用SearchUser方法获取到User对象
		author := dbm.SearchUser(u.AuthorId)

		/**
		TODO
		在这里需要通过token去获取到该用户是否有对该视频点赞
		*/

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = "http://10.34.151.198:8080/static" + u.PlayUrl
		video.CoverUrl = "http://10.34.151.198:8080/static" + u.CoverUrl
		video.IsFavorite = false
		if err != nil {
			fmt.Println("err = ", err)
			break
		}

		//为第numcount个视频赋值
		m[numcount] = video
		numcount++
	}
	return m
}

//获取某id视频的点赞数
func (mgr *manager) GetFavouriteCount(id int64) int64 {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	defer db.Close()
	//获取数据库该id的点赞数
	favouriteCount := db.QueryRow("SELECT COUNT(*) From user_video where video_id= ? AND user_video.favourite = 1", &id)
	var count int64
	favouriteCount.Scan(&count)
	return count
}

func (mgr *manager) GetVideoList() []Video {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	//defer db.Close()
	video_list, err := db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video order by create_time desc")
	if err != nil {
		return nil
	}
	m := make([]Video, 30)
	numcount := 0
	//获取视频的信息
	for video_list.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u VideoDB
		err := video_list.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video Video

		//调用SearchUser方法获取到User对象
		author := dbm.SearchUser(u.AuthorId)

		/**
		TODO
		在这里需要通过token去获取到该用户是否有对该视频点赞
		*/

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = "http://10.34.151.198:8080/static" + u.PlayUrl
		video.CoverUrl = "http://10.34.151.198:8080/static" + u.CoverUrl
		video.IsFavorite = false
		if err != nil {
			fmt.Println("err = ", err)
			break
		}

		//为第numcount个视频赋值
		m[numcount] = video
		numcount++
	}
	return m
}

// 更新用户的点赞操作
func (mgr *manager) UpdateUserFavorite(userId int64, videoId string, favourite string) {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	defer db.Close()

	video_Id, _ := strconv.ParseInt(videoId, 10, 8)
	favourInt, _ := strconv.ParseInt(favourite, 10, 8)

	//首先查询该用户是否与该视频有关系，如果没有则添加新用户，并将favourite赋值为1
	row := db.QueryRow("select * from user_video where user_id = ? and video_id = ?", userId, video_Id)

	var result Favourite
	err := row.Scan(&result.user_id, &result.video_id, &result.favourite, &result.comment)

	//已经建立关系
	if err == nil {
		if favourInt == 1 {
			_, err := db.Exec("update user_video set favourite = ? where user_id = ? and video_id = ?", 1, userId, videoId)
			if err != nil {
				fmt.Println("1在修改用户点赞操作时出现错误，err = ", err)
				return
			}
		} else if favourInt == 2 {
			if result.comment == "" {
				db.Exec("delete from user_video where user_id=? and video_id=?", userId, videoId)
			} else {
				_, err := db.Exec("update user_video set favourite = ? where user_id = ? and video_id = ?", 2, userId, videoId)
				if err != nil {
					fmt.Println("2在修改用户点赞操作时出现错误，err = ", err)
					return
				}
			}
		}

	} else {
		//添加用户和视频的关系
		_, err := db.Exec("insert into user_video values(?,?,1,?)", userId, videoId, "")
		if err != nil {
			fmt.Println("添加用户与视频的关系出错")
		}
	}

}
