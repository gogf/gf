// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gspath implements file index and search for folders.
//

package gspath

import (
	"runtime"
	"strings"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/text/gstr"
)

// updateCacheByPath adds all files under <path> recursively.
func (sp *SPath) updateCacheByPath(path string) {
	if sp.cache == nil {
		return
	}
	sp.addToCache(path, path)
}

// formatCacheName formats <name> with following rules:
// 1. The separator is unified to char '/'.
// 2. The name should be started with '/' (similar as HTTP URI).
func (sp *SPath) formatCacheName(name string) string {
	if runtime.GOOS != "linux" {
		name = gstr.Replace(name, "\\", "/")
	}
	return "/" + strings.Trim(name, "./")
}

// nameFromPath converts <filePath> to cache name.
func (sp *SPath) nameFromPath(filePath, rootPath string) string {
	name := gstr.Replace(filePath, rootPath, "")
	name = sp.formatCacheName(name)
	return name
}

// makeCacheValue formats <filePath> to cache value.
func (sp *SPath) makeCacheValue(filePath string, isDir bool) string {
	if isDir {
		return filePath + "_D_"
	}
	return filePath + "_F_"
}

// parseCacheValue parses cache value to file path and type.
func (sp *SPath) parseCacheValue(value string) (filePath string, isDir bool) {
	if value[len(value)-2 : len(value)-1][0] == 'F' {
		return value[:len(value)-3], false
	}
	return value[:len(value)-3], true
}

// addToCache adds an item to cache.
// If <filePath> is a directory, it also adds its all sub files/directories recursively
// to the cache.
func (sp *SPath) addToCache(filePath, rootPath string) {
	// Add itself firstly.
	idDir := gfile.IsDir(filePath)
	sp.cache.SetIfNotExist(
		sp.nameFromPath(filePath, rootPath), sp.makeCacheValue(filePath, idDir),
	)
	// If it's a directory, it adds its all sub files/directories recursively.
	if idDir {
		if files, err := gfile.ScanDir(filePath, "*", true); err == nil {
			//fmt.Println("gspath add to cache:", filePath, files)
			for _, path := range files {
				sp.cache.SetIfNotExist(sp.nameFromPath(path, rootPath), sp.makeCacheValue(path, gfile.IsDir(path)))
			}
		} else {
			//fmt.Errorf(err.Error())
		}
	}
}

// addMonitorByPath adds gfsnotify monitoring recursively.
// When the files under the directory are updated, the cache will be updated meanwhile.
// Note that since the listener is added recursively, if you delete a directory, the files (including the directory)
// under the directory will also generate delete events, which means it will generate N+1 events in total
// if a directory deleted and there're N files under it.
func (sp *SPath) addMonitorByPath(path string) {
	if sp.cache == nil {
		return
	}
	_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
		//glog.Debug(event.String())
		switch {
		case event.IsRemove():
			sp.cache.Remove(sp.nameFromPath(event.Path, path))

		case event.IsRename():
			if !gfile.Exists(event.Path) {
				sp.cache.Remove(sp.nameFromPath(event.Path, path))
			}

		case event.IsCreate():
			sp.addToCache(event.Path, path)
		}
	}, true)
}

// removeMonitorByPath removes gfsnotify monitoring of <path> recursively.
func (sp *SPath) removeMonitorByPath(path string) {
	if sp.cache == nil {
		return
	}
	_ = gfsnotify.Remove(path)
}
