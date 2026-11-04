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

Sets the after processing function, called after the main function executes. Note: This function is only called when the main function executes successfully; if the main function returns an error, this function will not be called.

### `WithOnError(f func(context.Context, error)) *Executor[T, R]`

Sets the error handler function, called when an error occurs. Note: This function is only called when the main function returns an error.

### `Do(ctx context.Context) (R, error)`

Executes the entire execution flow, returning the result and possible errors. Execution order is: before -> main -> after (only when main function succeeds) or before -> main -> onError (only when main function fails).

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

### 4. Error Handling

```go
// Execute specific handling when an error occurs
errorHandledExecutor := gexecutor.New[string, string]("input").
    WithMain(func(ctx context.Context, input string) (string, error) {
        // Main logic that may fail
        return "", fmt.Errorf("some error occurred")
    }).
    WithOnError(func(ctx context.Context, err error) {
        // Log error or perform error recovery logic
        log.Printf("Error occurred: %v", err)
    }).
    Do(context.Background())
```

### 5. No Input Scenario

When no input parameter is needed, you can use the empty struct `struct{}{}` instead:

```go
// When no input parameters are required
executor := gexecutor.New[struct{}, string](struct{}{}).  // Using empty struct as input type
    WithMain(func(ctx context.Context, input struct{}) (string, error) {
        // Main logic, no need to use input parameter
        return "result_without_input", nil
    }).
    WithAfter(func(ctx context.Context, result string) {
        // Post-processing
        fmt.Println("Result:", result)
    }).
    Do(context.Background())
```

### 6. Closure Handling Implicit Input

Closures can perfectly handle "implicit input" by capturing external variables rather than accessing data through function parameters. This approach allows functions to access variables from the outer scope:

```go
// External variables as implicit input
var (
    config = &Config{Timeout: 10, Retries: 3}
    logger = NewLogger()
)

// Using closures to capture external variables
executor := gexecutor.New[struct{}, bool](struct{}{}).
    WithBefore(func(ctx context.Context, input struct{}) {
        // Closure can access external config and logger variables
        logger.Info("Starting operation with config", config)
    }).
    WithMain(func(ctx context.Context, input struct{}) (bool, error) {
        // Main function can also access external variables via closure
        result := processWithConfig(config, ctx)
        return result, nil
    }).
    WithAfter(func(ctx context.Context, result bool) {
        // Post-processing can also access external variables
        logger.Info("Operation completed with result", result)
    }).
    WithOnError(func(ctx context.Context, err error) {
        // Error handling can also access external variables
        logger.Error("Operation failed with error", err, "using config", config)
    }).
    Do(context.Background())
```

Advantages of closure handling implicit input:
- **Avoid parameter passing**: No need to pass all dependencies as parameters to each function
- **State sharing**: Multiple functions can share the same external state
- **Interface simplification**: Function signatures remain clean without extra parameters for additional data
- **Flexibility**: Easily access any number of external variables in closures
- **Encapsulation**: External variables remain encapsulated and not exposed to other modules

## Notes

- The main function must be set, otherwise the `Do()` method will return an `ErrMainFuncNotSet` error
- If the passed context is `nil`, it will automatically use `context.Background()`
- All `WithXxx` methods are optional and can be chosen as needed
- Each `WithXxx` call returns a new executor instance, leaving the original instance unchanged
- The `WithAfter` function is only called after successful execution of the main function, and is skipped if the main function returns an error
- The `WithOnError` function is only called when the main function returns an error
- When no input parameter is needed, using `struct{}{}` as the input type is a Go idiom since it doesn't occupy memory space
- Closures can be used to handle implicit input by capturing external variables rather than using function parameters

## Error Handling

- `ErrMainFuncNotSet`: Returned when the main function is not set
- Other errors are returned by the main function