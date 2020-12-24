package http

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	Injector "github.com/shenyisyn/goft-ioc"
	"github.com/wike2019/wike_go/src/core/config"
	"log"
)

type DBConfig struct {
}
func NewDBConfig() *DBConfig {
	return &DBConfig{}
}

func(this *DBConfig) GormDB() *gorm.DB{
	dns:="root:roo1111111t@tcp(192.168.3.2:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	if configData := Injector.BeanFactory.Get((*config.SysConfig)(nil)); configData != nil {
		dns = configData.(*config.SysConfig).Server.Mysql
	}
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	return db
}