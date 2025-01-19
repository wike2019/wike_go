package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/wike2019/wike_go/model"
	ctl "github.com/wike2019/wike_go/pkg/service/ctl"
	"gorm.io/gorm"
)

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
	res, err := ctl.GetItem[*model.SysDictionary](this.DB.DB, Item, func(db *gorm.DB) *gorm.DB {
		return db.Preload("SysDictionaryDetails", "status = ?", 2)
	})
	ctl.Error(err, 400)
	c.Success("获取字典列表成功", res)
}
