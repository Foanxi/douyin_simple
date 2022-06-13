package global

import (
	"github.com/RaymondCode/simple-demo/config"
	"gorm.io/gorm"
)

var (
	Conf config.Config
	Db   *gorm.DB
	Jwt  config.Jwt
)
