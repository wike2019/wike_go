# 介绍

Wike Go 是一个基于 [Gin](https://github.com/gin-gonic/gin) 构建的**模块化 Go Web 框架**，通过 [Uber fx](https://github.com/uber-go/fx) 实现依赖注入，提供了一套开箱即用的企业级 Web 应用基础设施。

开发者只需关注业务逻辑，框架负责处理配置加载、中间件编排、缓存、定时任务、日志、鉴权等通用能力。

如果你已经准备在项目中使用 Wike Go，请访问 [快速开始](/fast/快速开始)。

## 核心架构

```
main.go
  └── core.God() → Provide(Config) → Mount(Router) → Run()
        │
        ├── core/          框架核心（DI、HTTP 服务、中间件、日志）
        ├── func/          功能模块（缓存、JWT、限流、熔断、重试等）
        ├── db/            数据库层（GORM + SQLite）
        ├── model/         数据模型（自定义时间类型、JSON 类型、定时任务模型）
        ├── utils/         工具集（分布式锁、线程安全 Map/Set、密码哈希、Redis 队列等）
        ├── cronJob/       定时任务调度器
        ├── fxTags/        fx 依赖注入标签辅助
        └── constData/     常量定义
```

## 技术栈

| 领域 | 选型 |
|------|------|
| HTTP 框架 | [Gin](https://github.com/gin-gonic/gin) |
| 依赖注入 | [Uber fx](https://github.com/uber-go/fx) |
| 数据库 ORM | [GORM](https://gorm.io) + SQLite |
| 缓存 | [go-cache](https://github.com/patrickmn/go-cache)（本地）/ [go-redis](https://github.com/redis/go-redis)（分布式） |
| 日志 | [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack)（日志轮转） |
| 配置管理 | [Viper](https://github.com/spf13/viper)（YAML） |
| 鉴权 | [golang-jwt](https://github.com/golang-jwt/jwt)（JWT 认证） |
| 定时任务 | [robfig/cron](https://github.com/robfig/cron) |
| 协程池 | [panjf2000/ants](https://github.com/panjf2000/ants) |

## 主要特性

### 请求处理管线

请求经过完整的中间件链路处理：

```
请求 → 优雅停机拦截 → 链路追踪 → 访问日志 → CORS → Panic 恢复
     → 请求体限制 → 超时控制 → 限流 → RBAC 鉴权 → 业务 Handler
```

### 缓存旁路模式

通过 `Cacher` 接口统一内存缓存和 Redis 缓存，`FindWithCallBack` 实现通用的 cache-aside 读取策略。

### 弹性设计

- **熔断器** — 基于 sony/gobreaker，60% 失败率触发熔断
- **重试机制** — 指数退避策略
- **分布式锁** — Redis 实现，支持自动续期
- **令牌桶限流** — 精细化流量控制

### 泛型工具集

充分利用 Go 泛型提升类型安全：

- `MapSync[T]` — 线程安全的泛型 Map
- 泛型 Set — 高效集合操作
- `InfoData[T]` — 泛型 JWT Claims

### 优雅停机

收到信号后标记 Reject 状态，新请求返回 503，等待存量请求处理完毕（10 秒超时），依次执行清理函数。

## 最小示例

```go
func main() {
    g := core.God()
    g.Provide(Config).MountWithEmpty(NewRouter).Run()
}
```

三行代码完成：创建框架实例 → 注入配置 → 挂载路由 → 启动服务。

业务路由只需实现两个方法：

- `Build(r *gin.RouterGroup, GCore *core.GCore)` — 注册路由
- `Path() string` — 定义路由前缀

## 设计理念

- **约定优于配置** — Viper 加载 YAML，所有参数都有合理默认值
- **模块独立** — 缓存、鉴权、限流、熔断等模块互不耦合，按需组合
- **DI 驱动** — 通过 fx 实现松耦合，便于测试和替换实现
- **中间件可组合** — 全局中间件和路由级中间件灵活编排
- **并发安全** — 工具层全面使用读写锁、协程池、分布式锁保障线程安全

## 适用场景

- RESTful API 服务
- 企业级后台管理系统
- 微服务架构中的单体服务
- 需要快速交付的 Go Web 项目


## 项目接口文档
https://wg.ng2-oa.com/
