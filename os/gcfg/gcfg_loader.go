// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

// Loader is a generic configuration manager that provides
// configuration loading, watching and management similar to Spring Boot's @ConfigurationProperties
type Loader[T any] struct {
	config        *Config                              // The configuration instance to watch
	propertyKey   string                               // The property key pattern to watch
	targetStruct  *T                                   // The target struct pointer to bind configuration to
	mutex         sync.RWMutex                         // Mutex for thread-safe operations
	onChange      func(T) error                        // Callback function when configuration changes
	converter     func(data any, target *T) error      // Optional custom converter function
	loadErrorFunc func(ctx context.Context, err error) // Optional error handling function for load failures
	reuse         bool                                 // reuse the same target struct, default is false to avoid data race
	watcherName   string                               // watcher name
}

// NewLoader creates a new Loader instance
// config: the configuration instance to watch for changes
// propertyKey: the property key pattern to watch (use "" or "." to watch all configuration)
// targetStruct: pointer to the struct that will receive the configuration values
func NewLoader[T any](config *Config, propertyKey string, targetStruct ...*T) *Loader[T] {
	if len(targetStruct) > 0 {
		return &Loader[T]{
			config:       config,
			propertyKey:  propertyKey,
			targetStruct: targetStruct[0],
			reuse:        false,
		}
	}
	return &Loader[T]{
		config:       config,
		propertyKey:  propertyKey,
		targetStruct: new(T),
		reuse:        false,
	}
}

// NewLoaderWithAdapter creates a new Loader instance
// adapter: the adapter instance to use for loading and watching configuration
// propertyKey: the property key pattern to watch (use "" or "." to watch all configuration)
// targetStruct: pointer to the struct that will receive the configuration values
func NewLoaderWithAdapter[T any](adapter Adapter, propertyKey string, targetStruct ...*T) *Loader[T] {
	return NewLoader(NewWithAdapter(adapter), propertyKey, targetStruct...)
}

// OnChange sets the callback function that will be called when configuration changes
// The callback function receives the updated configuration struct and can return an error
func (l *Loader[T]) OnChange(fn func(updated T) error) *Loader[T] {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.onChange = fn
	return l
}

// Load loads configuration from the config instance and binds it to the target struct
// The context is passed to the underlying configuration adapter
func (l *Loader[T]) Load(ctx context.Context) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Get configuration data
	var data *gvar.Var
	if l.propertyKey == "" || l.propertyKey == "." {
		// Get all configuration data
		configData, err := l.config.Data(ctx)
		if err != nil {
			if l.loadErrorFunc != nil {
				l.loadErrorFunc(ctx, err)
			}
			return err
		}
		data = gvar.New(configData)
	} else {
		// Get specific property
		configValue, err := l.config.Get(ctx, l.propertyKey)
		if err != nil {
			if l.loadErrorFunc != nil {
				l.loadErrorFunc(ctx, err)
			}
			return err
		}
		if configValue != nil {
			data = configValue
		} else {
			data = gvar.New(nil)
		}
	}

	// Use custom converter if provided, otherwise use default gconv.Scan
	if l.converter != nil && data != nil {
		if l.reuse {
			if err := l.converter(data.Val(), l.targetStruct); err != nil {
				if l.loadErrorFunc != nil {
					l.loadErrorFunc(ctx, err)
				}
				return err
			}
		} else {
			var newConfig T
			if err := l.converter(data.Val(), &newConfig); err != nil {
				if l.loadErrorFunc != nil {
					l.loadErrorFunc(ctx, err)
				}
				return err
			}
			l.targetStruct = &newConfig
		}
	} else {
		if data != nil {
			if l.reuse {
				if err := data.Scan(l.targetStruct); err != nil {
					if l.loadErrorFunc != nil {
						l.loadErrorFunc(ctx, err)
					}
					return err
				}
			} else {
				var newConfig T
				if err := data.Scan(&newConfig); err != nil {
					if l.loadErrorFunc != nil {
						l.loadErrorFunc(ctx, err)
					}
					return err
				}
				l.targetStruct = &newConfig
			}
		}
	}

	// Call change callback if exists
	if l.onChange != nil {
		return l.onChange(*l.targetStruct)
	}

	return nil
}

