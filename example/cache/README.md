# Cache 缓存服务示例

本项目演示了 `wike_go` 框架中缓存模块的使用方法，支持 **Memory（内存缓存）** 和 **Redis（分布式缓存）** 两种存储后端，通过统一的 `cache.Cacher` 接口进行操作。

## 目录结构

```
cache/
├── config/
│   └── config.go          # 配置中心，基于 viper 读取 config.yaml
├── controller/
│   └── cache.go           # 控制器，演示缓存的实际使用
├── db/
│   └── core.db            # SQLite 数据库文件
├── model/
│   └── sum.go             # 数据模型示例
├── redis/
│   └── redis.go           # Redis 客户端初始化
├── config.yaml            # 配置文件
├── go.mod                 # Go 模块定义
└── main.go                # 程序入口
```

## 核心概念

### 1. 缓存接口 `cache.Cacher`

框架提供统一的缓存接口，所有缓存后端（Memory / Redis）都实现该接口，使用时无需关心底层存储细节。

### 2. 缓存服务 `cache.Service`

`cache.Service` 是对 `cache.Cacher` 的高层封装，提供带回调的缓存查询方法 `FindWithCallBack`，实现"缓存未命中时自动执行回调并写入缓存"的模式。

### 3. 两种缓存后端

| 后端 | 构造函数 | 适用场景 |
|------|----------|----------|
| Memory | `memory.NewCache` | 单机部署、开发调试、无需持久化 |
| Redis | `redis.NewCache` | 分布式部署、需要持久化、多实例共享 |

## 配置文件

`config.yaml`:

```yaml
port: 8888
development: true
redis: "192.168.3.2:6379"
```

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `port` | 服务端口 | 8888 |
| `development` | 开发模式 | true |
| `redis` | Redis 地址 | - |

## 使用方法

### 步骤一：初始化 Redis 客户端

```go
// redis/redis.go
package redis

import (
    "context"
    "log"

    "github.com/redis/go-redis/v9"
    "github.com/spf13/viper"
)

func InitRedis(cfg *viper.Viper) *redis.Client {
    addr := cfg.GetString("redis")
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        log.Fatalf("failed to ping redis: %s", err.Error())
    }
    log.Println("redis connected:", addr)
    return client
}
```

### 步骤二：注册缓存提供者（main.go）

```go
package main

import (
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_app/controller"
    redisInit "github.com/wike2019/wike_app/redis"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/cache"
    "github.com/wike2019/wike_go/v2/func/cache/memory"
    "github.com/wike2019/wike_go/v2/func/cache/redis"
    "github.com/wike2019/wike_go/v2/fxTags"
)

func main() {
    g := core.God()
    g.GlobalUse(core.CustomRecover(g))

    // 注册 Memory 缓存实现，标记为 "cache_memory"
    g.Provide(fxTags.CreateInterFace(memory.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_memory"), nil))

    // 注册 Redis 缓存实现，标记为 "cache_redis"
    g.Provide(fxTags.CreateInterFace(redis.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_redis"), nil))

    // 提供配置和 Redis 客户端
    g.Provide(config.Config, redisInit.InitRedis).
        // 创建缓存服务，选择 Redis 作为底层存储
        Provide(fxTags.Create(cache.ServiceCache, "", fxTags.ParamList(fxTags.CreateTag("cache_redis")))).
        // 挂载控制器，同时注入 memory 和 redis 两种缓存
        Mount(controller.NewRouter3,
            fxTags.ParamList(
                fxTags.CreateTag("cache_memory"), fxTags.CreateTag("cache_redis"))).
        Run()
}
```

**关键说明：**

- `fxTags.CreateInterFace` 将具体实现注册为 `cache.Cacher` 接口，并通过 tag 区分不同后端
- `fxTags.CreateTag("cache_memory")` / `fxTags.CreateTag("cache_redis")` 用于依赖注入时的标识
- `cache.ServiceCache` 创建缓存服务时，通过 `ParamList` 指定使用哪个后端（此处选择 Redis）
- `Mount` 挂载控制器时可同时注入多个缓存实例

