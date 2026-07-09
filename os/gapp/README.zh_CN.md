# gapp

多服务器应用级别的生命周期管理。

## 介绍

`gapp` 包提供了统一的 `Server` 接口和 `App` 结构体，用于协调多个服务器的启动和关闭，包括优雅关闭的信号处理。同时支持通过 `Option` 机制在服务器启动前进行结构化的应用初始化。

## Context 传递

- `Boot`、`Start`、`Stop`、`Run` 均接收 `context.Context`，用于链路追踪与取消传播。
- 参数为 nil 时使用 `gctx.GetInitCtx()`，与框架内其他入口（如 `ghttp.Server.Start`）一致。
- 首次 `Boot` 成功后，其后 `Start`/`Stop` 若传入 nil，会沿用该次 Boot 所使用的规范化 context。
- `Start` 并发启动服务器时，若外层 context 先结束，已对成功启动的服务器按与原启动失败一致的策略回滚关闭。
- `Run` 在收到 OS 关闭信号 **或** root context 结束时触发优雅关闭；实际的 `Stop` 使用 `gctx.NeverDone(root)`，在保留链路元数据的同时，避免 teardown 被子 context 的超时或取消提前打断。

## Server 接口

```go
type Server interface {
    Start() error
    Stop(graceful bool) error
}
```

- `Start()`：以非阻塞方式启动服务器。
- `Stop(graceful)`：停止服务器。当 `graceful` 为 `true` 时，服务器等待正在处理的请求完成后再关闭。当 `graceful` 为 `false` 时，服务器被强制关闭。

## 适配器

该包为内置服务器类型提供了适配器构造函数：

| 构造函数 | 包装类型 |
|---|---|
| `NewHTTPServerAdapter(s *ghttp.Server)` | HTTP 服务器 |
| `NewTCPServerAdapter(s *gtcp.Server)` | TCP 服务器 |
| `NewUDPServerAdapter(s *gudp.Server)` | UDP 服务器 |

对于 gRPC 服务器，请使用 `contrib/rpc/grpcx` 包中的 `grpcx.NewGappServerAdapter(s *grpcx.GrpcServer)`。

## 启动 / 初始化

`Option` 类型允许在服务器启动前注册初始化逻辑，为数据库、缓存、外部连接等应用级别的初始化提供结构化的方式。

### Option 接口

```go
type Option interface {
    Apply(ctx context.Context, app *App) (func(ctx context.Context), error)
}
```

`Apply` 方法在 `Boot()` 期间执行。它可以返回一个可选的清理函数，该函数将在 `Stop()` 期间按注册的逆序调用。

### Option 构造函数

| 构造函数 | 描述 |
|---|---|
| `NewOption(f)` | 一次性初始化，无清理函数 |
| `NewOptionWithHook(f)` | 初始化并返回可选的清理函数 |

### 使用示例

```go
package main

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/os/gapp"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    httpServer := g.Server()
    httpServer.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("hello")
    })

    app := g.App(gapp.NewHTTPServerAdapter(httpServer))

    // 一次性初始化
    app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
        g.Log().Info(ctx, "应用初始化完成")
    }))

    // 带清理函数的初始化
    app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
        conn := connectToExternalService()
        return func(ctx context.Context) {
            conn.Close()
        }, nil
    }))

    app.Run(gctx.GetInitCtx())
}
```

### 生命周期顺序

1. `Boot(ctx)` -- 按注册顺序应用所有 Option（`nil` context 等价于 `GetInitCtx()`）
2. `Start(ctx)` -- 并发启动所有服务器（`Boot` 后若 `ctx` 为 nil，则沿用 Boot context）
3. `Run` 运行阶段阻塞直到关闭信号或 root context 结束
4. `Stop(ctx, graceful)` -- 按逆序执行清理函数，然后按逆序停止服务器

未显式调用 `Boot()` 时，`Start()` 会使用本次 `Start` 解析得到的 context 自动执行 `Boot`。

## 单次生命周期

`App` 面向 **单次进程生命周期**（典型的 `main()` 用法）：

- 首次成功的 `Stop` 会执行清理 hook 并停止所有服务器，之后的 `Stop` 调用为 no-op。
- `Stop` 之后再次 `Start` 可能对个别服务器成功，但 `App.Stop` 不会再次停止它们。需要新一轮生命周期时请创建新的 `App`。
- `contrib/scheduler/gjob` 中的任务服务器（`WorkerServer`、`CronServer`）在 `Stop` 后也会拒绝再次 `Start`。任务只能在 `Start` 之前添加；若服务器已启动或已停止，`Add` 会返回错误。

使用 `gapp.Run()` 时，请通过 adapter 注册服务器，并在 contrib 服务器上使用 `StartManaged()` 等非阻塞入口。不要对同一进程中的同一服务器再调用 legacy 阻塞 `Run()`（例如 `grpcx.GrpcServer.Run()` 或 ghttp admin 重启路径），否则可能注册冲突的 OS 信号处理器。

## 使用示例

```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/os/gapp"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    httpServer := g.Server()
    httpServer.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("hello")
    })

    app := g.App(gapp.NewHTTPServerAdapter(httpServer))
    app.Run(gctx.GetInitCtx())
}
```

## App 方法

| 方法 | 描述 |
|---|---|
| `New(servers ...Server)` | 创建新的 App，可选传入初始服务器 |
| `Add(servers ...Server)` | 向 App 添加服务器 |
| `Option(opts ...Option)` | 注册在 Boot 期间应用的 Option |
| `Boot(ctx context.Context) error` | 应用所有 Option；幂等操作 |
| `Booted() bool` | 返回 Boot 是否已成功调用 |
| `Start(ctx context.Context) error` | 并发启动服务器（`nil ctx` 复用 Boot context）；在等待启动完成时会响应 ctx |
| `Run(ctx context.Context) error` | Boot、Start、注册处理函数，阻塞直至信号或 ctx 结束；关闭失败时返回 Stop 错误 |
| `Stop(ctx context.Context, graceful bool) error` | 清理 hook 后以逆序停止服务器（`nil ctx` 复用 Boot context） |
| `Servers() []Server` | 返回已注册服务器的副本 |
