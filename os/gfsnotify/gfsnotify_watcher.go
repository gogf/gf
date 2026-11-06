// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

// Add monitors `path` with callback function `callbackFunc` to the watcher.
//
// The parameter `path` can be either a file or a directory path.
// The optional parameter `recursive` specifies whether monitoring the `path` recursively,
// which is true in default.
func (w *Watcher) Add(path string, callbackFunc func(event *Event), option ...WatchOption,
) (callback *Callback, err error) {
	return w.AddOnce("", path, callbackFunc, option...)
}

// ReAdd re-registers `path` with the underlying file system watcher while preserving
// all callbacks already bound to the path.
func (w *Watcher) ReAdd(path string) error {
	return w.watcher.Add(path)
}

// AddOnce monitors `path` with callback function `callbackFunc` only once using unique name
// `name` to the watcher. If AddOnce is called multiple times with the same `name` parameter,
// `path` is only added to monitor once.
//
// It returns error if it's called twice with the same `name`.
//
// The parameter `path` can be either a file or a directory path.
// The optional parameter `recursive` specifies whether monitoring the `path` recursively,
// which is true in default.
func (w *Watcher) AddOnce(name, path string, callbackFunc func(event *Event), option ...WatchOption,
) (callback *Callback, err error) {
	var watchOption = w.getWatchOption(option...)
	w.nameSet.AddIfNotExistFuncLock(name, func() bool {
		// Firstly add the path to watcher.
		//
		// A path can only be watched once; watching it more than once is a no-op and will
		// not return an error.
		callback, err = w.addWithCallbackFunc(
			name, path, callbackFunc, option...,
		)
		if err != nil {
			return false
		}

		// If it's recursive adding, it then adds all sub-folders to the monitor.
		// NOTE:
		// 1. It only recursively adds **folders** to the monitor, NOT files,
		//    because if the folders are monitored and their sub-files are also monitored.
		// 2. It bounds no callbacks to the folders, because it will search the callbacks
		//    from its parent recursively if any event produced.
		if fileIsDir(path) && !watchOption.NoRecursive {
			for _, subPath := range fileAllDirs(path) {
				if fileIsDir(subPath) {
					if watchAddErr := w.watcher.Add(subPath); watchAddErr != nil {
						err = gerror.Wrapf(
							err,
							`add watch failed for path "%s", err: %s`,
							subPath, watchAddErr.Error(),
						)
					} else {
						intlog.Printf(context.TODO(), "watcher adds monitor for: %s", subPath)
					}
				}
			}
		}
		if name == "" {
			return false
		}
		return true
	})
	return
}

func (w *Watcher) getWatchOption(option ...WatchOption) WatchOption {
	if len(option) > 0 {
		return option[0]
	}
	return WatchOption{}
}

// addWithCallbackFunc adds the path to underlying monitor, creates and returns a callback object.
// Very note that if it calls multiple times with the same `path`, the latest one will overwrite the previous one.
func (w *Watcher) addWithCallbackFunc(name, path string, callbackFunc func(event *Event), option ...WatchOption,
) (callback *Callback, err error) {
	var watchOption = w.getWatchOption(option...)
	// Check and convert the given path to absolute path.
	if realPath := fileRealPath(path); realPath == "" {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `"%s" does not exist`, path)
	} else {
		path = realPath
	}
	// Create callback object.
	callback = &Callback{
		Id:        callbackIdGenerator.Add(1),
		Func:      callbackFunc,
		Path:      path,
		name:      name,
		recursive: !watchOption.NoRecursive,
	}
	// Register the callback to watcher.
	w.callbacks.LockFunc(func(m map[string]*glist.TList[*Callback]) {
		list, ok := m[path]
		if !ok {
			list = glist.NewT[*Callback](true)
			m[path] = list
		}
		callback.elem = list.PushBack(callback)
	})
	// Add the path to underlying monitor.
	if err = w.watcher.Add(path); err != nil {
		err = gerror.Wrapf(err, `add watch failed for path "%s"`, path)
	} else {
		intlog.Printf(context.TODO(), "watcher adds monitor for: %s", path)
	}
	// Add the callback to global callback map.
	callbackIdMap.Set(callback.Id, callback)
	return
}

