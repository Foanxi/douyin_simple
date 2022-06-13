package Dao

import (
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
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

// GetLastVideoId 找到最后的视频编号
func (vdi *videoDaoImp) GetLastVideoId() int64 {
	var video model.Video
	global.Db.Table("video").Last(&video)
	return video.Id
}

// InsertVideo 插入新的视频信息
func (vdi *videoDaoImp) InsertVideo(authorId int64, playUrl string, coverUrl string, time string) bool {
	id := vdi.GetLastVideoId()
	atomic.AddInt64(&id, 1)
	global.Db.Exec("INSERT INTO video value (?,?,?,?,?,?,?)", &id, &authorId, &playUrl, &coverUrl, 0, 0, time)
	return true
}

// FavouriteByUserId 查找用户所喜爱的视频,id是当前用户的id
func (vdi *videoDaoImp) FavouriteByUserId(id int64) []_type.Video {
	var videos []model.Video
	global.Db.Raw("select id, author_id, play_url, cover_url, favourite_count, comment_count from video v1,favourite_video v2 where v2.user_id=? and v1.id = v2.video_id and v2.favourite=1", id).Scan(&videos)

	//查询对应id的用户所喜爱的视频
	m := make([]_type.Video, len(videos))

	//获取视频的信息
	for i := 0; i <= 30 && i < len(videos); i++ {
		result := videos[i]
		//数据库使用的临时存放结果集数据的对象

		//真正的Video对象
		video := _type.Video{
			Id:            result.Id,
			Author:        Udi.GetUserById(result.AuthorId),
			PlayUrl:       global.Conf.Ipconfig.Ip_url + "static" + result.PlayUrl,
			CoverUrl:      global.Conf.Ipconfig.Ip_url + "static" + result.CoverUrl,
			FavoriteCount: result.FavouriteCount,
			CommentCount:  result.CommentCount,
			IsFavorite:    Fdi.GetVideoFavouriteRelation(id, result.Id),
		}

		video.Author.IsFollow = Udi.GetUserRelation(result.AuthorId, id)

		//为第numCount个视频赋值
		m[i] = video
	}
	return m
}
func (vdi *videoDaoImp) GetVideoList(token string) []_type.Video {
	var videos []model.Video
	global.Db.Raw("select * from video order by create_time desc").Scan(&videos)

	m := make([]_type.Video, len(videos))

	user, _ := Udi.GetUserByToken(token)
	//获取视频的信息
	for i := 0; i <= 30 && i < len(videos); i++ {
		result := videos[i]
		//数据库使用的临时存放结果集数据的对象

		//真正的Video对象
		video := _type.Video{
			Id:            result.Id,
			Author:        Udi.GetUserById(result.AuthorId),
			PlayUrl:       global.Conf.Ipconfig.Ip_url + "static" + result.PlayUrl,
			CoverUrl:      global.Conf.Ipconfig.Ip_url + "static" + result.CoverUrl,
			FavoriteCount: result.FavouriteCount,
			CommentCount:  result.CommentCount,
			IsFavorite:    Fdi.GetVideoFavouriteRelation(user.Id, result.Id),
		}

		video.Author.IsFollow = Udi.GetUserRelation(result.AuthorId, user.Id)

		//为第numCount个视频赋值
		m[i] = video
	}
	return m
}
func (vdi *videoDaoImp) GetUserPublish(userId string) []_type.Video {
	var videos []model.Video
	id, _ := strconv.ParseInt(userId, 10, 8)
	global.Db.Raw("select * from video where author_id = ?", id).Scan(&videos)

	m := make([]_type.Video, len(videos))
	//获取视频的信息
	for i := 0; i < len(videos); i++ {
		result := videos[i]
		//数据库使用的临时存放结果集数据的对象

		//真正的Video对象
		video := _type.Video{
			Id:            result.Id,
			Author:        Udi.GetUserById(result.AuthorId),
			PlayUrl:       global.Conf.Ipconfig.Ip_url + "static" + result.PlayUrl,
			CoverUrl:      global.Conf.Ipconfig.Ip_url + "static" + result.CoverUrl,
			FavoriteCount: result.FavouriteCount,
			CommentCount:  result.CommentCount,
			IsFavorite:    Fdi.GetVideoFavouriteRelation(result.AuthorId, result.Id),
		}

		//为第 i 个视频赋值
		m[i] = video
	}
	return m
}
