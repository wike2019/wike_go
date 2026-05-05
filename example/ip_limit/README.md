# IP 限流中间件使用指南

本示例演示如何在 Wike Go 框架中使用 **基于 IP 的令牌桶限流中间件**（Rate Limiter）。框架通过 `core.RateLimiter()` 提供开箱即用的限流能力，基于客户端 IP 地址进行独立限流，防止单个 IP 的高频请求压垮服务。

## 核心概念

### 令牌桶算法

令牌桶（Token Bucket）是一种经典的限流算法：

```
┌─────────────────────────────┐
│         令牌桶               │
│                             │
│   以固定速率 (rate) 生成令牌   │
│   桶容量上限为 capacity       │
│                             │
│   ● ● ● ● ○ ○ ○ ○          │
│   (已有令牌)  (空位)          │
│                             │
└─────────────────────────────┘
         │
         ▼
    请求到达时取走一个令牌
    有令牌 → 放行
    无令牌 → 拒绝 (429 Too Many Requests)
```

**两个关键参数：**

| 参数 | 含义 | 作用 |
|------|------|------|
| `rate`（速率） | 每秒生成令牌的数量 | 控制持续请求的平均速率 |
| `capacity`（容量） | 桶中最多存放的令牌数 | 允许的突发请求量上限 |

### 基于 IP 的独立限流

框架为每个客户端 IP 维护独立的令牌桶，互不影响：

```
客户端 A (192.168.1.1) → 令牌桶 A [rate=1, cap=2]
客户端 B (192.168.1.2) → 令牌桶 B [rate=1, cap=2]
客户端 C (10.0.0.1)    → 令牌桶 C [rate=1, cap=2]
```

## 前置条件

- Go 1.20 及以上版本
- 已安装 Wike Go：`go get github.com/wike2019/wike_go/v2@v2.0.11`

## 项目结构

```
ip_limit/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── go.mod             # 依赖管理
├── main.go            # 入口文件（IP 限流示例）
└── README.md          # 本文档
```

## 使用方法

### 1. 通过依赖注入注册限流器

限流器通过框架的依赖注入系统（基于 `uber/fx`）进行管理。使用 `fxTags` 为限流器实例打标签，支持多个不同配置的限流器共存。

```go
import (
    "github.com/wike2019/wike_go/v2/func/rateLimiter"
    "github.com/wike2019/wike_go/v2/fxTags"
)

// 注册限流器到依赖注入容器，标签为 "limit"
g.Provide(fxTags.Create(rateLimiter.RateLimit, fxTags.CreateTag("limit"), nil))
```

**参数说明：**

| 参数 | 值 | 说明 |
|------|-----|------|
| 第 1 个参数 | `rateLimiter.RateLimit` | 限流器的构造函数 |
| 第 2 个参数 | `fxTags.CreateTag("limit")` | 为该实例设置依赖注入标签 |
| 第 3 个参数 | `nil` | 额外配置（此处无需额外配置） |

### 2. 在路由控制器中注入限流器

路由控制器通过构造函数接收限流器实例：

```go
type router struct {
    ctl.Controller
    GetLimit rateLimiter.GetLimit  // 限流器获取函数
}

func NewRouter(GetLimit rateLimiter.GetLimit) *router {
    return &router{GetLimit: GetLimit}
}
```

### 3. 在路由中使用限流中间件

在 `Build()` 方法中通过 `core.RateLimiter()` 注册限流中间件：

```go
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.Use(core.RateLimiter(this.GetLimit, 1, 2))
    r.GET("/healthz", this.healtzh)
}
```

**`core.RateLimiter()` 参数说明：**

| 参数位置 | 类型 | 说明 |
|---------|------|------|
| 第 1 个 | `rateLimiter.GetLimit` | 限流器实例（通过依赖注入获取） |
| 第 2 个 | `int` | 令牌生成速率（每秒生成的令牌数） |
| 第 3 个 | `int` | 令牌桶容量（桶中最多存放的令牌数） |

### 4. 挂载路由并启动

使用 `Mount` 挂载路由时，通过 `fxTags.ParamList` 指定依赖注入的标签匹配：

```go
func main() {
    g := core.God()
    // 注册限流器
    g.Provide(fxTags.Create(rateLimiter.RateLimit, fxTags.CreateTag("limit"), nil))
    // 挂载路由，注入标签为 "limit" 的限流器
    g.Provide(config.Config).Mount(NewRouter, fxTags.ParamList(fxTags.CreateTag("limit"))).Run()
}
```

> **注意：** 这里使用 `Mount` 而非 `MountWithEmpty`，因为 `NewRouter` 需要接收限流器参数。

## 完整示例代码

以下是本项目 `main.go` 的完整代码：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ctl"
    "github.com/wike2019/wike_go/v2/func/rateLimiter"
    "github.com/wike2019/wike_go/v2/fxTags"
)

type router struct {
    ctl.Controller
    GetLimit rateLimiter.GetLimit
}

func NewRouter(GetLimit rateLimiter.GetLimit) *router {
    return &router{GetLimit: GetLimit}
}

func (this *router) healtzh(context *gin.Context) {
    c := this.SetContext(context)
    c.OKWithMsg("hello world")
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    // 限流中间件：每秒生成 1 个令牌，桶容量为 2
    r.Use(core.RateLimiter(this.GetLimit, 1, 2))
    r.GET("/healthz", this.healtzh)
}

func (this router) Path() string {
    return "/"
}

