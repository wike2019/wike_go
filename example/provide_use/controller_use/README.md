# controller_use 使用指南

## 项目概述

`controller_use` 是基于 `wike_go/v2` 框架的 Web 应用示例项目，演示了如何使用该框架的**控制器（Controller）**机制来组织路由、处理请求、参数校验以及缓存集成。

核心依赖：
- **Gin** — HTTP 路由引擎
- **wike_go/v2** — 封装了 Gin 的上层框架，提供控制器抽象、依赖注入（基于 uber/fx）、缓存接口等
- **Viper** — 配置管理
- **go-redis** — Redis 客户端

## 项目结构

```
controller_use/
├── main.go              # 入口文件，注册依赖和控制器
├── config.yaml          # 配置文件
├── config/
│   └── config.go        # Viper 配置初始化
├── controller/
│   ├── first.go         # 控制器1：健康检查（GET）
│   ├── second.go        # 控制器2：求和计算（POST + 参数校验）
│   └── cache.go         # 控制器3：缓存演示（Memory + Redis）
├── model/
│   └── sum.go           # 数据模型 + 校验逻辑
├── redis/
│   └── redis.go         # Redis 客户端初始化
├── db/
│   └── core.db          # SQLite 数据库文件
└── logs/
    └── app.log          # 日志文件
```

## 核心概念

### 1. 控制器接口

每个控制器需要满足以下约定：

| 方法 | 作用 |
|------|------|
| `Build(r *gin.RouterGroup, GCore *core.GCore)` | 注册路由到路由组 |
| `Path() string` | 返回路由前缀路径 |

控制器结构体需要内嵌 `ctl.Controller`，它提供了 `SetContext` 方法来包装 Gin 的 Context，获得统一的响应能力（`c.OK()`、`c.OKWithMsg()` 等）。

### 2. 依赖注入

框架使用 `uber/fx` 进行依赖注入，通过 `fxTags` 包实现接口绑定和标签化注入。

### 3. 挂载方式

| 方法 | 适用场景 |
|------|----------|
| `MountWithEmpty(constructor)` | 控制器无外部依赖，构造函数无参数 |
| `Mount(constructor, params)` | 控制器有外部依赖，需要通过 `fxTags.ParamList` 指定注入参数 |

---

## 详细文件说明

### config.yaml — 配置文件

```yaml
port: 8888
development: true
redis: "192.168.3.2:6379"
```

### config/config.go — 配置初始化

```go
func Config() *viper.Viper
```

- 使用 Viper 读取 `config.yaml`
- 设置默认值：`port=8888`、`logPath=./logs/app.log`、`development=true`
- 返回 `*viper.Viper` 实例，供其他模块通过依赖注入获取

### redis/redis.go — Redis 初始化

```go
func InitRedis(cfg *viper.Viper) *redis.Client
```

- 从 Viper 配置中读取 `redis` 地址
- 创建 Redis 客户端并 Ping 验证连接
- 连接失败时 `log.Fatalf` 直接终止程序

### model/sum.go — 数据模型

```go
type ModelSum struct {
    Num1 float64
    Num2 float64
}

func (this *ModelSum) Check() error
```

- 定义求和请求的数据结构
- `Check()` 方法校验两个数字不能小于零，返回中文错误信息

---

## 三个控制器详解

### 控制器1：健康检查（first.go）

最简单的控制器，无外部依赖。

```go
type router struct {
    ctl.Controller  // 内嵌基础控制器
}

func NewRouter() *router {
    return &router{}
}

func (this *router) healtzh(context *gin.Context) {
    c := this.SetContext(context)  // 包装 context
    c.OKWithMsg("hello world")    // 返回带消息的成功响应
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/healthz", this.healtzh)  // 注册 GET 路由
}

func (this router) Path() string {
    return "/"  // 路由前缀
}
```

**请求示例：**
```bash
curl http://localhost:8888/healthz
```

---

### 控制器2：求和计算（second.go）

演示 POST 请求体解析 + 参数校验。

```go
type router2 struct {
    ctl.Controller
}

func NewRouter2() *router2 {
    return &router2{}
}

func (this *router2) PostData(context *gin.Context) {
    c := this.SetContext(context)
    data := &model.ModelSum{}
    json.Unmarshal(c.Data, &data)     // 从 c.Data 解析 JSON 请求体
    ctl.Error(data.Check(), 400)       // 校验失败则返回 400 错误
    c.OK(data.Num1 + data.Num2)        // 返回求和结果
}

func (this router2) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.POST("/sum", this.PostData)
}

func (this router2) Path() string {
    return "/"
}
```

**请求示例：**
```bash
# 正常请求
curl -X POST http://localhost:8888/sum \
  -H "Content-Type: application/json" \
  -d '{"Num1": 10, "Num2": 20}'
# 响应：30

# 校验失败
curl -X POST http://localhost:8888/sum \
  -H "Content-Type: application/json" \
  -d '{"Num1": -1, "Num2": 20}'
# 响应：400 错误，"第一个数不能小于零"
```

**关键流程：**
1. `this.SetContext(context)` — 包装 Gin Context，框架自动将请求体读入 `c.Data`（`[]byte`）
2. `json.Unmarshal(c.Data, &data)` — 手动反序列化 JSON
3. `ctl.Error(data.Check(), 400)` — 如果 `Check()` 返回非 nil error，框架自动 panic 并返回 400 状态码（由 `CustomRecover` 中间件捕获）
4. `c.OK(result)` — 返回成功响应

---

