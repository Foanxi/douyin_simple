package global

import (
	"database/sql"
	"github.com/RaymondCode/simple-demo/config"
)

var (
	Conf config.Config
	Db   *sql.DB
)
