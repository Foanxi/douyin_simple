package initalize

import (
	"database/sql"
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"time"
)

func Mysql() {
	m := global.Conf.Mysql
	var dsn = fmt.Sprintf("%s:%s@%s", m.Username, m.Password, m.Url)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql error: %s", err)
		return
	}

	// 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	db.SetMaxOpenConns(100)
	// 设置了连接可复用的最大时间
	db.SetConnMaxLifetime(time.Hour)

	global.Db = db
}
