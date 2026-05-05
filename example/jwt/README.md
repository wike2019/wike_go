# JWT 令牌工具

本节介绍 Wike Go 框架内置的 JWT（JSON Web Token）工具包，基于 Go 泛型实现，支持将任意结构体编码为 Token 并解析还原。

## 概述

`github.com/wike2019/wike_go/v2/func/jwt` 是一个轻量级的 JWT 工具包，底层基于 `golang-jwt/jwt/v4`，使用 Go 1.18+ 泛型特性，让你可以将**任意自定义结构体**作为 Token 的载荷（Payload），无需手动处理 `map[string]interface{}` 或类型断言。

### 核心特点

- **泛型支持** — `Create[T]` / `Parse[T]` 直接操作自定义结构体，类型安全
- **HS256 签名** — 使用 HMAC-SHA256 对称加密算法
- **过期时间** — 内置 `ExpiresAt` 支持，通过 `time.Duration` 灵活设置
- **签发者标识** — 支持 `Issuer` 字段标识令牌来源
- **零配置** — 不依赖框架启动流程，可独立使用

## 项目结构

```
jwt/
├── config/
│   └── config.go      # 配置中心（基于 Viper）
├── config.yaml        # 配置文件
├── db/
│   └── core.db        # SQLite 数据库文件
├── go.mod             # 依赖管理
├── go.sum             # 依赖校验
└── main.go            # 入口文件（JWT 使用示例）
```

## 核心 API

### Create — 生成 Token

```go
func Create[T any](info T, duration time.Duration, SECRET string, Issuer string) (string, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `info` | `T` | 要编码到 Token 中的自定义结构体数据 |
| `duration` | `time.Duration` | Token 有效期（如 `time.Hour * 24`） |
| `SECRET` | `string` | 签名密钥（用于 HS256 加密） |
| `Issuer` | `string` | 签发者标识（如服务名称） |

**返回值：**
- `string` — 签名后的 JWT Token 字符串
- `error` — 签名失败时返回错误

### Parse — 解析 Token

```go
func Parse[T any](token string, SECRET string) (*T, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `token` | `string` | 待解析的 JWT Token 字符串 |
| `SECRET` | `string` | 签名密钥（必须与生成时一致） |

**返回值：**
- `*T` — 解析出的自定义结构体指针
- `error` — Token 无效、已过期或密钥不匹配时返回错误

### InfoData — 内部载荷结构

```go
type InfoData[T any] struct {
    Core T                       // 用户自定义数据
    jwt.RegisteredClaims         // JWT 标准声明（ExpiresAt, Issuer 等）
}
```

## 完整示例

以下是当前项目 `main.go` 的完整代码：

```go
package main

import (
    "fmt"
    "time"

    "github.com/wike2019/wike_go/v2/func/jwt"
)

// 签名密钥
const key = "12345567"

// 自定义载荷结构体
type Info struct {
    Name string
}

func main() {
    // 构造载荷数据
    coreData := Info{
        Name: "我是wike",
    }

    // 生成 Token（有效期 10 秒，签发者 "wike"）
    token, err := jwt.Create[Info](coreData, time.Second*10, key, "wike")
    fmt.Println(token, err)

    // 解析 Token，还原为 Info 结构体
    info, err := jwt.Parse[Info](token, key)
    fmt.Println(info, err)
}
```

**运行输出：**

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDb3JlIjp7Ik5hbWUiOiLmiJHmmK93aWtlIn0sImlzcyI6Indpa2UiLCJleHAiOjE3MTQ4...} <nil>
&{我是wike} <nil>
```

## 使用场景 Demo

### 场景一：用户登录签发 Token

```go
type UserClaims struct {
    UserID   int64
    Username string
    Role     string
}

func Login(username, password string) (string, error) {
    // 验证用户名密码（省略）
    user := UserClaims{
        UserID:   1001,
        Username: username,
        Role:     "admin",
    }

    // 签发 Token，有效期 24 小时
    token, err := jwt.Create[UserClaims](user, time.Hour*24, "your-secret-key", "my-app")
    if err != nil {
        return "", fmt.Errorf("生成Token失败: %w", err)
    }
    return token, nil
}
```

### 场景二：中间件中解析 Token 鉴权

```go
func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        if tokenStr == "" {
            c.JSON(401, gin.H{"code": -1, "msg": "缺少Token"})
            c.Abort()
            return
        }

        // 解析 Token，还原用户信息
        userInfo, err := jwt.Parse[UserClaims](tokenStr, secret)
        if err != nil {
            c.JSON(401, gin.H{"code": -1, "msg": "Token无效或已过期"})
            c.Abort()
            return
        }

        // 将用户信息存入上下文，后续处理函数可直接使用
        c.Set("userID", userInfo.UserID)
        c.Set("username", userInfo.Username)
        c.Set("role", userInfo.Role)
        c.Next()
    }
}
```

### 场景三：不同业务使用不同载荷

```go
// API 密钥载荷
type APIKeyClaims struct {
    AppID       string
    Permissions []string
}

