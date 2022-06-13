package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	//var users []model.User
	db, err := gorm.Open(mysql.Open("root:19635588@tcp(localhost:3306)/douyin?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return
	}
	//result := model.User{
	//	Id:   3,
	//	Name: "胡",
	//}

	db.Exec("insert into user(id,name,password) values(?,?,?)", 5, "zhuang", "31")
	//db.Where("id = ?", 1).Find(&users)
	//fmt.Println("user = ", result.CommentId)
	//name := "Hu"
	//db.Where("Name = ?", name).Delete(&user)
}
