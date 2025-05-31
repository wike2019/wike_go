package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/common"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"github.com/wike2019/wike_go/server/model"
)

func (this CoreCtl) getApiParams(context *gin.Context) (*ctl.Page, *model.APISearch, *common.Empty, *ctl.Controller) {
	Item := &model.APISearch{}
	c := this.SetContext(context)
	c.PageInit()
	err := json.Unmarshal(c.Data, Item)
	ctl.Error(err, 400)
	return &c.Page, Item, &common.Empty{}, c
}
func (this CoreCtl) getApi(context *gin.Context) interface{} {
	res, err := this.service.GetApiService(this.getApiParams(context))
	ctl.Error(err, 400)
	return res
}

func (this CoreCtl) getApiFrontCommonParams(context *gin.Context) (*common.Empty, *model.APISearch, *common.Empty, *ctl.Controller) {
	c := this.SetContext(context)
	return &common.Empty{}, &model.APISearch{}, &common.Empty{}, c
}
func (this CoreCtl) getApiFront(context *gin.Context) interface{} {
	res, err := this.service.GetApiFrontService(this.getApiFrontCommonParams(context))
	ctl.Error(err, 400)
	return res
}
