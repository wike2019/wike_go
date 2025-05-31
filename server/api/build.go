package api

import (
	"github.com/gin-gonic/gin"
	model2 "github.com/wike2019/wike_go/common"
	"github.com/wike2019/wike_go/pkg/core"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"github.com/wike2019/wike_go/pkg/func/ctl/fastRouter"
	"github.com/wike2019/wike_go/pkg/os"
	"github.com/wike2019/wike_go/server/model"
	"github.com/wike2019/wike_go/server/service"
)

type CoreCtl struct {
	ctl.Controller
	DBInstance *core.CoreDb
	roleCtl    *core.RoleCtl
	service    *service.Service
}

func NewRouter(DB *core.CoreDb, roleCtl *core.RoleCtl, service *service.Service) *CoreCtl {
	return &CoreCtl{DBInstance: DB, roleCtl: roleCtl, service: service}
}
func (this CoreCtl) Path() string {
	return "/core"
}
func (this CoreCtl) CoreDb() *core.CoreDb {
	return this.DBInstance
}
func (this *CoreCtl) Build(r *gin.RouterGroup, GCore *core.GCore) {
	//t := r.Use(core.Root(), core.Authorizer(this.roleCtl.E))
	fastRouter.FastRouter[model2.IDSearch, model.SysDictionary, model2.Empty](this.DBInstance.DB, this.DBInstance, r, GCore, "/dictionary", "字典", "系统接口", "core")
	fastRouter.FastRouter[model2.IDSearch, model.SysDictionaryDetail, model2.Empty](this.DBInstance.DB, this.DBInstance, r, GCore, "/dictionaryDetail", "字典详细", "系统接口", "core")
	fastRouter.FastRouter[model2.IDSearch, model.Role, model2.Empty](this.DBInstance.DB, this.DBInstance, r, GCore, "/role", "角色", "系统接口", "core")

	this.SetDoc(this.getApiParams, this.service.GetApiService)
	GCore.PostWithRbac(r, this, "系统内部接口", "/api", this.getApi, "获取接口列表")

	this.SetDoc(this.getApiFrontCommonParams, this.service.GetApiFrontService)
	GCore.GetWithRbac(r, this, "系统内部接口", "/getApiFront", this.getApiFront, "获得接口列表")

	this.SetDoc(this.dictionaryListParams, this.service.DictionaryListService)
	GCore.PostWithRbac(r, this, "系统内部接口", "/dictionaryList", this.dictionaryList, "获取字典列表")
	this.SetDoc(this.dictionaryItemParams, this.service.DictionaryItemService)
	GCore.GetWithRbac(r, this, "系统内部接口", "/dictionaryItem", this.dictionaryItem, "模糊搜索字典")

	this.SetDoc(this.dictionarySearchParams, this.service.DictionarySearchService)
	GCore.GetWithRbac(r, this, "系统内部接口", "/dictionarySearch", this.dictionarySearch, "模糊搜索字典")
	this.SetDocRaw(nil, nil, nil, ctl.DataList[os.Server]{})
	GCore.GetWithRbac(r, this, "系统内部接口", "/systemInfo", this.SystemInfo, "获取服务器信息")

	this.SetDoc(this.roleListParams, this.service.RoleListService)
	GCore.GetWithRbac(r, this, "系统内部接口", "/roleList", this.roleList, "获取角色名称列表")

	this.SetDoc(this.ruleCreateParams, this.service.RuleCreateService)
	GCore.PostWithRbac(r, this, "系统内部接口", "/ruleCreate", this.ruleCreate, "创建规则")
	////
	this.SetDoc(this.ruleInfoParams, this.service.RuleInfoService)
	GCore.GetWithRbac(r, this, "系统内部接口", "/ruleInfo", this.ruleInfo, "获得接口权限详细")
	////
	this.SetEmpty()
	GCore.GetWithRbac(r, this, "系统内部接口", "/ruleSync", this.ruleSync, "同步权限")
	this.SetEmpty()
	GCore.GetWithRbac(r, this, "系统内部接口", "/log", this.log, "接口日志查询")
	this.SetEmptyWithHeader(&model2.Empty{})
	GCore.GetWithRbac(r, this, "系统内部接口", "/jobList", this.jobList, "定时任务查询")

	this.SetDocRaw(nil, nil, nil, ctl.DataList[[]model.Role]{})
	GCore.GetWithRbac(r, this, "系统内部接口", "/roleDataList", this.roleDataList, "获取角色列表")

}
