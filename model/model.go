package model

import "gorm.io/gorm"

type API struct {
	gorm.Model        // 主键
	Group      string `gorm:"size:300" search:"like"` // 路由分组
	Name       string `gorm:"size:300" search:"like"` // 路由名称
	Input      string `gorm:"size:2000"`              // 入参数
	Output     string `gorm:"size:2000"`              // 出参数
	Path       string `gorm:"size:300" search:"like"` // 路由路径
	Method     string `gorm:"size:100" search:"true"` // 请求方法
	Status     int    `gorm:"type:int" search:"true"`
}

func (this *API) TableName() string {
	return "apis"
}

type Job struct {
	gorm.Model        // 主键
	Name       string `gorm:"size:300" search:"like"` // 路由名称
	Cron       string `gorm:"size:2000"`
	Func       string `gorm:"size:2000"` // 入参数
}

func (this *Job) TableName() string {
	return "jobs"
}

type SysDictionary struct {
	Name   string `json:"name" form:"name" gorm:"column:name;comment:字典名（中）"`    // 字典名（中）
	Type   string `json:"type" form:"type" gorm:"column:type;comment:分类;unique"` // 分类，添加唯一索引
	Status int    `json:"status" form:"status" gorm:"column:status;comment:状态"`  // 状态
	Desc   string `json:"desc" form:"desc" gorm:"column:desc;comment:描述"`        // 描述
	gorm.Model
	SysDictionaryDetails []SysDictionaryDetail `json:"sysDictionaryDetails" form:"sysDictionaryDetails"`
}

func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}

type SysDictionaryDetail struct {
	Label           string `json:"label" form:"label" gorm:"column:label;comment:展示值"`                                  // 展示值
	Value           string `json:"value" form:"value" gorm:"column:value;comment:字典值"`                                  // 字典值
	Extend          string `json:"extend" form:"extend" gorm:"column:extend;comment:扩展值"`                               // 扩展值
	Status          int    `json:"status" form:"status" gorm:"column:status;comment:启用状态"`                              // 启用状态
	Sort            int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序标记"`                                    // 排序标记
	SysDictionaryID int    `json:"sysDictionaryID" form:"sysDictionaryID" gorm:"column:sys_dictionary_id;comment:关联标记"` // 关联标记
	gorm.Model
}

func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_details"
}
