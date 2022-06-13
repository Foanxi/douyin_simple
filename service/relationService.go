package service

import (
	"github.com/RaymondCode/simple-demo/Dao"
	_type "github.com/RaymondCode/simple-demo/type"
)

func RelationActionService(token string, favouriteUserId string, actionType string) bool {
	statue := Dao.Fdi.UserFavoriteUser(token, favouriteUserId, actionType)
	return statue
}

func FollowListService(useId string) []_type.User {
	var usersList []_type.User
	//进入service层
	usersList = Dao.Udi.GetAuthorById(useId)
	return usersList
}

func FollowerListService(useId string) []_type.User {
	var userFanList []_type.User
	userFanList = Dao.Udi.GetFanList(useId)
	return userFanList
}
