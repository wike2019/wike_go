# 中间件使用指南

本节演示如何在 Wike Go 框架中使用全局中间件。框架通过 `g.GlobalUse()` 方法注册中间件，所有中间件按注册顺序依次执行。

## 前置条件

- Go 1.20 及以上版本
- 已安装 Wike Go：`go get github.com/wike2019/wike_go/v2@v2.0.11`

## 项目结构

```
middleware_use/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── go.mod             # 依赖管理
└── main.go            # 入口文件（中间件注册示例）
```

## 中间件注册方式

使用 `g.GlobalUse()` 注册全局中间件，中间件会作用于所有路由：

```go
g := core.God()
g.GlobalUse(中间件函数)
```

多个中间件按注册顺序依次执行，建议按照以下顺序注册：

```
TraceId → Recover → 优雅关闭 → 日志 → Body限制 → 跨域 → 超时
```

## 内置中间件一览

| 中间件 | 方法 | 说明 |
|--------|------|------|
| TraceId | `core.AddTrace()` | 为每个请求生成唯一追踪 ID |
| Recover | `core.CustomRecover(g)` | 捕获 panic，防止服务崩溃 |
| 优雅关闭 | `core.Reject(g)` | 收到关闭信号后拒绝新请求 |
| 访问日志 | `core.AccessLog(g)` | 记录请求的方法、路径、耗时等信息 |
| Body 大小限制 | `core.LimitBodySize(size)` | 限制请求体大小，防止大文件攻击 |
| 跨域 CORS | `core.CORSMiddleware()` | 允许跨域请求 |
| 超时控制 | `core.TimeoutMiddleware(d)` | 请求超时自动取消 |

---

## 各中间件详解

### 1. TraceId 中间件

为每个请求自动生成唯一的 `X-Trace-Id`，贯穿整个请求链路，方便日志追踪和问题排查。

```go
g.GlobalUse(core.AddTrace())
```

**使用场景：**
- 在日志中通过 TraceId 关联同一请求的所有日志
- 微服务间传递 TraceId 实现分布式链路追踪
- 排查线上问题时，通过 TraceId 快速定位完整调用链

### 2. Recover 中间件

捕获处理函数中的 `panic`，返回 500 错误响应而不是让整个服务崩溃。需要传入框架实例 `g` 以便记录错误日志。

```go
g.GlobalUse(core.CustomRecover(g))
```

**行为：**
- 捕获 panic 后记录错误堆栈到日志
- 向客户端返回统一的 500 错误响应
- 服务继续运行，不会因单个请求的 panic 而宕机

### 3. 优雅关闭中间件

当服务收到关闭信号（如 `SIGTERM`）时，拒绝新的请求进入，同时等待已有请求处理完毕后再关闭。

```go
g.GlobalUse(core.Reject(g))
```

**行为：**
- 正常运行时：请求正常通过
- 收到关闭信号后：新请求返回 503 Service Unavailable
- 已在处理中的请求会继续执行直到完成

### 4. 访问日志中间件

记录每个请求的详细信息，包括请求方法、路径、状态码、耗时等，便于监控和分析。

```go
g.GlobalUse(core.AccessLog(g))
```

**记录内容：**
- 请求方法（GET、POST 等）
- 请求路径和查询参数
- 响应状态码
- 请求处理耗时
- 客户端 IP 地址

### 5. Body 大小限制中间件

限制请求体的最大大小，防止恶意用户发送超大请求体导致内存溢出。

```go
// 限制为 32MB（32 << 20 = 33554432 字节）
g.GlobalUse(core.LimitBodySize(32 << 20))
```

**常用大小设置：**

| 表达式 | 大小 | 适用场景 |
|--------|------|----------|
| `1 << 20` | 1 MB | 纯文本 API |
| `10 << 20` | 10 MB | 包含图片上传 |
| `32 << 20` | 32 MB | 通用场景（本示例） |
| `100 << 20` | 100 MB | 大文件上传 |

### 6. 跨域 CORS 中间件

允许浏览器跨域请求，解决前后端分离开发中的跨域问题。

```go
g.GlobalUse(core.CORSMiddleware())
```

**行为：**
- 自动处理 `OPTIONS` 预检请求
- 设置 `Access-Control-Allow-Origin` 等响应头
- 允许常见的请求方法和请求头

### 7. 超时控制中间件

为每个请求设置最大处理时间，超时后自动取消请求上下文，防止慢请求占用资源。

```go
import "time"

// 设置 10 秒超时
g.GlobalUse(core.TimeoutMiddleware(time.Second * 10))
```

**常用超时设置：**

| 时长 | 表达式 | 适用场景 |
|------|--------|----------|
| 3 秒 | `time.Second * 3` | 简单查询接口 |
| 10 秒 | `time.Second * 10` | 通用接口（本示例） |
| 30 秒 | `time.Second * 30` | 复杂计算或外部调用 |
| 60 秒 | `time.Minute` | 文件上传等耗时操作 |

---

## 完整示例

以下是本项目 `main.go` 的完整代码，展示了所有中间件的注册方式：

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

    // 1. TraceId — 请求链路追踪
    g.GlobalUse(core.AddTrace())
    // 2. Recover — 捕获 panic，防止服务崩溃
    g.GlobalUse(core.CustomRecover(g))
    // 3. 优雅关闭 — 关闭时拒绝新请求
    g.GlobalUse(core.Reject(g))
    // 4. 访问日志 — 记录请求信息
    g.GlobalUse(core.AccessLog(g))
    // 5. Body 大小限制 — 防止超大请求
    g.GlobalUse(core.LimitBodySize(32 << 20))
    // 6. 跨域 — 允许前端跨域访问
    g.GlobalUse(core.CORSMiddleware())
    // 7. 超时控制 — 10 秒超时
    g.GlobalUse(core.TimeoutMiddleware(time.Second * 10))

    g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
```

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

## 运行

```shell
go run main.go
```

访问 `http://localhost:8888/healthz`，返回：

```json
{
  "code": 0,
  "msg": "hello world"
}
```

## 按需选用中间件

并非所有中间件都是必须的，可以根据项目需求选择性注册：

**推荐始终启用：**
- `CustomRecover` — 防止 panic 导致服务崩溃
- `AccessLog` — 生产环境必备的请求日志

**按需启用：**
- `AddTrace` — 需要链路追踪时启用
- `Reject` — 需要优雅关闭时启用
- `LimitBodySize` — 有文件上传或大请求体时启用
- `CORSMiddleware` — 前后端分离部署时启用
- `TimeoutMiddleware` — 需要防止慢请求时启用

## 下一步

- [初始化项目](/guide/初始化项目.md) — 了解框架最小启动流程
- [控制器](/guide/控制器.md) — 深入学习路由控制器的用法
- [依赖注入](/guide/依赖注入.md) — 使用 Mount 注入数据库、缓存等依赖
