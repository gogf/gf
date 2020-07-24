// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gspath implements file index and search for folders.
//
// It searches file internally with high performance in order by the directory adding sequence.
// Note that:
// If caching feature enabled, there would be a searching delay after adding/deleting files.
package gspath

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"os"
	"sort"
	"strings"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

// SPath manages the path searching feature.
type SPath struct {
	paths *garray.StrArray // The searching directories array.
	cache *gmap.StrStrMap  // Searching cache map, it is not enabled if it's nil.
}

// SPathCacheItem is a cache item for searching.
type SPathCacheItem struct {
	path  string // Absolute path for file/dir.
	isDir bool   // Is directory or not.
}

var (
	// Path to searching object mapping, used for instance management.
	pathsMap = gmap.NewStrAnyMap(true)
)

// New creates and returns a new path searching manager.
func New(path string, cache bool) *SPath {
	sp := &SPath{
		paths: garray.NewStrArray(true),
	}
	if cache {
		sp.cache = gmap.NewStrStrMap(true)
	}
	if len(path) > 0 {
		if _, err := sp.Add(path); err != nil {
			//intlog.Print(err)
		}
	}
	return sp
}

// Get creates and returns a instance of searching manager for given path.
// The parameter <cache> specifies whether using cache feature for this manager.
// If cache feature is enabled, it asynchronously and recursively scans the path
// and updates all sub files/folders to the cache using package gfsnotify.
func Get(root string, cache bool) *SPath {
	if root == "" {
		root = "/"
	}
	return pathsMap.GetOrSetFuncLock(root, func() interface{} {
		return New(root, cache)
	}).(*SPath)
}

// Search searches file <name> under path <root>.
// The parameter <root> should be a absolute path. It will not automatically
// convert <root> to absolute path for performance reason.
// The optional parameter <indexFiles> specifies the searching index files when the result is a directory.
// For example, if the result <a> is a directory, and <indexFiles> is [index.html, main.html], it will also
// search [index.html, main.html] under <a>. It returns the absolute file path if any of them found,
// or else it returns <a>.
func Search(root string, name string, indexFiles ...string) (filePath string, isDir bool) {
	return Get(root, false).Search(name, indexFiles...)
}

// SearchWithCache searches file <name> under path <root> with cache feature enabled.
// The parameter <root> should be a absolute path. It will not automatically
// convert <root> to absolute path for performance reason.
// The optional parameter <indexFiles> specifies the searching index files when the result is a directory.
// For example, if the result <a> is a directory, and <indexFiles> is [index.html, main.html], it will also
// search [index.html, main.html] under <a>. It returns the absolute file path if any of them found,
// or else it returns <a>.
func SearchWithCache(root string, name string, indexFiles ...string) (filePath string, isDir bool) {
	return Get(root, true).Search(name, indexFiles...)
}

// Set deletes all other searching directories and sets the searching directory for this manager.
func (sp *SPath) Set(path string) (realPath string, err error) {
	realPath = gfile.RealPath(path)
	if realPath == "" {
		realPath, _ = sp.Search(path)
		if realPath == "" {
			realPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + path)
		}
	}
	if realPath == "" {
		return realPath, errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
	}
	// The set path must be a directory.
	if gfile.IsDir(realPath) {
		realPath = strings.TrimRight(realPath, gfile.Separator)
		if sp.paths.Search(realPath) != -1 {
			for _, v := range sp.paths.Slice() {
				sp.removeMonitorByPath(v)
			}
		}
		intlog.Print("paths clear:", sp.paths)
		sp.paths.Clear()
		if sp.cache != nil {
			sp.cache.Clear()
		}
		sp.paths.Append(realPath)
		sp.updateCacheByPath(realPath)
		sp.addMonitorByPath(realPath)
		return realPath, nil
	} else {
		return "", errors.New(path + " should be a folder")
	}
}

// Add adds more searching directory to the manager.
// The manager will search file in added order.
func (sp *SPath) Add(path string) (realPath string, err error) {
	realPath = gfile.RealPath(path)
	if realPath == "" {
		realPath, _ = sp.Search(path)
		if realPath == "" {
			realPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + path)
		}
	}
	if realPath == "" {
		return realPath, errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
	}
	// The added path must be a directory.
	if gfile.IsDir(realPath) {
		//fmt.Println("gspath:", realPath, sp.paths.Search(realPath))
		// It will not add twice for the same directory.
		if sp.paths.Search(realPath) < 0 {
			realPath = strings.TrimRight(realPath, gfile.Separator)
			sp.paths.Append(realPath)
			sp.updateCacheByPath(realPath)
			sp.addMonitorByPath(realPath)
		}
		return realPath, nil
	} else {
		return "", errors.New(path + " should be a folder")
	}
}

// Search searches file <name> in the manager.
// The optional parameter <indexFiles> specifies the searching index files when the result is a directory.
// For example, if the result <a> is a directory, and <indexFiles> is [index.html, main.html], it will also
// search [index.html, main.html] under <a>. It returns the absolute file path if any of them found,
// or else it returns <a>.
func (sp *SPath) Search(name string, indexFiles ...string) (filePath string, isDir bool) {
	// No cache enabled.
	if sp.cache == nil {
		sp.paths.LockFunc(func(array []string) {
			path := ""
			for _, v := range array {
				path = gfile.Join(v, name)
				if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
					path = gfile.Abs(path)
					// Security check: the result file path must be under the searching directory.
					if len(path) >= len(v) && path[:len(v)] == v {
						filePath = path
						isDir = stat.IsDir()
						break
					}
				}
			}
		})
		if len(indexFiles) > 0 && isDir {
			if name == "/" {
				name = ""
			}
			path := ""
			for _, file := range indexFiles {
				path = filePath + gfile.Separator + file
				if gfile.Exists(path) {
					filePath = path
					isDir = false
					break
				}
			}
		}
		return
	}
	// Using cache feature.
	name = sp.formatCacheName(name)
	if v := sp.cache.Get(name); v != "" {
		filePath, isDir = sp.parseCacheValue(v)
		if len(indexFiles) > 0 && isDir {
			if name == "/" {
				name = ""
			}
			for _, file := range indexFiles {
				if v := sp.cache.Get(name + "/" + file); v != "" {
					return sp.parseCacheValue(v)
				}
			}
		}
	}
	return
}

// Remove deletes the <path> from cache files of the manager.
// The parameter <path> can be either a absolute path or just a relative file name.
func (sp *SPath) Remove(path string) {
	if sp.cache == nil {
		return
	}
	if gfile.Exists(path) {
		for _, v := range sp.paths.Slice() {
			name := gstr.Replace(path, v, "")
			name = sp.formatCacheName(name)
			sp.cache.Remove(name)
		}
	} else {
		name := sp.formatCacheName(path)
		sp.cache.Remove(name)
	}
}

// Paths returns all searching directories.
func (sp *SPath) Paths() []string {
	return sp.paths.Slice()
}

// AllPaths returns all paths cached in the manager.
func (sp *SPath) AllPaths() []string {
	if sp.cache == nil {
		return nil
	}
	paths := sp.cache.Keys()
	if len(paths) > 0 {
		sort.Strings(paths)
	}
	return paths
}

// Size returns the count of the searching directories.
func (sp *SPath) Size() int {
	return sp.paths.Len()
}
