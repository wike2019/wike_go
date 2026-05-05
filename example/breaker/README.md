# 熔断器（Circuit Breaker）

本节介绍如何在 Wike Go 项目中使用熔断器保护外部服务调用。当下游服务不可用时，熔断器会自动切断请求，避免雪崩效应，并在服务恢复后自动放行。

## 核心概念

熔断器有三种状态：

```
  ┌──────────┐    连续失败达到阈值    ┌──────────┐
  │  Closed  │ ──────────────────→ │   Open   │
  │ (正常放行) │                      │ (拒绝请求) │
  └──────────┘                      └──────────┘
       ↑                                 │
       │          Timeout 超时后          │
       │                                 ↓
       │                           ┌───────────┐
       └────── 探测请求成功 ──────── │ Half-Open │
                                   │ (试探放行)  │
                                   └───────────┘
```

| 状态 | 行为 |
|------|------|
| **Closed（关闭）** | 正常状态，所有请求正常通过。当失败率达到阈值时，切换到 Open |
| **Open（打开）** | 熔断状态，所有请求立即被拒绝，不会调用下游服务。经过 Timeout 时间后，切换到 Half-Open |
| **Half-Open（半开）** | 试探状态，允许有限数量的请求通过。如果成功则恢复到 Closed，失败则回到 Open |

## 依赖

框架基于 [sony/gobreaker](https://github.com/sony/gobreaker) 封装，提供了开箱即用的熔断器工厂方法。

```shell
go get github.com/sony/gobreaker
go get github.com/wike2019/wike_go/v2@v2.0.11
```

## 项目结构

```
breaker/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── db/
│   └── core.db        # SQLite 数据库文件
├── config.yaml        # 配置文件
├── go.mod             # 依赖管理
├── main.go            # 入口文件（熔断器示例）
└── README.md
```

## API 说明

### `breaker.NewCircuitBreaker`

快速创建一个熔断器实例，使用框架内置的默认触发策略。

```go
func NewCircuitBreaker(
    name        string,         // 熔断器名称，用于日志和监控标识
    MaxRequests uint32,         // Half-Open 状态下允许通过的最大请求数
    Interval    time.Duration,  // Closed 状态下的统计周期，周期结束后计数器清零
    Timeout     time.Duration,  // Open 状态持续时间，超时后进入 Half-Open
) *gobreaker.CircuitBreaker
```

**默认触发策略（ReadyToTrip）：**

```go
// 当请求数 >= 6 且失败率 >= 60% 时触发熔断
var DefaultToTrip = func(counts gobreaker.Counts) bool {
    failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
    return counts.Requests >= 6 && failureRatio >= 0.6
}
```

### `breaker.NewCircuitBreakerWithSettings`

使用完全自定义的 `gobreaker.Settings` 创建熔断器，适用于需要自定义触发策略、状态变更回调等高级场景。

```go
func NewCircuitBreakerWithSettings(settings gobreaker.Settings) *gobreaker.CircuitBreaker
```

### `CircuitBreaker.Execute`

执行受熔断器保护的函数。如果熔断器处于 Open 状态，直接返回错误而不执行函数体。

```go
func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error)
```

## 参数详解

| 参数 | 类型 | 说明 | 示例值 |
|------|------|------|--------|
| `name` | `string` | 熔断器名称，用于区分不同的熔断器实例 | `"http-breaker"` |
| `MaxRequests` | `uint32` | Half-Open 状态下允许通过的最大探测请求数。探测成功则恢复，失败则重新熔断 | `5` |
| `Interval` | `time.Duration` | Closed 状态下的统计窗口周期。每个周期结束后，失败计数器会重置 | `10 * time.Second` |
| `Timeout` | `time.Duration` | Open 状态的持续时间。超时后自动进入 Half-Open 状态进行探测 | `3 * time.Second` |

## 使用示例

### 基础用法：保护 HTTP 请求

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sony/gobreaker"
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/breaker"
    "github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
    ctl.Controller
    Breaker *gobreaker.CircuitBreaker
}

