package fastRouter

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/common"
	"github.com/wike2019/wike_go/pkg/core"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"gorm.io/gorm"
)

func FastRouter[Q common.IDSearcher, T any, H any](DB *gorm.DB, DBInstance *core.CoreDb, r *gin.RouterGroup, GCore *core.GCore, path string, name string, groupName string, prefix string) {
	res := &FastInstance[Q, T, H]{
		DB:         DB,
		path:       path,
		name:       name,
		prefix:     prefix,
		groupName:  groupName,
		DBInstance: DBInstance,
	}
	res.Build(r, GCore)
}

type FastInstance[Q common.IDSearcher, T any, H any] struct {
	ctl.Controller
	DB         *gorm.DB
	path       string
	name       string
	groupName  string
	prefix     string
	DBInstance *core.CoreDb
}

func (this *FastInstance[Q, T, H]) Build(r *gin.RouterGroup, GCore *core.GCore) {
	var query Q
	var data T
	var header H
	this.SetDocRaw(common.Empty{}, data, header, ctl.DataList[string]{})
	GCore.PostWithRbac(r, this, this.groupName, this.path+"/create", this.Create, "添加"+this.name)
	this.SetDocRaw(query, data, header, ctl.DataList[string]{})
	GCore.PostWithRbac(r, this, this.groupName, this.path+"/update", this.Update, "修改"+this.name)
	this.SetDocRaw(query, common.Empty{}, header, ctl.DataList[string]{})
	GCore.DelWithRbac(r, this, this.groupName, this.path+"/delete", this.Del, "删除"+this.name)
}

func (this *FastInstance[Q, T, H]) Path() string {
	if this.prefix[0] != '/' {
		this.prefix = "/" + this.prefix
	}
	return this.prefix
}

func (this *FastInstance[Q, T, H]) Create(context *gin.Context) interface{} {
	c := this.SetContext(context)
	var data T
	err := json.Unmarshal(c.Data, &data)
	ctl.Error(err, 400)
	err = ctl.CreateItem(this.DB, &data)
	ctl.Error(err, 400)
	return "添加" + this.name + "成功"
}
func (this *FastInstance[Q, T, H]) Update(context *gin.Context) interface{} {
	var query Q
	var data T
	err := context.ShouldBindQuery(&query)
	ctl.Error(err, 400)
	c := this.SetContext(context)
	err = json.Unmarshal(c.Data, &data)
	ctl.Error(err, 400)
	id := query.GetID()
	err = ctl.UpdateItem[T](this.DB, data, id)
	ctl.Error(err, 400)
	return "修改" + this.name + "成功"
}

func (this *FastInstance[Q, T, H]) Del(context *gin.Context) interface{} {
	var query Q
	err := context.ShouldBindQuery(&query)
	ctl.Error(err, 400)
	this.SetContext(context)
	id := query.GetID()
	err = ctl.DeleteItem[T](this.DB, id)
	ctl.Error(err, 400)
	return "删除" + this.name + "成功"
}
