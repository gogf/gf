# gexecutor

`gexecutor` 是 GoFrame 框架中的一个通用执行器组件，它提供了一种优雅的方式来执行带有前置和后置处理逻辑的函数。它使用泛型设计，支持任意输入和输出类型。

## 特性

- **泛型支持**：支持任意输入类型 `T` 和输出类型 `R`
- **链式调用**：支持链式配置执行器
- **前置/后置处理**：支持在主函数执行前后执行自定义逻辑
- **上下文感知**：支持传递上下文参数
- **错误处理**：内置错误处理机制

## 安装

```bash
go get github.com/gogf/gf/v2
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/gogf/gf/v2/util/gexecutor"
)

func main() {
    // 创建一个执行器，输入为整数，输出为字符串
    executor := gexecutor.New[int, string](42)
    
    // 配置执行器：前置处理 -> 主函数 -> 后置处理
    result, err := executor.
        WithBefore(func(ctx context.Context, input int) {
            fmt.Printf("Before: input is %d\n", input)
        }).
        WithMain(func(ctx context.Context, input int) (string, error) {
            return fmt.Sprintf("processed_%d", input), nil
        }).
        WithAfter(func(ctx context.Context, result string) {
            fmt.Printf("After: result is %s\n", result)
        }).
        Do(context.Background())
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Final result: %s\n", result)
}
```

## API 接口

### `New[T any, R any](input T) *Executor[T, R]`

创建一个新的执行器实例。

### `WithMain(f func(context.Context, T) (R, error)) *Executor[T, R]`

设置主执行函数。

### `WithBefore(f func(context.Context, T)) *Executor[T, R]`

设置前置处理函数，在主函数执行前调用。

### `WithAfter(f func(context.Context, R)) *Executor[T, R]`

设置后置处理函数，在主函数执行后调用。

### `Do(ctx context.Context) (R, error)`

执行整个执行流程，返回结果和可能的错误。

## 使用场景

### 1. 业务逻辑执行

```go
// 执行订单处理逻辑
orderExecutor := gexecutor.New[*Order, *ProcessedOrder](order).
    WithBefore(validateOrder).
    WithMain(processOrder).
    WithAfter(notifyOrderProcessed).
    Do(context.Background())
```

### 2. 数据处理管道

```go
// 数据清洗和转换
dataProcessor := gexecutor.New[[]string, []string](rawData).
    WithBefore(logProcessing).
    WithMain(cleanAndTransform).
    WithAfter(saveProcessedData).
    Do(ctx)
```

### 3. 模板复用

由于每次 `WithXxx` 调用都会返回新实例，可以创建执行器模板并安全复用：

```go
// 创建通用模板
templateExecutor := gexecutor.New[int, string](0).
    WithBefore(func(ctx context.Context, input int) {
        // 通用前置处理
    }).
    WithAfter(func(ctx context.Context, result string) {
        // 通用后置处理
    })

// 基于模板创建特定的执行器
doubler := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
    return fmt.Sprintf("doubled: %d", input*2), nil
})

tripler := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
    return fmt.Sprintf("tripled: %d", input*3), nil
})
```

## 注意事项

- 主函数必须设置，否则 `Do()` 方法会返回 `ErrMainFuncNotSet` 错误
- 如果传入的上下文为 `nil`，会自动使用 `context.Background()`
- 所有 `WithXxx` 方法都是可选的，可以根据需要选择使用
- 每次 `WithXxx` 调用都会返回一个新的执行器实例，不影响原始实例

## 错误处理

- `ErrMainFuncNotSet`：当主函数未设置时返回此错误
- 其他错误由主函数返回