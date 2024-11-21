package core

type API struct {
	ID     uint   `gorm:"primaryKey"` // 主键
	Group  string `gorm:"size:300"`   // 路由分组
	Name   string `gorm:"size:300"`   // 路由名称
	Input  string `gorm:"size:2000"`  // 入参数
	Output string `gorm:"size:2000"`  // 出参数
	Path   string `gorm:"size:300"`   // 路由路径
	Method string `gorm:"size:100"`   // 请求方法
}