func NewRouter() *router {
    return &router{
        // 创建熔断器：名称 "breaker"，半开状态最多放行 5 个请求，
        // 统计窗口 10 秒，熔断后 3 秒进入半开状态
        Breaker: breaker.NewCircuitBreaker("breaker", 5, time.Second*10, time.Second*3),
    }
}

// Job 定义受熔断器保护的业务逻辑
func (this *router) Job() (interface{}, error) {
    resp, err := http.Get("https://baidu.com")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
    }
    return "Request succeeded", nil
}

// healthz 处理函数，通过熔断器执行业务逻辑
func (this *router) healthz(context *gin.Context) {
    c := this.SetContext(context)
    data, err := this.Breaker.Execute(this.Job)
    ctl.Error(err, 500)
    c.OK(data)
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
    r.GET("/healthz", this.healthz)
}

func (this router) Path() string {
    return "/"
}

func main() {
    g := core.God()
    g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
```

### 自定义触发策略

当默认的触发策略不满足需求时，使用 `NewCircuitBreakerWithSettings` 完全自定义：

```go
import (
    "github.com/sony/gobreaker"
    "github.com/wike2019/wike_go/v2/func/breaker"
    "time"
)

cb := breaker.NewCircuitBreakerWithSettings(gobreaker.Settings{
    Name:        "custom-breaker",
    MaxRequests: 3,                    // 半开状态最多放行 3 个请求
    Interval:    30 * time.Second,     // 统计窗口 30 秒
    Timeout:     10 * time.Second,     // 熔断 10 秒后进入半开
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // 自定义策略：连续失败 3 次即触发熔断
        return counts.ConsecutiveFailures >= 3
    },
    OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
        // 状态变更回调，可用于日志记录或告警
        fmt.Printf("熔断器 [%s] 状态变更: %s → %s\n", name, from, to)
    },
})
```

### 在控制器中使用多个熔断器

不同的外部服务可以使用独立的熔断器实例，互不影响：

```go
type router struct {
    ctl.Controller
    PaymentBreaker *gobreaker.CircuitBreaker  // 支付服务熔断器
    OrderBreaker   *gobreaker.CircuitBreaker  // 订单服务熔断器
}

func NewRouter() *router {
    return &router{
        PaymentBreaker: breaker.NewCircuitBreaker("payment", 3, time.Second*30, time.Second*10),
        OrderBreaker:   breaker.NewCircuitBreaker("order", 5, time.Second*10, time.Second*5),
    }
}
```

## 运行

```shell
go run main.go
```

访问 `http://localhost:8888/healthz`，正常情况下返回：

```json
{
  "code": 0,
  "data": "Request succeeded"
}
```

当下游服务不可用、熔断器触发后，返回 500 错误：

```json
{
  "code": 500,
  "msg": "circuit breaker is open"
}
```

## 最佳实践

1. **为每个外部服务创建独立的熔断器** — 避免一个服务故障导致所有熔断器打开
2. **合理设置 Timeout** — 太短会导致频繁探测，太长会延迟恢复。建议根据下游服务的平均恢复时间设置
3. **合理设置 MaxRequests** — Half-Open 状态下的探测请求数不宜过多，避免对刚恢复的服务造成压力
4. **利用 OnStateChange 回调** — 在状态变更时记录日志或发送告警，便于运维监控
5. **结合重试机制** — 熔断器保护的是调用方，可以在 Job 函数内部加入有限次重试

## 下一步

- [中间件](/guide/中间件.md) — 了解如何在中间件层面集成熔断器
- [依赖注入](/guide/依赖注入.md) — 通过 DI 容器管理熔断器实例
- [sony/gobreaker 文档](https://github.com/sony/gobreaker) — 了解更多高级配置选项
