# 定时任务 (Cron Task)

本节介绍如何在 Wike Go 框架中使用内置的定时任务功能，包括任务的创建、停止与重启。

## 概述

Wike Go 框架基于 `robfig/cron` 封装了定时任务能力，通过链式调用 `.Cron()` 即可注册定时任务，并支持通过任务名称动态控制任务的启停。

## 项目结构

```
cron_task/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── go.mod             # 依赖管理
├── go.sum             # 依赖校验
└── main.go            # 入口文件（定时任务示例）
```

## 核心 API

### Cron 方法签名

```go
func (g *God) Cron(spec string, cmd func(), name string, immediately bool) *God
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `spec` | `string` | Cron 表达式（支持 6 位，含秒级精度） |
| `cmd` | `func()` | 定时执行的函数 |
| `name` | `string` | 任务名称（用于后续启停控制） |
| `immediately` | `bool` | 是否在注册时立即执行一次 |

### Cron 表达式格式（6 位）

```
┌──────────── 秒 (0-59)
│ ┌────────── 分 (0-59)
│ │ ┌──────── 时 (0-23)
│ │ │ ┌────── 日 (1-31)
│ │ │ │ ┌──── 月 (1-12)
│ │ │ │ │ ┌── 星期 (0-6, 0=周日)
│ │ │ │ │ │
* * * * * *
```

常用表达式示例：

| 表达式 | 说明 |
|--------|------|
| `* * * * * *` | 每秒执行 |
| `0 * * * * *` | 每分钟执行 |
| `0 0 * * * *` | 每小时执行 |
| `0 0 0 * * *` | 每天零点执行 |
| `0 30 9 * * 1-5` | 工作日每天 9:30 执行 |
| `*/5 * * * * *` | 每 5 秒执行 |
| `0 */10 * * * *` | 每 10 分钟执行 |

### 任务控制方法

通过 `GCore` 实例动态控制已注册的定时任务：

```go
// 停止指定名称的定时任务
GCore.StopTask("任务名称")

// 重启指定名称的定时任务
GCore.RestartTask("任务名称")
```

## 完整示例

```go
package main

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ctl"
)

// 定义路由控制器
type router struct {
    ctl.Controller
}

func NewRouter() *router {
    return &router{}
}

func (this *router) healtzh(context *gin.Context) {
    c := this.SetContext(context)
    c.OKWithMsg("hello world")
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/healthz", this.healtzh)
}

func (this router) Path() string {
    return "/"
}

func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewRouter).
        Cron("* * * * * *", func() {
            fmt.Println(time.Now())
        }, "定时任务一", true).
        Invokes(func(GCore *core.GCore) {
            go func() {
                time.Sleep(5 * time.Second)
                GCore.StopTask("定时任务一")
                fmt.Println("已经停止，30秒后再开启")
                time.Sleep(30 * time.Second)
                GCore.RestartTask("定时任务一")
            }()
        }).
        Run()
}
```

## 启动流程解析

```
core.God() → Provide(配置) → MountWithEmpty(路由) → Cron(表达式, 函数, 名称, 立即执行) → Invokes(回调) → Run()
```

| 方法 | 作用 |
|------|------|
| `core.God()` | 创建框架实例 |
| `Provide(fn...)` | 注册配置到 DI 容器 |
| `MountWithEmpty(ctor)` | 挂载路由控制器 |
| `Cron(spec, cmd, name, immediately)` | 注册定时任务 |
| `Invokes(fn)` | 注册启动后回调（可访问 GCore） |
| `Run()` | 启动 HTTP 服务和定时任务 |

## 使用场景 Demo

### 场景一：每分钟清理过期缓存

```go
g.Provide(config.Config).
    MountWithEmpty(NewRouter).
    Cron("0 * * * * *", func() {
        fmt.Println("清理过期缓存...")
        // 执行缓存清理逻辑
    }, "缓存清理", false).
    Run()
```

### 场景二：多个定时任务并存

```go
g.Provide(config.Config).
    MountWithEmpty(NewRouter).
    Cron("*/10 * * * * *", func() {
        fmt.Println("每10秒：检查队列消息")
    }, "队列检查", false).
    Cron("0 0 * * * *", func() {
        fmt.Println("每小时：同步数据")
    }, "数据同步", false).
    Cron("0 0 2 * * *", func() {
        fmt.Println("每天凌晨2点：生成报表")
    }, "报表生成", false).
    Run()
```

### 场景三：根据业务条件动态启停任务

```go
g.Provide(config.Config).
    MountWithEmpty(NewRouter).
    Cron("*/5 * * * * *", func() {
        fmt.Println("执行数据采集...")
    }, "数据采集", true).
    Invokes(func(GCore *core.GCore) {
        // 在路由中通过 GCore 控制任务
        // 例如：提供 HTTP 接口来启停任务
    }).
    Run()
```

在路由控制器中控制任务：

```go
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/healthz", this.healtzh)

    // 停止定时任务
    r.POST("/task/stop", func(c *gin.Context) {
        GCore.StopTask("数据采集")
        c.JSON(200, gin.H{"msg": "任务已停止"})
    })

    // 重启定时任务
    r.POST("/task/restart", func(c *gin.Context) {
        GCore.RestartTask("数据采集")
        c.JSON(200, gin.H{"msg": "任务已重启"})
    })
}
```

## 配置文件

`config.yaml`：

```yaml
port: 8888
development: true
```

## 运行

```shell
go run main.go
```

启动后：
- HTTP 服务监听 `http://localhost:8888`
- 定时任务按照 Cron 表达式自动执行
- 可通过 `GCore.StopTask()` / `GCore.RestartTask()` 动态控制任务

## 依赖

```
github.com/wike2019/wike_go/v2 v2.0.16
github.com/gin-gonic/gin v1.12.0
github.com/spf13/viper v1.21.0
```
