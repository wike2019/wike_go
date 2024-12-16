package ctl

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

type Controller struct {
	PageTime
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
func (r *Controller) SetDoc(query interface{}, body interface{}, header interface{}, output interface{}) {
	r.clear()
	r.query = query
	r.body = body
	r.header = header
	r.output = output
}
func (r *Controller) GetInnerData() (interface{}, interface{}, interface{}, interface{}) {
	defer r.clear()
	return r.query, r.body, r.header, r.output
}
func (r *Controller) SetContext(ctx *gin.Context) *Controller {
	clone := new(Controller)
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		Error(errors.New("参数错误"+err.Error()), 400)
		return nil
	}
	clone.Data = body
	err = json.Unmarshal(clone.Data, clone)

	if err != nil {
		Error(errors.New("参数错误"+err.Error()), 400)
		return nil
	}
	queryStr := ctx.Query("page")
	countStr := ctx.Query("count")
	// 转换为 int
	if queryStr != "" {
		page, err := strconv.Atoi(queryStr)
		Error(err, 400)
		clone.CurPage = page
	}
	if countStr != "" {
		count, err := strconv.Atoi(countStr)
		Error(err, 400)
		clone.Count = count
	}
	clone.Context = ctx
	if clone.CurPage == 0 {
		clone.CurPage = 1
	}
	if clone.Count == 0 {
		clone.Count = 10
	}
	clone.Offset = (clone.CurPage - 1) * clone.Count
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

	result.SumPage = r.SumPage
	result.SumCount = r.SumCount
	result.CurPage = r.CurPage
	result.Count = r.Count
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
