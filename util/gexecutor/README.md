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

设置后置处理函数，在主函数执行后调用。注意：只有在主函数执行成功时才会调用此函数；如果主函数返回错误，则不会调用此函数。

### `WithOnError(f func(context.Context, error)) *Executor[T, R]`

设置错误处理函数，当发生错误时调用。注意：仅在主函数返回错误时才会调用此函数。

### `Do(ctx context.Context) (R, error)`

执行整个执行流程，返回结果和可能的错误。执行顺序为：before -> main -> after（仅当main函数成功时）或 before -> main -> onError（仅当main函数失败时）。

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

### 4. 错误处理

```go
// 在发生错误时执行特定处理
errorHandledExecutor := gexecutor.New[string, string]("input").
    WithMain(func(ctx context.Context, input string) (string, error) {
        // 可能出错的主逻辑
        return "", fmt.Errorf("some error occurred")
    }).
    WithOnError(func(ctx context.Context, err error) {
        // 记录错误日志或执行错误恢复逻辑
        log.Printf("Error occurred: %v", err)
    }).
    Do(context.Background())
```

### 5. 无输入情况处理

当不需要输入参数时，可以使用空结构体 `struct{}{}` 来代替：

```go
// 当不需要任何输入参数时
executor := gexecutor.New[struct{}, string](struct{}{}).  // 使用空结构体作为输入类型
    WithMain(func(ctx context.Context, input struct{}) (string, error) {
        // 主逻辑，不需要使用 input 参数
        return "result_without_input", nil
    }).
    WithAfter(func(ctx context.Context, result string) {
        // 后置处理
        fmt.Println("Result:", result)
    }).
    Do(context.Background())
```

### 6. 闭包处理隐式输入

闭包可以完美处理"隐式输入"，即通过捕获外部变量的方式访问数据，而不是通过函数参数。这种方式使得可以在执行器函数中访问外部作用域的变量：

```go
// 外部变量作为隐式输入
var (
    config = &Config{Timeout: 10, Retries: 3}
    logger = NewLogger()
)

// 使用闭包捕获外部变量
executor := gexecutor.New[struct{}, bool](struct{}{}).
    WithBefore(func(ctx context.Context, input struct{}) {
        // 闭包可以访问外部的 config 和 logger 变量
        logger.Info("Starting operation with config", config)
    }).
    WithMain(func(ctx context.Context, input struct{}) (bool, error) {
        // 主函数同样可以通过闭包访问外部变量
        result := processWithConfig(config, ctx)
        return result, nil
    }).
    WithAfter(func(ctx context.Context, result bool) {
        // 后置处理也可以访问外部变量
        logger.Info("Operation completed with result", result)
    }).
    WithOnError(func(ctx context.Context, err error) {
        // 错误处理同样可以访问外部变量
        logger.Error("Operation failed with error", err, "using config", config)
    }).
    Do(context.Background())
```

闭包处理隐式输入的优势：
- **避免参数传递**：无需将所有依赖项作为参数传递给每个函数
- **状态共享**：多个函数可以共享相同的外部状态
- **简化接口**：函数签名保持简洁，不必为了传递额外数据而增加参数
- **灵活性**：可以在闭包中轻松访问任意数量的外部变量
- **封装性**：外部变量仍然保持封装，不会暴露给其他模块

## 注意事项

- 主函数必须设置，否则 `Do()` 方法会返回 `ErrMainFuncNotSet` 错误
- 如果传入的上下文为 `nil`，会自动使用 `context.Background()`
- 所有 `WithXxx` 方法都是可选的，可以根据需要选择使用
- 每次 `WithXxx` 调用都会返回一个新的执行器实例，不影响原始实例
- `WithAfter` 函数只在主函数成功执行后才调用，如果主函数返回错误则不会调用
- `WithOnError` 函数只在主函数返回错误时才调用
- 当不需要输入参数时，使用 `struct{}{}` 作为输入类型是 Go 语言的惯用法，因为它不占用内存空间
- 闭包可用于处理隐式输入，通过捕获外部变量而非函数参数来访问数据

## 错误处理

- `ErrMainFuncNotSet`：当主函数未设置时返回此错误
- 其他错误由主函数返回