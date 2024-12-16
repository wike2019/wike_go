package ctl

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Controller struct {
	PageTime
	*gin.Context
	query  interface{}
	body   interface{}
	header interface{}
	output interface{}
}

func (r *Controller) clear() {
	r.query = nil
	r.body = nil
	r.header = nil
	r.output = nil
}
func (r *Controller) SetDoc(query interface{}, body interface{}, header interface{}, output interface{}) {
	r.clear()
	r.query = query
	r.body = body
	r.header = header
	r.output = output
}
func (r *Controller) GetInnerData() (interface{}, interface{}, interface{}, interface{}) {
	return r.query, r.body, r.header, r.output
}
func (r *Controller) SetContext(ctx *gin.Context) *Controller {
	clone := new(Controller)
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		Error(errors.New("参数错误"), 400)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
	err = ctx.ShouldBind(clone)
	if err != nil {
		Error(errors.New("参数错误"), 400)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
	clone.Context = ctx

	if clone.Page.CurPage == 0 {
		clone.Page.CurPage = 1
	}
	if clone.Page.Count == 0 {
		clone.Page.Count = 10
	}
	clone.Page.Offset = (clone.Page.CurPage - 1) * clone.Page.Count

	return clone
}
func (r *Controller) Success(msg string, data interface{}) {
	r.JSON(http.StatusOK, gin.H{
		"code":     Success,
		"msg":      msg,
		"data":     data,
		"trace_id": r.GetString("trace_id"),
	})
}
func (r *Controller) List(msg string, list interface{}) {
	if list == nil {
		list = make([]interface{}, 0)
	}
	result := PageResult.Get().(*Pagination)
	defer func() {
		result.Reset()
		PageResult.Put(result)
	}()
	result.Page = r.Page
	result.List = list
	r.Success(msg, result)
}

func (r *Controller) Failed(code errorCode, msg string) {
	errMsg := codeMsg[code] + ": " + msg
	if code != Success {
		r.Set("error_code", int(code))
		r.Set("error_msg", msg)
	}
	r.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":     code,
		"msg":      errMsg,
		"data":     nil,
		"trace_id": r.GetString("trace_id"),
	})
}
