// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfsnotify provides a platform-independent interface for file system notifications.
package gfsnotify

import (
	"context"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcache"
)

// Watcher is the monitor for file changes.
type Watcher struct {
	watcher   *fsnotify.Watcher // Underlying fsnotify object.
	events    *gqueue.Queue     // Used for internal event management.
	cache     *gcache.Cache     // Used for repeated event filter.
	nameSet   *gset.StrSet      // Used for AddOnce feature.
	callbacks *gmap.StrAnyMap   // Path(file/folder) to callbacks mapping.
	closeChan chan struct{}     // Used for watcher closing notification.
}

// Callback is the callback function for Watcher.
type Callback struct {
	Id        int                // Unique id for callback object.
	Func      func(event *Event) // Callback function.
	Path      string             // Bound file path (absolute).
	name      string             // Registered name for AddOnce.
	elem      *glist.Element     // Element in the callbacks of watcher.
	recursive bool               // Is bound to sub-path recursively or not.
}

// Event is the event produced by underlying fsnotify.
type Event struct {
	event   fsnotify.Event // Underlying event.
	Path    string         // Absolute file path.
	Op      Op             // File operation.
	Watcher *Watcher       // Parent watcher.
}

// WatchOption holds the option for watching.
type WatchOption struct {
	// NoRecursive explicitly specifies no recursive watching.
	// Recursive watching will also watch all its current and following created subfolders and sub-files.
	//
	// Note that the recursive watching is enabled in default.
	NoRecursive bool
}

// Op is the bits union for file operations.
type Op uint32

// internalPanic is the custom panic for internal usage.
type internalPanic string

const (
	CREATE Op = 1 << iota
	WRITE
	REMOVE
	RENAME
	CHMOD
)

const (
	repeatEventFilterDuration               = time.Millisecond // Duration for repeated event filter.
	callbackExitEventPanicStr internalPanic = "exit"           // Custom exit event for internal usage.
)

var (
	mu                  sync.Mutex                // Mutex for concurrent safety of defaultWatcher.
	defaultWatcher      *Watcher                  // Default watcher.
	callbackIdMap       = gmap.NewIntAnyMap(true) // Global callback id to callback function mapping.
	callbackIdGenerator = gtype.NewInt()          // Atomic id generator for callback.
)

// New creates and returns a new watcher.
// Note that the watcher number is limited by the file handle setting of the system.
// Example: fs.inotify.max_user_instances system variable in linux systems.
//
// In most case, you can use the default watcher for usage instead of creating one.
func New() (*Watcher, error) {
	w := &Watcher{
		cache:     gcache.New(),
		events:    gqueue.New(),
		nameSet:   gset.NewStrSet(true),
		closeChan: make(chan struct{}),
		callbacks: gmap.NewStrAnyMap(true),
	}
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		w.watcher = watcher
	} else {
		intlog.Printf(context.TODO(), "New watcher failed: %v", err)
		return nil, err
	}
	go w.watchLoop()
	go w.eventLoop()
	return w, nil
}

// Add monitors `path` using default watcher with callback function `callbackFunc`.
//
// The parameter `path` can be either a file or a directory path.
// The optional parameter `recursive` specifies whether monitoring the `path` recursively, which is true in default.
func Add(path string, callbackFunc func(event *Event), option ...WatchOption) (callback *Callback, err error) {
	w, err := getDefaultWatcher()
	if err != nil {
		return nil, err
	}
	return w.Add(path, callbackFunc, option...)
}

// AddOnce monitors `path` using default watcher with callback function `callbackFunc` only once using unique name `name`.
//
// If AddOnce is called multiple times with the same `name` parameter, `path` is only added to monitor once.
// It returns error if it's called twice with the same `name`.
//
// The parameter `path` can be either a file or a directory path.
// The optional parameter `recursive` specifies whether monitoring the `path` recursively, which is true in default.
func AddOnce(name, path string, callbackFunc func(event *Event), option ...WatchOption) (callback *Callback, err error) {
	w, err := getDefaultWatcher()
	if err != nil {
		return nil, err
	}
	return w.AddOnce(name, path, callbackFunc, option...)
}

// Remove removes all monitoring callbacks of given `path` from watcher recursively.
func Remove(path string) error {
	w, err := getDefaultWatcher()
	if err != nil {
		return err
	}
	return w.Remove(path)
}

// RemoveCallback removes specified callback with given id from watcher.
func RemoveCallback(callbackID int) error {
	w, err := getDefaultWatcher()
	if err != nil {
		return err
	}
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackID); r != nil {
		callback = r.(*Callback)
	}
	if callback == nil {
		return gerror.NewCodef(gcode.CodeInvalidParameter, `callback for id %d not found`, callbackID)
	}
	w.RemoveCallback(callbackID)
	return nil
}

// Exit is only used in the callback function, which can be used to remove current callback
// of itself from the watcher.
func Exit() {
	panic(callbackExitEventPanicStr)
}

// getDefaultWatcher creates and returns the default watcher.
// This is used for lazy initialization purpose.
func getDefaultWatcher() (*Watcher, error) {
	mu.Lock()
	defer mu.Unlock()
	if defaultWatcher != nil {
		return defaultWatcher, nil
	}
	var err error
	defaultWatcher, err = New()
	return defaultWatcher, err
}
