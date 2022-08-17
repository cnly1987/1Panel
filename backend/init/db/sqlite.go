package db

import (
	"github.com/1Panel-dev/1Panel/global"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteGorm() *gorm.DB {
	s := global.CONF.Sqlite
	if db, err := gorm.Open(sqlite.Open(s.Dsn()), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		return db
	}
}