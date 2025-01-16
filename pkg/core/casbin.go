package core

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"log"
	"os"
)

type Rule struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type AggregatedRule struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Roles  []string `json:"roles"`
}

const modelconfig = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// 鉴权系统
func NewEnforcer(db *CoreDb) *casbin.Enforcer {
	_, err := os.Stat("model.conf")
	if os.IsNotExist(err) {
		os.WriteFile("model.conf", []byte(modelconfig), os.ModePerm)
	}
	adapter, err := gormadapter.NewAdapterByDB(db.DB)
	if err != nil {
		log.Fatalf("Failed to create NewAdapterByDB: %v", err)
	}

	// 2. 创建 Casbin Enforcer，加载模型和策略
	enforcer, err := casbin.NewEnforcer("model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}

	// 3. 加载策略（如果已经存在）
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}
	return enforcer
}

// 核心鉴权系统
type RoleCtl struct {
	E   *casbin.Enforcer
	zap *zap.Logger
}

func NewCtl(e *casbin.Enforcer, zap *zap.Logger) *RoleCtl {
	return &RoleCtl{E: e, zap: zap}
}

// 添加规则
func (this *RoleCtl) AddRule(role string, path string, method string) {
	_, err := this.E.AddPolicy(role, path, method)
	if err != nil {
		this.zap.Error("添加鉴权规则失败:" + err.Error())
	}
	err = this.E.SavePolicy()
	if err != nil {
		this.zap.Error("保存鉴权规则失败:" + err.Error())
	}
}

// 添加角色
func (this *RoleCtl) AddRole(role string, parentRole string) {
	_, err := this.E.AddGroupingPolicy(role, parentRole)
	if err != nil {
		this.zap.Error("添加继承关系失败:" + err.Error())
	}
	err = this.E.SavePolicy()
	if err != nil {
		this.zap.Error("保存继承规则失败:" + err.Error())
	}
}

func (this *RoleCtl) GetAllParentRoles(role string) []string {
	parentRoles := make(map[string]bool) // 使用 map 防止重复
	var dfs func(r string)

	dfs = func(r string) {
		roles, _ := this.E.GetRolesForUser(r) // 获取直接父角色
		for _, parentRole := range roles {
			if !parentRoles[parentRole] { // 防止重复访问
				parentRoles[parentRole] = true
				dfs(parentRole) // 递归获取父角色
			}
		}
	}
	dfs(role)
	// 将 map 转为 slice 返回
	result := make([]string, 0, len(parentRoles))
	for r := range parentRoles {
		result = append(result, r)
	}
	return result
}

func (this *RoleCtl) GetRulesForRole(role string) []Rule {
	// 使用 GetFilteredPolicy 按照第一个字段（角色/主体）筛选规则
	res := make([]Rule, 0)
	rules, _ := this.E.GetFilteredPolicy(0, role)
	for _, item := range rules {
		res = append(res, Rule{
			Role:   item[0],
			Path:   item[1],
			Method: item[2],
		})
	}
	return res
}

func (this *RoleCtl) GetRulesForInheritRole(role string) []AggregatedRule {
	list := this.GetAllParentRoles(role)
	res := make([]Rule, 0)
	for _, item := range list {
		res = append(res, this.GetRulesForRole(item)...)
	}
	return this.AggregateRules(res)
}
func (this *RoleCtl) AggregateRules(rules []Rule) []AggregatedRule {
	// 使用 map 去重，key 是 path + method，value 是角色数组
	aggregatedMap := make(map[string]AggregatedRule)

	for _, rule := range rules {
		// 组合 Path 和 Method 作为 key
		key := rule.Path + "|" + rule.Method

		// 检查是否已经存在
		if aggregatedRule, exists := aggregatedMap[key]; exists {
			// 如果已存在，将 Role 添加到角色数组中（避免重复）
			exists := false
			for _, r := range aggregatedRule.Roles {
				if r == rule.Role {
					exists = true
					break
				}
			}
			if !exists {
				aggregatedRule.Roles = append(aggregatedRule.Roles, rule.Role)
			}
			aggregatedMap[key] = aggregatedRule
		} else {
			// 如果不存在，新建 AggregatedRule
			aggregatedMap[key] = AggregatedRule{
				Path:   rule.Path,
				Method: rule.Method,
				Roles:  []string{rule.Role},
			}
		}
	}

	// 将 map 转换为 slice
	result := make([]AggregatedRule, 0, len(aggregatedMap))
	for _, aggregatedRule := range aggregatedMap {
		result = append(result, aggregatedRule)
	}

	return result
}

func (this *RoleCtl) DeleteRulesForRole(role string) error {
	// 使用 RemoveFilteredPolicy 删除与指定角色相关的规则
	this.E.RemoveFilteredPolicy(0, role)
	// 保存更改到持久化存储
	err := this.E.SavePolicy()
	return err
}

func (this *RoleCtl) DeleteRoleInheritance(childRole string, parentRole string) error {
	// 删除该角色的所有继承关系
	this.E.RemoveGroupingPolicy(childRole, parentRole)
	// 保存到持久化存储
	err := this.E.SavePolicy()
	return err
}
