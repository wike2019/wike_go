package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/model"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type router2 struct {
	ctl.Controller
}

// 注册第二个控制器
func NewRouter2() *router2 {
	return &router2{}
}

func (this *router2) PostData(context *gin.Context) {
	c := this.SetContext(context)
	data := &model.ModelSum{}
	json.Unmarshal(c.Data, &data)
	ctl.Error(data.Check(), 400)
	c.OK(data.Num1 + data.Num2)
}
func (this router2) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.POST("/sum", this.PostData)
}
func (this router2) Path() string {
	return "/"
}