### 控制器3：缓存演示（cache.go）

演示通过依赖注入使用 Memory 缓存和 Redis 缓存。

```go
type router3 struct {
    ctl.Controller
    redisCache  cache.Cacher   // Redis 缓存实例
    memoryCache cache.Cacher   // 内存缓存实例
}

func NewRouter3(memoryCache cache.Cacher, redisCache cache.Cacher) *router3 {
    return &router3{redisCache: redisCache, memoryCache: memoryCache}
}

func (this *router3) cache(context *gin.Context) {
    c := this.SetContext(context)

    // 内存缓存操作
    this.memoryCache.Set("test", "1", 0)
    str := ""
    this.memoryCache.Get("test", &str)
    fmt.Println(str)  // 输出: 1

    // Redis 缓存操作
    this.redisCache.Set("test2", "2", 0)
    str2 := ""
    this.redisCache.Get("test2", &str2)
    fmt.Println(str2)  // 输出: 2

    c.OK("ok")
}

func (this router3) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.POST("/cache", this.cache)
}

func (this router3) Path() string {
    return "/"
}
```

**请求示例：**
```bash
curl -X POST http://localhost:8888/cache
```

**缓存接口 `cache.Cacher`：**
```go
Set(key string, value interface{}, ttl time.Duration)
Get(key string, dest interface{})
```

- `ttl` 传 `0` 表示永不过期
- `Get` 通过指针回写值

---

## main.go — 启动入口详解

```go
func main() {
    g := core.God()                              // 1. 创建框架核心实例
    g.GlobalUse(core.CustomRecover(g))           // 2. 注册全局 panic 恢复中间件

    // 3. 注册缓存接口实现（带标签）
    g.Provide(fxTags.CreateInterFace(
        memory.NewCache,                          // 构造函数
        new(cache.Cacher),                        // 接口类型
        fxTags.CreateTag("cache_memory"),          // 注入标签
        nil,
    ))
    g.Provide(fxTags.CreateInterFace(
        redis.NewCache,
        new(cache.Cacher),
        fxTags.CreateTag("cache_redis"),
        nil,
    ))

    // 4. 注册配置和 Redis，挂载控制器
    g.Provide(config.Config, redisInit.InitRedis).
        MountWithEmpty(controller.NewRouter).      // 控制器1：无依赖
        MountWithEmpty(controller.NewRouter2).     // 控制器2：无依赖
        Mount(controller.NewRouter3,               // 控制器3：有依赖
            fxTags.ParamList(
                fxTags.CreateTag("cache_memory"),
                fxTags.CreateTag("cache_redis"),
            ),
        ).
        Run()                                      // 5. 启动服务
}
```

**启动流程：**

```
core.God() → 创建实例
    ↓
GlobalUse() → 注册全局中间件
    ↓
Provide() → 注册依赖（Config、Redis、Cache 接口）
    ↓
MountWithEmpty() / Mount() → 挂载控制器
    ↓
Run() → 启动 HTTP 服务（监听 config.yaml 中的 port）
```

---

## 快速开始 Demo

### 1. 创建一个新控制器

```go
// controller/my_controller.go
package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ctl"
)

type myRouter struct {
    ctl.Controller
}

func NewMyRouter() *myRouter {
    return &myRouter{}
}

func (this *myRouter) hello(context *gin.Context) {
    c := this.SetContext(context)
    c.OKWithMsg("你好，世界！")
}

func (this myRouter) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/hello", this.hello)
}

func (this myRouter) Path() string {
    return "/api"  // 路由前缀为 /api，完整路径为 /api/hello
}
```

### 2. 在 main.go 中注册

```go
g.Provide(config.Config, redisInit.InitRedis).
    MountWithEmpty(controller.NewRouter).
    MountWithEmpty(controller.NewRouter2).
    MountWithEmpty(controller.NewMyRouter).  // 新增
    Mount(controller.NewRouter3,
        fxTags.ParamList(
            fxTags.CreateTag("cache_memory"),
            fxTags.CreateTag("cache_redis"),
        ),
    ).
    Run()
```

### 3. 测试

```bash
curl http://localhost:8888/api/hello
# 响应：你好，世界！
```

---

## API 汇总

| 方法 | 路径 | 功能 | 控制器 |
|------|------|------|--------|
| GET | `/healthz` | 健康检查 | first.go |
| POST | `/sum` | 两数求和（带校验） | second.go |
| POST | `/cache` | 缓存读写演示 | cache.go |

---

## 关键 API 速查

| API | 说明 |
|-----|------|
| `core.God()` | 创建框架核心实例 |
| `g.GlobalUse(middleware)` | 注册全局中间件 |
| `g.Provide(constructors...)` | 注册依赖到 DI 容器 |
| `g.MountWithEmpty(constructor)` | 挂载无依赖的控制器 |
| `g.Mount(constructor, params)` | 挂载有依赖的控制器 |
| `fxTags.CreateInterFace(fn, iface, tag, nil)` | 将实现绑定到接口（带标签） |
| `fxTags.ParamList(tags...)` | 声明控制器构造函数的注入参数列表 |
| `fxTags.CreateTag(name)` | 创建依赖注入标签 |
| `this.SetContext(ctx)` | 包装 Gin Context，获取框架增强能力 |
| `c.OK(data)` | 返回成功响应（带数据） |
| `c.OKWithMsg(msg)` | 返回成功响应（带消息） |
| `ctl.Error(err, code)` | 若 err 非 nil，panic 并返回指定状态码 |
