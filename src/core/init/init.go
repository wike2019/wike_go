package init

import (
	Injector "github.com/shenyisyn/goft-ioc"
	"github.com/wike2019/wike_go/src/core/config"
	"github.com/wike2019/wike_go/src/core/sql"
)
// 默认加载 加载配置文件 和 Db解析 用于数据库适配
func init()  {
	initConfig := Config.InitConfig()
	Injector.BeanFactory.Set(initConfig) // add global into (new)BeanFactory
	Injector.BeanFactory.Set(sql.NewGPAUtil())
}