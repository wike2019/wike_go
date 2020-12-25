package Config


//系统配置 默必须配置端口 mysql地址 redis地址
import (
"gopkg.in/yaml.v2"
	"log"
)

//配置解析

type ServerConfig struct {
	Port int32
	Name string
	Mysql string
	Redis string
}

//系统配置
type SysConfig struct {
	Server *ServerConfig
}

func (this *SysConfig) Name() string {
	return "SysConfig"
}

func NewSysConfig() *SysConfig {
	return &SysConfig{Server: &ServerConfig{Port: 8080, Name: "wikeGo微服务框架",Mysql:"root:root@tcp(192.168.3.2:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",Redis:"192.168.3.2:6379"}}
}
func InitConfig() *SysConfig {
	config := NewSysConfig()
	if b := LoadConfigFile(); b != nil {
		err := yaml.Unmarshal(b, config)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config

}

