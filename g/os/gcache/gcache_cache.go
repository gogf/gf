// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcache

import (
    "sync/atomic"
    "unsafe"
)

// 缓存对象。
// 底层只有一个缓存对象，如果需要提高并发性能，可新增缓存对象无锁哈希表，用键名做固定分区。
type Cache struct {
    *memCache // 底层缓存对象
}

// Cache对象按照缓存键名首字母做了分组
func New() *Cache {
    return &Cache {
        memCache : newMemCache(),
    }
}

// 清空缓存中的所有数据
func (c *Cache) Clear() {
    atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.memCache)), unsafe.Pointer(newMemCache()))
}