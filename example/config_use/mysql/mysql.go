package mysql

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDB struct {
	*gorm.DB
}

func InitMysql(cfg *viper.Viper) *MysqlDB {
	dsn := cfg.GetString("mysql")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql: " + err.Error())
	}
	return &MysqlDB{DB: db}
}