// MustLoad is like Load but panics if there is an error
func (l *Loader[T]) MustLoad(ctx context.Context) {
	if err := l.Load(ctx); err != nil {
		panic(err)
	}
}

// Watch starts watching for configuration changes and automatically updates the target struct
// name: the name of the watcher, which is used to identify this watcher
// This method sets up a watcher that will call Load() when configuration changes are detected
func (l *Loader[T]) Watch(ctx context.Context, name string) error {
	if name == "" {
		return gerror.New("Watcher name cannot be empty")
	}
	adapter := l.config.GetAdapter()
	if watcherAdapter, ok := adapter.(WatcherAdapter); ok {
		watcherAdapter.AddWatcher(name, func(ctx context.Context) {
			// Reload configuration when change is detected
			if err := l.Load(ctx); err != nil {
				// Use the configured error handler if available, otherwise execute default logging
				if l.loadErrorFunc != nil {
					l.loadErrorFunc(ctx, err)
				} else {
					// Default logging using intlog (internal logging for development)
					intlog.Errorf(ctx, "Configuration load failed in watcher %s: %v", name, err)
				}
			}
		})
		l.watcherName = name
		return nil
	}
	return gerror.New("Watcher adapter not found")
}

// MustWatch is like Watch but panics if there is an error
func (l *Loader[T]) MustWatch(ctx context.Context, name string) {
	if err := l.Watch(ctx, name); err != nil {
		panic(err)
	}
}

// MustLoadAndWatch is a convenience method that calls MustLoad and MustWatch
func (l *Loader[T]) MustLoadAndWatch(ctx context.Context, name string) {
	l.MustLoad(ctx)
	l.MustWatch(ctx, name)
}

// Get returns the current configuration struct
// This method is thread-safe and returns a copy of the current configuration
func (l *Loader[T]) Get() T {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return *l.targetStruct
}

// GetPointer returns a pointer to the current configuration struct
// This method is thread-safe and returns a pointer to the current configuration
// The returned pointer is safe for read operations but should not be modified
func (l *Loader[T]) GetPointer() *T {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.targetStruct
}

// SetConverter sets a custom converter function that will be used during Load operations
// The converter function receives the source data and the target struct pointer
func (l *Loader[T]) SetConverter(converter func(data any, target *T) error) *Loader[T] {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.converter = converter
	return l
}

// SetLoadErrorHandler sets an error handling function that will be called when Load operations fail
func (l *Loader[T]) SetLoadErrorHandler(errorFunc func(ctx context.Context, err error)) *Loader[T] {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.loadErrorFunc = errorFunc
	return l
}

// SetReuseTargetStruct sets whether to reuse the same target struct or create a new one on updates
func (l *Loader[T]) SetReuseTargetStruct(reuse bool) *Loader[T] {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.reuse = reuse
	return l
}

// StopWatch stops watching for configuration changes and removes the associated watcher
func (l *Loader[T]) StopWatch(ctx context.Context) (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.watcherName == "" {
		return false, gerror.New("No watcher name specified")
	}
	adapter := l.config.GetAdapter()
	if watcherAdapter, ok := adapter.(WatcherAdapter); ok {
		watcherAdapter.RemoveWatcher(l.watcherName)
		l.watcherName = ""
		return true, nil
	}
	return false, gerror.New("Watcher adapter not found")
}

// IsWatching returns true if the loader is currently watching for configuration changes
func (l *Loader[T]) IsWatching() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if l.watcherName == "" {
		return false
	}
	adapter := l.config.GetAdapter()
	if watcherAdapter, ok := adapter.(WatcherAdapter); ok {
		return watcherAdapter.IsWatching(l.watcherName)
	}
	return false
}