// Close closes the watcher.
func (w *Watcher) Close() {
	close(w.closeChan)
	if err := w.watcher.Close(); err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	// Wait for the watch loop to stop pushing events before closing the queue.
	w.watchLoopWg.Wait()
	w.events.Close()
}

// Remove removes watching and all callbacks associated with the `path` recursively.
// Note that, it's recursive in default if given `path` is a directory.
func (w *Watcher) Remove(path string) error {
	var (
		err          error
		subPaths     []string
		removedPaths = make([]string, 0)
	)
	removedPaths = append(removedPaths, path)
	if fileIsDir(path) {
		subPaths, err = fileScanDir(path, "*", true)
		if err != nil {
			return err
		}
		removedPaths = append(removedPaths, subPaths...)
	}

	for _, removedPath := range removedPaths {
		// remove the callbacks of the path.
		if value := w.callbacks.Remove(removedPath); value != nil {
			for {
				if item := value.PopFront(); item != nil {
					callbackIdMap.Remove(item.Id)
				} else {
					break
				}
			}
		}
		// remove the monitor of the path from underlying monitor.
		if watcherRemoveErr := w.watcher.Remove(removedPath); watcherRemoveErr != nil {
			err = gerror.Wrapf(
				err,
				`remove watch failed for path "%s", err: %s`,
				removedPath, watcherRemoveErr.Error(),
			)
		}
	}
	return err
}

// RemoveCallback removes callback with given callback id from watcher.
//
// It removes the underlying watch when no callback remains on the path.
// Subdirectory watches added for recursive monitoring are removed together with
// that path, unless a callback on the subdirectory or a recursive ancestor still needs them.
func (w *Watcher) RemoveCallback(callbackID int) {
	callback := callbackIdMap.Get(callbackID)
	if callback == nil {
		return
	}
	callbackIdMap.Remove(callbackID)
	if callback.name != "" {
		w.nameSet.Remove(callback.name)
	}
	w.callbacks.LockFunc(func(m map[string]*glist.TList[*Callback]) {
		list, ok := m[callback.Path]
		if !ok {
			return
		}
		list.Remove(callback.elem)
		if list.Len() > 0 {
			return
		}
		delete(m, callback.Path)
		pathIsDir := fileIsDir(callback.Path)
		if pathIsDir && hasRecursiveCallbackIn(m, callback.Path) {
			return
		}
		if err := w.watcher.Remove(callback.Path); err != nil {
			intlog.Errorf(context.TODO(), `remove watch failed for path "%s": %v`, callback.Path, err)
		}
		if pathIsDir {
			w.removeDescendantWatches(m, callback.Path)
		}
	})
}

// The caller must hold the callbacks map lock.
func (w *Watcher) removeDescendantWatches(m map[string]*glist.TList[*Callback], path string) {
	for _, watchedPath := range w.watcher.WatchList() {
		if !isDescendantPath(path, watchedPath) {
			continue
		}
		if _, exists := m[watchedPath]; exists {
			continue
		}
		if hasRecursiveCallbackIn(m, watchedPath) {
			continue
		}
		if err := w.watcher.Remove(watchedPath); err != nil {
			intlog.Errorf(context.TODO(), `remove watch failed for path "%s": %v`, watchedPath, err)
		}
	}
}

func isDescendantPath(root, candidate string) bool {
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return false
	}
	parentPrefix := ".." + string(filepath.Separator)
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, parentPrefix)
}

// The caller must hold the callbacks map lock.
func hasRecursiveCallbackIn(m map[string]*glist.TList[*Callback], path string) bool {
	for dirPath := fileDir(path); ; {
		if list, ok := m[dirPath]; ok {
			for _, callback := range list.FrontAll() {
				if callback.recursive {
					return true
				}
			}
		}
		parentDirPath := fileDir(dirPath)
		if parentDirPath == dirPath {
			return false
		}
		dirPath = parentDirPath
	}
}
