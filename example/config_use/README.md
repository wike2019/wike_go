# config_use 使用指南

基于 `wike_go/v2` 框架的完整示例，展示如何通过**依赖注入**组织配置、数据库和路由。

## 目录结构

```
config_use/
├── config.yaml          # 配置文件
├── config/
│   └── config.go        # 配置中心（viper 读取 yaml）
├── mysql/
│   └── mysql.go         # MySQL 初始化（从 viper 读取 DSN）
├── redis/
│   └── redis.go         # Redis 初始化（从 viper 读取地址）
├── main.go              # 入口：依赖注入 + 路由挂载 + 启动
└── go.mod
```

## 配置文件

### config.yaml

```yaml
port: 8888
development: true
mysql: "root:root@tcp(192.168.3.2:3310)/test?charset=utf8mb4&parseTime=True&loc=Local"
redis: "192.168.3.2:6379"
logPath: "./logs/app.log"
```

## 各模块代码

### config/config.go — 配置中心

```go
package config

import (
	"log"

	"github.com/spf13/viper"
)

func Config() *viper.Viper {
	viper.SetDefault("port", "8888")
	viper.SetDefault("logPath", "./logs/app.log")
	viper.SetDefault("development", true)
	viper.SetConfigFile("config.yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err.Error())
	}
	return viper.GetViper()
}
```

### mysql/mysql.go — MySQL 初始化

```go
package mysql

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDB struct {
	*gorm.DB
}

func InitMysql(cfg *viper.Viper) *MysqlDB {
	dsn := cfg.GetString("mysql")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql: " + err.Error())
	}
	return &MysqlDB{DB: db}
}
```

### redis/redis.go — Redis 初始化

```go
package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisDB struct {
	*redis.Client
}

func InitRedis(cfg *viper.Viper) *RedisDB {
	addr := cfg.GetString("redis")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to ping redis: %s", err.Error())
	}
	log.Println("redis connected:", addr)
	return &RedisDB{Client: client}
}
```

### main.go — 入口

```go
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_app/mysql"
	"github.com/wike2019/wike_app/redis"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
	ctl.Controller
	db    *mysql.MysqlDB
	redis *redis.RedisDB
}

func NewRouter(db *mysql.MysqlDB, redis *redis.RedisDB) *router {
	return &router{db: db, redis: redis}
}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.OKWithMsg("hello world")
}

func (this *router) doSomething(context *gin.Context) {
	c := this.SetContext(context)
	// 这里可以操作 mysql 和 redis
	fmt.Println(this.db)
	fmt.Println(this.redis)
	c.OKWithMsg("hello world")
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/healthz", this.healtzh)
	r.GET("/doSomething", this.doSomething)
}

func (this router) Path() string {
	return "/"
}

func main() {
	g := core.God()
	g.Provide(config.Config, mysql.InitMysql, redis.InitRedis).
		MountWithEmpty(NewRouter).
		Run()
}
```

### go.mod

```
module github.com/wike2019/wike_app

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/redis/go-redis/v9 v9.19.0
	github.com/spf13/viper v1.21.0
	github.com/wike2019/wike_go/v2 v2.0.11
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)
```

## 运行方式

```bash
# 初始化依赖
go mod tidy

# 确保 config.yaml 中的 mysql 和 redis 地址可达，然后启动
go run main.go
```

启动后访问：

- `GET http://localhost:8888/healthz` — 健康检查
- `GET http://localhost:8888/doSomething` — 业务接口（可操作 MySQL 和 Redis）

## 框架运行机制

```
core.God()
  │
  ├── Provide(config.Config, mysql.InitMysql, redis.InitRedis)
  │     │
  │     │  框架内部（uber/dig）自动解析依赖链：
  │     │  1. config.Config() → 返回 *viper.Viper
  │     │  2. mysql.InitMysql(cfg *viper.Viper) → 返回 *MysqlDB
  │     │  3. redis.InitRedis(cfg *viper.Viper) → 返回 *RedisDB
  │     │
  ├── MountWithEmpty(NewRouter)
  │     │
  │     │  NewRouter(db *MysqlDB, redis *RedisDB) → 自动注入上面的实例
  │     │  调用 Build() 注册路由，Path() 指定路由前缀
  │     │
  └── Run()
        │
        启动 Gin HTTP 服务，监听 config.yaml 中的 port
```

## 关键概念

| 概念 | 说明 |
|------|------|
| **依赖注入** | `Provide` 注册的函数，参数由框架自动注入（底层 `uber/fx` + `uber/dig`） |
| **路由接口** | 实现 `Build` 和 `Path` 两个方法即可挂载路由组 |
| **Controller 基类** | 嵌入 `ctl.Controller`，通过 `SetContext` 获得封装后的上下文，提供 `OKWithMsg` 等便捷响应方法 |
| **配置读取** | 任何 `Provide` 的初始化函数都可以接收 `*viper.Viper` 参数来读取配置 |
