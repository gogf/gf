// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gmlock

import (
    "gitee.com/johng/gf/g/container/gtype"
    "sync"
)

// 互斥锁对象
type Mutex struct {
    mu     sync.RWMutex
    wid    *gtype.Int64        // 当前Lock产生的唯一id(主要用于计时Unlock的校验)
    rcount *gtype.Int          // RLock次数
    wcount *gtype.Int          // Lock次数
}

// 创建一把内存锁使用的底层RWMutex
func NewMutex() *Mutex {
    return &Mutex{
        wid    : gtype.NewInt64(),
        rcount : gtype.NewInt(),
        wcount : gtype.NewInt(),
    }
}

func (l *Mutex) Lock() {
    l.wcount.Add(1)
    l.mu.Lock()
    l.wid.Add(1)
}

// 安全的Unlock
func (l *Mutex) Unlock() {
    if l.wcount.Val() > 0 {
        if l.wcount.Add(-1) >= 0 {
            l.mu.Unlock()
        } else {
            // 标准库这里会panic
            l.wcount.Add(1)
        }
    }
}

func (l *Mutex) RLock() {
    l.rcount.Add(1)
    l.mu.RLock()
}

// 安全的RUnlock
func (l *Mutex) RUnlock() {
    if l.rcount.Val() > 0 {
        if l.rcount.Add(-1) >= 0 {
            l.mu.RUnlock()
        } else {
            // 标准库这里会panic
            l.rcount.Add(1)
        }
    }
}

// 不阻塞Lock
func (l *Mutex) TryLock() bool {
    // 初步读写次数检查, 但无法保证原子性
    if l.wcount.Val() == 0 && l.rcount.Val() == 0 {
        // 第二次检查, 保证原子操作
        if l.wcount.Add(1) == 1 {
            l.mu.Lock()
            l.wid.Add(1)
            return true
        } else {
            l.wcount.Add(-1)
        }
    }
    return false
}

// 不阻塞RLock
func (l *Mutex) TryRLock() bool {
    // 只要不存在写锁
    if l.wcount.Val() == 0 {
        l.rcount.Add(1)
        l.mu.RLock()
        return true
    }
    return false
}

