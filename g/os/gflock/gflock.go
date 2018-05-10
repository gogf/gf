// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件锁.
package gflock

import (
    "sync"
    "github.com/theckman/go-flock"
    "gitee.com/johng/gf/g/os/gfile"
)

// 文件锁
type Locker struct {
    mu    sync.RWMutex
    flock *flock.Flock
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
