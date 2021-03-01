package main

import (
	_ "github.com/go-Sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type DBConfig struct {
}
func NewDBConfig() *DBConfig {
	return &DBConfig{}
}

func(this *DBConfig) GormDB() *gorm.DB{
	dns:="root:root@tcp(192.168.3.2:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	return db
}

