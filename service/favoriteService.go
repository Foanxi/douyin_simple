package service

import (
	"github.com/RaymondCode/simple-demo/Dao"
	_type "github.com/RaymondCode/simple-demo/type"
)

func FavoriteActionService(token string, videoId string, actionType string) bool {
	userId := Dao.Udi.GerAllUser()[token].Id
	Dao.Fdi.UpdateUserFavorite(userId, videoId, actionType)
	if _, exist := Dao.Udi.GerAllUser()[token]; exist {
		return true
	} else {
		return false
	}
}

func FavoriteListService(token string) []_type.Video {
	//调用FavouriteByUserId查出该用户喜爱的视频列表
	userId, _ := Dao.Udi.GerAllUser()[token]
	videosList := Dao.Vdi.FavouriteByUserId(userId.Id)
	return videosList
}
