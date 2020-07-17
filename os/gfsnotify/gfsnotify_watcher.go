// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/intlog"

	"github.com/gogf/gf/container/glist"
)

// Add monitors <path> with callback function <callbackFunc> to the watcher.
// The optional parameter <recursive> specifies whether monitoring the <path> recursively,
// which is true in default.
func (w *Watcher) Add(path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	return w.AddOnce("", path, callbackFunc, recursive...)
}

// AddOnce monitors <path> with callback function <callbackFunc> only once using unique name
// <name> to the watcher. If AddOnce is called multiple times with the same <name> parameter,
// <path> is only added to monitor once.
// It returns error if it's called twice with the same <name>.
//
// The optional parameter <recursive> specifies whether monitoring the <path> recursively,
// which is true in default.
func (w *Watcher) AddOnce(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	w.nameSet.AddIfNotExistFuncLock(name, func() bool {
		// Firstly add the path to watcher.
		callback, err = w.addWithCallbackFunc(name, path, callbackFunc, recursive...)
		if err != nil {
			return false
		}
		// If it's recursive adding, it then adds all sub-folders to the monitor.
		// NOTE:
		// 1. It only recursively adds **folders** to the monitor, NOT files,
		//    because if the folders are monitored and their sub-files are also monitored.
		// 2. It bounds no callbacks to the folders, because it will search the callbacks
		//    from its parent recursively if any event produced.
		if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
			for _, subPath := range fileAllDirs(path) {
				if fileIsDir(subPath) {
					if err := w.watcher.Add(subPath); err != nil {
						intlog.Error(err)
					} else {
						intlog.Printf("watcher adds monitor for: %s", subPath)
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

// addWithCallbackFunc adds the path to underlying monitor, creates and returns a callback object.
func (w *Watcher) addWithCallbackFunc(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	// Check and convert the given path to absolute path.
	if t := fileRealPath(path); t == "" {
		return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
	} else {
		path = t
	}
	// Create callback object.
	callback = &Callback{
		Id:        callbackIdGenerator.Add(1),
		Func:      callbackFunc,
		Path:      path,
		name:      name,
		recursive: true,
	}
	if len(recursive) > 0 {
		callback.recursive = recursive[0]
	}
	// Register the callback to watcher.
	w.callbacks.LockFunc(func(m map[string]interface{}) {
		list := (*glist.List)(nil)
		if v, ok := m[path]; !ok {
			list = glist.New(true)
			m[path] = list
		} else {
			list = v.(*glist.List)
		}
		callback.elem = list.PushBack(callback)
	})
	// Add the path to underlying monitor.
	if err := w.watcher.Add(path); err != nil {
		intlog.Error(err)
	} else {
		intlog.Printf("watcher adds monitor for: %s", path)
	}
	// Add the callback to global callback map.
	callbackIdMap.Set(callback.Id, callback)

	//intlog.Print("addWithCallbackFunc", name, path, callback.recursive)
	return
}

// Close closes the watcher.
func (w *Watcher) Close() {
	w.events.Close()
	if err := w.watcher.Close(); err != nil {
		intlog.Error(err)
	}
	close(w.closeChan)
}

// Remove removes monitor and all callbacks associated with the <path> recursively.
func (w *Watcher) Remove(path string) error {
	// Firstly remove the callbacks of the path.
	if r := w.callbacks.Remove(path); r != nil {
		list := r.(*glist.List)
		for {
			if r := list.PopFront(); r != nil {
				callbackIdMap.Remove(r.(*Callback).Id)
			} else {
				break
			}
		}
	}
	// Secondly remove monitor of all sub-files which have no callbacks.
	if subPaths, err := fileScanDir(path, "*", true); err == nil && len(subPaths) > 0 {
		for _, subPath := range subPaths {
			if w.checkPathCanBeRemoved(subPath) {
				if err := w.watcher.Remove(subPath); err != nil {
					intlog.Error(err)
				}
			}
		}
	}
	// Lastly remove the monitor of the path from underlying monitor.
	return w.watcher.Remove(path)
}

// checkPathCanBeRemoved checks whether the given path have no callbacks bound.
func (w *Watcher) checkPathCanBeRemoved(path string) bool {
	// Firstly check the callbacks in the watcher directly.
	if v := w.callbacks.Get(path); v != nil {
		return false
	}
	// Secondly check its parent whether has callbacks.
	dirPath := fileDir(path)
	if v := w.callbacks.Get(dirPath); v != nil {
		for _, c := range v.(*glist.List).FrontAll() {
			if c.(*Callback).recursive {
				return false
			}
		}
		return false
	}
	// Recursively check its parent.
	parentDirPath := ""
	for {
		parentDirPath = fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if v := w.callbacks.Get(parentDirPath); v != nil {
			for _, c := range v.(*glist.List).FrontAll() {
				if c.(*Callback).recursive {
					return false
				}
			}
			return false
		}
		dirPath = parentDirPath
	}
	return true
}

// RemoveCallback removes callback with given callback id from watcher.
func (w *Watcher) RemoveCallback(callbackId int) {
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackId); r != nil {
		callback = r.(*Callback)
	}
	if callback != nil {
		if r := w.callbacks.Get(callback.Path); r != nil {
			r.(*glist.List).Remove(callback.elem)
		}
		callbackIdMap.Remove(callbackId)
		if callback.name != "" {
			w.nameSet.Remove(callback.name)
		}
	}
}
