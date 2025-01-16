package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/model"
	"github.com/wike2019/wike_go/pkg/core"
	"github.com/wike2019/wike_go/pkg/service/ctl"
)

type CoreCtl struct {
	ctl.Controller
	DB *core.CoreDb
}

func NewRouter(DB *core.CoreDb) *CoreCtl {
	return &CoreCtl{DB: DB}
}
func (this CoreCtl) Path() string {
	return "/core"
}
func (this *CoreCtl) Build(r *gin.RouterGroup, GCore *core.GCore) {
	this.SetDoc(ctl.Page{}, model.API{}, nil, ctl.PageDoc[model.API]())
	GCore.PostWithRbac(r, this, "系统内部接口", "/api", this.getApi, "获取接口列表")
	GCore.GetWithRbac(r, this, "系统内部接口", "/menu", this.getMenu, "获取菜单列表")
	this.SetDoc(nil, nil, nil, ctl.HttpDocList[model.SysDictionary]{})
	GCore.GetWithRbac(r, this, "系统内部接口", "/dictionaryList", this.dictionaryList, "获取字典列表")
	GCore.PostWithRbac(r, this, "系统内部接口", "/dictionaryCreate", this.dictionaryCreate, "添加字典")
	GCore.PostWithRbac(r, this, "系统内部接口", "/dictionaryDelete", this.dictionaryDelete, "删除字典")
	GCore.PostWithRbac(r, this, "系统内部接口", "/dictionaryUpdate", this.dictionaryUpdate, "修改字典")
	GCore.PostWithRbac(r, this, "系统内部接口", "/dictionaryItem", this.dictionaryItem, "模糊搜索字典")
	GCore.GetWithRbac(r, this, "系统内部接口", "/systemInfo", this.SystemInfo, "获取服务器信息")
}