### 步骤三：在控制器中使用缓存

```go
package controller

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/cache"
    "github.com/wike2019/wike_go/v2/func/ctl"
)

type router3 struct {
    ctl.Controller
    redisCache  cache.Cacher      // Redis 缓存实例
    memoryCache cache.Cacher      // Memory 缓存实例
    cache       *cache.Service    // 缓存服务（高层封装）
}

func NewRouter3(memoryCache cache.Cacher, redisCache cache.Cacher, cache *cache.Service) *router3 {
    return &router3{redisCache: redisCache, memoryCache: memoryCache, cache: cache}
}

func (this router3) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.POST("/cache", this.cacheData)
}

func (this *router3) cacheData(context *gin.Context) {
    c := this.SetContext(context)
    key := "healtzh"

    // 使用 FindWithCallBack：缓存命中直接返回，未命中则执行回调并缓存结果
    data := cache.FindWithCallBack[gin.H](key, time.Second*60, this.cache, func() gin.H {
        return gin.H{"data": "form cache"}
    })
    c.OK(data)
}

func (this router3) Path() string {
    return "/"
}
```

## 核心 API

### `cache.FindWithCallBack[T]`

```go
func FindWithCallBack[T any](key string, ttl time.Duration, service *cache.Service, callback func() T) T
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `key` | `string` | 缓存键名 |
| `ttl` | `time.Duration` | 缓存过期时间 |
| `service` | `*cache.Service` | 缓存服务实例 |
| `callback` | `func() T` | 缓存未命中时的数据获取回调 |

**执行逻辑：**

1. 根据 `key` 查询缓存
2. 命中 → 直接返回缓存数据（反序列化为类型 `T`）
3. 未命中 → 执行 `callback()` 获取数据 → 写入缓存（设置 TTL）→ 返回数据

## 使用 Demo

### Demo 1：使用 Redis 缓存服务（推荐）

通过 `cache.Service` + `FindWithCallBack` 实现自动缓存管理：

```go
func (this *router3) getUserInfo(context *gin.Context) {
    c := this.SetContext(context)
    userID := c.Param("id")

    // 缓存 5 分钟，未命中时从数据库查询
    user := cache.FindWithCallBack[User]("user:"+userID, time.Minute*5, this.cache, func() User {
        // 这里执行实际的数据库查询
        return queryUserFromDB(userID)
    })
    c.OK(user)
}
```

### Demo 2：直接使用 Cacher 接口

如果需要更细粒度的控制，可以直接操作 `cache.Cacher` 接口：

```go
func (this *router3) manualCache(context *gin.Context) {
    c := this.SetContext(context)

    // 使用 memory 缓存
    this.memoryCache.Set("key1", "value1", time.Minute*10)
    val, _ := this.memoryCache.Get("key1")

    // 使用 redis 缓存
    this.redisCache.Set("key2", "value2", time.Hour)
    val2, _ := this.redisCache.Get("key2")

    c.OK(gin.H{"memory": val, "redis": val2})
}
```

### Demo 3：切换缓存后端

只需修改 `main.go` 中 `ServiceCache` 的参数即可切换底层存储：

```go
// 使用 Redis 作为缓存服务存储
Provide(fxTags.Create(cache.ServiceCache, "", fxTags.ParamList(fxTags.CreateTag("cache_redis"))))

// 切换为 Memory 作为缓存服务存储（无需 Redis，适合开发调试）
Provide(fxTags.Create(cache.ServiceCache, "", fxTags.ParamList(fxTags.CreateTag("cache_memory"))))
```

## 依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/wike2019/wike_go/v2` | v2.0.11 | 核心框架（含缓存模块） |
| `github.com/gin-gonic/gin` | v1.12.0 | HTTP 框架 |
| `github.com/redis/go-redis/v9` | v9.19.0 | Redis 客户端 |
| `github.com/spf13/viper` | v1.21.0 | 配置管理 |

## 运行

```bash
# 确保 Redis 服务已启动
# 修改 config.yaml 中的 redis 地址

go run main.go
```

服务启动后，访问 `POST /cache` 即可测试缓存功能。
