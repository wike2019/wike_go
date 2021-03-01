package Bean

import (
	Injector "github.com/shenyisyn/goft-ioc"
	_ "github.com/wike2019/wike_go/src/Core/Init"
	"reflect"
	"sync"
)

//BeanFactory 容器 用于依赖注入
var Instance *BeanFactory //单例
var once2 sync.Once       //锁


type  BeanFactory struct {
	ExprData map[string]interface{} //保存对象 给表达式使用
}


func ( this *BeanFactory) ApplyAll() { //循环依赖注入
	for t, v := range Injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			Injector.BeanFactory.Apply(v.Interface())
		}
	}
}

func New() *BeanFactory {
	once2.Do(func() {
		Instance = &BeanFactory{ExprData: make(map[string]interface{})}
	})
	return Instance
}


func (this *BeanFactory) Beans(beans ...Bean) *BeanFactory { //注入对象
	for _, bean := range beans {
		this.ExprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return this
}
func (this *BeanFactory) Config(cfgs ...interface{}) *BeanFactory { //根据方法注入反射
	Injector.BeanFactory.Config(cfgs...)
	return this
}

func (this *BeanFactory) Get(name interface{}) interface{} { //根据方法注入反射
	return  Injector.BeanFactory.Get(name)

}