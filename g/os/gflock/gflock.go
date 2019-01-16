// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gflock implements a thread-safe sync.Locker interface for file locking.
// 
// 文件锁.
package gflock

import (
    "sync"
    "gitee.com/johng/gf/third/github.com/theckman/go-flock"
    "gitee.com/johng/gf/g/os/gfile"
)

// 文件锁
type Locker struct {
    mu    sync.RWMutex // 用于外部接口调用的互斥锁(阻塞机制)
    flock *flock.Flock // 底层文件锁对象
}

// 创建文件锁
func New(file string) *Locker {
    dir := gfile.TempDir() + gfile.Separator + "gflock"
    if !gfile.Exists(dir) {
        gfile.Mkdir(dir)
    }
    path := dir + gfile.Separator + file
    lock := flock.NewFlock(path)
    return &Locker{
        flock : lock,
    }
}

func (l *Locker) Path() string {
    return l.flock.Path()
}

// 当前文件锁是否处于锁定状态(Lock)
func (l *Locker) IsLocked() bool {
    return l.flock.Locked()
}

// 尝试Lock文件，如果失败立即返回
func (l *Locker) TryLock() bool {
    ok, _ := l.flock.TryLock()
    if ok {
        l.mu.Lock()
    }
    return ok
}

// 尝试RLock文件，如果失败立即返回
func (l *Locker) TryRLock() bool {
    ok, _ := l.flock.TryRLock()
    if ok {
        l.mu.RLock()
    }
    return ok
}

func (l *Locker) Lock() {
    l.mu.Lock()
    l.flock.Lock()
}

func (l *Locker) UnLock() {
    l.flock.Unlock()
    l.mu.Unlock()
}

func (l *Locker) RLock() {
    l.mu.RLock()
    l.flock.RLock()
}

func (l *Locker) RUnlock() {
    l.flock.Unlock()
    l.mu.RUnlock()
}
