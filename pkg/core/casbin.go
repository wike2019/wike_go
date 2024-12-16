package core

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"log"
	"os"
)

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
func (this *RoleCtl) AddRule(role string, prefix string, path string, method string) {
	_, err := this.E.AddPolicy(role, prefix+path, method)
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
