package initialize

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func Mysql() {
	m := global.Conf.Mysql
	var dsn = fmt.Sprintf("%s:%s@%s", m.Username, m.Password, m.Url)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		fmt.Printf("mysql error: %s", err)
		return
	}

	sqlDb, _ := db.DB()

	// 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDb.SetMaxOpenConns(100)
	// 设置了连接可复用的最大时间
	sqlDb.SetConnMaxLifetime(time.Hour)

	global.Db = db
}
