package controller

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

// 时间戳转换成日期
func ToDate(timestamp int64) string {
	// 将时间戳转换为 time.Time 对象
	t := time.Unix(timestamp, 0) // 第二个参数是纳秒，这里设置为 0

	// 将时间格式化为字符串
	// 格式化字符串 "2006-01-02 15:04:05" 定义了输出格式（这里是年-月-日 时:分:秒）
	// 这个特殊的日期 "2006-01-02 15:04:05" 在 Go 中用于表示格式
	return t.Format("2006-01-02 15:04:05")
}

type Page struct {
	SumPage  int `json:"sumPage" swaggerignore:"true"`  // 总页数
	SumCount int `json:"sumCount" swaggerignore:"true"` // 总条数
	CurPage  int `json:"page" form:"page"`              // 当前页
	Offset   int `json:"-"`                             // 起始量
	Count    int `json:"count" form:"count"`            // 每页返回数量
}

type Pagination struct {
	List interface{} `json:"list"`
	Page
}

// PageTime 适用 列表-分页-时间范围
type PageTime struct {
	Page
	StartTime int64 `json:"startTime" form:"startTime"` // 开始时间
	EndTime   int64 `json:"endTime" form:"endTime"`     // 结束时间
}

var PageResult *sync.Pool
var ErrorResult *sync.Pool

func init() {
	PageResult = &sync.Pool{
		New: func() interface{} {
			return &Pagination{}
		},
	}
	ErrorResult = &sync.Pool{
		New: func() interface{} {
			return &StatusError{}
		},
	}
}

// FormatPage 格式化分页数据-起始量
func (r *Page) FormatPage(ctx *gin.Context) {
	ctx.ShouldBind(r)
	if r.CurPage == 0 {
		r.CurPage = 1
	}
	if r.Count == 0 {
		r.Count = 10
	}
	r.Offset = (r.CurPage - 1) * r.Count
}

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
}

func (*Controller) Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":     Success,
		"msg":      msg,
		"data":     data,
		"trace_id": ctx.GetString("trace_id"),
	})
}
func (r *Controller) List(ctx *gin.Context, msg string, list interface{}) {
	if list == nil {
		list = make([]interface{}, 0)
	}
	result := PageResult.Get().(*Pagination)
	defer PageResult.Put(result)
	result.Page = r.Page
	result.List = list
	r.Success(ctx, msg, result)
}

func (*Controller) Failed(ctx *gin.Context, code errorCode, msg string) {
	errMsg := codeMsg[code] + ": " + msg
	if code != Success {
		ctx.Set("error_code", int(code))
		ctx.Set("error_msg", msg)
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":     code,
		"msg":      errMsg,
		"data":     nil,
		"trace_id": ctx.GetString("trace_id"),
	})
}

type StatusError struct {
	Code int
	Msg  string
}

func (this *StatusError) Error() string {
	return this.Msg
}

func Error(code int, msg string) {
	result := ErrorResult.Get().(*StatusError)
	defer ErrorResult.Put(result)
	result.Msg = msg
	result.Code = code
	panic(result)
}
