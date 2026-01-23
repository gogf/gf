# gexecutor

`gexecutor` is a generic executor component in the GoFrame framework that provides an elegant way to execute functions with before and after processing logic. It uses generic design to support arbitrary input and output types.

## Features

- **Generic Support**: Supports arbitrary input type `T` and output type `R`
- **Chainable API**: Supports chainable configuration of the executor
- **Before/After Processing**: Supports custom logic execution before and after the main function
- **Context Aware**: Supports passing context parameters
- **Error Handling**: Built-in error handling mechanism

## Installation

```bash
go get github.com/gogf/gf/v2
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/gogf/gf/v2/util/gexecutor"
)

func main() {
    // Create an executor with integer input and string output
    executor := gexecutor.New[int, string](42)
    
    // Configure the executor: before processing -> main function -> after processing
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

## API Interface

### `New[T any, R any](input T) *Executor[T, R]`

Creates a new executor instance.

### `WithMain(f func(context.Context, T) (R, error)) *Executor[T, R]`

Sets the main execution function.

### `WithBefore(f func(context.Context, T)) *Executor[T, R]`

Sets the before processing function, called before the main function executes.

### `WithAfter(f func(context.Context, R)) *Executor[T, R]`

Sets the after processing function, called after the main function executes.

### `Do(ctx context.Context) (R, error)`

Executes the entire execution flow, returning the result and possible errors.

## Usage Scenarios

### 1. Business Logic Execution

```go
// Execute order processing logic
orderExecutor := gexecutor.New[*Order, *ProcessedOrder](order).
    WithBefore(validateOrder).
    WithMain(processOrder).
    WithAfter(notifyOrderProcessed).
    Do(context.Background())
```

### 2. Data Processing Pipeline

```go
// Data cleaning and transformation
dataProcessor := gexecutor.New[[]string, []string](rawData).
    WithBefore(logProcessing).
    WithMain(cleanAndTransform).
    WithAfter(saveProcessedData).
    Do(ctx)
```

### 3. Template Reuse

Since each `WithXxx` call returns a new instance, you can create executor templates and reuse them safely:

```go
// Create a general template
templateExecutor := gexecutor.New[int, string](0).
    WithBefore(func(ctx context.Context, input int) {
        // General before processing
    }).
    WithAfter(func(ctx context.Context, result string) {
        // General after processing
    })

// Create specific executors based on the template
doubler := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
    return fmt.Sprintf("doubled: %d", input*2), nil
})

tripler := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
    return fmt.Sprintf("tripled: %d", input*3), nil
})
```

## Notes

- The main function must be set, otherwise the `Do()` method will return an `ErrMainFuncNotSet` error
- If the passed context is `nil`, it will automatically use `context.Background()`
- All `WithXxx` methods are optional and can be chosen as needed
- Each `WithXxx` call returns a new executor instance, leaving the original instance unchanged

## Error Handling

- `ErrMainFuncNotSet`: Returned when the main function is not set
- Other errors are returned by the main function