package Init

import (
	Injector "github.com/shenyisyn/goft-ioc"
	"github.com/wike2019/wike_go/src/Core/Sql"
)
// 默认加载 加载配置文件 和 Db解析 用于数据库适配
func init()  {
	Injector.BeanFactory.Set(Sql.NewGPAUtil())
}