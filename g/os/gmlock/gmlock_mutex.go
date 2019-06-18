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

// The high level Mutex.
// It wraps the sync.RWMutex to implements more rich features.
type Mutex struct {
    mu     sync.RWMutex
    wid    *gtype.Int64 // Unique id, used for multiple safely Unlock.
    number *gtype.Int   // Locking number.
                        //   0: writing lock false;
                        //   1: writing lock true;
                        // >=2: reading lock;
}

// NewMutex creates and returns a new mutex.
func NewMutex() *Mutex {
    return &Mutex{
        wid    : gtype.NewInt64(),
	    number : gtype.NewInt(),
    }
}

// Lock locks mutex for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (m *Mutex) Lock() {
	m.mu.Lock()
	m.number.Set(1)
	m.wid.Add(1)
}

// Unlock unlocks the write lock.
// It is safe to be called multiple times.
func (m *Mutex) Unlock() {
	if m.number.Cas(1, 0) {
		m.mu.Unlock()
	}
}

// RLock locks mutex for reading.
// If the mutex is already locked for writing,
// It blocks until the lock is available.
func (m *Mutex) RLock() {
    m.mu.RLock()
	m.number.Add(2)
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if mutex is not locked for reading
// on entry to RUnlock.
// It is safe to be called multiple times.
func (m *Mutex) RUnlock() {
    if n := m.number.Val(); n >= 2 && n%2 == 0 {
        if n := m.number.Add(-2); n >= 0 && n%2 == 0 {
            m.mu.RUnlock()
        } else {
            m.number.Add(2)
        }
    }
}

// TryLock tries locking the mutex for writing.
// It returns true if success, or if there's a write/read lock on the mutex,
// it returns false.
func (m *Mutex) TryLock() bool {
    if m.number.Cas(0, 1) {
        m.mu.Lock()
        m.wid.Add(1)
        return true
    }
    return false
}

// TryRLock tries locking the mutex for reading.
// It returns true if success, or if there's a writing lock on the mutex, it returns false.
//
// Note that it's not an atomic operation.
//
// TODO It should be an atomic operation.
func (m *Mutex) TryRLock() bool {
    // There must be no writing lock on the mutex.
    if n := m.number.Val(); n%2 == 0 {
        m.mu.RLock()
	    m.number.Add(2)
        return true
    }
    return false
}

// TryLockFunc tries locking the mutex for writing with given callback function <f>.
// it returns true if success, or if there's a write/read lock on the mutex,
// it returns false.
//
// It releases the lock after <f> is executed.
func (m *Mutex) TryLockFunc(f func()) bool {
	if m.TryLock() {
		defer m.Unlock()
		f()
		return true
	}
	return false
}

// TryRLockFunc tries locking the mutex for reading with given callback function <f>.
// It returns true if success, or if there's a write lock on the mutex, it returns false.
//
// It releases the lock after <f> is executed.
func (m *Mutex) TryRLockFunc(f func()) bool {
	if m.TryRLock() {
		defer m.RUnlock()
		f()
		return true
	}
	return false
}

// LockFunc locks the mutex for writing with given callback function <f>.
// If there's a write/read lock the mutex, it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
func (m *Mutex) LockFunc(f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

// RLockFunc locks the mutex for reading with given callback function <f>.
// If there's a write lock the mutex, it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
func (m *Mutex) RLockFunc(f func()) {
	m.RLock()
	defer m.RUnlock()
	f()
}
