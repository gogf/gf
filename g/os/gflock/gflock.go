// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gflock implements a concurrent-safe sync.Locker interface for file locking.
package gflock

import (
    "sync"
    "github.com/gogf/gf/third/github.com/theckman/go-flock"
    "github.com/gogf/gf/g/os/gfile"
)

// File locker.
type Locker struct {
    mu    sync.RWMutex // 用于外部接口调用的互斥锁(阻塞机制)
    flock *flock.Flock // 底层文件锁对象
}

// New creates and returns a new file locker with given <file>.
// The parameter <file> usually is a absolute file path.
func New(file string) *Locker {
    dir := gfile.TempDir() + gfile.Separator + "gflock"
    if !gfile.Exists(dir) {
        _ = gfile.Mkdir(dir)
    }
    path := dir + gfile.Separator + file
    lock := flock.NewFlock(path)
    return &Locker{
        flock : lock,
    }
}

// Path returns the file path of the locker.
func (l *Locker) Path() string {
    return l.flock.Path()
}

// IsLocked returns whether the locker is locked.
func (l *Locker) IsLocked() bool {
    return l.flock.Locked()
}

// TryLock tries get the writing lock of the locker.
// It returns true if success, or else returns false immediately.
func (l *Locker) TryLock() bool {
    ok, _ := l.flock.TryLock()
    if ok {
        l.mu.Lock()
    }
    return ok
}

// TryRLock tries get the reading lock of the locker.
// It returns true if success, or else returns false immediately.
func (l *Locker) TryRLock() bool {
    ok, _ := l.flock.TryRLock()
    if ok {
        l.mu.RLock()
    }
    return ok
}

func (l *Locker) Lock() (err error) {
    l.mu.Lock()
    err = l.flock.Lock()
    return
}

func (l *Locker) UnLock() (err error) {
    err = l.flock.Unlock()
    l.mu.Unlock()
    return
}

func (l *Locker) RLock() (err error) {
    l.mu.RLock()
    err = l.flock.RLock()
    return
}

func (l *Locker) RUnlock() (err error) {
    err = l.flock.Unlock()
    l.mu.RUnlock()
    return
}
