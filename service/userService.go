package service

import (
	"github.com/RaymondCode/simple-demo/Dao"
	"github.com/RaymondCode/simple-demo/jwt"
	_type "github.com/RaymondCode/simple-demo/type"
	"strconv"
	"sync/atomic"
)

func RegisterService(username string, password string) _type.UserLoginResponse {
	var response _type.UserLoginResponse
	token := username + password
	if _, exist := Dao.Udi.GetAllUser()[token]; exist {
		response.Response = _type.Response{StatusCode: 1, StatusMsg: "User already exist"}
		return response
	} else {
		var userIdSequence = Dao.Udi.GetLastId()
		atomic.AddInt64(&userIdSequence, 1)
		newUser := _type.User{
			Id:            userIdSequence,
			Password:      password,
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		Dao.Udi.AddUser(newUser)
		response.Response = _type.Response{StatusCode: 1, StatusMsg: "User already exist"}
		response.Response = _type.Response{StatusCode: 0}
		response.UserId = Dao.Udi.GetLastId()
		response.Token, _ = jwt.GenerateToke(username, response.UserId)
		return response
	}
}
func LoginService(username string, password string) _type.UserLoginResponse {
	var response _type.UserLoginResponse
	token := username + password
	if user, exist := Dao.Udi.GetAllUser()[token]; exist {
		response.Response = _type.Response{StatusCode: 0}
		response.UserId = user.Id
		response.Token, _ = jwt.GenerateToke(username, response.UserId)
		return response
	} else {
		response.Response = _type.Response{StatusCode: 1, StatusMsg: "User doesn't exist"}
		return response
	}
}

func UserInfoService(token string, authorId string) _type.User {
	user, _ := Dao.Udi.GetUserByToken(token)
	id, _ := strconv.ParseInt(authorId, 10, 8)
	user.IsFollow = Dao.Udi.GetUserRelation(id, user.Id)
	return user
}
