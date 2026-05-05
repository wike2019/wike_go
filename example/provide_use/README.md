# Wike Go 依赖注入使用指南

本项目是 `wike_go` 框架依赖注入（DI）功能的完整示例，涵盖了从基础注入到高级 Tag 标签、接口注入、缓存等全部用法。

## 目录

- [项目结构](#项目结构)
- [前置条件](#前置条件)
- [快速开始](#快速开始)
- [核心概念](#核心概念)
  - [God 实例与链式调用](#god-实例与链式调用)
  - [Config 配置注入](#config-配置注入)
  - [Provide 依赖注册](#provide-依赖注册)
  - [Supply 直接值注入](#supply-直接值注入)
  - [Invokes 立即执行](#invokes-立即执行)
  - [MountWithEmpty 无依赖控制器挂载](#mountwithempty-无依赖控制器挂载)
  - [Mount + fxTags 带标签的依赖注入](#mount--fxtags-带标签的依赖注入)
  - [接口注入 CreateInterFace](#接口注入-createinterface)
- [完整 Demo](#完整-demo)
  - [Demo 1：基础控制器（无外部依赖）](#demo-1基础控制器无外部依赖)
  - [Demo 2：同类型多实例注入（fxTags）](#demo-2同类型多实例注入fxtags)
  - [Demo 3：接口注入（Memory + Redis 缓存）](#demo-3接口注入memory--redis-缓存)
  - [Demo 4：Supply 直接注入已有对象](#demo-4supply-直接注入已有对象)
- [API 速查表](#api-速查表)
- [常见问题](#常见问题)

---

## 项目结构

```
provide_use/
├── main.go              # 入口文件，展示全部 DI 用法
├── config.yaml          # 配置文件
├── config/
│   └── config.go        # 配置中心（基于 Viper）
├── model/
│   ├── config.go        # NoParam 配置模型（提供 Dog、Cat 默认实例）
│   ├── dog.go           # Dog 模型
│   └── cat.go           # Cat、CatSupply 模型
├── provide/
│   ├── first.go         # ProvideExampleFirst —— 第一个 Animal 提供者
│   └── second.go        # ProvideExampleSecond —— 第二个 Animal 提供者（含聚合方法）
├── controller/
│   ├── router.go        # 基础路由控制器（注入 Dog + Cat）
│   ├── routerDog.go     # Tag 路由控制器（注入两个 Animal）
│   └── cache.go         # 缓存路由控制器（注入 Memory + Redis 缓存）
├── redis/
│   └── redis.go         # Redis 客户端初始化
├── db/
│   └── core.db          # SQLite 数据库文件
├── go.mod
└── go.sum
```

## 前置条件

- Go 1.20+
- Redis 服务（可选，用于缓存示例）
- 安装框架依赖：

```shell
go get github.com/wike2019/wike_go/v2@v2.0.11
```

## 快速开始

```shell
# 1. 进入项目目录
cd provide_use

# 2. 安装依赖
go mod tidy

# 3. 修改 config.yaml 中的 redis 地址（如不需要 Redis 可去掉相关代码）
# redis: "127.0.0.1:6379"

# 4. 运行
go run main.go

# 5. 测试接口
curl http://localhost:8888/animal       # 基础控制器
curl http://localhost:8888/dog/healthz  # Tag 注入控制器
curl -X POST http://localhost:8888/cache # 缓存控制器
```

---

## 核心概念

### God 实例与链式调用

`core.God()` 是框架的入口，返回一个支持链式调用的实例。整个应用的依赖注册、控制器挂载、启动都通过链式调用完成：

```go
g := core.God()
g.Config(&model.NoParam{}).   // 注册配置模型
    Provide(config.Config).   // 注册依赖
    MountWithEmpty(NewRouter). // 挂载控制器
    Run()                     // 启动服务
```

### Config 配置注入

`Config()` 方法接收一个配置模型，框架会自动调用该模型上的方法来生成依赖对象。

```go
// model/config.go
type NoParam struct{}

// 框架会自动调用这些方法，将返回值注入到 DI 容器中
func (this *NoParam) Dog() *Dog {
    return &Dog{Name: "我是一只狗"}
}
func (this *NoParam) Cat() *Cat {
    return &Cat{Name: "我是一只猫"}
}
```

```go
// main.go
g.Config(&model.NoParam{})
// 此时 DI 容器中已有 *model.Dog 和 *model.Cat 两个实例
```

> **注意**：同一类型只能有一个 Config 方法。如果定义了两个返回 `*Cat` 的方法（如 `Cat()` 和 `Cat2()`），框架会因为无法区分而异常退出。

### Provide 依赖注册

`Provide()` 将一个构造函数注册到 DI 容器。框架会自动解析函数参数，从容器中查找匹配的依赖并注入。

```go
// config/config.go —— 无参数的 Provide
func Config() *viper.Viper {
    // ...加载配置
    return viper.GetViper()
}

// redis/redis.go —— 依赖 *viper.Viper 的 Provide
func InitRedis(cfg *viper.Viper) *redis.Client {
    addr := cfg.GetString("redis")
    client := redis.NewClient(&redis.Options{Addr: addr})
    // ...
    return client
}
```

```go
// main.go —— 多个 Provide 可以链式调用，也可以在一个 Provide 中传入多个函数
g.Provide(config.Config, redisInit.InitRedis)
```

依赖解析是自动的：`InitRedis` 需要 `*viper.Viper`，框架会从容器中找到 `Config()` 的返回值并注入。

### Supply 直接值注入

`Supply()` 将一个**已经创建好的对象**直接放入 DI 容器，无需构造函数。

```go
CatSupply := &model.CatSupply{Name: "先初始化完成再注入"}

g.Supply(CatSupply)
// 此时 DI 容器中有 *model.CatSupply 实例
```

适用场景：
- 需要在框架启动前手动初始化的对象
- 外部传入的配置或连接
- 测试时注入 mock 对象

### Invokes 立即执行

`Invokes()` 注册一个在框架启动时**立即执行**的函数。函数参数会从 DI 容器中自动注入。

```go
g.Invokes(func(r *model.CatSupply) {
    fmt.Println(r.Name, "这个是立即执行的，这里的代码不能阻塞 必须开协程")
})
```

> **重要**：`Invokes` 中的代码会在启动阶段同步执行，**不能阻塞**。如果有耗时操作，必须开启 goroutine。

### MountWithEmpty 无依赖控制器挂载

当控制器的构造函数**不需要额外的 Tag 参数**时，使用 `MountWithEmpty`：

```go
// controller/router.go
func NewRouter(dog *model.Dog, cat *model.Cat) *router {
    return &router{dog: dog, cat: cat}
}

// main.go
g.MountWithEmpty(controller.NewRouter)
```

框架会自动从 DI 容器中找到 `*model.Dog` 和 `*model.Cat` 并注入到 `NewRouter`。

### Mount + fxTags 带标签的依赖注入

当 DI 容器中存在**同一类型的多个实例**时，框架无法区分该注入哪个。此时需要使用 `fxTags` 打标签。

**问题场景**：两个 `*provide.Animal` 实例，框架不知道选哪个。

**解决方案**：

```go
// 第一步：注册时打标签
g.Provide(fxTags.Create(provide.ProvideExampleFirst, fxTags.CreateTag("dogFirst"), nil)).
  Provide(fxTags.Create(provide.ProvideExampleSecond, fxTags.CreateTag("dogSecond"), nil))

// 第二步：挂载时指定参数标签顺序
g.Mount(controller.NewRouterDog,
    fxTags.ParamList(
        fxTags.CreateTag("dogFirst"),
        fxTags.CreateTag("dogSecond"),
    ))
```

`fxTags.Create` 参数说明：

| 参数 | 说明 |
|------|------|
| 第 1 个 | 构造函数 |
| 第 2 个 | 输出标签（`fxTags.CreateTag("name")`） |
| 第 3 个 | 输入标签（`nil` 表示无特殊输入标签） |

`fxTags.ParamList` 按**构造函数参数顺序**指定每个参数对应的标签。

### 接口注入 CreateInterFace

当需要注入**接口类型**而非具体类型时，使用 `fxTags.CreateInterFace`。这是实现多态和策略模式的关键。

```go
// 注册 memory 缓存实现，标记为 cache.Cacher 接口
g.Provide(fxTags.CreateInterFace(
    memory.NewCache,           // 构造函数
    new(cache.Cacher),         // 接口类型指针
    fxTags.CreateTag("cache_memory"), // 输出标签
    nil,                       // 输入标签
))

// 注册 redis 缓存实现，同样标记为 cache.Cacher 接口
g.Provide(fxTags.CreateInterFace(
    redis.NewCache,
    new(cache.Cacher),
    fxTags.CreateTag("cache_redis"),
    nil,
))

// 挂载控制器时按标签注入
g.Mount(controller.NewRouter3,
    fxTags.ParamList(
        fxTags.CreateTag("cache_memory"),
        fxTags.CreateTag("cache_redis"),
    ))
```

控制器中通过接口使用，无需关心具体实现：

```go
type router3 struct {
    ctl.Controller
    redisCache  cache.Cacher  // 接口类型
    memoryCache cache.Cacher  // 接口类型
}

func (this *router3) cache(context *gin.Context) {
    this.memoryCache.Set("test", "1", 0)  // 使用内存缓存
    this.redisCache.Set("test2", "2", 0)  // 使用 Redis 缓存
}
```

---

## 完整 Demo

### Demo 1：基础控制器（无外部依赖）

最简单的用法，通过 `Config` 自动生成依赖，`MountWithEmpty` 挂载控制器。

**model/config.go**
```go
type NoParam struct{}

func (this *NoParam) Dog() *Dog {
    return &Dog{Name: "我是一只狗"}
}
func (this *NoParam) Cat() *Cat {
    return &Cat{Name: "我是一只猫"}
}
```

**controller/router.go**
```go
type router struct {
    ctl.Controller
    dog *model.Dog
    cat *model.Cat
}

func NewRouter(dog *model.Dog, cat *model.Cat) *router {
    return &router{dog: dog, cat: cat}
}

func (this *router) healtzh(context *gin.Context) {
    c := this.SetContext(context)
    c.OK(this.dog.Name + this.cat.Name)
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/animal", this.healtzh)
}

func (this router) Path() string {
    return "/"
}
```

**main.go（最小启动）**
```go
func main() {
    g := core.God()
    g.Config(&model.NoParam{}).
        Provide(config.Config).
        MountWithEmpty(controller.NewRouter).
        Run()
}
```

**测试**：`GET /animal` → 返回 `"我是一只狗我是一只猫"`

---

### Demo 2：同类型多实例注入（fxTags）

当同一类型有多个实例时，必须使用 Tag 区分。

**provide/first.go**
```go
func ProvideExampleFirst(dog *model.Dog) *Animal {
    return &Animal{dog: dog}
}
```

**provide/second.go**
```go
type Animal struct {
    dog *model.Dog
}

func (this *Animal) Show() string {
    return this.dog.Name + "这个是聚合方法"
}

func ProvideExampleSecond(dog *model.Dog) *Animal {
    return &Animal{dog: dog}
}
```

**controller/routerDog.go**
```go
type routerDog struct {
    ctl.Controller
    Animal1 *provide.Animal
    Animal2 *provide.Animal
}

func NewRouterDog(dog *provide.Animal, dog2 *provide.Animal) *routerDog {
    return &routerDog{Animal1: dog, Animal2: dog2}
}

func (this *routerDog) healtzh(context *gin.Context) {
    c := this.SetContext(context)
    c.OK(this.Animal1.Show() + this.Animal2.Show())
}

func (this routerDog) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/healthz", this.healtzh)
}

func (this routerDog) Path() string {
    return "/dog"
}
```

**main.go 注册方式**
```go
g.Provide(fxTags.Create(provide.ProvideExampleSecond, fxTags.CreateTag("dogSecond"), nil)).
  Provide(fxTags.Create(provide.ProvideExampleFirst, fxTags.CreateTag("dogFirst"), nil)).
  Mount(controller.NewRouterDog,
      fxTags.ParamList(
          fxTags.CreateTag("dogFirst"),
          fxTags.CreateTag("dogSecond"),
      ))
```

**测试**：`GET /dog/healthz` → 返回两个 Animal 的聚合结果

---

### Demo 3：接口注入（Memory + Redis 缓存）

基于接口的依赖注入，实现同一接口的不同实现并存。

**controller/cache.go**
```go
type router3 struct {
    ctl.Controller
    redisCache  cache.Cacher
    memoryCache cache.Cacher
}

func NewRouter3(memoryCache cache.Cacher, redisCache cache.Cacher) *router3 {
    return &router3{redisCache: redisCache, memoryCache: memoryCache}
}

func (this *router3) cache(context *gin.Context) {
    c := this.SetContext(context)

    // 使用内存缓存
    this.memoryCache.Set("test", "1", 0)
    str := ""
    this.memoryCache.Get("test", &str)

    // 使用 Redis 缓存
    this.redisCache.Set("test2", "2", 0)
    str2 := ""
    this.redisCache.Get("test2", &str2)

    c.OK("ok")
}

func (this router3) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.POST("/cache", this.cache)
}

func (this router3) Path() string {
    return "/"
}
```

**main.go 注册方式**
```go
g.Provide(fxTags.CreateInterFace(memory.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_memory"), nil)).
  Provide(fxTags.CreateInterFace(redis.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_redis"), nil)).
  Mount(controller.NewRouter3,
      fxTags.ParamList(
          fxTags.CreateTag("cache_memory"),
          fxTags.CreateTag("cache_redis"),
      ))
```

**测试**：`POST /cache` → 同时操作内存缓存和 Redis 缓存

---

### Demo 4：Supply 直接注入已有对象

将已创建的对象直接注入容器，配合 `Invokes` 立即使用。

```go
func main() {
    // 在框架启动前手动创建对象
    CatSupply := &model.CatSupply{Name: "先初始化完成再注入"}

    g := core.God()
    g.Supply(CatSupply).
      Invokes(func(r *model.CatSupply) {
          // 框架启动时立即执行，r 就是上面 Supply 的对象
          fmt.Println(r.Name, "这个是立即执行的，这里的代码不能阻塞 必须开协程")
      }).
      Run()
}
```

---

## API 速查表

| 方法 | 说明 | 使用场景 |
|------|------|----------|
| `core.God()` | 创建框架实例 | 入口，必须调用 |
| `.Config(model)` | 注册配置模型，自动调用其方法生成依赖 | 提供基础模型（Dog、Cat 等） |
| `.Provide(fn...)` | 注册构造函数到 DI 容器 | 配置、数据库、Redis 等服务初始化 |
| `.Supply(obj)` | 直接注入已有对象 | 手动创建的对象、外部传入的依赖 |
| `.Invokes(fn)` | 注册启动时立即执行的函数 | 初始化逻辑、启动检查（不可阻塞） |
| `.MountWithEmpty(ctor)` | 挂载无 Tag 依赖的控制器 | 控制器参数类型在容器中唯一 |
| `.Mount(ctor, params)` | 挂载带 Tag 的控制器 | 控制器参数类型在容器中有多个实例 |
| `.Run()` | 启动 HTTP 服务 | 链式调用的最后一步 |
| `fxTags.Create(fn, outTag, inTag)` | 为 Provide 的构造函数打输出标签 | 同类型多实例区分 |
| `fxTags.CreateInterFace(fn, iface, outTag, inTag)` | 接口类型的 Provide + 打标签 | 同接口多实现区分 |
| `fxTags.CreateTag(name)` | 创建标签 | 配合 Create / CreateInterFace / ParamList |
| `fxTags.ParamList(tags...)` | 按参数顺序指定标签列表 | 配合 Mount 使用 |

## 常见问题

### 1. 启动报错：多个同类型实例冲突

**错误原因**：DI 容器中存在两个相同类型的对象，框架不知道注入哪个。

**解决方案**：使用 `fxTags.Create` 为每个实例打标签，使用 `Mount` + `fxTags.ParamList` 指定注入顺序。

```go
// 错误写法
g.Provide(provide.ProvideExampleFirst).
  Provide(provide.ProvideExampleSecond)  // 两个都返回 *Animal，冲突！

// 正确写法
g.Provide(fxTags.Create(provide.ProvideExampleFirst, fxTags.CreateTag("dogFirst"), nil)).
  Provide(fxTags.Create(provide.ProvideExampleSecond, fxTags.CreateTag("dogSecond"), nil))
```

### 2. Invokes 中的代码阻塞导致服务无法启动

`Invokes` 在启动阶段同步执行。如果有耗时操作（如监听消息队列），必须开启 goroutine：

```go
g.Invokes(func(r *model.CatSupply) {
    go func() {
        // 耗时操作放在 goroutine 中
    }()
})
```

### 3. Config 模型中定义了多个返回相同类型的方法

```go
// 错误：两个方法都返回 *Cat，框架无法区分
func (this *NoParam) Cat() *Cat { ... }
func (this *NoParam) Cat2() *Cat { ... }  // 会导致启动失败
```

每种类型在 Config 模型中只能有一个工厂方法。如果需要多个同类型实例，使用 `fxTags` 方案。

### 4. Mount 和 MountWithEmpty 的区别

| 方法 | 适用场景 | 是否需要 fxTags |
|------|----------|----------------|
| `MountWithEmpty` | 控制器参数类型在容器中唯一 | 不需要 |
| `Mount` | 控制器参数类型在容器中有多个实例 | 需要 ParamList |

### 5. Create 和 CreateInterFace 的区别

| 方法 | 适用场景 | 额外参数 |
|------|----------|----------|
| `fxTags.Create` | 注入具体类型（如 `*Animal`） | 无 |
| `fxTags.CreateInterFace` | 注入接口类型（如 `cache.Cacher`） | 需要 `new(接口)` 指定接口类型 |

---

## 配置文件

**config.yaml**

```yaml
port: 8888           # 服务端口
development: true    # 开发模式
redis: "127.0.0.1:6379"  # Redis 地址
```

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `port` | HTTP 服务监听端口 | `8888` |
| `development` | 开发模式，输出详细日志 | `true` |
| `logPath` | 日志文件路径 | `./logs/app.log` |
| `redis` | Redis 连接地址 | 无 |

## 技术栈

- **Web 框架**：[Gin](https://github.com/gin-gonic/gin)
- **依赖注入**：[Uber Fx](https://github.com/uber-go/fx)（由 wike_go 封装）
- **配置管理**：[Viper](https://github.com/spf13/viper)
- **缓存**：内存缓存（go-cache）+ Redis（go-redis）
- **ORM**：GORM（可选）
