# 工具集与实用功能 (Other)

本节汇总展示 Wike Go 框架提供的各类实用工具和功能模块，涵盖并发控制、数据结构、加密、重试、一致性哈希、系统信息、Redis 队列/分布式锁、MD5 校验、对象拷贝等。

## 概述

Wike Go 框架在 `utils` 和 `func` 包下内置了大量开箱即用的工具函数，帮助开发者快速实现常见的业务需求，无需引入额外第三方库。

## 项目结构

```
other/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── go.mod             # 依赖管理
├── go.sum             # 依赖校验
└── main.go            # 入口文件（工具集使用示例）
```

## 功能目录

| 功能模块 | 包路径 | 说明 |
|----------|--------|------|
| 并发安全 Map | `utils.NewMap[T]()` | 泛型线程安全 Map |
| 密码加密 | `utils.PasswordHash / PasswordVerify` | bcrypt 密码哈希与验证 |
| 集合去重 | `utils.NewSet()` | 支持任意类型的去重集合 |
| 随机数据 | `utils.RandomString / GetRandomNum` | 随机字符串和随机数生成 |
| 并发池 | `ants_service.NewPool(n)` | 基于 ants 的 goroutine 池 |
| 一致性哈希 | `hash.New()` | 一致性哈希环，用于负载均衡 |
| 系统信息 | `os.InitOS / InitCPU / InitRAM / InitDisk` | 获取系统运行时信息 |
| 时间解析 | `utils.ParseDuration` | 支持天级别的时间解析 |
| Redis 队列 | `utils.NewRedisQueue(client)` | 基于 Redis 的消息队列 |
| Redis 分布式锁 | `utils.NewRedisLock(client)` | 基于 Redis 的分布式锁 |
| MD5 工具 | `utils.MD5V / CheckMd5` | MD5 计算与校验 |
| 对象拷贝 | `utils.CopyProperties` | 结构体字段自动拷贝 |
| 重试机制 | `retry.NewRetry(fn)` | 可配置次数和延迟的重试器 |

---

## 1. 并发安全 Map

基于泛型的线程安全 Map，适用于多 goroutine 并发读写场景。

```go
import "github.com/wike2019/wike_go/v2/utils"

// 创建一个 value 类型为 string 的并发安全 Map
mapstring := utils.NewMap[string]()

// 写入
mapstring.Set("111", "2222")
mapstring.Set("222", "3333")

// 读取
value := mapstring.Get("111")  // "2222"

// 获取所有 key 和 value
keys := mapstring.Keys()       // ["111", "222"]
values := mapstring.Values()   // ["2222", "3333"]
```

### API

| 方法 | 说明 |
|------|------|
| `NewMap[T]()` | 创建泛型并发安全 Map |
| `Set(key, value)` | 设置键值对 |
| `Get(key)` | 获取值 |
| `Keys()` | 获取所有 key |
| `Values()` | 获取所有 value |

---

## 2. 密码加密与验证

基于 bcrypt 的密码哈希，适用于用户注册和登录场景。

```go
import "github.com/wike2019/wike_go/v2/utils"

// 注册时：对明文密码加密
hashed, err := utils.PasswordHash("123456")
// hashed = "$2a$10$..."

// 登录时：验证密码是否匹配
isValid := utils.PasswordVerify("123456", hashed)
// isValid = true

isValid = utils.PasswordVerify("wrong_password", hashed)
// isValid = false
```

### API

| 方法 | 说明 |
|------|------|
| `PasswordHash(password string) (string, error)` | 对明文密码进行 bcrypt 加密 |
| `PasswordVerify(password, hash string) bool` | 验证明文密码与哈希是否匹配 |

---

## 3. 集合去重 (Set)

支持任意类型的去重集合，相同值只保留一份。

```go
import "github.com/wike2019/wike_go/v2/utils"

// 字符串去重
set := utils.NewSet()
set.Add("111")
set.Add("222")
set.Add("333")
set.Add("222")  // 重复，不会添加

var str []string
set.List(&str)
// str = ["111", "222", "333"]

// 结构体去重（按值比较）
type Test struct {
    Name string
}

set2 := utils.NewSet()
set2.Add(&Test{Name: "a1"})
set2.Add(&Test{Name: "a1"})  // 相同值，去重
set2.Add(&Test{Name: "a2"})

var obj []*Test
set2.List(&obj)
// obj 长度为 2: [{Name:"a1"}, {Name:"a2"}]
```

