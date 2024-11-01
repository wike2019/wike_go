package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/pkg/core"
	"log"
)

// 配置中心
func MyViper(core *core.GCore) *viper.Viper {
	//fmt.Println(core)
	viper.SetDefault("port", "8888")
	viper.SetDefault("logPath", "./logs/app.log")
	viper.SetDefault("development", true)
	viper.SetConfigFile("config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")      // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")           // 还可以在工作目录中查找配置
	err := viper.ReadInConfig()        // 查找并读取配置文件
	if err != nil {                    // 处理读取配置文件的错误
		log.Fatalf("Fatal error config file: %s \n", err.Error())
	}
	//业务参数
	viper.SetDefault("Timeout", 3000)
	core.Stop(func() error {
		fmt.Println("这里是viper清理")
		return nil
	})
	return viper.GetViper()
}
