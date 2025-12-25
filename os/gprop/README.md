# gprop - Generic Configuration Manager

`gprop` is a generic configuration manager that binds configuration values to Go structs, similar to Spring Boot's `@ConfigurationProperties`.

## Overview

`gprop` provides a clean API for binding configuration data (from files, environment variables, or other configuration sources) to Go structs. It supports:
- Generic type-safe configuration binding
- Configuration change watching
- Thread-safe operations
- Callback functions support
- Custom converters
- Custom error handling

## Installation

```bash
go get github.com/gogf/gf/v2
```

## Quick Start

Here's a simple example showing how to use `gprop` to bind configuration to a struct:

```golang
package main

import (
    "context"
    "fmt"
    
    "github.com/gogf/gf/v2/os/gprop"
)

type ServerConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
    Name string `json:"name"`
}

func main() {
    ctx := context.Background()
    
    // Create configuration data
    content := `{"host":"localhost","port":8080,"name":"test-server"}`
    
    // Create target struct instance
    var config ServerConfig
    
    // Create configurator from content
    configurator, err := gprop.FromContent[ServerConfig](content, "", &config)
    if err != nil {
        panic(err)
    }
    
    // Load configuration
    err = configurator.Load(ctx)
    if err != nil {
        panic(err)
    }
    
    // Use configuration
    fmt.Printf("Server: %s:%d (%s)\n", config.Host, config.Port, config.Name)
}
```

## Real-World Usage Example

Here's an example of how `gprop` is used in a real-world application for permission management:

```go
package permission

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gprop"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

var (
	ruleConfig             RuleConfig
	RuleConfigConfigurator *gprop.Configurator[RuleConfig]
)

// API defines the API entity structure
type API struct {
	Module  string
	Path    string
	Method  string
	Summary string
	Struct  string
	Auth    bool
}

type RuleConfig struct {
	Rules *gmap.KVMap[string, API] `json:"rules" yaml:"rules"`
}

func StartRules(ctx context.Context) {
	RuleConfigConfigurator = gprop.New[RuleConfig](g.Cfg("rule"), "rules", &ruleConfig)
	RuleConfigConfigurator.SetConverter(func(data any, target *RuleConfig) error {
		m := gmap.NewKVMap[string, API](false)
		var apis []API
		err := gconv.Scan(data, &apis)
		if err != nil {
			return err
		}
		for _, api := range apis {
			key := fmt.Sprintf("method=%s;path=%s;auth=%v", api.Method, api.Path, api.Auth)
			m.Set(strings.ToUpper(key), api)
		}
		target.Rules = m
		return nil
	})
	RuleConfigConfigurator.MustLoad(ctx)
	RuleConfigConfigurator.MustWatch(ctx, "rule-config-watcher")
}

func DoAuth(method string, path string) bool {
	key := fmt.Sprintf("method=%s;path=%s;auth=%v", method, path, true)
	return RuleConfigConfigurator.Get().Rules.Contains(strings.ToUpper(key))
}
```

In this example:
- A custom configuration structure is defined to manage API permissions
- A custom converter is used to transform the raw configuration data into a specialized map structure
- The configuration is loaded and watched for changes
- The loaded configuration is used to check if an API requires authentication

## Core Features

### 1. Basic Configuration Binding

The simplest way to bind configuration to a struct:

```go
// Create configurator
configurator := gprop.New(config, "server", &configStruct)

// Load configuration
err := configurator.Load(ctx)
```

### 2. Loading Configuration from File

```go
var config ServerConfig

// Create configurator from file
configurator, err := gprop.FromFile[ServerConfig]("config.json", "server", &config)
if err != nil {
    panic(err)
}

err = configurator.Load(ctx)
```

### 3. Loading Configuration from Content String

```go
content := `{"host":"localhost","port":8080}`
var config ServerConfig

configurator, err := gprop.FromContent[ServerConfig](content, "server", &config)
if err != nil {
    panic(err)
}

err = configurator.Load(ctx)
```

### 4. Watching Configuration Changes

`gprop` supports watching configuration changes and automatically updating:

```go
// Set configuration change callback
configurator.OnChange(func(updated ServerConfig) error {
    fmt.Printf("Configuration updated: %+v\n", updated)
    return nil
})

// Start watching configuration changes
err := configurator.Watch(ctx, "my-watcher")
```

### 5. Getting Current Configuration

Use the Get() method to get a copy of the current configuration:

```go
currentConfig := configurator.Get()
```

### 6. Custom Converter

You can set a custom converter to handle the conversion from configuration data to struct:

```go
configurator.SetConverter(func(data interface{}, target *ServerConfig) error {
    // Custom conversion logic
    // For example: convert from map to struct
    if m, ok := data.(map[string]interface{}); ok {
        target.Host = m["host"].(string)
        target.Port = int(m["port"].(int64))
        target.Name = m["name"].(string)
    }
    return nil
})

err := configurator.Load(ctx)
```

### 7. Custom Error Handling

You can set a custom error handling function to be called when configuration loading fails:

```go
configurator.SetLoadErrorHandler(func(ctx context.Context, err error) {
    // Custom error handling logic, e.g., logging
    fmt.Printf("Configuration loading failed: %v\n", err)
})

err := configurator.Load(ctx)
```

### 8. MustLoad Method

If you want to panic when configuration loading fails, use the MustLoad method:

```go
// Load configuration, panic if it fails
configurator.MustLoad(ctx)
```

## API Reference

### `New[T any](config *gcfg.Config, propertyKey string, targetStruct *T) *Configurator[T]`

Creates a new `Configurator` instance.

Parameters:
- `config`: The configuration instance to watch for changes
- `propertyKey`: The property key pattern to watch (use "" or "." to watch all configuration)
- `targetStruct`: Pointer to the struct that will receive the configuration values

### `Load(ctx context.Context) error`

Loads configuration from the config instance and binds it to the target struct.
The context is passed to the underlying configuration adapter.

### `MustLoad(ctx context.Context)`

Similar to Load but panics if there is an error.

### `Watch(ctx context.Context, name string) error`

Starts watching for configuration changes and automatically updates the target struct.
name: the name of the watcher, which is used to identify this watcher.
This method sets up a watcher that will call Load() when configuration changes are detected.

### `OnChange(fn func(updated T) error)`

Sets the callback function that will be called when configuration changes.
The callback function receives the updated configuration struct and can return an error.

### `Get() T`

Returns the current configuration struct.
This method is thread-safe and returns a copy of the current configuration.

### `SetConverter(converter func(data any, target *T) error)`

Sets a custom converter function that will be used during Load operations.
The converter function receives the source data and the target struct pointer.

### `SetLoadErrorHandler(errorFunc func(ctx context.Context, err error))`

Sets an error handling function that will be called when Load operations fail.

### `FromFile[T any](filePath string, propertyKey string, targetStruct *T) (*Configurator[T], error)`

Creates a Configurator from a file path.

### `FromContent[T any](content string, propertyKey string, targetStruct *T) (*Configurator[T], error)`

Creates a Configurator from content string.

## Use Cases

- Application configuration management
- Dynamic configuration updates
- Microservice configuration center integration
- Environment-specific configuration management
- Scenarios requiring custom configuration conversion logic
- Configuration systems requiring detailed error handling

## Notes

- All operations are thread-safe
- Configuration watchers require underlying configuration adapter support
- Custom converters and error handlers are optional
- Uses gvar module for unified value processing
- Watch method returns whether the watcher was successfully added, not whether configuration loading succeeded
- When configuration loading fails, if there is an error handler, it executes, otherwise it fails silently
- Function parameters are clearly named to improve code readability