### API

| 方法 | 说明 |
|------|------|
| `NewSet()` | 创建新的集合 |
| `Add(item)` | 添加元素（自动去重） |
| `List(dest)` | 将集合导出为切片 |

---

## 4. 随机数据生成

```go
import "github.com/wike2019/wike_go/v2/utils"

// 生成指定长度的随机字符串
str := utils.RandomString(10)  // 如 "aB3kLm9xPq"

// 生成指定范围的随机整数
num := utils.GetRandomNum(1, 99)  // 1 到 99 之间的随机数
```

### API

| 方法 | 说明 |
|------|------|
| `RandomString(length int) string` | 生成指定长度的随机字符串 |
| `GetRandomNum(min, max int) int` | 生成 [min, max] 范围内的随机整数 |

---

## 5. 并发池 (Goroutine Pool)

基于 ants 封装的协程池，限制并发 goroutine 数量，防止资源耗尽。

```go
import "github.com/wike2019/wike_go/v2/func/ants_service"

// 1. 创建池，限制最多 10 个 goroutine 并发
pool, err := ants_service.NewPool(10)

// 2. 设置总任务数（内部调用 WaitGroup.Add）
pool.SetTotal(100)

// 3. 提交任务
for i := 0; i < 100; i++ {
    pool.Submit(func() error {
        // 业务逻辑
        return nil
    })
}

// 4. 等待所有任务完成
pool.Wait()

// 5. 检查结果
fmt.Println(pool.Ok, pool.Fail, pool.Error())
```

### API

| 方法 | 说明 |
|------|------|
| `NewPool(size int) (*Pool, error)` | 创建指定大小的协程池 |
| `SetTotal(n int)` | 设置总任务数 |
| `Submit(fn func() error)` | 提交任务到池中 |
| `Wait()` | 阻塞等待所有任务完成 |
| `Ok` | 成功任务数 |
| `Fail` | 失败任务数 |
| `Error()` | 返回所有错误 |

---

## 6. 一致性哈希 (Consistent Hash)

用于分布式系统中的负载均衡，将请求按 key 路由到固定节点。

```go
import "github.com/wike2019/wike_go/v2/func/hash"

h := hash.New()

// 添加节点（比如服务器地址）
h.Add("192.168.1.1:8080")
h.Add("192.168.1.2:8080")
h.Add("192.168.1.3:8080")

// 根据 key 获取对应的节点
node, err := h.Get("user_123")
if err != nil {
    panic(err)
}
fmt.Println(node) // 输出某个节点地址

// 相同的 key 总是路由到相同的节点
node2, _ := h.Get("user_123")
// node == node2
```

### API

| 方法 | 说明 |
|------|------|
| `New() *Hash` | 创建一致性哈希环 |
| `Add(node string)` | 添加节点 |
| `Get(key string) (string, error)` | 根据 key 获取对应节点 |

### 使用场景

- 分布式缓存路由（如 Redis 集群分片）
- 微服务负载均衡
- 分布式任务调度

---

## 7. 系统信息获取

获取当前服务器的运行时信息，适用于监控面板和健康检查。

```go
import "github.com/wike2019/wike_go/v2/func/os"

// 获取 Go 运行时信息：系统类型、CPU 核数、编译器、Go 版本、goroutine 数
fmt.Println(os.InitOS())

// 获取每个 CPU 核心的使用率百分比和物理核心数
fmt.Println(os.InitCPU())

// 获取内存已用 MB、总量 MB、使用百分比
fmt.Println(os.InitRAM())

// 获取指定挂载点的磁盘用量（MB/GB）和使用百分比
fmt.Println(os.InitDisk("./"))
```

### API

