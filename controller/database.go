package controller

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB(db *sql.DB) {
	var err error
	db, err = sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}

func searchUser(userid string) User {
	row := db.QueryRow("select * from user where id=?", userid)
	if row == nil {
		fmt.Print("查询失败")
		return User{}
	}
	var u User
	err := row.Scan(&u.Id, &u.Name, &u.Password)
	if err != nil {
		fmt.Print("添加至结构体失败")
	}
	return u
}

func FindAllUser() []User {
	db, _ = sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	u := User{}
	users := make([]User, 0)
	rows, _ := db.Query("select * from user")
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Password)
		if err != nil {
			return nil
		}
		users = append(users, u)
	}
	return users
}

func AddUser(id int64, username string, password string) bool {
	db, err := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	_, err = db.Exec("INSERT INTO user(Id,Name,Password)VALUES (?,?,?)", &id, &username, &password)
	if err != nil {
		return false
	}
	return true
}

func GetLastId() int64 {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	var id int64
	db.QueryRow("select id from user order by id desc limit 1").Scan(&id)
	return id
}

func GerAllUser() map[string]User {
	db, _ := sql.Open("mysql", "root:19635588@tcp(127.0.0.1:3306)/douyin?charset=utf8")
	m := make(map[string]User)
	rows, _ := db.Query("select id, name, password, followcount, followercount, isfollow from user")
	fmt.Print(rows)
	var u User
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Password, &u.FollowCount, &u.FollowerCount, &u.IsFollow)
		if err != nil {
			return nil
		}
		token := u.Name + u.Password
		m[token] = u
	}
	defer db.Close()
	return m
}
