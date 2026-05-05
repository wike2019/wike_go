# 路由级中间件使用指南

本示例演示如何在 Wike Go 框架中使用 **路由级中间件**（Route-Level Middleware）。与全局中间件（`g.GlobalUse()`）不同，路由级中间件通过 `r.Use()` 注册在路由组内部，**仅对该调用之后注册的路由生效**，实现更精细的中间件控制。

## 核心概念

### 全局中间件 vs 路由级中间件

| 特性 | 全局中间件 | 路由级中间件 |
|------|-----------|-------------|
| 注册方式 | `g.GlobalUse(middleware)` | `r.Use(middleware)` |
| 作用范围 | 所有路由 | 仅 `r.Use()` 之后注册的路由 |
| 注册位置 | `main()` 函数中 | `Build()` 方法中 |
| 典型场景 | 日志、Recover、CORS | 超时控制、鉴权、限流 |

### 关键规则：注册顺序决定生效范围

在 `Build()` 方法中，`r.Use()` **只对其后面注册的路由生效**，之前注册的路由不受影响。这是 Gin 框架路由组的核心行为。

```
r.GET("/a", handlerA)       ← 不受中间件影响
r.Use(someMiddleware)       ← 注册路由级中间件
r.GET("/b", handlerB)       ← 受中间件影响
r.GET("/c", handlerC)       ← 受中间件影响
```

## 前置条件

- Go 1.20 及以上版本
- 已安装 Wike Go：`go get github.com/wike2019/wike_go/v2@v2.0.11`

## 项目结构

```
middleware_router_use/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── logs/
│   └── app.log        # 应用日志
├── go.mod             # 依赖管理
├── main.go            # 入口文件（路由级中间件示例）
└── README.md          # 本文档
```

## 使用方法

### 1. 定义路由控制器

路由控制器需要实现三个要素：

- 嵌入 `ctl.Controller` 获得上下文管理能力
- 实现 `Build(r *gin.RouterGroup, GCore *core.GCore)` 方法注册路由和中间件
- 实现 `Path() string` 方法定义路由前缀

```go
type router struct {
    ctl.Controller
}

func NewRouter() *router {
    return &router{}
}
```

### 2. 在 Build 中注册路由级中间件

在 `Build()` 方法中，通过 `r.Use()` 在路由之间插入中间件，控制哪些路由受中间件约束：

```go
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    // 这个路由不受超时中间件影响
    r.GET("/healthz", this.healtzh)

    // 注册路由级中间件：1 秒超时
    r.Use(core.TimeoutMiddleware(time.Second * 1))

    // 这个路由受超时中间件影响
    r.GET("/timeOut", this.healtzh)
}
```

### 3. 定义路由前缀

```go
func (this router) Path() string {
    return "/"
}
```

### 4. 启动服务

```go
func main() {
    g := core.God()
    g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
```

## 完整示例代码

以下是本项目 `main.go` 的完整代码：

```go
package main

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
    ctl.Controller
}

func NewRouter() *router {
    return &router{}
}

func (this *router) healtzh(context *gin.Context) {
    c := this.SetContext(context)
    time.Sleep(time.Second * 2) // 模拟耗时操作（2 秒）
    c.OKWithMsg("hello world")
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    // 此路由无超时限制，正常返回（虽然 handler 会 sleep 2 秒）
    r.GET("/healthz", this.healtzh)

    // 注册路由级超时中间件：1 秒超时
    r.Use(core.TimeoutMiddleware(time.Second * 1))

    // 此路由受超时中间件约束，handler sleep 2 秒 > 超时 1 秒，触发超时
    r.GET("/timeOut", this.healtzh)
}

func (this router) Path() string {
    return "/"
}

func main() {
    g := core.God()
    g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
```

## 运行与验证

### 启动服务

```shell
cd middleware_router_use
go run main.go
```

### 测试路由行为差异

**请求 `/healthz`（无超时中间件）：**

```shell
curl http://localhost:8888/healthz
```

等待约 2 秒后正常返回：

```json
{
  "code": 0,
  "msg": "hello world"
}
```

**请求 `/timeOut`（受 1 秒超时中间件约束）：**

```shell
curl http://localhost:8888/timeOut
```

由于 handler 需要 2 秒，但超时限制为 1 秒，请求会在 1 秒后超时返回错误：

```json
{
  "code": -1,
  "msg": "request timeout"
}
```

### 行为对比

| 路由 | 超时中间件 | handler 耗时 | 结果 |
|------|-----------|-------------|------|
| `GET /healthz` | 无 | 2 秒 | 正常返回 `hello world` |
| `GET /timeOut` | 1 秒 | 2 秒 | 超时返回错误 |

这就是路由级中间件的核心价值：**同一个路由组内，不同路由可以有不同的中间件策略**。

## 配置文件

`config.yaml`：

```yaml
port: 8888
development: true
```

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `port` | 服务监听端口 | `8888` |
| `development` | 开发模式，开启后输出更详细的日志 | `true` |

## 配置中心

`config/config.go` 基于 Viper 实现配置读取：

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

## 进阶用法

### 多个路由级中间件叠加

可以在 `Build()` 中多次调用 `r.Use()` 实现分层中间件：

```go
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    // 公开路由 — 无任何额外中间件
    r.GET("/public", this.publicHandler)

    // 注册鉴权中间件
    r.Use(authMiddleware())

    // 需要登录的路由
    r.GET("/profile", this.profileHandler)
    r.GET("/settings", this.settingsHandler)

    // 再叠加管理员权限中间件
    r.Use(adminMiddleware())

    // 需要管理员权限的路由
    r.GET("/admin/users", this.adminUsersHandler)
}
```

效果：

```
/public          → 无额外中间件
/profile         → authMiddleware
/settings        → authMiddleware
/admin/users     → authMiddleware + adminMiddleware
```

### 结合全局中间件使用

路由级中间件与全局中间件可以同时使用，执行顺序为：

```
全局中间件 → 路由级中间件 → Handler
```

```go
func main() {
    g := core.God()

    // 全局中间件：所有路由生效
    g.GlobalUse(core.AddTrace())
    g.GlobalUse(core.CustomRecover(g))
    g.GlobalUse(core.AccessLog(g))

    // 路由级中间件在各 Router 的 Build() 中按需注册
    g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
```

## 适用场景

| 场景 | 推荐方式 | 说明 |
|------|---------|------|
| 日志、Recover、CORS | 全局中间件 | 所有请求都需要 |
| 接口鉴权 | 路由级中间件 | 部分接口公开，部分需要登录 |
| 超时控制 | 路由级中间件 | 不同接口超时要求不同 |
| 限流 | 路由级中间件 | 不同接口限流策略不同 |
| 权限分级 | 路由级中间件 | 普通用户 / 管理员分层控制 |

## 下一步

- 查看 `middleware_use/` — 了解全局中间件的用法
- 查看 `provide_use/` — 了解依赖注入的用法
