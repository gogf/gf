// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gspath implements file index and search for folders.
//

package gspath

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfsnotify"
	"github.com/gogf/gf/g/text/gstr"
	"runtime"
	"strings"
)

// 递归添加目录下的文件
func (sp *SPath) updateCacheByPath(path string) {
	if sp.cache == nil {
		return
	}
	sp.addToCache(path, path)
}

// 格式化name返回符合规范的缓存名称，分隔符号统一为'/'，且前缀必须以'/'开头(类似HTTP URI).
func (sp *SPath) formatCacheName(name string) string {
	if runtime.GOOS != "linux" {
		name = gstr.Replace(name, "\\", "/")
	}
	return "/" + strings.Trim(name, "./")
}

// 根据path计算出对应的缓存name, dirPath为检索根目录路径
func (sp *SPath) nameFromPath(filePath, rootPath string) string {
	name := gstr.Replace(filePath, rootPath, "")
	name = sp.formatCacheName(name)
	return name
}

// 按照一定数据结构生成缓存的数据项字符串
func (sp *SPath) makeCacheValue(filePath string, isDir bool) string {
	if isDir {
		return filePath + "_D_"
	}
	return filePath + "_F_"
}

// 按照一定数据结构解析数据项字符串
func (sp *SPath) parseCacheValue(value string) (filePath string, isDir bool) {
	if value[len(value)-2 : len(value)-1][0] == 'F' {
		return value[:len(value)-3], false
	}
	return value[:len(value)-3], true
}

// 添加path到缓存中(递归)
func (sp *SPath) addToCache(filePath, rootPath string) {
	// 首先添加自身
	idDir := gfile.IsDir(filePath)
	sp.cache.SetIfNotExist(sp.nameFromPath(filePath, rootPath), sp.makeCacheValue(filePath, idDir))
	// 如果添加的是目录，那么需要递归添加
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

// 添加文件目录监控(递归)，当目录下的文件有更新时，会同时更新缓存。
// 这里需要注意的点是，由于添加监听是递归添加的，那么假如删除一个目录，那么该目录下的文件(包括目录)也会产生一条删除事件，总共会产生N条事件。
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

// 删除监听(递归)
func (sp *SPath) removeMonitorByPath(path string) {
	if sp.cache == nil {
		return
	}
	_ = gfsnotify.Remove(path)
}
