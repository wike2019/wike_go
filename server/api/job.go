package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"github.com/wike2019/wike_go/server/model"
)

func (this CoreCtl) jobList(context *gin.Context) interface{} {
	c := this.SetContext(context)
	Item := model.Job{}
	data, err := ctl.ListItemAll[model.Job](this.DBInstance.DB, Item, nil)
	ctl.Error(err, 400)
	res := ctl.Item("获取日志成功", data, c)
	return res
}
