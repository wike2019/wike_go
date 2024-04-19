package casbinInit

import (
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"os"
)

// 鉴权系统
func New() *casbin.Enforcer {
	os.WriteFile("policy.csv", []byte(""), os.ModePerm)
	e, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		panic("Casbin初始化失败" + err.Error())
	}
	return e
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
