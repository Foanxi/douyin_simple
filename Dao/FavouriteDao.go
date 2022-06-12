package Dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	_type "github.com/RaymondCode/simple-demo/type"
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
	var result _type.Favourite
	//首先查询该用户是否与该视频有关系，如果没有则添加新用户，并将favourite赋值为1
	_ = global.Db.QueryRow("select * from favouriter_video where user_id = ? and video_id = ?", userId, video_Id).Scan(&result.UserId, &result.VideoId, &result.Favourite)

	//已经建立关系
	if result.UserId != 0 {
		if favourInt == 1 {
			_, err := global.Db.Exec("update favouriter_video set favourite = ? where user_id = ? and video_id = ?", 1, userId, videoId)
			if err != nil {
				fmt.Println("在修改用户点赞操作为1时出现错误，err = ", err)
			}
			_, err = global.Db.Exec("update video set favourite_count = favourite_count+1 where id=?", video_Id)
		} else if favourInt == 2 {
			_, err := global.Db.Exec("delete from favouriter_video where user_id=? and video_id=?", userId, videoId)
			if err != nil {
				fmt.Println("在删除关系的时候出错，err = ", err)
			}
			_, err = global.Db.Exec("update video set favourite_count = favourite_count-1 where id=?", video_Id)
			if err != nil {
				fmt.Println("更新视频表点赞减1时出错")
			}
		}
	} else {
		//添加用户和视频的关系
		_, err := global.Db.Exec("insert into favouriter_video values(?,?,1)", userId, videoId)
		if err != nil {
			fmt.Println("添加用户与视频的关系出错")
		}
		_, err = global.Db.Exec("update video set favourite_count = favourite_count+1 where id=?", video_Id)
		if err != nil {
			fmt.Println("更新视频表点赞加1时出错")
		}
	}
}

func (Fdi *favoriteDaoImp) UserFavoriteUser(token string, favouriteUserId string, actionType string) bool {
	//对数据进行预处理
	favouriteType, _ := strconv.ParseInt(actionType, 10, 8)
	id := Udi.GerAllUser()[token].Id

	//进行查询获取用户
	rows, err := global.Db.Query("select * from user where Id = ?", id)
	defer rows.Close()

	if err != nil {
		return false
	}
	favouiteId, _ := strconv.ParseInt(favouriteUserId, 10, 8)
	if rows.Next() {
		if favouiteId != id {
			//关注
			if favouriteType == 1 {
				_, err := global.Db.Exec("insert into author_fans(author_id,favourite_id) values(?,?)", favouriteUserId, id)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return true
			} else {
				// 取消关注
				_, err := global.Db.Exec("delete from author_fans where author_id = ? and  favourite_id = ? ", favouriteUserId, id)
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
func (Fdi *favoriteDaoImp) GetVideoFavouriteRelation(userId int64, videoId int64) bool {
	rows, err := global.Db.Query("select * from favouriter_video where user_id = ? and video_id =? and favourite = 1", userId, videoId)
	defer rows.Close()
	if err != nil {
		return false
	}
	return rows.Next()
}
