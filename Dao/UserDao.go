package Dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/initialize"
	"github.com/RaymondCode/simple-demo/jwt"
	_type "github.com/RaymondCode/simple-demo/type"
)

var Db = global.Db

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

//	获取全部用户
func (mgr *userDaoImp) GetAllUser() map[string]_type.User {
	if global.Db == nil {
		InitDB()
	}
	//创建一个map存放User对象，并通过token作为键获取对应的User对象
	m := make(map[string]_type.User)

	rows, _ := global.Db.Query("select id, name, password, followcount, followercount from user")
	defer rows.Close()

	var u _type.User
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Password, &u.FollowCount, &u.FollowerCount)
		if err != nil {
			return nil
		}
		token := u.Name + u.Password
		m[token] = u
	}
	return m
}

// 添加新用户
func (mgr *userDaoImp) AddUser(user _type.User) bool {
	_, err := global.Db.Exec("INSERT INTO user(Id,Name,Password,FollowCount,FollowerCount)VALUES (?,?,?,?,?)", &user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)
	if err != nil {
		fmt.Println("添加新用户时出错")
		return false
	}
	return true
}

//获取用户的最后一个id
func (mgr *userDaoImp) GetLastId() int64 {
	var id int64
	err := global.Db.QueryRow("select id from user order by id desc limit 1").Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

//返回指定id的用户
func (mgr *userDaoImp) GetUserById(id int64) _type.User {
	var user _type.User
	err := global.Db.QueryRow("select * from user where id = ?", id).Scan(&user.Id, &user.Name, &user.Password, &user.FollowCount, &user.FollowerCount)

	if err != nil {
		fmt.Println("返回指定id的用户时出错")
	}
	return user
}

// GetUserRelation id 表示关注作者的id，token表示当前用户的令牌
func (mgr *userDaoImp) GetUserRelation(authorId int64, favouriteId int64) bool {
	rows, err := global.Db.Query("select * from author_fans where author_id = ? and favourite_id = ?", authorId, favouriteId)
	if err != nil {
		return false
	}
	return rows.Next()
}

//
func (mgr *userDaoImp) SearchUser(userid int64) _type.User {
	row := global.Db.QueryRow("select id, name, followcount, followercount from user where id=?", &userid)
	if row == nil {
		fmt.Print("查询失败")
		return _type.User{}
	}
	var u _type.User
	err := row.Scan(&u.Id, &u.Name, &u.FollowCount, &u.FollowerCount)
	if err != nil {
		fmt.Print("添加至结构体失败")
	}
	return u
}

// GetAuthorById 获取全部关注人
func (mgr *userDaoImp) GetAuthorById(userId string) []_type.User {
	rows, err := global.Db.Query("SELECT author_id, Name, FollowCount, FollowerCount FROM author_fans,user where favourite_id = ? and author_id = user.Id", userId)
	if err != nil {
		fmt.Println("查询关注时出错，err = ", err)
	}
	authorList := make([]_type.User, 20)
	var numCount int64
	for rows.Next() {
		var user _type.User
		rows.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount)
		user.IsFollow = true
		authorList[numCount] = user
		numCount++
	}
	authorList = authorList[:numCount]
	return authorList
}

func (mgr *userDaoImp) GetFanList(userId string) []_type.User {
	rows, err := global.Db.Query("SELECT Id, Name, FollowCount, FollowerCount FROM author_fans,user where author_id = ? and favourite_id = user.Id", userId)
	var count int64
	global.Db.QueryRow("select count(*) from author_fans,user where author_id = user.Id and author_id = ?", userId).Scan(&count)
	if err != nil {
		fmt.Println("查询粉丝出错，err = ", err)
	}
	authorList := make([]_type.User, count)
	var numCount int64
	for rows.Next() {
		var user _type.User
		rows.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount)
		user.IsFollow = true
		authorList[numCount] = user
		numCount++
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
