// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"context"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/internal/intlog"
)

// watchLoop starts the loop for event listening from underlying inotify monitor.
func (w *Watcher) watchLoop() {
	go func() {
		for {
			select {
			// Close event.
			case <-w.closeChan:
				return

			// Event listening.
			case ev := <-w.watcher.Events:
				// Filter the repeated event in custom duration.
				_, err := w.cache.SetIfNotExist(
					context.Background(),
					ev.String(),
					func(ctx context.Context) (value interface{}, err error) {
						w.events.Push(&Event{
							event:   ev,
							Path:    ev.Name,
							Op:      Op(ev.Op),
							Watcher: w,
						})
						return struct{}{}, nil
					}, repeatEventFilterDuration,
				)
				if err != nil {
					intlog.Errorf(context.TODO(), `%+v`, err)
				}

			case err := <-w.watcher.Errors:
				intlog.Errorf(context.TODO(), `%+v`, err)
			}
		}
	}()
}

// eventLoop is the core event handler.
func (w *Watcher) eventLoop() {
	go func() {
		for {
			if v := w.events.Pop(); v != nil {
				event := v.(*Event)
				// If there's no any callback of this path, it removes it from monitor.
				callbacks := w.getCallbacks(event.Path)
				if len(callbacks) == 0 {
					w.watcher.Remove(event.Path)
					continue
				}
				switch {
				case event.IsRemove():
					// It should check again the existence of the path.
					// It adds it back to the monitor if it still exists.
					if fileExists(event.Path) {
						// It adds the path back to monitor.
						// We need no worry about the repeat adding.
						if err := w.watcher.Add(event.Path); err != nil {
							intlog.Errorf(context.TODO(), `%+v`, err)
						} else {
							intlog.Printf(context.TODO(), "fake remove event, watcher re-adds monitor for: %s", event.Path)
						}
						// Change the event to RENAME, which means it renames itself to its origin name.
						event.Op = RENAME
					}

				case event.IsRename():
					// It should check again the existence of the path.
					// It adds it back to the monitor if it still exists.
					// Especially Some editors might do RENAME and then CHMOD when it's editing file.
					if fileExists(event.Path) {
						// It might lost the monitoring for the path, so we add the path back to monitor.
						// We need no worry about the repeat adding.
						if err := w.watcher.Add(event.Path); err != nil {
							intlog.Errorf(context.TODO(), `%+v`, err)
						} else {
							intlog.Printf(context.TODO(), "fake rename event, watcher re-adds monitor for: %s", event.Path)
						}
						// Change the event to CHMOD.
						event.Op = CHMOD
					}

				case event.IsCreate():
					// =========================================
					// Note that it here just adds the path to monitor without any callback registering,
					// because its parent already has the callbacks.
					// =========================================
					if fileIsDir(event.Path) {
						// If it's a folder, it then does adding recursively to monitor.
						for _, subPath := range fileAllDirs(event.Path) {
							if fileIsDir(subPath) {
								if err := w.watcher.Add(subPath); err != nil {
									intlog.Errorf(context.TODO(), `%+v`, err)
								} else {
									intlog.Printf(context.TODO(), "folder creation event, watcher adds monitor for: %s", subPath)
								}
							}
						}
					} else {
						// If it's a file, it directly adds it to monitor.
						if err := w.watcher.Add(event.Path); err != nil {
							intlog.Errorf(context.TODO(), `%+v`, err)
						} else {
							intlog.Printf(context.TODO(), "file creation event, watcher adds monitor for: %s", event.Path)
						}
					}
				}
				// Calling the callbacks in order.
				for _, callback := range callbacks {
					go func(callback *Callback) {
						defer func() {
							if err := recover(); err != nil {
								switch err {
								case callbackExitEventPanicStr:
									w.RemoveCallback(callback.Id)
								default:
									panic(err)
								}
							}
						}()
						callback.Func(event)
					}(callback)
				}
			} else {
				break
			}
		}
	}()
}

// getCallbacks searches and returns all callbacks with given `path`.
// It also searches its parents for callbacks if they're recursive.
func (w *Watcher) getCallbacks(path string) (callbacks []*Callback) {
	// Firstly add the callbacks of itself.
	if v := w.callbacks.Get(path); v != nil {
		for _, v := range v.(*glist.List).FrontAll() {
			callback := v.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}
	// Secondly searches its direct parent for callbacks.
	// It is special handling here, which is the different between `recursive` and `not recursive` logic
	// for direct parent folder of `path` that events are from.
	dirPath := fileDir(path)
	if v := w.callbacks.Get(dirPath); v != nil {
		for _, v := range v.(*glist.List).FrontAll() {
			callback := v.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}
	// Lastly searches all the parents of directory of `path` recursively for callbacks.
	for {
		parentDirPath := fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if v := w.callbacks.Get(parentDirPath); v != nil {
			for _, v := range v.(*glist.List).FrontAll() {
				callback := v.(*Callback)
				if callback.recursive {
					callbacks = append(callbacks, callback)
				}
			}
		}
		dirPath = parentDirPath
	}
	return
}
