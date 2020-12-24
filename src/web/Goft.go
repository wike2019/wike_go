package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-ioc"
	"github.com/wike2019/wike_go/src/core/config"
	"github.com/wike2019/wike_go/src/core/ioc"
	"strings"
	"sync"
)

//web 相关操作

var innerRouter *GoftTree // inner tree node . backup httpmethod and path
var routerOnce sync.Once

func getInnerRouter() *GoftTree {
	routerOnce.Do(func() {
		innerRouter = NewGoftTree()
	})
	return innerRouter
}

type Goft struct {
	*gin.Engine
	 g            *gin.RouterGroup // 保存 group对象
	currentGroup string // temp-var for group string
}

func Ignite() *Goft {
	g := &Goft{Engine: gin.New(),}
	g.Use(ErrorHandler()) //强迫加载的异常处理中间件

	return g
}
func (this *Goft) Launch() {
	var port int32 = 8080
	if configData := Injector.BeanFactory.Get((*config.SysConfig)(nil)); configData != nil {
		port = configData.(*config.SysConfig).Server.Port
	}
	this.Run(fmt.Sprintf(":%d", port))
}

func (this *Goft) Handle(httpMethod, relativePath string, handler interface{}) *Goft {
	if h := Convert(handler); h != nil {
		getInnerRouter().addRoute(httpMethod, this.getPath(relativePath), h) // for future
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}
func (this *Goft) getPath(relativePath string) string {
	g := "/" + this.currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}
func (this *Goft) HandleWithFairing(httpMethod, relativePath string, handler interface{}, fairings ...Fairing) *Goft {
	if h := Convert(handler); h != nil {
		getInnerRouter().addRoute(httpMethod, this.getPath(relativePath), fairings) //for future
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

// 注册中间件
func (this *Goft) Attach(f ...Fairing) *Goft {
	for _, f1 := range f {
		Injector.BeanFactory.Set(f1)
	}
	getFairingHandler().AddFairing(f...)
	return this
}



func (this *Goft) Mount(group string, classes ...IClass) *Goft {
	this.g = this.Group(group)
	for _, class := range classes {
		this.currentGroup = group
		class.Build(this)
		//this.beanFactory.inject(class)
		ioc.New().Beans(class)
	}
	return this
}

