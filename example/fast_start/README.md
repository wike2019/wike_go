# 初始化项目

本节将从零开始创建一个 Wike Go 项目，帮助你理解框架的最小启动流程。

## 前置条件

- Go 1.20 及以上版本
- 已安装 Wike Go：`go get github.com/wike2019/wike_go/v2@v2.0.11`

## 项目结构

一个最小的 Wike Go 项目只需要以下文件：

```
my_project/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── go.mod             # 依赖管理
└── main.go            # 入口文件
```

::: tip 说明
框架遵循**约定优于配置**的理念，所有参数都有合理的默认值。最小启动只需要一个配置文件和一个入口文件。
:::

## 第一步：初始化 Go 模块

```shell
mkdir my_project && cd my_project
go mod init my_project
go get github.com/wike2019/wike_go/v2@v2.0.11
```

## 第二步：编写配置文件

创建 `config.yaml`，这是框架启动所需的最基础配置：

```yaml
port: 8888
development: true
```

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `port` | 服务监听端口 | `8888` |
| `development` | 开发模式，开启后会输出更详细的日志 | `true` |

## 第三步：编写配置中心

创建 `config/config.go`，使用 Viper 加载配置文件：

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

::: info 为什么需要配置中心？
配置中心通过 `Provide` 注入到框架的依赖注入容器中，框架内部的日志、HTTP 服务、中间件等模块都会从中读取配置。`SetDefault` 确保即使配置文件缺少某些字段，程序也能正常运行。
:::

## 第四步：编写入口文件

创建 `main.go`，定义一个路由控制器并启动服务：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ctl"
    "my_project/config"
)

// 定义控制器，嵌入 ctl.Controller 获得统一的响应方法
type HelloRouter struct {
    ctl.Controller
}

func NewHelloRouter() *HelloRouter {
    return &HelloRouter{}
}

// Path 返回路由前缀
func (h HelloRouter) Path() string {
    return "/"
}

// Build 注册具体路由
func (h HelloRouter) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/healthz", h.healthz)
}

// 处理函数
func (h *HelloRouter) healthz(ctx *gin.Context) {
    c := h.SetContext(ctx)
    c.OKWithMsg("hello world")
}

func main() {
    g := core.God()
    g.Provide(config.Config).MountWithEmpty(NewHelloRouter).Run()
}
```

### 核心启动流程

整个启动过程只有一行链式调用：

```
core.God() → Provide(配置) → MountWithEmpty(路由) → Run()
```

| 方法 | 作用 |
|------|------|
| `core.God()` | 创建框架实例 |
| `Provide(fn...)` | 注册依赖到 DI 容器（配置、数据库等） |
| `MountWithEmpty(ctor)` | 挂载无外部依赖的路由控制器 |
| `Run()` | 启动 HTTP 服务 |

### 控制器接口

每个路由控制器需要实现两个方法：

```go
type Controller interface {
    Build(r *gin.RouterGroup, GCore *core.GCore)  // 注册路由
    Path() string                                  // 路由前缀
}
```

通过嵌入 `ctl.Controller`，你可以直接使用以下响应方法：

| 方法 | 说明 |
|------|------|
| `SetContext(ctx)` | 初始化上下文，每个处理函数的第一行调用 |
| `OK(data)` | 返回成功响应 + 数据 |
| `OKWithMsg(msg)` | 返回成功响应 + 消息 |
| `OKWithList(list)` | 返回分页列表响应 |
| `Failed(code, msg)` | 返回错误响应 |

## 第五步：运行

```shell
go run main.go
```

浏览器访问 `http://localhost:8888/healthz`，看到如下输出即表示项目初始化成功：

```json
{
  "code": 0,
  "msg": "hello world"
}
```

## 挂载多个控制器

实际项目中通常有多个路由模块，可以通过链式调用挂载：

```go
func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewUserRouter).
        MountWithEmpty(NewOrderRouter).
        MountWithEmpty(NewProductRouter).
        Run()
}
```

::: tip 提示
当控制器需要依赖注入（如数据库、缓存等）时，使用 `Mount(ctor, params)` 替代 `MountWithEmpty`。详见 [依赖注入](/guide/依赖注入.md) 章节。
:::

## 下一步

- [使用自定义配置文件](/guide/使用自定义配置文件.md) — 了解更多配置选项
- [控制器](/guide/控制器.md) — 深入学习路由控制器的用法
- [依赖注入](/guide/依赖注入.md) — 使用 Mount 注入数据库、缓存等依赖
