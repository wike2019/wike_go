package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 权限注册函数
func (this *GCore) GetWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler gin.HandlerFunc, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, Input(query, body, header), Output(output), group.Path()+path, http.MethodGet)
	r.GET(path, handler)

}
func (this *GCore) PostWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler gin.HandlerFunc, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, Input(query, body, header), Output(output), group.Path()+path, http.MethodPost)
	r.POST(path, handler)
}
func (this *GCore) DelWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler gin.HandlerFunc, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, Input(query, body, header), Output(output), group.Path()+path, http.MethodDelete)
	r.DELETE(path, handler)
}
func (this *GCore) PutWithRbac(r gin.IRoutes, group Controller, groupName string, path string, handler gin.HandlerFunc, name string) {
	query, body, header, output := group.GetInnerData()
	this.db.ApiTable(name, groupName, Input(query, body, header), Output(output), group.Path()+path, http.MethodPut)
	r.PUT(path, handler)
}
