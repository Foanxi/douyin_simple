package Dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	_type "github.com/RaymondCode/simple-demo/type"
	"strconv"
	"sync/atomic"
)

type videoDaoImp struct {
}

var Vdi VideoDaoImp = &videoDaoImp{}

type VideoDaoImp interface {
	GetLastVideoId() int64
	InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool
	FavouriteByUserId(id int64) []_type.Video
	GetVideoList(token string) []_type.Video
	GetUserPublish(userId string) []_type.Video
}

//找到最后的视频编号
func (vdi *videoDaoImp) GetLastVideoId() int64 {
	var id int64
	err := global.Db.QueryRow("select id from video order by id desc limit 1").Scan(&id)
	if err != nil {
		fmt.Println("获取视频最后id的时候失败")
	}
	if id == 0 {
		return 1
	}
	return id
}

// InsertVideo 插入新的视频信息
func (vdi *videoDaoImp) InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool {
	id := vdi.GetLastVideoId()
	atomic.AddInt64(&id, 1)
	_, err := global.Db.Exec("INSERT INTO video(id,author_id,play_url,cover_url,favourite_count,comment_count,create_time) value (?,?,?,?,?,?,?)", &id, &authorId, &playUrl, &coverUrl, 0, 0, time)
	if err != nil {
		return false
	}
	return true
}

// FavouriteByUserId 查找用户所喜爱的视频,id是当前用户的id
func (vdi *videoDaoImp) FavouriteByUserId(id int64) []_type.Video {
	var count int64
	global.Db.QueryRow("select count(*) from video v1,favouriter_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id).Scan(&count)
	if count > 30 {
		count = 30
	}

	//查询对应id的用户所喜爱的视频
	rows, _ := global.Db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video v1,favouriter_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id)
	m := make([]_type.Video, count)
	defer rows.Close()

	//获取视频的信息
	numCount := 0
	for rows.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u _type.VideoDB
		_ = rows.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video _type.Video

		//调用SearchUser方法获取到User对象
		author := Udi.GetUserById(u.AuthorId)

		author.IsFollow = Udi.GetUserRelation(u.AuthorId, id)

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = Fdi.GetVideoFavouriteRelation(id, u.Id)

		//为第numCount个视频赋值
		m[numCount] = video
		numCount++
	}
	return m
}
func (vdi *videoDaoImp) GetVideoList(token string) []_type.Video {
	var count int64
	global.Db.QueryRow("select count(*) from video").Scan(&count)
	if count > 30 {
		count = 30
	}
	videoList, err := global.Db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video order by create_time desc")
	defer videoList.Close()

	if err != nil {
		return nil
	}
	m := make([]_type.Video, count)
	numCount := 0

	user, _ := Udi.GetUserByToken(token)
	//获取视频的信息
	for videoList.Next() {
		//数据库使用的临时存放结果集数据的对象
		var u _type.VideoDB
		err := videoList.Scan(&u.Id, &u.AuthorId, &u.PlayUrl, &u.CoverUrl, &u.FavoriteCount, &u.CommentCount)

		//真正的Video对象
		var video _type.Video

		//调用SearchUser方法获取到User对象
		author := Udi.GetUserById(u.AuthorId)

		//调用GetUserRelation获取到当前用户是否有关注视频作者
		author.IsFollow = Udi.GetUserRelation(author.Id, user.Id)

		//为video对象赋值
		video.Id = u.Id
		video.Author = author
		video.CommentCount = u.CommentCount
		video.FavoriteCount = u.FavoriteCount
		video.PlayUrl = global.Conf.Ipconfig.Ip_url + "static" + u.PlayUrl
		video.CoverUrl = global.Conf.Ipconfig.Ip_url + "static" + u.CoverUrl
		video.IsFavorite = Fdi.GetVideoFavouriteRelation(user.Id, u.Id)
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
func (vdi *videoDaoImp) GetUserPublish(userId string) []_type.Video {
	var count int64

	id, _ := strconv.ParseInt(userId, 10, 8)
	global.Db.QueryRow("select count(*) from video where author_id=?", id).Scan(&count)
	if count > 30 {
		count = 30
	}
	videoList, err := global.Db.Query("select id, author_id, play_url, cover_url, favourite_count, comment_count from video where author_id = ?", id)
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
		author := Udi.GetUserById(u.AuthorId)

		row, _ := global.Db.Query("select * from favouriter_video where user_id = ? and video_id =? and favourite = 1", id, u.Id)

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
