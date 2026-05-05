# 优雅停机 (Graceful Shutdown)

本节介绍如何在 Wike Go 框架中使用 `.Stop()` 实现服务的优雅停机，确保在进程退出前完成资源清理工作。

## 概述

在生产环境中，服务可能因部署更新、手动重启或系统信号（如 `SIGINT`、`SIGTERM`）而终止。如果直接强制退出，可能导致：

- 数据库连接未正确关闭，造成连接泄漏
- 正在处理的请求被中断，返回异常响应
- 缓存数据未持久化，导致数据丢失
- 文件句柄未释放，临时文件残留

Wike Go 框架通过 `.Stop()` 方法提供了优雅停机钩子，让你在服务关闭前执行自定义的清理逻辑。

## 项目结构

```
stop/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── go.mod             # 依赖管理
├── go.sum             # 依赖校验
└── main.go            # 入口文件（优雅停机示例）
```

## 核心 API

### Stop 方法签名

```go
func (g *God) Stop(fn func() error) *God
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `fn` | `func() error` | 停机时执行的清理函数，返回 error 表示清理是否成功 |

**特点：**
- 支持链式调用，可注册多个 Stop 钩子
- 在收到系统终止信号时自动触发
- 清理函数按注册顺序依次执行
- 返回 `*God` 实例，可继续链式调用其他方法

## 完整示例

以下是当前项目 `main.go` 的完整代码：

```go
package main

import (
    "fmt"

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
    g.Provide(config.Config).MountWithEmpty(NewRouter).Stop(func() error {
        fmt.Println("这里做全局清理")
        return nil
    }).Run()
}
```

## 启动流程解析

```
core.God() → Provide(配置) → MountWithEmpty(路由) → Stop(清理函数) → Run()
```

| 方法 | 作用 |
|------|------|
| `core.God()` | 创建框架实例 |
| `Provide(fn...)` | 注册配置到 DI 容器 |
| `MountWithEmpty(ctor)` | 挂载路由控制器 |
| `Stop(fn)` | 注册优雅停机钩子 |
| `Run()` | 启动 HTTP 服务 |

## 生命周期

```
启动阶段                          运行阶段              停机阶段
   │                                │                    │
   ▼                                ▼                    ▼
God() → Provide → Mount → Run() → 处理请求... → 收到终止信号
                                                         │
                                                         ▼
                                                  停止接收新请求
                                                         │
                                                         ▼
                                                  等待进行中的请求完成
                                                         │
                                                         ▼
                                                  执行 Stop 钩子
                                                         │
                                                         ▼
                                                    进程退出
```

## 使用场景 Demo

### 场景一：关闭数据库连接

```go
func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewRouter).
        Stop(func() error {
            fmt.Println("正在关闭数据库连接...")
            // 假设 db 是全局数据库实例
            // sqlDB, _ := db.DB()
            // return sqlDB.Close()
            return nil
        }).
        Run()
}
```

### 场景二：注册多个清理钩子

```go
func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewRouter).
        Stop(func() error {
            fmt.Println("1. 关闭数据库连接")
            // db.Close()
            return nil
        }).
        Stop(func() error {
            fmt.Println("2. 关闭 Redis 连接")
            // redisClient.Close()
            return nil
        }).
        Stop(func() error {
            fmt.Println("3. 刷新日志缓冲区")
            // logger.Sync()
            return nil
        }).
        Run()
}
```

### 场景三：结合定时任务的清理

```go
func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewRouter).
        Cron("*/10 * * * * *", func() {
            fmt.Println("执行定时数据同步...")
        }, "数据同步", false).
        Stop(func() error {
            fmt.Println("停机前：执行最后一次数据同步")
            // syncData()
            fmt.Println("停机前：关闭数据库连接")
            // db.Close()
            return nil
        }).
        Run()
}
```

### 场景四：清理临时文件和缓存

```go
func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewRouter).
        Stop(func() error {
            fmt.Println("清理临时上传文件...")
            err := os.RemoveAll("/tmp/uploads")
            if err != nil {
                fmt.Printf("清理临时文件失败: %v\n", err)
                return err
            }
            fmt.Println("临时文件清理完成")
            return nil
        }).
        Run()
}
```

### 场景五：通知外部服务下线

```go
func main() {
    g := core.God()
    g.Provide(config.Config).
        MountWithEmpty(NewRouter).
        Stop(func() error {
            fmt.Println("通知注册中心：服务下线")
            // 从服务注册中心注销
            // registry.Deregister(serviceID)

            fmt.Println("通知监控系统：服务即将停止")
            // 发送告警或通知
            // alertManager.Send("服务正在停机维护")

            return nil
        }).
        Run()
}
```

## 常见清理任务清单

| 清理项 | 说明 | 优先级 |
|--------|------|--------|
| 数据库连接 | 关闭 MySQL/PostgreSQL/SQLite 连接池 | 高 |
| Redis 连接 | 关闭 Redis 客户端连接 | 高 |
| 消息队列 | 关闭 Kafka/RabbitMQ 生产者和消费者 | 高 |
| 日志刷新 | 将缓冲区中的日志写入磁盘 | 高 |
| 临时文件 | 清理上传的临时文件 | 中 |
| 服务注销 | 从注册中心（Consul/Etcd）注销 | 中 |
| 定时任务 | 等待当前执行中的任务完成 | 中 |
| 外部通知 | 通知监控系统或告警平台 | 低 |
| 统计上报 | 上报最后一批统计数据 | 低 |

## 注意事项

1. **清理函数应尽快完成** — 避免在 Stop 钩子中执行耗时过长的操作，否则可能被系统强制终止
2. **错误处理** — Stop 函数返回 error 时应记录日志，便于排查停机异常
3. **资源释放顺序** — 注册多个 Stop 钩子时，注意依赖关系（如先关闭业务层，再关闭数据库层）
4. **幂等性** — 清理函数应具备幂等性，避免重复执行导致异常

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

启动后按 `Ctrl+C` 发送终止信号，观察控制台输出清理日志：

```
^C
这里做全局清理
```

## 依赖

```
github.com/wike2019/wike_go/v2 v2.0.16
github.com/gin-gonic/gin v1.12.0
github.com/spf13/viper v1.21.0
```
