package service

import (
	"github.com/RaymondCode/simple-demo/Dao"
	_type "github.com/RaymondCode/simple-demo/type"
)

func FavoriteActionService(token string, videoId string, actionType string) bool {
	user, _ := Dao.Udi.GetUserByToken(token)
	Dao.Fdi.UpdateUserFavorite(user.Id, videoId, actionType)
	if user.Id != 0 {
		return true
	} else {
		return false
	}
}

func FavoriteListService(token string) []_type.Video {
	//调用FavouriteByUserId查出该用户喜爱的视频列表
	user, _ := Dao.Udi.GetUserByToken(token)
	videosList := Dao.Vdi.FavouriteByUserId(user.Id)
	return videosList
}
