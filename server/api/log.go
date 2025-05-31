package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"github.com/wike2019/wike_go/pkg/func/memorylog"
)

func (this CoreCtl) log(context *gin.Context) interface{} {
	c := this.SetContext(context)
	res := ctl.Item("获取日志成功", memorylog.LogInfo.All(), c)
	return res
}
