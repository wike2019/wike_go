package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/pkg/doc"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	"net/http"
	"reflect"
)

// 权限注册函数
func (this *GCore) GetWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler func(c *gin.Context) interface{}, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, doc.Input(query, body, header), doc.Output(output), group.Path()+path, http.MethodGet)
	r.GET(path, WrapCtxHandler(handler))

}
func (this *GCore) PostWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler func(c *gin.Context) interface{}, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, doc.Input(query, body, header), doc.Output(output), group.Path()+path, http.MethodPost)
	r.POST(path, WrapCtxHandler(handler))
}
func (this *GCore) DelWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler func(c *gin.Context) interface{}, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, doc.Input(query, body, header), doc.Output(output), group.Path()+path, http.MethodDelete)
	r.DELETE(path, WrapCtxHandler(handler))
}
func (this *GCore) PutWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler func(c *gin.Context) interface{}, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, doc.Input(query, body, header), doc.Output(output), group.Path()+path, http.MethodPut)
	r.PUT(path, WrapCtxHandler(handler))
}
func WrapCtxHandler(fn func(c *gin.Context) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := fn(c)
		if reflect.TypeOf(result).Kind() == reflect.String {
			str := fmt.Sprintf("%v", result)
			ctl.OKWithMsg(c, str)
		} else {
			ctl.OK(c, result)
		}
	}
}
