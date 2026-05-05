# Bloom Filter（布隆过滤器）

## 简介

`bloomfilter` 包是对 [hugh2632/bloomfilter](https://github.com/hugh2632/bloomfilter) 的轻量封装，提供了一个简洁易用的布隆过滤器实现。布隆过滤器是一种空间效率极高的概率型数据结构，用于判断一个元素是否存在于集合中。

**核心特点：**

- 判断元素**不存在**时，结果是 100% 准确的
- 判断元素**存在**时，有一定的误判率（false positive）
- 空间效率远高于传统的 Set / Map 数据结构
- 不支持删除操作

**典型应用场景：**

- 缓存穿透防护：在查询数据库前先用布隆过滤器判断 key 是否存在，避免大量无效查询打到数据库
- 去重判断：URL 去重、爬虫已访问页面判断、消息去重
- 黑名单过滤：垃圾邮件地址、恶意 IP 快速过滤

---

## API 说明

### 结构体

```go
type Bloom struct {
    Bloom bloomfilter.IFilter
}
```

内部持有一个 `bloomfilter.IFilter` 接口实例，底层使用内存存储的位数组。

### 方法一览

| 方法 | 签名 | 说明 |
|------|------|------|
| `NewBloom` | `func NewBloom() *Bloom` | 创建默认布隆过滤器（10240 字节位数组 + 默认哈希函数） |
| `Add` | `func (this *Bloom) Add(key string)` | 将元素添加到过滤器中 |
| `Exists` | `func (this *Bloom) Exists(key string) bool` | 判断元素是否可能存在。返回 `false` 表示一定不存在，返回 `true` 表示可能存在 |
| `Clear` | `func (this *Bloom) Clear() error` | 清空过滤器中的所有数据，重置为初始状态 |

---

## 使用示例

### 基础用法

```go
package main

import (
    "fmt"
    "github.com/wike2019/wike_go/v2/func/bloomfilter"
)

func main() {
    // 1. 创建布隆过滤器
    bf := bloomfilter.NewBloom()

    // 2. 添加元素
    bf.Add("apple")
    bf.Add("banana")
    bf.Add("cherry")

    // 3. 判断元素是否存在
    fmt.Println(bf.Exists("apple"))  // true  — 已添加，一定返回 true
    fmt.Println(bf.Exists("grape"))  // false — 未添加，大概率返回 false

    // 4. 清空过滤器
    err := bf.Clear()
    if err != nil {
        fmt.Println("清空失败:", err)
        return
    }
    fmt.Println(bf.Exists("apple"))  // false — 清空后不再存在
}
```

### 缓存穿透防护

```go
package main

import (
    "fmt"
    "github.com/wike2019/wike_go/v2/func/bloomfilter"
)

// 模拟数据库中已有的 key
var dbKeys = []string{"user:1001", "user:1002", "user:1003", "order:5001", "order:5002"}

func main() {
    // 初始化布隆过滤器，预热数据库中已有的 key
    bf := bloomfilter.NewBloom()
    for _, key := range dbKeys {
        bf.Add(key)
    }

    // 模拟查询请求
    queryKeys := []string{"user:1001", "user:9999", "order:5001", "order:0000"}

    for _, key := range queryKeys {
        if !bf.Exists(key) {
            // 布隆过滤器判断不存在 → 一定不存在，直接返回，不查数据库
            fmt.Printf("[%s] 布隆过滤器拦截：key 不存在，跳过数据库查询\n", key)
            continue
        }
        // 布隆过滤器判断可能存在 → 继续查数据库确认
        fmt.Printf("[%s] 布隆过滤器放行：可能存在，查询数据库...\n", key)
    }
}
```

**输出：**

```
[user:1001] 布隆过滤器放行：可能存在，查询数据库...
[user:9999] 布隆过滤器拦截：key 不存在，跳过数据库查询
[order:5001] 布隆过滤器放行：可能存在，查询数据库...
[order:0000] 布隆过滤器拦截：key 不存在，跳过数据库查询
```

### URL 去重（爬虫场景）

```go
package main

import (
    "fmt"
    "github.com/wike2019/wike_go/v2/func/bloomfilter"
)

func main() {
    visited := bloomfilter.NewBloom()

    urls := []string{
        "https://example.com/page1",
        "https://example.com/page2",
        "https://example.com/page1", // 重复
        "https://example.com/page3",
        "https://example.com/page2", // 重复
    }

    for _, url := range urls {
        if visited.Exists(url) {
            fmt.Printf("[跳过] 已访问: %s\n", url)
            continue
        }
        visited.Add(url)
        fmt.Printf("[抓取] 新页面: %s\n", url)
    }
}
```

**输出：**

```
[抓取] 新页面: https://example.com/page1
[抓取] 新页面: https://example.com/page2
[跳过] 已访问: https://example.com/page1
[抓取] 新页面: https://example.com/page3
[跳过] 已访问: https://example.com/page2
```

---

## 注意事项

1. **误判率（False Positive）**：布隆过滤器存在一定的误判率。`Exists` 返回 `true` 时，元素**可能**不存在；返回 `false` 时，元素**一定**不存在。业务逻辑中需要对 `true` 的情况做二次确认。

2. **不支持删除**：标准布隆过滤器不支持删除元素。如果需要删除功能，考虑使用 Counting Bloom Filter。

3. **容量规划**：默认使用 `10240` 字节（约 81920 bit）的位数组。当存储的元素数量过多时，误判率会显著上升。如果需要存储大量元素，建议调整位数组大小或使用多个过滤器实例。

4. **线程安全**：当前实现未做并发保护。在多 goroutine 环境下使用时，需要自行加锁（`sync.Mutex` 或 `sync.RWMutex`）。
