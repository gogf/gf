// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock

import (
    "github.com/gogf/gf/g/container/gtype"
    "sync"
)

// The mutex.
type Mutex struct {
    mu     sync.RWMutex
    wid    *gtype.Int64        // Unique id for this mutex.
    rcount *gtype.Int          // RLock count.
    wcount *gtype.Int          // Lock count.
}

// NewMutex creates and returns a new mutex.
func NewMutex() *Mutex {
    return &Mutex{
        wid    : gtype.NewInt64(),
        rcount : gtype.NewInt(),
        wcount : gtype.NewInt(),
    }
}

// Lock locks mutex for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (l *Mutex) Lock() {
    l.wcount.Add(1)
    l.mu.Lock()
    l.wid.Add(1)
}

// Unlock unlocks the write lock.
// It is safe to be called multiple times.
func (l *Mutex) Unlock() {
    if l.wcount.Val() > 0 {
        if l.wcount.Add(-1) >= 0 {
            l.mu.Unlock()
        } else {
            l.wcount.Add(1)
        }
    }
}

// RLock locks mutex for reading.
func (l *Mutex) RLock() {
    l.rcount.Add(1)
    l.mu.RLock()
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
// It is safe to be called multiple times.
func (l *Mutex) RUnlock() {
    if l.rcount.Val() > 0 {
        if l.rcount.Add(-1) >= 0 {
            l.mu.RUnlock()
        } else {
            l.rcount.Add(1)
        }
    }
}

// TryLock tries locking the mutex with write lock,
// it returns true if success, or if there's a write/read lock the mutex,
// it returns false.
func (l *Mutex) TryLock() bool {
    // The first check, but it cannot ensure the atomicity.
    if l.wcount.Val() == 0 && l.rcount.Val() == 0 {
        // The second check, it can ensure the atomicity.
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

// TryRLock tries locking the mutex with read lock.
// It returns true if success, or if there's a write lock on mutex, it returns false.
func (l *Mutex) TryRLock() bool {
    // There must be no write lock on mutex.
    if l.wcount.Val() == 0 {
        l.rcount.Add(1)
        l.mu.RLock()
        return true
    }
    return false
}

