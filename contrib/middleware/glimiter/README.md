# Rate Limiter for GoFrame

GoFrame 的速率限制器实现，支持基于内存和 Redis 的限速策略。

## 特性

- ✅ **多种实现**: 支持内存（Memory）和 Redis 两种存储后端
- ✅ **滑动窗口**: Redis 使用精确的滑动窗口算法
- ✅ **并发安全**: 内存版使用 CAS 原子操作，Redis 版使用 Lua 脚本
- ✅ **灵活配置**: 支持自定义键生成、错误处理
- ✅ **标准化**: 符合 RFC 6585 规范的响应头
- ✅ **高性能**: 内存版无额外开销，Redis 版原子操作

## 安装

```bash
go get github.com/gogf/gf/contrib/middleware/glimiter/v2@latest
```

## 快速开始

### 1. 基础使用

```go
import "github.com/gogf/gf/contrib/middleware/glimiter/v2"

// 创建限速器：每分钟最多 100 次请求
limiter := glimiter.NewMemoryLimiter(100, time.Minute)

// 检查是否允许请求
allowed, err := limiter.Allow(ctx, "user:123")
if !allowed {
    // 请求被限速
}
```

### 2. HTTP 中间件

#### 按 IP 限速

```go
s := g.Server()
limiter := glimiter.NewMemoryLimiter(100, time.Minute)

s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(glimiter.MiddlewareByIP(limiter))
    
    group.GET("/users", handler)
})
```

#### 按 API Key 限速

```go
limiter := glimiter.NewMemoryLimiter(1000, time.Hour)

s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(glimiter.MiddlewareByAPIKey(limiter, "X-API-Key"))
    
    group.GET("/data", handler)
})
```

#### 自定义限速逻辑

```go
limiter := glimiter.NewMemoryLimiter(50, time.Minute)

s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(glimiter.Middleware(glimiter.MiddlewareConfig{
        Limiter: limiter,
        KeyFunc: func(r *ghttp.Request) string {
            // 自定义 key：结合 IP 和 User-Agent
            return r.GetClientIp() + ":" + r.UserAgent()
        },
        ErrorHandler: func(r *ghttp.Request) {
            r.Response.WriteStatus(429)
            r.Response.WriteJson(g.Map{
                "error": "Rate limit exceeded",
                "retry": time.Now().Add(time.Minute).Unix(),
            })
        },
    }))
    
    group.GET("/resource", handler)
})
```

### 3. 使用 Redis 限速器

```go
import (
    "github.com/gogf/gf/contrib/middleware/glimiter/v2"
    "github.com/gogf/gf/v2/database/gredis"
)

// 创建 Redis 连接
redis, err := gredis.New(&gredis.Config{
    Address: "127.0.0.1:6379",
})

// 创建 Redis 限速器
limiter := glimiter.NewRedisLimiter(redis, 100, time.Minute)

// 使用方式与内存限速器相同
allowed, err := limiter.Allow(ctx, "user:123")
```

## 核心接口

### Limiter 接口

```go
type Limiter interface {
    // 棜查是否允许单个请求
    Allow(ctx context.Context, key string) (bool, error)
    
    // 检查是否允许 N 个请求
    AllowN(ctx context.Context, key string, n int) (bool, error)
    
    // 阻塞直到允许请求
    Wait(ctx context.Context, key string) error
    
    // 获取限制配置
    GetLimit() int
    GetWindow() time.Duration
    
    // 获取剩余配额
    GetRemaining(ctx context.Context, key string) (int, error)
    
    // 重置限制
    Reset(ctx context.Context, key string) error
}
```

## 使用场景

### 1. 多层限速

针对不同时间窗口设置多层限制，防止突发流量和长期滥用：

```go
// 第一层：防突发（每秒）
burstLimiter := glimiter.NewMemoryLimiter(10, time.Second)

// 第二层：常规限制（每分钟）
normalLimiter := glimiter.NewMemoryLimiter(100, time.Minute)

// 第三层：长期限制（每小时）
hourlyLimiter := glimiter.NewMemoryLimiter(1000, time.Hour)

s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(
        glimiter.MiddlewareByIP(burstLimiter),
        glimiter.MiddlewareByIP(normalLimiter),
        glimiter.MiddlewareByIP(hourlyLimiter),
    )
    
    group.GET("/search", handler)
})
```

### 2. 路由级限速

不同的 API 路由使用不同的限速策略：