func main() {
    g := core.God()
    // 注册限流器到依赖注入容器
    g.Provide(fxTags.Create(rateLimiter.RateLimit, fxTags.CreateTag("limit"), nil))
    // 挂载路由并注入限流器
    g.Provide(config.Config).Mount(NewRouter, fxTags.ParamList(fxTags.CreateTag("limit"))).Run()
}
```

## 运行与验证

### 启动服务

```shell
cd ip_limit
go run main.go
```

### 测试限流效果

本示例配置为 `rate=1, capacity=2`，即每秒生成 1 个令牌，桶最多容纳 2 个令牌。

**正常请求（桶中有令牌）：**

```shell
curl http://localhost:8888/healthz
```

返回：

```json
{
  "code": 0,
  "msg": "hello world"
}
```

**快速连续请求（耗尽令牌后被限流）：**

```shell
# 快速发送多个请求，模拟高频访问
for i in {1..5}; do
    echo "--- 请求 $i ---"
    curl -s http://localhost:8888/healthz
    echo ""
done
```

前 2 个请求正常返回（桶容量为 2），之后的请求被限流：

```
--- 请求 1 ---
{"code":0,"msg":"hello world"}
--- 请求 2 ---
{"code":0,"msg":"hello world"}
--- 请求 3 ---
{"code":-1,"msg":"rate limit exceeded"}
--- 请求 4 ---
{"code":-1,"msg":"rate limit exceeded"}
--- 请求 5 ---
{"code":-1,"msg":"rate limit exceeded"}
```

**等待令牌恢复后再次请求：**

```shell
# 等待 2 秒，令牌桶会补充令牌
sleep 2
curl http://localhost:8888/healthz
```

返回正常：

```json
{
  "code": 0,
  "msg": "hello world"
}
```

### 行为总结

| 场景 | 行为 |
|------|------|
| 桶中有令牌 | 请求正常通过，消耗 1 个令牌 |
| 桶中无令牌 | 返回 429 限流错误 |
| 等待一段时间 | 令牌按 rate 速率恢复，请求恢复正常 |

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

## 常用限流参数配置

根据不同业务场景选择合适的 `rate` 和 `capacity`：

| 场景 | rate（每秒令牌数） | capacity（桶容量） | 说明 |
|------|-------------------|-------------------|------|
| 普通 API | 10 | 20 | 允许每秒 10 次请求，突发最多 20 次 |
| 登录接口 | 1 | 3 | 严格限制，防止暴力破解 |
| 文件上传 | 2 | 5 | 适度限制，防止资源滥用 |
| 公开查询 | 50 | 100 | 宽松限制，保证可用性 |
| 本示例 | 1 | 2 | 演示用，每秒 1 次，突发最多 2 次 |

## 进阶用法

### 不同路由使用不同限流策略

可以注册多个限流器实例，为不同路由组配置不同的限流参数：

```go
func main() {
    g := core.God()

    // 注册两个不同标签的限流器
    g.Provide(fxTags.Create(rateLimiter.RateLimit, fxTags.CreateTag("strict"), nil))
    g.Provide(fxTags.Create(rateLimiter.RateLimit, fxTags.CreateTag("loose"), nil))

    // 分别注入到不同的路由控制器
    g.Provide(config.Config).
        Mount(NewAuthRouter, fxTags.ParamList(fxTags.CreateTag("strict"))).
        Mount(NewApiRouter, fxTags.ParamList(fxTags.CreateTag("loose"))).
        Run()
}
```

在各自的 `Build()` 中使用不同参数：

```go
// 登录路由 — 严格限流
func (this authRouter) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.Use(core.RateLimiter(this.GetLimit, 1, 3))   // 每秒 1 次，突发 3 次
    r.POST("/login", this.login)
}

// API 路由 — 宽松限流
func (this apiRouter) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.Use(core.RateLimiter(this.GetLimit, 50, 100)) // 每秒 50 次，突发 100 次
    r.GET("/users", this.listUsers)
}
```

### 结合路由级中间件实现部分路由限流

限流中间件也遵循路由级中间件的规则 — 只对 `r.Use()` 之后注册的路由生效：

```go
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    // 健康检查不限流
    r.GET("/healthz", this.healtzh)

    // 以下路由启用限流
    r.Use(core.RateLimiter(this.GetLimit, 10, 20))
    r.GET("/api/data", this.getData)
    r.POST("/api/submit", this.submit)
}
```

### Mount vs MountWithEmpty

| 方法 | 适用场景 | 说明 |
|------|---------|------|
| `MountWithEmpty(NewRouter)` | 路由构造函数无参数 | 无需依赖注入 |
| `Mount(NewRouter, fxTags.ParamList(...))` | 路由构造函数有参数 | 通过标签匹配注入依赖 |

本示例使用 `Mount` 是因为 `NewRouter` 需要接收 `rateLimiter.GetLimit` 参数。

## 依赖注入流程图

```
┌──────────────────────────────────────────────────────┐
│ main()                                               │
│                                                      │
│  g.Provide(fxTags.Create(                            │
│      rateLimiter.RateLimit,  ← 限流器构造函数          │
│      fxTags.CreateTag("limit"), ← 标签 "limit"       │
│      nil                                             │
│  ))                                                  │
│       │                                              │
│       ▼                                              │
│  g.Mount(NewRouter,                                  │
│      fxTags.ParamList(                               │
│          fxTags.CreateTag("limit") ← 匹配标签         │
│      )                                               │
│  )                                                   │
│       │                                              │
│       ▼                                              │
│  NewRouter(GetLimit) ← 自动注入限流器实例               │
│       │                                              │
│       ▼                                              │
│  Build() 中使用 core.RateLimiter(this.GetLimit, 1, 2) │
└──────────────────────────────────────────────────────┘
```

## 下一步

- 查看 `middleware_use/` — 了解全局中间件的用法
- 查看 `middleware_router_use/` — 了解路由级中间件的用法
- 查看 `provide_use/` — 了解依赖注入的更多用法
