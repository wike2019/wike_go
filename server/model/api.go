package model

import "gorm.io/gorm"

// api接口模型
type API struct {
	gorm.Model        // 主键
	Input      string `gorm:"size:2000"` // 入参数
	Output     string `gorm:"size:2000"` // 出参数
	Status     int    `gorm:"type:int"`
	APISearch
}

// 可以用于搜索的字段
type APISearch struct {
	Group  string `gorm:"size:300" search:"like"` // 路由分组
	Name   string `gorm:"size:300" search:"like"` // 路由名称
	Path   string `gorm:"size:300" search:"like"` // 路由路径
	Method string `gorm:"size:100" search:"true"` // 请求方法
}

func (this *API) TableName() string {
	return "apis"
}
func (this *APISearch) TableName() string {
	return "apis"
}

// 子路有
type ChildAPi struct {
	Value  uint   `json:"value"`
	Label  string `json:"label"`
	Path   string
	Method string
}

// 前端路由列表
type APIFront struct {
	Value    uint       `json:"value"`
	Label    string     `json:"label"`
	Children []ChildAPi `json:"children"`
}

// 前端路由权限列表
type APIDataForRole struct {
	APIList       []Rule `json:"apiList"`
	APIParentList []Rule `json:"apiParentList"`
	ApiIDs        []uint `json:"apiIDs"`
}

// 接口分组结构
type APIGroup struct {
	Group string
	APIs  []API
}
