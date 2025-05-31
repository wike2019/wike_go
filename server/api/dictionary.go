package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/common"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"github.com/wike2019/wike_go/server/model"
)

func (this CoreCtl) dictionaryListParams(context *gin.Context) (*ctl.Page, *model.SysDictionarySearch, *common.Empty, *ctl.Controller) {
	Item := &model.SysDictionarySearch{}
	c := this.SetContext(context)
	c.PageInit()
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	return &c.Page, Item, &common.Empty{}, c
}

func (this CoreCtl) dictionaryList(context *gin.Context) interface{} {
	data, err := this.service.DictionaryListService(this.dictionaryListParams(context))
	ctl.Error(err, 400)
	return data
}

func (this CoreCtl) dictionaryItemParams(context *gin.Context) (*common.IDSearch, *common.Empty, *common.Empty, *ctl.Controller) {
	query := &common.IDSearch{}
	c := this.SetContext(context)
	err := c.ShouldBindQuery(query)
	ctl.Error(err, 400)
	return query, &common.Empty{}, &common.Empty{}, c
}
func (this CoreCtl) dictionaryItem(context *gin.Context) interface{} {
	data, err := this.service.DictionaryItemService(this.dictionaryItemParams(context))
	ctl.Error(err, 400)
	return data
}

func (this CoreCtl) dictionarySearchParams(context *gin.Context) (*common.KeyWordsSearch, *common.Empty, *common.Empty, *ctl.Controller) {
	query := &common.KeyWordsSearch{}
	c := this.SetContext(context)
	err := c.ShouldBindQuery(query)
	ctl.Error(err, 400)
	return query, &common.Empty{}, &common.Empty{}, c
}
func (this CoreCtl) dictionarySearch(context *gin.Context) interface{} {
	data, err := this.service.DictionarySearchService(this.dictionarySearchParams(context))
	ctl.Error(err, 400)
	return data
}
