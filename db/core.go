package db

import (
	"github.com/glebarez/sqlite"
	"github.com/wike2019/wike_go/v2/server/model"
	"gorm.io/gorm"
)

type CoreDb struct {
	DB *gorm.DB
}

func InitDb() *CoreDb {
	dbSqlite, err := gorm.Open(sqlite.Open("./db/core.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = dbSqlite.AutoMigrate(&model.CronJob{})
	if err != nil {
		panic("failed to migrate database")
	}

	return &CoreDb{
		DB: dbSqlite,
	}
}