| 方法 | 说明 |
|------|------|
| `InitOS()` | 系统类型、CPU 核数、编译器、Go 版本、goroutine 数 |
| `InitCPU()` | CPU 核心使用率和物理核心数 |
| `InitRAM()` | 内存已用/总量/使用百分比 |
| `InitDisk(path string)` | 指定路径的磁盘用量和使用百分比 |

---

## 8. 友好时间解析

支持天级别（`d`）的时间字符串解析，标准库 `time.ParseDuration` 不支持天。

```go
import "github.com/wike2019/wike_go/v2/utils"

duration := utils.ParseDuration("1d5h20m")
// 等价于 24h + 5h + 20m = 29h20m
```

---

## 9. Redis 消息队列

基于 Redis List 实现的轻量级消息队列，支持任意结构体的序列化/反序列化。

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/wike2019/wike_go/v2/utils"
)

type Order struct {
    ID     int
    UserID int
    Amount float64
}

client := redis.NewClient(&redis.Options{
    Addr: "192.168.3.2:6379",
})

queue := utils.NewRedisQueue(client)

// ========== 生产者：推送消息 ==========
order := Order{ID: 1001, UserID: 42, Amount: 99.9}
err := queue.Push("order_queue", order)
if err != nil {
    log.Fatal(err)
}
fmt.Println("消息已推送")

// ========== 消费者：阻塞等待并消费消息 ==========
var received Order
err = queue.Pop("order_queue", &received)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("收到订单: ID=%d, 金额=%.1f\n", received.ID, received.Amount)
```

### API

| 方法 | 说明 |
|------|------|
| `NewRedisQueue(client) *RedisQueue` | 创建 Redis 队列实例 |
| `Push(key string, data interface{}) error` | 推送消息到队列 |
| `Pop(key string, dest interface{}) error` | 阻塞消费消息 |

---

## 10. Redis 分布式锁

基于 Redis 实现的分布式互斥锁，适用于防止并发重复操作。

```go
import (
    "context"
    "github.com/redis/go-redis/v9"
    "github.com/wike2019/wike_go/v2/utils"
)

client := redis.NewClient(&redis.Options{
    Addr: "192.168.3.2:6379",
})

lock := utils.NewRedisLock(client)

// 加锁，key="order_lock"，TTL=10秒
lockData, err := lock.Lock("order_lock", 10, nil)
if err != nil {
    log.Fatal("加锁失败:", err)
}
fmt.Println("加锁成功")

// ===== 执行需要互斥的业务逻辑 =====
fmt.Println("正在处理订单...")
// 比如：扣库存、创建订单等

// 释放锁
err = lockData.Release(context.Background())
if err != nil {
    log.Println("释放锁失败:", err)
}
fmt.Println("锁已释放")
```

### API

| 方法 | 说明 |
|------|------|
| `NewRedisLock(client) *RedisLock` | 创建分布式锁实例 |
| `Lock(key string, ttl int, ctx) (*Lock, error)` | 加锁，返回锁对象 |
| `lockData.Release(ctx) error` | 释放锁 |

### 使用场景

- 防止订单重复支付
- 库存扣减互斥
- 定时任务防重复执行

---

## 11. MD5 工具

用于数据完整性校验，常见于文件分片上传场景。

```go
import "github.com/wike2019/wike_go/v2/utils"

// 1. 计算 MD5 值
data := []byte("hello world")
md5Str := utils.MD5V(data)
fmt.Println(md5Str) // 输出: 5eb63bbbe01eeed093cb22bb8f5acdc3

// 2. 校验分片完整性（文件分片上传场景）
chunk := []byte("这是文件的第一个分片内容")
expectedMd5 := utils.MD5V(chunk) // 上传前客户端计算的 MD5

// 服务端收到分片后校验
ok := utils.CheckMd5(chunk, expectedMd5)
fmt.Println("分片校验:", ok) // true

// 如果数据被篡改或传输损坏
ok = utils.CheckMd5([]byte("损坏的数据"), expectedMd5)
fmt.Println("分片校验:", ok) // false
```

### API

| 方法 | 说明 |
|------|------|
| `MD5V(data []byte) string` | 计算字节数组的 MD5 值 |
| `CheckMd5(data []byte, expected string) bool` | 校验数据的 MD5 是否与期望值一致 |

---

## 12. 对象属性拷贝

将一个结构体的字段值拷贝到另一个结构体，自动匹配同名字段，忽略目标中不存在的字段。适用于 PO → VO 转换。

```go
import "github.com/wike2019/wike_go/v2/utils"

