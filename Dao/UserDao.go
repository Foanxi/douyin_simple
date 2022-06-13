package Dao

import (
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/initialize"
	"github.com/RaymondCode/simple-demo/jwt"
	"github.com/RaymondCode/simple-demo/model"
	_type "github.com/RaymondCode/simple-demo/type"
)

type userDaoImp struct {
}

func InitDB() {
	initialize.LoadConfig()
	initialize.Mysql()
}

var Udi UserDaoImp = &userDaoImp{}

type UserDaoImp interface {
	GetUserById(userid int64) _type.User
	AddUser(user _type.User) bool
	GetLastId() int64
	GetAllUser() map[string]_type.User
	GetUserRelation(authorId int64, favouriteId int64) bool
	SearchUser(userid int64) _type.User
	GetAuthorById(userId string) []_type.User
	GetFanList(userId string) []_type.User
	GetUserByToken(token string) (_type.User, error)
}

// GetAllUser 获取全部用户
func (mgr *userDaoImp) GetAllUser() map[string]_type.User {
	if global.Db == nil {
		InitDB()
	}
	//创建一个map存放User对象，并通过token作为键获取对应的User对象
	m := make(map[string]_type.User)

	var users []model.User

	global.Db.Raw("select * from user").Scan(&users)

	for i := 0; i < len(users); i++ {
		result := users[i]
		user := _type.User{
			Id:            result.Id,
			Name:          result.Name,
			Password:      result.Password,
			FollowCount:   result.FollowCount,
			FollowerCount: result.FollowerCount,
		}
		token := user.Name + user.Password
		m[token] = user
	}
	return m
}

// AddUser 添加新用户
func (mgr *userDaoImp) AddUser(user _type.User) bool {
	global.Db.Exec("INSERT INTO user VALUES (?,?,?,?,?)", &user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)
	return true
}

// GetLastId 获取用户的最后一个id
func (mgr *userDaoImp) GetLastId() int64 {
	var user model.User
	global.Db.Last(&user)
	return user.Id
}

// GetUserById 返回指定id的用户
func (mgr *userDaoImp) GetUserById(id int64) _type.User {

	var userModel model.User
	global.Db.Raw("select * from user where id = ?", id).Scan(&userModel)
	user := _type.User{
		Id:            userModel.Id,
		Password:      userModel.Password,
		Name:          userModel.Name,
		FollowCount:   userModel.FollowCount,
		FollowerCount: userModel.FollowerCount,
	}
	return user
}

// GetUserRelation authorId 表示关注作者的id，favouriteId 表示当前用户的id
func (mgr *userDaoImp) GetUserRelation(authorId int64, favouriteId int64) bool {
	var favouriteUser model.FavouriteUser
	global.Db.Raw("select * from favourite_user where author_id = ? and favourite_id = ?", authorId, favouriteId).Scan(&favouriteUser)
	return favouriteUser.AuthorId == authorId
}

// SearchUser 查询一个用户
func (mgr *userDaoImp) SearchUser(userid int64) _type.User {
	var user model.User
	global.Db.Raw("select * from user where id=?", userid).Scan(&user)
	u := _type.User{
		Id:            user.Id,
		Password:      user.Password,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
	}
	return u
}

// GetAuthorById 获取全部关注人
func (mgr *userDaoImp) GetAuthorById(userId string) []_type.User {
	var users []model.User
	global.Db.Raw("SELECT * FROM favourite_user,user where favourite_id = ? and author_id = user.Id", userId).Scan(&users)
	authorList := make([]_type.User, len(users))
	for i := 0; i < len(users); i++ {
		result := users[i]
		user := _type.User{
			Id:            result.Id,
			Password:      result.Password,
			Name:          result.Name,
			FollowCount:   result.FollowCount,
			FollowerCount: result.FollowerCount,
			IsFollow:      true,
		}
		authorList[i] = user
	}
	return authorList
}

func (mgr *userDaoImp) GetFanList(userId string) []_type.User {
	var users []model.User
	global.Db.Raw("SELECT * FROM favourite_user,user where author_id = ? and favourite_id = user.Id", userId).Scan(&users)
	authorList := make([]_type.User, len(users))
	for i := 0; i < len(users); i++ {
		result := users[i]
		user := _type.User{
			Id:            result.Id,
			Name:          result.Name,
			FollowCount:   result.FollowCount,
			FollowerCount: result.FollowerCount,
			IsFollow:      true,
		}
		authorList[i] = user
	}
	return authorList
}

func (mgr *userDaoImp) GetUserByToken(token string) (_type.User, error) {
	var user _type.User

	uid, err := jwt.GetUidByToken(token)
	if err != nil {
		return user, err

	}

	user = Udi.SearchUser(uid)

	return user, nil
}
