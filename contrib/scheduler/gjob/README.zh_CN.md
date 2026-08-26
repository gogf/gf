# gjob

GoFrame 应用生命周期管理中的后台任务服务器。

## 介绍

`gjob` 包提供了两种 `gapp.Server` 实现，用于后台任务处理：

- **WorkerServer**：管理在独立协程中运行的长驻后台任务
- **CronServer**：管理按 cron 表达式定时执行的调度任务

两种服务器类型均实现了 `gapp.Server` 接口，可以注册到 `gapp.App` 进行统一的生命周期管理。

## 基于配置的启动方式

`NewServersFromConfig` 辅助函数从应用配置文件中读取任务配置并创建相应的服务器实例。当你希望将任务定义与处理代码解耦时，推荐使用此方式。

### 配置格式

```yaml
scheduler:
  job:
    - name: my-worker       # 任务唯一名称
      type: worker           # 任务类型："worker" 或 "cron"
      enable: true           # 是否启用
    - name: my-cron          # 任务唯一名称
      type: cron             # 任务类型："worker" 或 "cron"
      enable: true           # 是否启用
      spec: "*/2 * * * * *"  # cron 表达式（仅 cron 类型需要）
```

### 使用示例

```go
package main

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    gjob "github.com/gogf/gf/contrib/scheduler/gjob/v2"
    "github.com/gogf/gf/v2/os/gapp"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    httpServer := g.Server()
    httpServer.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("ok")
    })

    handlers := gjob.HandlerMap{
        "my-worker": gjob.WorkerHandler(func(ctx context.Context) func() {
            go doBackgroundWork(ctx)
            return func() {
                g.Log().Info(ctx, "worker cleaned up")
            }
        }),
        "my-cron": gjob.CronHandler(func(ctx context.Context) error {
            return syncData(ctx)
        }),
    }

    jobServers := gjob.NewServersFromConfig(context.Background(), handlers)
    app := g.App(gapp.NewHTTPServerAdapter(httpServer))
    app.Add(jobServers...)
    app.Run(gctx.GetInitCtx())
}
```

被禁用（`enable: false`）的任务或 `HandlerMap` 中没有对应处理函数的任务会自动跳过。

## WorkerServer

WorkerServer 管理在独立协程中运行的后台任务。每个任务接收一个上下文，该上下文在服务器停止时会被取消，任务还可以返回一个可选的清理函数。

### WorkerHandler

```go
type WorkerHandler func(ctx context.Context) func()
```

处理器接收一个上下文，并返回一个可选的清理函数。清理函数在服务器关闭、上下文取消后调用。

### 使用示例

```go
package main

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    gjob "github.com/gogf/gf/contrib/scheduler/gjob/v2"
    "github.com/gogf/gf/v2/os/gapp"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    workerSrv := gjob.NewWorkerServer(gctx.GetInitCtx(),
        gjob.WorkerTask{
            Name: "order-worker",
            Handler: func(ctx context.Context) func() {
                // 开始后台任务。
                go processOrders(ctx)
                return func() {
                    // 关闭时清理资源。
                    g.Log().Info(ctx, "order worker cleaned up")
                }
            },
        },
    )

    app := g.App(workerSrv)
    app.Run(gctx.GetInitCtx())
}
```

## CronServer

CronServer 管理通过 `gcron` 按 cron 表达式定时执行的任务。任务以单例模式注册，即如果上一次执行尚未完成，则跳过本次触发。

### CronHandler

```go
type CronHandler func(ctx context.Context) error
```

### 使用示例

```go
package main

import (
    "context"

    gjob "github.com/gogf/gf/contrib/scheduler/gjob/v2"
    "github.com/gogf/gf/v2/os/gapp"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    cronSrv := gjob.NewCronServer(gctx.GetInitCtx(),
        gjob.CronTask{
            Name: "data-sync",
            Spec: "0 */5 * * * *", // 每5分钟执行一次
            Handler: func(ctx context.Context) error {
                return syncData(ctx)
            },
        },
    )

    app := g.App(cronSrv)
    app.Run(gctx.GetInitCtx())
}
```

## 组合使用

两种服务器类型可以在同一应用中组合使用：

```go
package main

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    gjob "github.com/gogf/gf/contrib/scheduler/gjob/v2"
    "github.com/gogf/gf/v2/os/gapp"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    httpServer := g.Server()
    httpServer.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("ok")
    })

    workerSrv := gjob.NewWorkerServer(gctx.GetInitCtx(),
        gjob.WorkerTask{
            Name:    "event-worker",
            Handler: eventWorkerHandler,
        },
    )

    cronSrv := gjob.NewCronServer(gctx.GetInitCtx(),
        gjob.CronTask{
            Name:    "cleanup",
            Spec:    "0 0 2 * * *", // 每天凌晨2点
            Handler: cleanupHandler,
        },
    )

    app := g.App(
        gapp.NewHTTPServerAdapter(httpServer),
        workerSrv,
        cronSrv,
    )
    app.Run(gctx.GetInitCtx())
}
```

## API 参考

### WorkerServer

| 方法 | 说明 |
|---|---|
| `NewWorkerServer(ctx context.Context, tasks ...WorkerTask)` | 使用生命周期 context 创建 WorkerServer |
| `Add(tasks ...WorkerTask) error` | 在 Start 前添加任务；服务器已启动或已停止时返回错误 |
| `Start() error` | 并发启动所有任务 |
| `Stop(graceful bool) error` | 停止所有任务；取消 context 并等待 goroutine 退出 |

### CronServer

| 方法 | 说明 |
|---|---|
| `NewCronServer(ctx context.Context, tasks ...CronTask)` | 使用生命周期 context 创建 CronServer |
| `Add(tasks ...CronTask) error` | 在 Start 前添加任务；服务器已启动或已停止时返回错误 |
| `Start() error` | 注册并启动所有 cron 任务 |
| `Stop(graceful bool) error` | 停止 cron 调度器 |

### 配置辅助

| 方法 | 说明 |
|---|---|
| `NewServersFromConfig(ctx, HandlerMap)` | 从 `scheduler.job` 配置创建服务器 |
