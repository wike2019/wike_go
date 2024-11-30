package controller

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type errorCode int

// 核心控制器基础方法
const (
	Success      errorCode = 200
	Failed       errorCode = 500
	ParamError   errorCode = 400
	NotFound     errorCode = 404
	UnAuthorized errorCode = 401
)

var codeMsg = map[errorCode]string{
	Success:      "正常",
	Failed:       "系统异常",
	ParamError:   "参数错误",
	NotFound:     "记录不存在",
	UnAuthorized: "未授权",
}

type LocalTimeInt int64

func (t LocalTimeInt) MarshalJSON() ([]byte, error) {
	tTime := time.Unix(int64(t), 0)
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if &t == nil || tTime.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

// UnmarshalJSON 为 LocalTimeInt 类型自定义 JSON 反序列化方法
func (lti *LocalTimeInt) UnmarshalJSON(data []byte) error {
	// 将 JSON 中的数字解析为 int64
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}

	// 将 int64 的时间戳转换为 time.Time，然后更新 LocalTimeInt 的值
	*lti = LocalTimeInt(timestamp)
	return nil
}

// ToTime 方法将 LocalTimeInt 类型转换为 time.Time
func (lti LocalTimeInt) ToTime() time.Time {
	return time.Unix(int64(lti), 0)
}

type LocalTime time.Time

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *LocalTime) String() string {
	// 如果时间 null 那么我们需要把返回的值进行修改
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s", time.Time(*t).Format("2006-01-02 15:04:05"))
}

func (t *LocalTime) IsZero() bool {
	return time.Time(*t).IsZero()
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")

	local, _ := time.LoadLocation("Asia/Shanghai")

	t1, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, local)
	*t = LocalTime(t1)
	return err
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if &t == nil || t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

type Page struct {
	SumPage  int `json:"sumPage" swaggerignore:"true" desc:"总页数"`  // 总页数
	SumCount int `json:"sumCount" swaggerignore:"true" desc:"总条数"` // 总条数
	CurPage  int `json:"page" form:"page" desc:"当前页"   `           // 当前页
	Offset   int `json:"-"`                                        // 起始量
	Count    int `json:"count" form:"count" desc:"每页返回数量"`         // 每页返回数量
}

type Pagination struct {
	List interface{} `json:"list"`
	Page Page
}

func (this *Pagination) Reset() {
	this.List = nil
	this.Page.SumPage = 0
	this.Page.SumCount = 0
	this.Page.CurPage = 0
	this.Page.Offset = 0
	this.Page.Count = 0
}

// PageTime 适用 列表-分页-时间范围
type PageTime struct {
	Page      Page
	StartTime int64 `json:"startTime" form:"startTime"` // 开始时间
	EndTime   int64 `json:"endTime" form:"endTime"`     // 结束时间
}

var PageResult *sync.Pool

func init() {
	PageResult = &sync.Pool{
		New: func() interface{} {
			return &Pagination{}
		},
	}
}

// FormatPage 格式化分页数据-起始量

// FormatTotal 格式化分页数据-总页数
func (p *Page) FormatTotal(total int64) {
	p.SumCount = int(total)
	if p.Count > 0 {
		p.SumPage = int(math.Ceil(float64(p.SumCount) / float64(p.Count)))
	} else {
		p.SumPage = 1
	}
}

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
	clone.Context = ctx
	ctx.ShouldBindBodyWith(clone, binding.JSON)
	if clone.Page.CurPage == 0 {
		clone.Page.CurPage = 1
	}
	if clone.Page.Count == 0 {
		clone.Page.Count = 10
	}
	clone.Page.Offset = (clone.Page.CurPage - 1) * clone.Page.Count
	return clone
}
func (r *Controller) ShouldBindJson(v interface{}) error {
	return r.ShouldBindBodyWith(v, binding.JSON)
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

type StatusError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func Error(err error, code int) {
	if err != nil {
		panic(StatusError{Msg: err.Error(), Code: code})
	}
}

type PageDocList[T any] struct {
	List T `json:"list"`
	Page Page
}
type HttpReturn struct {
	Code    errorCode   `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id"`
}

type HttpDoc[T any] struct {
	Code    errorCode `json:"code"`
	Msg     string    `json:"msg"`
	Data    T         `json:"data"`
	TraceId string    `json:"trace_id" desc:"追踪id"`
}

func PageDoc[T any]() HttpDoc[PageDocList[T]] {
	return HttpDoc[PageDocList[T]]{}
}
