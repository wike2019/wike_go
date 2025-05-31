package core

import (
	"github.com/glebarez/sqlite"
	"github.com/wike2019/wike_go/server/model"
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
	err = dbSqlite.AutoMigrate(&model.API{}, &model.SysDictionary{}, &model.SysDictionaryDetail{}, &model.Job{}, &model.Rule{}, &model.Role{})
	if err != nil {
		panic("failed to migrate database")
	}
	dbSqlite.Model(model.API{}).Where("1=1").Update("status", 2)
	dbSqlite.Model(model.Job{}).Where("1=1").Delete(&model.Job{})
	return &CoreDb{
		DB: dbSqlite,
	}
}

func (this *CoreDb) ApiTable(name string, group string, input string, output string, path string, method string) *CoreDb {
	res := &model.API{
		APISearch: model.APISearch{
			Name:   name,
			Group:  group,
			Path:   path,
			Method: method,
		},
		Input:  input,
		Output: output,

		Status: 1,
	}
	historyApi := &model.API{}
	err := this.DB.Where("method=? and path=?", method, path).First(historyApi).Error
	if err != nil {
		this.DB.Create(res)
	} else {
		historyApi.Input = input
		historyApi.Output = output
		historyApi.Name = name
		historyApi.Group = group
		historyApi.Status = 1
		this.DB.Save(historyApi)
	}

	return this
}

func (this *CoreDb) GetData() map[string]model.APIGroup {
	var list []model.API
	this.DB.Find(&list)

	rs := make(map[string]model.APIGroup)

	for _, item := range list {
		// 检查分组是否存在
		group, ok := rs[item.Group]
		if !ok {
			// 如果分组不存在，初始化
			group = model.APIGroup{
				Group: item.Group,
				APIs:  make([]model.API, 0),
			}
		}

		// 修改分组中的 API 列表
		group.APIs = append(group.APIs, item)

		// 将修改后的分组放回 map
		rs[item.Group] = group
	}

	return rs

}