// 数据库模型（包含敏感字段）
type UserPO struct {
    ID       int
    Name     string
    Password string
    Age      int
    Email    string
}

// 返回给前端的 VO（不含 Password）
type UserVO struct {
    ID    int
    Name  string
    Age   int
    Email string
}

user := UserPO{
    ID:       1,
    Name:     "张三",
    Password: "secret123",
    Age:      25,
    Email:    "zhangsan@example.com",
}

// 拷贝到 VO，自动忽略 Password 字段（VO 中没有）
var vo UserVO
err := utils.CopyProperties(&vo, user)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%+v\n", vo)
// 输出: {ID:1 Name:张三 Age:25 Email:zhangsan@example.com}
```

### API

| 方法 | 说明 |
|------|------|
| `CopyProperties(dest, src interface{}) error` | 将 src 的同名字段拷贝到 dest |

---

## 13. 重试机制 (Retry)

可配置重试次数和延迟的重试器，适用于网络请求、外部服务调用等不稳定操作。

```go
import (
    "time"
    "github.com/wike2019/wike_go/v2/func/retry"
)

// 定义需要重试的函数
func unstableJob() error {
    // 可能失败的操作，如 HTTP 请求、数据库连接等
    return nil
}

// 创建重试器
re := retry.NewRetry(unstableJob)

// 设置最大重试次数
re.SetTimes(5)

// 设置每次重试的间隔
re.SetDelay(time.Second * 1)

// 执行（失败会自动重试）
err := re.Do()
if err != nil {
    fmt.Println("重试5次后仍然失败:", err)
}
```

### API

| 方法 | 说明 |
|------|------|
| `NewRetry(fn func() error) *Retry` | 创建重试器 |
| `SetTimes(n int)` | 设置最大重试次数 |
| `SetDelay(d time.Duration)` | 设置重试间隔 |
| `Do() error` | 执行函数，失败自动重试 |

### 使用场景

- HTTP 请求超时重试
- 数据库连接失败重试
- 第三方 API 调用重试
- 消息发送失败重试

---

## 完整示例代码

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ants_service"
    "github.com/wike2019/wike_go/v2/func/ctl"
    "github.com/wike2019/wike_go/v2/func/hash"
    "github.com/wike2019/wike_go/v2/func/os"
    "github.com/wike2019/wike_go/v2/func/retry"
    "github.com/wike2019/wike_go/v2/utils"
)

func main() {
    // 并发安全 Map
    mapstring := utils.NewMap[string]()
    mapstring.Set("111", "2222")
    fmt.Println(mapstring.Get("111"))

    // 密码加密
    hashed, _ := utils.PasswordHash("123456")
    fmt.Println(utils.PasswordVerify("123456", hashed))

    // 集合去重
    set := utils.NewSet()
    set.Add("aaa")
    set.Add("bbb")
    set.Add("aaa")

    // 随机数据
    fmt.Println(utils.RandomString(10))
    fmt.Println(utils.GetRandomNum(1, 99))

    // 并发池
    pool, _ := ants_service.NewPool(10)
    pool.SetTotal(50)
    for i := 0; i < 50; i++ {
        pool.Submit(func() error { return nil })
    }
    pool.Wait()
    fmt.Println(pool.Ok, pool.Fail)

    // 一致性哈希
    h := hash.New()
    h.Add("node1")
    h.Add("node2")
    node, _ := h.Get("key1")
    fmt.Println(node)

    // 系统信息
    fmt.Println(os.InitOS())
    fmt.Println(os.InitRAM())

    // 时间解析
    fmt.Println(utils.ParseDuration("1d5h20m"))

    // 启动服务
    g := core.God()
    g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
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

## 依赖

```
github.com/wike2019/wike_go/v2 v2.0.16
github.com/gin-gonic/gin v1.12.0
github.com/redis/go-redis/v9 v9.19.0
github.com/spf13/viper v1.21.0
```
