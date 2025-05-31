package model

import "gorm.io/gorm"

type SysDictionary struct {
	gorm.Model
	SysDictionarySearch
	Desc                 string                `json:"desc" form:"desc" gorm:"column:desc;comment:描述"` // 描述
	SysDictionaryDetails []SysDictionaryDetail `json:"sysDictionaryDetails" form:"sysDictionaryDetails"`
}
type SysDictionarySearch struct {
	Name   string `json:"name" form:"name" gorm:"column:name;comment:字典名（中）" search:"like"`    // 字典名（中）
	Type   string `json:"type" form:"type" gorm:"column:type;comment:分类;unique" search:"like"` // 分类，添加唯一索引
	Status int    `json:"status" form:"status" gorm:"column:status;comment:状态" search:"true"`  // 状态
}

func (this *SysDictionarySearch) TableName() string {
	return "sys_dictionaries"
}

func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}

type SysDictionaryDetail struct {
	gorm.Model
	Label           string `json:"label" form:"label" gorm:"column:label;comment:展示值"`                                  // 展示值
	Value           string `json:"value" form:"value" gorm:"column:value;comment:字典值"`                                  // 字典值
	Extend          string `json:"extend" form:"extend" gorm:"column:extend;comment:扩展值"`                               // 扩展值
	Status          int    `json:"status" form:"status" gorm:"column:status;comment:启用状态"`                              // 启用状态
	Sort            int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序标记"`                                    // 排序标记
	SysDictionaryID int    `json:"sysDictionaryID" form:"sysDictionaryID" gorm:"column:sys_dictionary_id;comment:关联标记"` // 关联标记
	Desc            string `json:"desc" form:"desc" gorm:"column:desc;comment:描述"`                                      // 描述
}

func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_details"
}
