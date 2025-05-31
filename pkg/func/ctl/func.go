package ctl

import (
	"fmt"
	"reflect"
)

func (r *Controller) getFuncMateInfo(input interface{}, output interface{}) {
	// 获取函数的类型信息
	funcVal := reflect.ValueOf(output)
	funcType := funcVal.Type()

	funcInVal := reflect.ValueOf(input)
	funcInType := funcInVal.Type()

	// 确保传入的确实是一个函数
	if funcType.Kind() != reflect.Func {
		panic("类型错误")
	}
	if funcInType.Kind() != reflect.Func {
		panic("类型错误")
	}

	// 获取函数的参数类型和返回值类型
	if funcType.NumOut() != 2 {
		panic(fmt.Sprintf("%s-%d", "数据不合法", funcType.NumOut()))
	}
	if funcInType.NumOut() != 4 {
		panic(fmt.Sprintf("%s-%d", "数据不合法", funcInType.NumOut()))
	}
	// 获取参数类型（假设为结构体A）
	paramType := funcInType.Out(0).Elem()
	bodyType := funcInType.Out(1).Elem()
	headerType := funcInType.Out(2).Elem()
	// 获取返回值类型（假设为结构体B）
	returnType := funcType.Out(0).Elem()

	// 使用反射创建参数类型和返回值类型的实例
	paramInstance := reflect.New(paramType)
	bodyInstance := reflect.New(bodyType)
	headerInstance := reflect.New(headerType)
	returnInstance := reflect.New(returnType)
	r.query = paramInstance.Interface()
	r.body = bodyInstance.Interface()
	r.header = headerInstance.Interface()
	r.output = returnInstance.Interface()
}
