package fxTags

import (
	"fmt"

	"go.uber.org/fx"
)

func Create(fn interface{}, resultTag string, param []string) interface{} {
	if len(param) > 0 {
		return fx.Annotate(fn, fx.ParamTags(param...), fx.ResultTags(resultTag))
	} else {
		return fx.Annotate(fn, fx.ResultTags(resultTag))
	}
}
func CreateInterFace(fn interface{}, I interface{}, resultTag string, param []string) interface{} {
	if len(param) > 0 {
		return fx.Annotate(fn, fx.As(I), fx.ParamTags(param...), fx.ResultTags(resultTag))
	} else {
		return fx.Annotate(fn, fx.As(I), fx.ResultTags(resultTag))
	}
}
func ParamList(param ...string) []string {
	return param
}
func CreateTag(tag string) string {
	return fmt.Sprintf(` name:"%s" `, tag)
}
func CreateGroup(tag string) string {
	return fmt.Sprintf(` group:"%s" `, tag)
}
