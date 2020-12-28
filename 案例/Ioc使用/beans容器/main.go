package main

import (
	"github.com/shenyisyn/goft-expr/src/expr"
	"github.com/wike2019/wike_go/src/core/Ioc"
)

func main()  {

	task:=&Task{}
	task.Show(10)//执行一次
	expr.BeanExpr("Task.Show(20)", Ioc.New().ExprData)//使用表达式执行一次

}
