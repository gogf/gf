// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件缓存.
package gfcache

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gfsnotify"
)

// 设置容量大小(MB)
func (c *Cache) SetCap(cap int) {
    c.cap.Set(cap)
}

// 获得缓存容量大小(byte)
func (c *Cache) GetCap() int {
    return c.cap.Val()
}

// 获得已缓存的文件大小(byte)
func (c *Cache) GetSize() int {
    return c.size.Val()
}

// 获得文件内容 string
func (c *Cache) GetContents(path string) string {
    return string(c.GetBinContents(path))
}

// 获得文件内容 []byte
func (c *Cache) GetBinContents(path string) []byte {
    v := c.cache.Get(path)
    if v != nil {
        return v.([]byte)
    }
    b := gfile.GetBinContents(path)
    if b != nil && (c.cap.Val() == 0 || c.size.Val() < c.cap.Val()) {
        c.addMonitor(path)
        c.cache.Set(path, b, 0)
        c.size.Add(len(b))
    }
    return b
}

// 添加文件监控
func (c *Cache) addMonitor(path string) {
    // 防止多goroutine同时调用
    if c.cache.Get(path) != nil {
        return
    }
    gfsnotify.Add(path, func(event *gfsnotify.Event) {
        //glog.Debug("gfcache:", event)
        r := c.cache.Get(path).([]byte)
        // 是否删除
        if event.IsRemove() {
            c.cache.Remove(path)
            c.size.Add(-len(r))
        }
        // 更新缓存内容
        if c.cap.Val() == 0 || c.size.Val() < c.cap.Val() {
            b := gfile.GetBinContents(path)
            if b != nil {
                dif := len(b) - len(r)
                c.cache.Set(path, b, 0)
                c.size.Add(dif)
            }
        }
    })
}