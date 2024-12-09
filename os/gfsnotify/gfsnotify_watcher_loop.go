// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"context"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

// watchLoop starts the loop for event listening from underlying inotify monitor.
func (w *Watcher) watchLoop() {
	for {
		select {
		// close event.
		case <-w.closeChan:
			return

		// event listening.
		case ev, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			// filter the repeated event in custom duration.
			var cacheFunc = func(ctx context.Context) (value interface{}, err error) {
				w.events.Push(&Event{
					event:   ev,
					Path:    ev.Name,
					Op:      Op(ev.Op),
					Watcher: w,
				})
				return struct{}{}, nil
			}
			_, err := w.cache.SetIfNotExist(
				context.Background(),
				ev.String(),
				cacheFunc,
				repeatEventFilterDuration,
			)
			if err != nil {
				intlog.Errorf(context.TODO(), `%+v`, err)
			}

		// error occurs in underlying watcher.
		case err := <-w.watcher.Errors:
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
	}
}

// eventLoop is the core event handler.
func (w *Watcher) eventLoop() {
	var (
		err error
		ctx = context.TODO()
	)
	for {
		if v := w.events.Pop(); v != nil {
			event := v.(*Event)
			// If there's no any callback of this path, it removes it from monitor,
			// as a path watching without callback is meaningless.
			callbacks := w.getCallbacksForPath(event.Path)
			if len(callbacks) == 0 {
				_ = w.watcher.Remove(event.Path)
				continue
			}

			switch {
			case event.IsRemove():
				// It should check again the existence of the path.
				// It adds it back to the monitor if it still exists.
				if fileExists(event.Path) {
					// A watch will be automatically removed if the watched path is deleted or
					// renamed.
					//
					// It here adds the path back to monitor.
					// We need no worry about the repeat adding.
					if err = w.watcher.Add(event.Path); err != nil {
						intlog.Errorf(ctx, `%+v`, err)
					} else {
						intlog.Printf(
							ctx,
							"fake remove event, watcher re-adds monitor for: %s",
							event.Path,
						)
					}
					// Change the event to RENAME, which means it renames itself to its origin name.
					event.Op = RENAME
				}

			case event.IsRename():
				// It should check again the existence of the path.
				// It adds it back to the monitor if it still exists.
				// Especially Some editors might do RENAME and then CHMOD when it's editing file.
				if fileExists(event.Path) {
					// A watch will be automatically removed if the watched path is deleted or
					// renamed.
					//
					// It might lose the monitoring for the path, so we add the path back to monitor.
					// We need no worry about the repeat adding.
					if err = w.watcher.Add(event.Path); err != nil {
						intlog.Errorf(ctx, `%+v`, err)
					} else {
						intlog.Printf(
							ctx,
							"fake rename event, watcher re-adds monitor for: %s",
							event.Path,
						)
					}
					// Change the event to CHMOD.
					event.Op = CHMOD
				}

			case event.IsCreate():
				// =================================================================================
				// Note that it here just adds the path to monitor without any callback registering,
				// because its parent already has the callbacks.
				// =================================================================================
				if w.checkRecursiveWatchingInCreatingEvent(event.Path) {
					// It handles only folders, watching folders also watching its sub files.
					for _, subPath := range fileAllDirs(event.Path) {
						if fileIsDir(subPath) {
							if err = w.watcher.Add(subPath); err != nil {
								intlog.Errorf(ctx, `%+v`, err)
							} else {
								intlog.Printf(
									ctx,
									"folder creation event, watcher adds monitor for: %s",
									subPath,
								)
							}
						}
					}
				}
			}
			// Calling the callbacks in multiple goroutines.
			for _, callback := range callbacks {
				go w.doCallback(event, callback)
			}
		} else {
			break
		}
	}
}

// checkRecursiveWatchingInCreatingEvent checks and returns whether recursive adding given `path` to watcher
// in creating event.
func (w *Watcher) checkRecursiveWatchingInCreatingEvent(path string) bool {
	if !fileIsDir(path) {
		return false
	}
	var (
		parentDirPath string
		dirPath       = path
	)
	for {
		parentDirPath = fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if callbackItem := w.callbacks.Get(parentDirPath); callbackItem != nil {
			for _, node := range callbackItem.(*glist.List).FrontAll() {
				callback := node.(*Callback)
				if callback.recursive {
					return true
				}
			}
		}
		dirPath = parentDirPath
	}
	return false
}

func (w *Watcher) doCallback(event *Event, callback *Callback) {
	defer func() {
		if exception := recover(); exception != nil {
			switch exception {
			case callbackExitEventPanicStr:
				w.RemoveCallback(callback.Id)
			default:
				if e, ok := exception.(error); ok {
					panic(gerror.WrapCode(gcode.CodeInternalPanic, e))
				}
				panic(exception)
			}
		}
	}()
	callback.Func(event)
}

// getCallbacksForPath searches and returns all callbacks with given `path`.
//
// It also searches its parents for callbacks if they're recursive.
func (w *Watcher) getCallbacksForPath(path string) (callbacks []*Callback) {
	// Firstly add the callbacks of itself.
	if item := w.callbacks.Get(path); item != nil {
		for _, node := range item.(*glist.List).FrontAll() {
			callback := node.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}
	// ============================================================================================================
	// Secondly searches its direct parent for callbacks.
	//
	// Note that it is SPECIAL handling here, which is the different between `recursive` and `not recursive` logic
	// for direct parent folder of `path` that events are from.
	// ============================================================================================================
	dirPath := fileDir(path)
	if item := w.callbacks.Get(dirPath); item != nil {
		for _, node := range item.(*glist.List).FrontAll() {
			callback := node.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}

	// Lastly searches all the parents of directory of `path` recursively for callbacks.
	for {
		parentDirPath := fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if item := w.callbacks.Get(parentDirPath); item != nil {
			for _, node := range item.(*glist.List).FrontAll() {
				callback := node.(*Callback)
				if callback.recursive {
					callbacks = append(callbacks, callback)
				}
			}
		}
		dirPath = parentDirPath
	}
	return
}
