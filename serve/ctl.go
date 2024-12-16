package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/wike2019/wike_go/model"
	"github.com/wike2019/wike_go/pkg/core"
	ctl "github.com/wike2019/wike_go/pkg/service/ctl"
	"gorm.io/gorm"
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
func (this CoreCtl) getApi(context *gin.Context) {
	Item := &model.API{}
	c := this.SetContext(context)
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	data, total, err := ctl.ListItem[*model.API](this.DB.DB, c.Offset, c.Count, Item, nil)
	ctl.Error(err, 400)
	c.FormatTotal(total)
	c.List("接口列表查询成功", data)
}

func (this CoreCtl) getMenu(context *gin.Context) {
	//c := this.SetContext(context)
	str := `
[
  {
    "name": "Dashboard",
    "path": "/dashboard",
	"component":"/dashboard/VIndex",
    "children": [
      {
        "name": "Overview",
        "path": "overview",
		"component":"/dashboard/VOverview"
      },
      {
        "name": "Stats",
        "path": "stats",
		"component":"/dashboard/VStats"
      }
    ]
  },
  {
    "name": "Settings",
    "path": "/settings/VIndex",
    "children": [],
	"component":"/settings/VIndex"
  }
]

`
	context.Header("Content-Type", "application/json")
	context.String(200, str)
}

func (this CoreCtl) dictionaryList(context *gin.Context) {
	Item := &model.SysDictionary{}
	c := this.SetContext(context)
	data, err := ctl.ListItemAll[*model.SysDictionary](this.DB.DB, Item, func(db *gorm.DB) *gorm.DB {
		return db.Preload("SysDictionaryDetails")
	})
	ctl.Error(err, 400)
	c.Success("获取字典列表成功", data)
}

func (this CoreCtl) dictionaryCreate(context *gin.Context) {
	Item := &model.SysDictionary{}
	c := this.SetContext(context)
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	err = ctl.CreateItem[*model.SysDictionary](this.DB.DB, Item)
	ctl.Error(err, 400)
	c.Success("添加字典成功", Item)
}

func (this CoreCtl) dictionaryDelete(context *gin.Context) {
	Item := &model.SysDictionary{}
	c := this.SetContext(context)
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	err = ctl.DeleteItem[*model.SysDictionary](this.DB.DB, Item)
	ctl.Error(err, 400)
	c.Success("删除字典成功", Item)
}
func (this CoreCtl) dictionaryUpdate(context *gin.Context) {
	Item := &model.SysDictionary{}
	c := this.SetContext(context)
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	err = ctl.UpdateItem[*model.SysDictionary](this.DB.DB, Item)
	ctl.Error(err, 400)
	c.Success("修改字典成功", Item)
}

func (this CoreCtl) dictionaryItem(context *gin.Context) {
	Item := &model.SysDictionary{}
	c := this.SetContext(context)
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	err = ctl.GetItem[*model.SysDictionary](this.DB.DB, Item, func(db *gorm.DB) *gorm.DB {
		return db.Preload("SysDictionaryDetails", "status = ?", 2)
	})
	ctl.Error(err, 400)
	c.Success("获取字典列表成功", Item)
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
}
