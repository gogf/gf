// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/internal/intlog"
)

// startWatchLoop starts the loop for event listening fro underlying inotify monitor.
func (w *Watcher) startWatchLoop() {
	go func() {
		for {
			select {
			// Close event.
			case <-w.closeChan:
				return

			// Event listening.
			case ev := <-w.watcher.Events:
				// Filter the repeated event in custom duration.
				w.cache.SetIfNotExist(ev.String(), func() (interface{}, error) {
					w.events.Push(&Event{
						event:   ev,
						Path:    ev.Name,
						Op:      Op(ev.Op),
						Watcher: w,
					})
					return struct{}{}, nil
				}, repeatEventFilterDuration)

			case err := <-w.watcher.Errors:
				intlog.Error(err)
			}
		}
	}()
}

// getCallbacks searches and returns all callbacks with given <path>.
// It also searches its parent for callbacks if they're recursive.
func (w *Watcher) getCallbacks(path string) (callbacks []*Callback) {
	// Firstly add the callbacks of itself.
	if v := w.callbacks.Get(path); v != nil {
		for _, v := range v.(*glist.List).FrontAll() {
			callback := v.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}
	// Secondly searches its parent for callbacks.
	dirPath := fileDir(path)
	if v := w.callbacks.Get(dirPath); v != nil {
		for _, v := range v.(*glist.List).FrontAll() {
			callback := v.(*Callback)
			if callback.recursive {
				callbacks = append(callbacks, callback)
			}
		}
	}
	// Lastly searches the parent recursively for callbacks.
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

// startEventLoop is the core event handler.
func (w *Watcher) startEventLoop() {
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
							intlog.Error(err)
						} else {
							intlog.Printf("fake remove event, watcher re-adds monitor for: %s", event.Path)
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
							intlog.Error(err)
						} else {
							intlog.Printf("fake rename event, watcher re-adds monitor for: %s", event.Path)
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
									intlog.Error(err)
								} else {
									intlog.Printf("folder creation event, watcher adds monitor for: %s", subPath)
								}
							}
						}
					} else {
						// If it's a file, it directly adds it to monitor.
						if err := w.watcher.Add(event.Path); err != nil {
							intlog.Error(err)
						} else {
							intlog.Printf("file creation event, watcher adds monitor for: %s", event.Path)
						}
					}

				}
				// Calling the callbacks in order.
				for _, v := range callbacks {
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
					}(v)
				}
			} else {
				break
			}
		}
	}()
}
