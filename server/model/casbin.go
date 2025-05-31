package model

import "gorm.io/gorm"

// 路由规则
type Rule struct {
	gorm.Model
	RuleSearch
	Path   string `json:"path"`
	Method string `json:"method"`
	APIId  uint   `json:"apiId"`
}

func (this *Rule) TableName() string {
	return "sys_rule"
}

// 路由搜索字段
type RuleSearch struct {
	Role string `json:"role" search:"true"`
}

func (this *RuleSearch) TableName() string {
	return "sys_rule"
}

// 根据接口获取路由角色列表
type AggregatedRule struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Roles  []string `json:"roles"`
}

// 角色列表
type Role struct {
	gorm.Model
	Parent string `json:"parent"`
	Child  string `json:"children"`
}
type RoleService struct {
	List  []Role
	Names []string
}

func (this *Role) TableName() string {
	return "sys_role"
}

// 添加规则
type RuleInput struct {
	Role  string `json:"role"`
	APIId []uint `json:"apiId"`
}

// 添加角色
type RoleInput struct {
	Role string `form:"role" search:"true"`
}
