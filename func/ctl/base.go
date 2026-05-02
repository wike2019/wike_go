package ctl

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	Page
	*gin.Context
	Data []byte
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

func (r *Controller) OKWithList(list interface{}) {
	r.JSON(http.StatusOK, gin.H{
		"data":    list,
		"page":    r.Page,
		"traceId": r.GetString("trace_id"),
		"code":    Success,
	})
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

func (r *Controller) IsExist(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
