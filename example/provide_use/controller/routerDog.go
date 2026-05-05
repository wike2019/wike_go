package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/provide"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type routerDog struct {
	ctl.Controller
	Animal1 *provide.Animal
	Animal2 *provide.Animal
}

func NewRouterDog(dog *provide.Animal, dog2 *provide.Animal) *routerDog {
	return &routerDog{Animal1: dog, Animal2: dog2}
}
func (this *routerDog) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.OK(this.Animal1.Show() + this.Animal2.Show())
}

func (this routerDog) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/healthz", this.healtzh)
}
func (this routerDog) Path() string {
	return "/dog"
}