// 生成 API 密钥（有效期 30 天）
func GenerateAPIKey(appID string, perms []string) (string, error) {
    claims := APIKeyClaims{
        AppID:       appID,
        Permissions: perms,
    }
    return jwt.Create[APIKeyClaims](claims, time.Hour*24*30, "api-secret", "api-gateway")
}

// 邮箱验证载荷
type EmailVerifyClaims struct {
    Email  string
    Action string
}

// 生成邮箱验证链接 Token（有效期 1 小时）
func GenerateEmailToken(email string) (string, error) {
    claims := EmailVerifyClaims{
        Email:  email,
        Action: "verify",
    }
    return jwt.Create[EmailVerifyClaims](claims, time.Hour, "email-secret", "user-service")
}
```

### 场景四：结合 Wike Go 框架使用

```go
package main

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/wike2019/wike_app/config"
    "github.com/wike2019/wike_go/v2/core"
    "github.com/wike2019/wike_go/v2/func/ctl"
    "github.com/wike2019/wike_go/v2/func/jwt"
)

const secret = "my-secret-key"

type UserClaims struct {
    UserID   int64
    Username string
}

type authRouter struct {
    ctl.Controller
}

func NewAuthRouter() *authRouter {
    return &authRouter{}
}

func (r authRouter) Path() string {
    return "/auth"
}

func (r authRouter) Build(rg *gin.RouterGroup, GCore *core.GCore) {
    rg.POST("/login", r.login)
    rg.GET("/profile", r.profile)
}

// 登录接口 — 签发 Token
func (r *authRouter) login(c *gin.Context) {
    ctx := r.SetContext(c)
    claims := UserClaims{UserID: 1, Username: "wike"}
    token, err := jwt.Create[UserClaims](claims, time.Hour*24, secret, "my-app")
    if err != nil {
        ctx.Failed(-1, "Token生成失败")
        return
    }
    ctx.OK(gin.H{"token": token})
}

// 需要鉴权的接口 — 解析 Token
func (r *authRouter) profile(c *gin.Context) {
    ctx := r.SetContext(c)
    tokenStr := c.GetHeader("Authorization")
    user, err := jwt.Parse[UserClaims](tokenStr, secret)
    if err != nil {
        ctx.Failed(-1, "Token无效或已过期")
        return
    }
    ctx.OK(gin.H{
        "userID":   user.UserID,
        "username": user.Username,
    })
}

func main() {
    g := core.God()
    g.Provide(config.Config).MountWithEmpty(NewAuthRouter).Run()
}
```

## Token 生命周期

```
Create[T](data, duration, secret, issuer)
    │
    ▼
┌─────────────────────────────┐
│  Header: {"alg":"HS256"}    │
│  Payload: {Core: T, Claims} │
│  Signature: HMAC-SHA256     │
└─────────────────────────────┘
    │
    ▼
  Token 字符串 (eyJhbGci...)
    │
    ▼
Parse[T](token, secret)
    │
    ├── 验证签名 ──── 失败 → error("Token不合法")
    ├── 验证过期 ──── 过期 → error("token is expired")
    └── 验证通过 ──── 返回 *T（原始结构体指针）
```

## 常用有效期参考

| 场景 | 推荐有效期 | 代码 |
|------|-----------|------|
| 短期会话 | 15 分钟 | `time.Minute * 15` |
| 用户登录 | 24 小时 | `time.Hour * 24` |
| 记住登录 | 7 天 | `time.Hour * 24 * 7` |
| API 密钥 | 30 天 | `time.Hour * 24 * 30` |
| 邮箱验证 | 1 小时 | `time.Hour` |
| 密码重置 | 30 分钟 | `time.Minute * 30` |

## 安全建议

1. **密钥管理** — 不要将密钥硬编码在代码中，应使用环境变量或配置文件管理
2. **密钥强度** — 生产环境使用至少 32 字符的随机字符串作为密钥
3. **有效期** — 根据业务场景设置合理的过期时间，避免 Token 长期有效
4. **传输安全** — Token 应通过 HTTPS 传输，避免中间人攻击
5. **存储安全** — 客户端避免将 Token 存储在 localStorage，推荐使用 HttpOnly Cookie

## 运行

```shell
go run main.go
```

## 依赖

```
github.com/wike2019/wike_go/v2 v2.0.16
github.com/golang-jwt/jwt/v4 v4.5.2    # 底层 JWT 库（间接依赖）
```
