package Web

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Core/Sql"
	"reflect"
	"sync"
)
// 定义返回解析器
var responderList []Responder
var once_resp_list sync.Once

func get_responder_list() []Responder {
	once_resp_list.Do(func() {
		responderList = []Responder{(StringResponder)(nil),
			(JsonResponder)(nil),
			(ViewResponder)(nil),
			(SqlResponder)(nil),
			(SqlQueryResponder)(nil),
		}
	})
	return responderList
}

type Responder interface {
	RespondTo() gin.HandlerFunc
}

func Convert(handler interface{}) gin.HandlerFunc {
	h_ref := reflect.ValueOf(handler)
	for _, r := range get_responder_list() {
		r_ref := reflect.TypeOf(r)
		if h_ref.Type().ConvertibleTo(r_ref) {
			return h_ref.Convert(r_ref).Interface().(Responder).RespondTo()
		}
	}
	return nil
}

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, getFairingHandler().handlerFairing(this, context).(string))
	}
}

type Json interface{}
type JsonResponder func(*gin.Context) Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, getFairingHandler().handlerFairing(this, context))
	}
}

type SqlQueryResponder func(*gin.Context) Sql.Query

func (this SqlQueryResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		getQuery := getFairingHandler().handlerFairing(this, context).(Sql.Query)
		ret, err := Sql.QueryForMapsByInterface(getQuery)
		if err != nil {
			panic(err)
		}
		context.JSON(200, ret)
	}
}

type SqlResponder func(*gin.Context) Sql.SimpleQuery

func (this SqlResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		getSql := getFairingHandler().handlerFairing(this, context).(Sql.SimpleQuery)
		ret, err := Sql.QueryForMaps(string(getSql), nil, []interface{}{}...)
		if err != nil {
			panic(err)
		}
		context.JSON(200, ret)
	}
}

// Deprecated: 暂时不提供View的解析
type View string
type ViewResponder func(*gin.Context) View

func (this ViewResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.HTML(200, string(this(context))+".html", context.Keys)
	}
}
