package ioc

// 接口所有可注入对象 Name 返回的的名字必须为类名 且类名不能重复
type Bean interface {
	Name() string
}
