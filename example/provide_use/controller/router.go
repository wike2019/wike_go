package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/model"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
	ctl.Controller
	dog *model.Dog
	cat *model.Cat
}

func NewRouter(dog *model.Dog, cat *model.Cat) *router {
	return &router{dog: dog, cat: cat}
}
func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.OK(this.dog.Name + this.cat.Name)
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/animal", this.healtzh)
}
func (this router) Path() string {
	return "/"
}
