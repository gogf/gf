# gjob

Background job servers for GoFrame application lifecycle management.

## Introduction

The `gjob` package provides two `gapp.Server` implementations for background job processing:

- **WorkerServer**: manages long-running background worker tasks that run in their own goroutines
- **CronServer**: manages scheduled cron tasks that run on specified intervals

Both server types implement the `gapp.Server` interface and can be registered with `gapp.App` for unified lifecycle management.

## Configuration-Based Setup

The `NewServersFromConfig` helper reads job configurations from the application config file and creates the appropriate server instances. This is the recommended way to set up job servers when you want to decouple task definitions from handler code.

### Configuration Format

```yaml
scheduler:
  job:
    - name: my-worker       # Unique task name
      type: worker           # Task type: "worker" or "cron"
      enable: true           # Whether to enable this task
    - name: my-cron          # Unique task name
      type: cron             # Task type: "worker" or "cron"
      enable: true           # Whether to enable this task
      spec: "*/2 * * * * *"  # Cron expression (cron type only)
```

### Usage

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

Tasks that are disabled (`enable: false`) or have no matching handler in the `HandlerMap` are automatically skipped.

## WorkerServer

WorkerServer manages background tasks that run in their own goroutines. Each task receives a context that is cancelled when the server stops, and can return an optional cleanup function.

### WorkerHandler

```go
type WorkerHandler func(ctx context.Context) func()
```

The handler receives a context and returns an optional cleanup function. The cleanup is called after the context is cancelled during server shutdown.

### Usage

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
                // Start background work.
                go processOrders(ctx)
                return func() {
                    // Cleanup resources on shutdown.
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

CronServer manages tasks that run on a cron schedule using `gcron`. Tasks are registered as singleton jobs, meaning concurrent invocations are skipped if the previous one is still running.

### CronHandler

```go
type CronHandler func(ctx context.Context) error
```

### Usage

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
            Spec: "0 */5 * * * *", // Every 5 minutes
            Handler: func(ctx context.Context) error {
                return syncData(ctx)
            },
        },
    )

    app := g.App(cronSrv)
    app.Run(gctx.GetInitCtx())
}
```

## Combined Usage

Both server types can be used together in a single application:

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
            Spec:    "0 0 2 * * *", // Daily at 2am
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

## API Reference

### WorkerServer

| Method | Description |
|---|---|
| `NewWorkerServer(ctx context.Context, tasks ...WorkerTask)` | Creates a new WorkerServer with lifecycle context |
| `Add(tasks ...WorkerTask) error` | Adds tasks before Start; returns error if server already started or stopped |
| `Start() error` | Starts all tasks concurrently |
| `Stop(graceful bool) error` | Stops all tasks; cancels context and waits for goroutines to finish |

### CronServer

| Method | Description |
|---|---|
| `NewCronServer(ctx context.Context, tasks ...CronTask)` | Creates a new CronServer with lifecycle context |
| `Add(tasks ...CronTask) error` | Adds tasks before Start; returns error if server already started or stopped |
| `Start() error` | Registers and starts all cron tasks |
| `Stop(graceful bool) error` | Stops the cron scheduler |

### Config Helper

| Method | Description |
|---|---|
| `NewServersFromConfig(ctx, HandlerMap)` | Creates servers from `scheduler.job` config |
