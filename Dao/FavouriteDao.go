package Dao

import (
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"strconv"
)

type favoriteDaoImp struct {
}

var Fdi FavouriteDaoImp = &favoriteDaoImp{}

type FavouriteDaoImp interface {
	UpdateUserFavorite(userId int64, videoId string, favourite string)
	UserFavoriteUser(userId string, favouriteUserId string, actionType string) bool
	GetVideoFavouriteRelation(userId int64, videoId int64) bool
}

// UpdateUserFavorite 更新用户的点赞操作
func (Fdi *favoriteDaoImp) UpdateUserFavorite(userId int64, videoId string, favourite string) {
	video_Id, _ := strconv.ParseInt(videoId, 10, 8)
	favourInt, _ := strconv.ParseInt(favourite, 10, 8)
	var result model.FavouriteVideo
	//首先查询该用户是否与该视频有关系，如果没有则添加新用户，并将favourite赋值为1
	global.Db.Raw("select * from favourite_video where user_id = ? and video_id = ?", userId, video_Id).Scan(&result)

	//已经建立关系
	if result.UserId != 0 {
		if favourInt == 1 {
			global.Db.Exec("update favourite_video set favourite = ? where user_id = ? and video_id = ?", 1, userId, videoId)
			global.Db.Exec("update video set favourite_count = favourite_count+1 where id=?", video_Id)
		} else if favourInt == 2 {
			global.Db.Exec("delete from favourite_video where user_id=? and video_id=?", userId, videoId)
			global.Db.Exec("update video set favourite_count = favourite_count-1 where id=?", video_Id)
		}
	} else {
		//添加用户和视频的关系
		global.Db.Exec("insert into favourite_video values(?,?,1)", userId, videoId)
		global.Db.Exec("update video set favourite_count = favourite_count+1 where id=?", video_Id)
	}
}

func (Fdi *favoriteDaoImp) UserFavoriteUser(token string, favouriteUserId string, actionType string) bool {
	//对数据进行预处理
	favouriteType, _ := strconv.ParseInt(actionType, 10, 8)
	user, _ := Udi.GetUserByToken(token)

	var users []model.User

	//进行查询获取用户
	global.Db.Raw("select * from user where Id = ?", user.Id).Scan(&users)
	favouiteId, _ := strconv.ParseInt(favouriteUserId, 10, 8)
	if users[0].Id != 0 {
		if favouiteId != user.Id {
			//关注
			if favouriteType == 1 {
				global.Db.Exec("insert into favourite_user(author_id,favourite_id) values(?,?)", favouriteUserId, user.Id)
				return true
			} else {
				// 取消关注
				global.Db.Exec("delete from favourite_user where author_id = ? and  favourite_id = ? ", favouriteUserId, user.Id)
				return true
			}
		}
	}
	return false
}
func (Fdi *favoriteDaoImp) GetVideoFavouriteRelation(userId int64, videoId int64) bool {
	var result model.FavouriteVideo
	global.Db.Raw("select * from favourite_video where user_id = ? and video_id =? and favourite = 1", userId, videoId).Scan(&result)
	return result.UserId == userId
}
