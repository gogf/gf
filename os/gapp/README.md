# gapp

Application-level lifecycle management for multiple servers.

## Introduction

The `gapp` package provides a unified `Server` interface and an `App` struct that coordinates the startup and shutdown of multiple servers, including signal handling for graceful shutdown. It also supports an `Option` mechanism for structured application initialization before servers start.

## Context propagation

- `Boot`, `Start`, `Stop`, and `Run` take a `context.Context` for tracing and cancellation.
- A nil argument is normalized to `gctx.GetInitCtx()`, matching framework defaults such as `ghttp.Server.Start`.
- After the first successful `Boot`, subsequent `Start`/`Stop` calls with a nil argument reuse the normalized context captured during `Boot`.
- During concurrent server startup in `Start`, if the resolved context completes first, partially started servers are rolled back similar to startup failure cleanup.
- `Run` shuts down gracefully when receiving an OS shutdown signal **or** when the root context completes; `Stop` receives `gctx.NeverDone(root)` so trace metadata survives while teardown itself is not shortened by parent cancellation.

## Server Interface

```go
type Server interface {
    Start() error
    Stop(graceful bool) error
}
```

- `Start()`: Starts the server in a non-blocking way.
- `Stop(graceful)`: Stops the server. When `graceful` is `true`, the server waits for in-flight requests to complete before shutting down. When `graceful` is `false`, the server is forcibly closed.

## Adapters

The package provides adapter constructors for built-in server types:

| Constructor | Wraps |
|---|---|
| `NewHTTPServerAdapter(s *ghttp.Server)` | HTTP server |
| `NewTCPServerAdapter(s *gtcp.Server)` | TCP server |
| `NewUDPServerAdapter(s *gudp.Server)` | UDP server |

For gRPC server, use `grpcx.NewGappServerAdapter(s *grpcx.GrpcServer)` from the `contrib/rpc/grpcx` package.

## Boot / Initialization

The `Option` type allows registering initialization logic that runs before servers start. This provides a structured way to set up databases, caches, external connections, and other application-level concerns.

### Option Interface

```go
type Option interface {
    Apply(ctx context.Context, app *App) (func(ctx context.Context), error)
}
```

The `Apply` method runs during `Boot()`. It can return an optional cleanup function that will be called during `Stop()` in reverse registration order.

### Option Constructors

| Constructor | Description |
|---|---|
| `NewOption(f)` | One-shot initialization, no cleanup |
| `NewOptionWithHook(f)` | Initialization with optional cleanup function |

### Usage

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

    // One-shot initialization
    app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
        g.Log().Info(ctx, "application initialized")
    }))

    // Initialization with cleanup
    app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
        conn := connectToExternalService()
        return func(ctx context.Context) {
            conn.Close()
        }, nil
    }))

    app.Run(gctx.GetInitCtx())
}
```

### Lifecycle Order

1. `Boot(ctx)` -- Applies all Options in registration order (`nil ctx` behaves like `GetInitCtx()`)
2. `Start(ctx)` -- Starts all servers concurrently (`nil ctx` after Boot reuses Boot context)
3. (running under `Run`) -- Blocks until shutdown signal or root ctx completion
4. `Stop(ctx, graceful)` -- Runs cleanup hooks in reverse order, then stops servers in reverse order

If `Boot()` is not called explicitly, `Start()` calls it automatically using the resolved `Start` context.

## Single lifecycle

`App` is designed for a **single process lifecycle** (typical `main()` usage):

- The first successful `Stop` runs cleanup hooks and stops all servers. Later `Stop` calls are no-ops.
- `Start` after `Stop` may succeed for individual servers, but `App.Stop` will not stop them again. Create a new `App` for another round.
- Job servers in `contrib/scheduler/gjob` (`WorkerServer`, `CronServer`) also reject `Start` after `Stop`. Add tasks only before `Start`; `Add` returns an error if called after the server has started or stopped.

When using `gapp.Run()`, register servers through adapters and call `StartManaged()`-style APIs on contrib servers. Do **not** call legacy blocking `Run()` on the same servers (for example `grpcx.GrpcServer.Run()` or ghttp admin restart paths), or you may register competing OS signal handlers.

## Usage

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

## App Methods

| Method | Description |
|---|---|
| `New(servers ...Server)` | Creates a new App with optional initial servers |
| `Add(servers ...Server)` | Adds servers to the App |
| `Option(opts ...Option)` | Registers Options to apply during Boot |
| `Boot(ctx context.Context) error` | Applies all Options; idempotent |
| `Booted() bool` | Returns whether Boot has been called successfully |
| `Start(ctx context.Context) error` | Starts servers concurrently (`nil ctx` reuses Boot context); honors ctx during startup waits |
| `Run(ctx context.Context) error` | Boots, starts, registers handlers, blocks until shutdown signal or ctx completes; returns Stop error on shutdown failure |
| `Stop(ctx context.Context, graceful bool) error` | Cleanup hooks then reverse server shutdown (`nil ctx` reuses Boot context) |
| `Servers() []Server` | Returns a copy of the registered servers |
