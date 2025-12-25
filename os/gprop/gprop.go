// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gprop provides a generic configuration manager that binds
// configuration values to Go structs, similar to Spring Boot's @ConfigurationProperties.
package gprop

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcfg"
)

// Configurator is a generic configuration manager that provides
// configuration loading, watching and management similar to Spring Boot's @ConfigurationProperties
type Configurator[T any] struct {
	config        *gcfg.Config                         // The configuration instance to watch
	propertyKey   string                               // The property key pattern to watch
	targetStruct  *T                                   // The target struct pointer to bind configuration to
	mutex         sync.RWMutex                         // Mutex for thread-safe operations
	onChange      func(T) error                        // Callback function when configuration changes
	converter     func(data any, target *T) error      // Optional custom converter function
	loadErrorFunc func(ctx context.Context, err error) // Optional error handling function for load failures
}

// New creates a new Configurator instance
// config: the configuration instance to watch for changes
// propertyKey: the property key pattern to watch (use "" or "." to watch all configuration)
// targetStruct: pointer to the struct that will receive the configuration values
func New[T any](config *gcfg.Config, propertyKey string, targetStruct *T) *Configurator[T] {
	return &Configurator[T]{
		config:       config,
		propertyKey:  propertyKey,
		targetStruct: targetStruct,
	}
}

// OnChange sets the callback function that will be called when configuration changes
// The callback function receives the updated configuration struct and can return an error
func (c *Configurator[T]) OnChange(fn func(updated T) error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.onChange = fn
}

// Load loads configuration from the config instance and binds it to the target struct
// The context is passed to the underlying configuration adapter
func (c *Configurator[T]) Load(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Get configuration data
	var data *gvar.Var
	if c.propertyKey == "" || c.propertyKey == "." {
		// Get all configuration data
		configData, err := c.config.Data(ctx)
		if err != nil {
			if c.loadErrorFunc != nil {
				c.loadErrorFunc(ctx, err)
			}
			return err
		}
		data = gvar.New(configData)
	} else {
		// Get specific property
		configValue, err := c.config.Get(ctx, c.propertyKey)
		if err != nil {
			if c.loadErrorFunc != nil {
				c.loadErrorFunc(ctx, err)
			}
			return err
		}
		if configValue != nil {
			data = configValue
		} else {
			data = gvar.New(nil)
		}
	}

	// Create a temporary variable to hold the new configuration
	var newConfig T
	// Use custom converter if provided, otherwise use default gconv.Scan
	if c.converter != nil && data != nil {
		if err := c.converter(data.Val(), &newConfig); err != nil {
			if c.loadErrorFunc != nil {
				c.loadErrorFunc(ctx, err)
			}
			return err
		}
	} else {
		// Convert and bind to temporary variable first
		if data != nil {
			if err := data.Scan(&newConfig); err != nil {
				if c.loadErrorFunc != nil {
					c.loadErrorFunc(ctx, err)
				}
				return err
			}
		}
	}

	// Update the target struct
	*c.targetStruct = newConfig

	// Call change callback if exists
	if c.onChange != nil {
		return c.onChange(newConfig)
	}

	return nil
}

// MustLoad is like Load but panics if there is an error
func (c *Configurator[T]) MustLoad(ctx context.Context) {
	if err := c.Load(ctx); err != nil {
		panic(err)
	}
}

// Watch starts watching for configuration changes and automatically updates the target struct
// name: the name of the watcher, which is used to identify this watcher
// This method sets up a watcher that will call Load() when configuration changes are detected
func (c *Configurator[T]) Watch(ctx context.Context, name string) error {
	adapter := c.config.GetAdapter()
	if watcherAdapter, ok := adapter.(gcfg.WatcherAdapter); ok {
		watcherAdapter.AddWatcher(name, func(ctx context.Context) {
			// Reload configuration when change is detected
			if err := c.Load(ctx); err != nil {
				// Use the configured error handler if available, otherwise execute default logging
				if c.loadErrorFunc != nil {
					c.loadErrorFunc(ctx, err)
				} else {
					// Default logging using intlog (internal logging for development)
					intlog.Errorf(ctx, "Configuration load failed in watcher %s: %v", name, err)
				}
			}
		})
		return nil
	}
	return gerror.New("Watcher adapter not found")
}

// MustWatch is like Watch but panics if there is an error
func (c *Configurator[T]) MustWatch(ctx context.Context, name string) {
	if err := c.Watch(ctx, name); err != nil {
		panic(err)
	}
}

// MustLoadAndWatch is a convenience method that calls MustLoad and MustWatch
func (c *Configurator[T]) MustLoadAndWatch(ctx context.Context, name string) {
	c.MustLoad(ctx)
	c.MustWatch(ctx, name)
}

// Get returns the current configuration struct
// This method is thread-safe and returns a copy of the current configuration
func (c *Configurator[T]) Get() T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return *c.targetStruct
}

// SetConverter sets a custom converter function that will be used during Load operations
// The converter function receives the source data and the target struct pointer
func (c *Configurator[T]) SetConverter(converter func(data any, target *T) error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.converter = converter
}

// SetLoadErrorHandler sets an error handling function that will be called when Load operations fail
func (c *Configurator[T]) SetLoadErrorHandler(errorFunc func(ctx context.Context, err error)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.loadErrorFunc = errorFunc
}

// FromFile creates a Configurator from a file path
// filePath: the path to the configuration file
// propertyKey: the property key pattern to watch (use "" or "." to watch all configuration)
// targetStruct: pointer to the struct that will receive the configuration values
func FromFile[T any](filePath string, propertyKey string, targetStruct *T) (*Configurator[T], error) {
	adapter, err := gcfg.NewAdapterFile(filePath)
	if err != nil {
		return nil, err
	}
	config := gcfg.NewWithAdapter(adapter)
	return New(config, propertyKey, targetStruct), nil
}

// FromContent creates a Configurator from content string
// content: the configuration content as a string
// propertyKey: the property key pattern to watch (use "" or "." to watch all configuration)
// targetStruct: pointer to the struct that will receive the configuration values
func FromContent[T any](content string, propertyKey string, targetStruct *T) (*Configurator[T], error) {
	adapter, err := gcfg.NewAdapterContent(content)
	if err != nil {
		return nil, err
	}
	config := gcfg.NewWithAdapter(adapter)
	return New(config, propertyKey, targetStruct), nil
}
