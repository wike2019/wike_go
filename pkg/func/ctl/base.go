package ctl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/common"
	"io"
	"net/http"
	"strconv"
)

type Controller struct {
	Page
	*gin.Context
	query  interface{}
	body   interface{}
	header interface{}
	output interface{}
	Data   []byte
}

func (r *Controller) clear() {
	r.query = nil
	r.body = nil
	r.header = nil
	r.output = nil
	r.Data = nil
}
func (r *Controller) SetDoc(input interface{}, output interface{}) {
	r.clear()
	r.getFuncMateInfo(input, output)

}
func (r *Controller) SetDocRaw(paramInstance interface{}, bodyInstance interface{}, headerInstance interface{}, returnInstance interface{}) {
	r.clear()
	r.query = paramInstance
	r.body = bodyInstance
	r.header = headerInstance
	r.output = returnInstance
}
func (r *Controller) SetEmpty() {
	r.clear()
	r.SetDocRaw(common.Empty{}, common.Empty{}, common.Empty{}, DataList[string]{})
}
func (r *Controller) SetEmptyWithHeader(header interface{}) {
	r.clear()
	r.SetDocRaw(common.Empty{}, common.Empty{}, header, DataList[string]{})
}
func (r *Controller) GetInnerData() (interface{}, interface{}, interface{}, interface{}) {
	defer r.clear()
	return r.query, r.body, r.header, r.output
}
func (r *Controller) PageInit() {
	queryStr := r.Context.Query("page")
	countStr := r.Context.Query("count")
	// 转换为 int
	if queryStr != "" {
		page, err := strconv.Atoi(queryStr)
		Error(err, 400)
		r.CurPage = page
	}
	if countStr != "" {
		count, err := strconv.Atoi(countStr)
		Error(err, 400)
		r.Count = count
	}
	if r.CurPage == 0 {
		r.CurPage = 1
	}
	if r.Count == 0 {
		r.Count = 10
	}
	r.Offset = (r.CurPage - 1) * r.Count
}
func (r *Controller) SetContext(ctx *gin.Context) *Controller {
	clone := new(Controller)
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		Error(errors.New("参数错误"+err.Error()), 400)
		return nil
	}
	clone.Data = body
	clone.Context = ctx
	return clone
}
func (r *Controller) OK(data interface{}) {
	r.JSON(http.StatusOK, data)
}
func (r *Controller) OKWithMsg(msg string) {
	r.JSON(http.StatusOK, gin.H{
		"msg":     msg,
		"traceId": r.GetString("trace_id"),
		"code":    Success,
	})
}

func OK(r *gin.Context, data interface{}) {
	r.JSON(http.StatusOK, data)
}
func OKWithMsg(r *gin.Context, msg string) {
	r.JSON(http.StatusOK, gin.H{
		"msg":     msg,
		"traceId": r.GetString("trace_id"),
		"code":    Success,
	})
}
func List[T any](msg string, list []T, context *Controller) *PageList[T] {
	return &PageList[T]{
		Data:    list,
		Page:    context.Page,
		Msg:     msg,
		TraceId: context.GetString("trace_id"),
		Code:    Success,
	}
}
func Item[T any](msg string, list T, context *Controller) *DataList[T] {
	return &DataList[T]{
		Data:    list,
		Msg:     msg,
		TraceId: context.GetString("trace_id"),
		Code:    Success,
	}
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
