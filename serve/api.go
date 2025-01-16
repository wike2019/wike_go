package serve

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/model"
	"github.com/wike2019/wike_go/pkg/service/ctl"
)

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
