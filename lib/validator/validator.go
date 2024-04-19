package validator

//gin > 1.4.0

//将验证器错误翻译成中文

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	//注册翻译器
	zh := zh.New()
	uni = ut.New(zh, zh)

	trans, _ = uni.GetTranslator("zh")

	//获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)
	//注册翻译器
	zh_translations.RegisterDefaultTranslations(validate, trans)
}

// Translate 翻译错误信息
func Translate(err error, v interface{}, tagType string) map[string][]string {

	var result = make(map[string][]string)

	errors := err.(validator.ValidationErrors)

	for _, err := range errors {
		// 获取字段的 reflect.Value
		field, _ := reflect.TypeOf(v).Elem().FieldByName(err.StructField())
		// 从 reflect.Value 中获取 json 标签
		jsonTag := field.Tag.Get(tagType)
		result[jsonTag] = append(result[jsonTag], strings.Replace(err.Translate(trans), err.Field(), jsonTag, -1))
	}
	return result
}
func TranslateJson(err error, v interface{}) map[string][]string {
	return Translate(err, v, "json")
}
func TranslateForm(err error, v interface{}) map[string][]string {
	return Translate(err, v, "form")
}
