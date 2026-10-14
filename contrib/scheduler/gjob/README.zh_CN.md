# gjob

GoFrame 应用生命周期管理中的后台任务服务器。

## 介绍

`gjob`包提供了两种`gapp.Server`实现，用于后台任务处理：

- **WorkerServer**：管理在独立协程中运行的长驻后台任务
- **CronServer**：管理按`cron`表达式定时执行的调度任务

两种服务器类型均实现了`gapp.Server`接口，可以注册到`gapp.App`进行统一的生命周期管理。

## 基于配置的启动方式

`NewServersFromConfig`辅助函数从应用配置文件中读取任务配置并创建相应的服务器实例。当你希望将任务定义与处理代码解耦时，推荐使用此方式。

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

    jobServers, err := gjob.NewServersFromConfig(context.Background(), handlers)
    if err != nil {
        g.Log().Fatal(context.Background(), err)
    }
    app := g.App(gapp.NewHTTPServerAdapter(httpServer))
    app.Add(jobServers...)
    app.Run(gctx.GetInitCtx())
}
```

被禁用（`enable: false`）的任务，或`HandlerMap`中没有对应处理函数的任务，会被跳过。读取`scheduler.job`失败、未知任务类型、缺少`cron`的`spec`、处理函数为 nil 或类型不匹配时，会返回错误，避免进程在这些任务缺失的情况下仍正常启动。直接构造的服务器也会在执行任何任务前校验全部处理函数。

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

CronServer 管理通过`gcron`按`cron`表达式定时执行的任务。任务以单例模式注册，即如果上一次执行尚未完成，则跳过本次触发。

### CronHandler

```go
type CronHandler func(ctx context.Context) error
```

### 使用示例

```go
package main

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    gjob "github.com/gogf/gf/contrib/scheduler/gjob/v2"
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
| `NewWorkerServer(ctx context.Context, tasks ...WorkerTask)` | 使用生命周期`context`创建 WorkerServer |
| `Add(tasks ...WorkerTask) error` | 在服务器未运行时添加任务；服务器已启动时返回错误 |
| `Start() error` | 并发启动所有任务；Stop 返回后可再次调用，Stop 进行中返回错误 |
| `Stop(graceful bool) error` | 取消`context`并等待 goroutine 和清理完成；并发调用等待同一次停止完成 |

### CronServer

| 方法 | 说明 |
|---|---|
| `NewCronServer(ctx context.Context, tasks ...CronTask)` | 使用生命周期`context`创建 CronServer |
| `Add(tasks ...CronTask) error` | 在服务器未运行时添加任务；服务器已启动时返回错误 |
| `Start() error` | 校验处理函数、注册全部任务后再启动 cron；注册失败不执行任务；Stop 返回后可再次调用 |
| `Stop(graceful bool) error` | 停止调度器并关闭所有任务条目；并发调用等待同一次停止完成 |

### 配置辅助

| 方法 | 说明 |
|---|---|
| `NewServersFromConfig(ctx, HandlerMap) ([]gapp.Server, error)` | 从`scheduler.job`配置创建服务器；非法配置返回错误 |

`scheduler.job`必须是任务对象列表。配置缺失、`null`或空列表时不创建服务器；标量值、列表位置的对象以及列表中的非对象元素都会返回配置错误。
