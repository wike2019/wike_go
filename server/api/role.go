package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/common"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"github.com/wike2019/wike_go/server/model"
)

//
//	func (this CoreCtl) roleCreate(context *gin.Context) {
//		Item := &model2.Role{}
//		c := this.SetContext(context)
//		err := json.Unmarshal(c.Data, Item)
//		ctl.Error(err, 400)
//		if Item.Parent == "" {
//			Item.Parent = "nologin"
//		}
//		err = ctl.CreateItem[*model2.Role](this.DB.DB, Item)
//		ctl.Error(err, 400)
//		c.Success("添加角色成功", Item)
//	}

func (this CoreCtl) roleListParams(context *gin.Context) (*model.Role, *common.Empty, *common.Empty, *ctl.Controller) {
	query := &model.Role{}
	c := this.SetContext(context)
	err := c.ShouldBindQuery(query)
	ctl.Error(err, 400)
	return query, &common.Empty{}, &common.Empty{}, c
}
func (this CoreCtl) roleList(context *gin.Context) interface{} {
	list, err := this.service.RoleListService(this.roleListParams(context))
	ctl.Error(err, 400)
	return list
}

func (this CoreCtl) roleDataList(context *gin.Context) interface{} {
	c := this.SetContext(context)
	list, err := ctl.ListItemAll[*model.Role](this.DBInstance.DB, model.Role{}, nil)
	ctl.Error(err, 400)
	data := ctl.Item[[]*model.Role]("获取api列表成功", list, c)
	return data
}

//
//	func (this CoreCtl) roleDelete(context *gin.Context) {
//		Item := &common.IDSearch{}
//		c := this.SetContext(context)
//		err := c.ShouldBindQuery(Item)
//		ctl.Error(err, 400)
//		err = ctl.DeleteItem[*model2.Role](this.DB.DB, Item.ID)
//		ctl.Error(err, 400)
//		c.Success("删除角色成功", Item)
//	}

func (this CoreCtl) ruleCreateParams(context *gin.Context) (*common.Empty, *model.RuleInput, *common.Empty, *ctl.Controller) {
	Item := &model.RuleInput{}
	c := this.SetContext(context)
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	return &common.Empty{}, Item, &common.Empty{}, c
}
func (this CoreCtl) ruleCreate(context *gin.Context) interface{} {
	data, err := this.service.RuleCreateService(this.ruleCreateParams(context))
	ctl.Error(err, 400)
	return data
}

func (this CoreCtl) ruleSync(context *gin.Context) interface{} {
	c := this.SetContext(context)
	err := this.DBInstance.DB.Exec("delete from casbin_rule").Error
	ctl.Error(err, 400)
	search := model.Rule{}
	data, err := ctl.ListItemAll[model.Rule](this.DBInstance.DB, search, nil)
	for _, item := range data {
		this.roleCtl.AddRule(item.Role, item.Path, item.Method)
	}
	search2 := model.Role{}
	data2, err := ctl.ListItemAll[model.Role](this.DBInstance.DB, search2, nil)
	for _, item := range data2 {
		this.roleCtl.AddRole(item.Child, item.Parent)
	}
	return ctl.Item[common.Empty]("信息同步成功", common.Empty{}, c)
}

func (this CoreCtl) ruleInfoParams(context *gin.Context) (*model.RoleInput, *common.Empty, *common.Empty, *ctl.Controller) {
	query := &model.RoleInput{}
	c := this.SetContext(context)
	err := c.Context.ShouldBindQuery(query)
	ctl.Error(err, 400)
	return query, &common.Empty{}, &common.Empty{}, c
}
func (this CoreCtl) ruleInfo(context *gin.Context) interface{} {
	res, err := this.service.RuleInfoService(this.ruleInfoParams(context))
	ctl.Error(err, 400)
	return res
}
