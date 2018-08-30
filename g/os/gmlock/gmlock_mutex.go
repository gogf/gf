// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gmlock

import (
    "sync"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gtime"
)

// 互斥锁对象
type Mutex struct {
    mu     sync.RWMutex
    rid    *gtype.Int64        // 当前RLock产生的唯一id(主要用于计时RUnlock的校验)
    wid    *gtype.Int64        // 当前Lock产生的唯一id(主要用于计时Unlock的校验)
    rcount *gtype.Int          // RLock次数
    wcount *gtype.Int          // Lock次数
}

// 创建一把内存锁使用的底层RWMutex
func NewMutex() *Mutex {
    return &Mutex{
        rid    : gtype.NewInt64(),
        wid    : gtype.NewInt64(),
        rcount : gtype.NewInt(),
        wcount : gtype.NewInt(),
    }
}

// 不阻塞Lock
func (l *Mutex) TryLock() bool {
    if l.wcount.Val() == 0 && l.rcount.Val() == 0 {
        l.Lock()
        return true
    }
    return false
}

func (l *Mutex) Lock() {
    l.wcount.Add(1)
    l.mu.Lock()
    l.wid.Set(gtime.Nanosecond())
}

// 安全的Unlock
func (l *Mutex) Unlock() {
    if l.wcount.Val() > 0 {
        l.mu.Unlock()
        l.wcount.Add(-1)
    }
}

// 不阻塞RLock
func (l *Mutex) TryRLock() bool {
    if l.wcount.Val() == 0 {
        l.RLock()
        return true
    }
    return false
}

func (l *Mutex) RLock() {
    l.rcount.Add(1)
    l.mu.RLock()
    l.rid.Set(gtime.Nanosecond())
}

// 安全的RUnlock
func (l *Mutex) RUnlock() {
    if l.rcount.Val() > 0 {
        l.mu.RUnlock()
        l.rcount.Add(-1)
    }
}
