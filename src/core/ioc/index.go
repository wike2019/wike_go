package ioc

import (
	Injector "github.com/shenyisyn/goft-ioc"
	_ "github.com/wike2019/wike_go/src/core/init"
	"reflect"
	"sync"
)
//Ioc 容器 用于依赖注入
var Instance *Ioc   //单例
var once2 sync.Once //锁


type  Ioc struct {
	ExprData map[string]interface{} //保存对象 给表达式使用
}


func ( this *Ioc) ApplyAll() { //循环依赖注入
	for t, v := range Injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			Injector.BeanFactory.Apply(v.Interface())
		}
	}
}

func New() *Ioc {
	once2.Do(func() {
		Instance = &Ioc{ExprData: make(map[string]interface{})}
	})
	return Instance
}


func (this *Ioc) Beans(beans ...Bean) *Ioc { //注入对象
	for _, bean := range beans {
		this.ExprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return this
}
func (this *Ioc) Config(cfgs ...interface{}) *Ioc { //根据方法注入反射
	Injector.BeanFactory.Config(cfgs...)
	return this
}