```go
s := g.Server()

// 公开 API：宽松限制
s.Group("/public", func(group *ghttp.RouterGroup) {
    publicLimiter := glimiter.NewMemoryLimiter(100, time.Minute)
    group.Middleware(glimiter.MiddlewareByIP(publicLimiter))
    group.GET("/info", handler)
})

// 认证 API：中等限制
s.Group("/auth", func(group *ghttp.RouterGroup) {
    authLimiter := glimiter.NewMemoryLimiter(5, time.Minute)
    group.Middleware(glimiter.MiddlewareByIP(authLimiter))
    group.POST("/login", handler)
})

// 敏感操作：严格限制
s.Group("/admin", func(group *ghttp.RouterGroup) {
    adminLimiter := glimiter.NewMemoryLimiter(10, time.Hour)
    group.Middleware(glimiter.MiddlewareByIP(adminLimiter))
    group.POST("/delete", handler)
})
```

### 3. 按用户限速

```go
limiter := glimiter.NewMemoryLimiter(1000, time.Hour)

middleware := glimiter.MiddlewareByUser(limiter, func(r *ghttp.Request) string {
    // 从上下文中获取用户 ID
    user := r.GetCtxVar("user").String()
    return user
})

s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(middleware)
    group.GET("/profile", handler)
})
```

### 4. 直接使用限速器

不通过中间件，在业务代码中直接使用：

```go
limiter := glimiter.NewMemoryLimiter(10, time.Minute)

func ProcessTask(ctx context.Context, taskID string) error {
    // 检查是否允许处理
    allowed, err := limiter.Allow(ctx, "task:"+taskID)
    if err != nil {
        return err
    }
    
    if !allowed {
        return errors.New("rate limit exceeded")
    }
    
    // 执行任务
    return doTask(taskID)
}
```

### 5. 等待配额可用

```go
limiter := glimiter.NewMemoryLimiter(5, time.Second)

func SendRequest(ctx context.Context) error {
    // 阻塞等待，直到配额可用
    if err := limiter.Wait(ctx, "api-call"); err != nil {
        return err
    }
    
    // 发送请求
    return makeAPICall()
}
```

## 响应头

限速中间件会自动设置以下 HTTP 响应头：

| 响应头 | 说明 |
|--------|------|
| `X-RateLimit-Limit` | 时间窗口内的最大请求数 |
| `X-RateLimit-Remaining` | 剩余可用请求数 |
| `X-RateLimit-Reset` | 限速重置时间（Unix 时间戳） |

## 最佳实践

### 1. 合理设置时间窗口

- **短时间窗口**（1 分钟内）：适合内存限速器，性能最优
- **长时间窗口**（1 小时以上）：建议使用 Redis 限速器，支持分布式

### 2. 使用多层限速

结合不同时间窗口的限速策略，既能防止突发流量，又能限制长期滥用：

- 第一层：秒级限速，防止突发攻击
- 第二层：分钟级限速，常规使用限制
- 第三层：小时级限速，长期配额管理

### 3. 区分不同场景

根据 API 的敏感程度和重要性设置不同的限速策略：

- **公开 API**：宽松限制，提供良好用户体验
- **认证 API**：中等限制，防止暴力破解
- **敏感操作**：严格限制，保护关键功能

### 4. 提供友好的错误信息

自定义错误处理器，告知用户何时可以重试：

```go
ErrorHandler: func(r *ghttp.Request) {
    r.Response.WriteStatus(429)
    r.Response.WriteJson(g.Map{
        "error": "Rate limit exceeded",
        "message": "You have exceeded the rate limit. Please try again later.",
        "retry_after": limiter.GetWindow().Seconds(),
    })
}
```

### 5. 监控和告警

在生产环境中，建议监控限速器的使用情况：

```go
// 定期检查剩余配额
remaining, _ := limiter.GetRemaining(ctx, key)
if remaining < 10 {
    // 发送告警
    log.Warn("Rate limit nearly exhausted", "key", key, "remaining", remaining)
}
```

## 性能考虑

### 内存限速器

- **优点**: 极高性能，无网络开销
- **缺点**: 单机限速，不适合分布式环境
- **适用**: 单体应用、短时间窗口

### Redis 限速器

- **优点**: 支持分布式，数据持久化
- **缺点**: 有网络延迟开销
- **适用**: 分布式应用、长时间窗口

## 并发安全

### 内存限速器

使用 CAS（Compare-And-Swap）原子操作确保并发安全：

### Redis 限速器

使用 Lua 脚本确保原子操作：


## 常见问题

### Q: 如何在分布式环境中使用？

A: 使用 `RedisLimiter` 替代 `MemoryLimiter`，确保所有服务实例共享同一个 Redis。

### Q: 如何实现动态调整限速配置？

A: 可以在运行时创建新的 `Limiter` 实例并更新中间件配置。

### Q: 时间窗口如何工作？

A: 
- **内存限速器**: 使用 gcache 的自动过期功能
- **Redis 限速器**: 使用滑动窗口算法，精确控制时间范围

### Q: 如何避免限速器成为性能瓶颈？

A: 
1. 短时间窗口使用内存限速器
2. 合理设置限速配额
3. 使用多层限速而非单一严格限制