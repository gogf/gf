// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmutex implements graceful concurrent-safe mutex with more rich features.
package gmutex

import (
	"sync"
	"sync/atomic"
)

// Mutex is a high level Mutex, which implements more rich features for mutex.
type Mutex struct {
	rwMutex sync.RWMutex
	wLocked int32
	rLocked int64
}

// New creates and returns a new mutex.
func New() *Mutex {
	return &Mutex{}
}

// Lock locks the mutex for writing purpose.
// If the mutex is already locked by another goroutine for reading or writing,
// it blocks until the lock is available.
func (m *Mutex) Lock() {
	m.rwMutex.Lock()
	atomic.StoreInt32(&m.wLocked, 1)
}

// Unlock unlocks writing lock on the mutex.
// It is safe to be called multiple times even there's no locks.
func (m *Mutex) Unlock() {
	if atomic.CompareAndSwapInt32(&m.wLocked, 1, 0) {
		m.rwMutex.Unlock()
	}
}

// TryLock tries locking the mutex for writing purpose.
// It returns true immediately if success, or if there's a write/reading lock on the mutex,
// it returns false immediately.
func (m *Mutex) TryLock() bool {
	locked := m.rwMutex.TryLock()
	if locked {
		atomic.StoreInt32(&m.wLocked, 1)
	}
	return locked
}

// RLock locks mutex for reading purpose.
// If the mutex is already locked for writing,
// it blocks until the lock is available.
func (m *Mutex) RLock() {
	m.rwMutex.RLock()
	atomic.AddInt64(&m.rLocked, 1)
}

// RUnlock unlocks the reading lock on the mutex.
// It is safe to be called multiple times even there's no locks.
func (m *Mutex) RUnlock() {
	for {
		rlocked := atomic.LoadInt64(&m.rLocked)
		if rlocked > 0 {
			if atomic.CompareAndSwapInt64(&m.rLocked, rlocked, rlocked-1) {
				m.rwMutex.RUnlock()
				return
			}
		} else {
			return
		}
	}
}

// TryRLock tries locking the mutex for reading purpose.
// It returns true immediately if success, or if there's a writing lock on the mutex,
// it returns false immediately.
func (m *Mutex) TryRLock() bool {
	locked := m.rwMutex.TryRLock()
	if locked {
		atomic.AddInt64(&m.rLocked, 1)
	}
	return locked
}

// IsLocked checks whether the mutex is locked with writing or reading lock.
// Note that the result might be changed after it's called,
// so it cannot be the criterion for atomic operations.
func (m *Mutex) IsLocked() bool {
	return m.IsWLocked() || m.IsRLocked()
}

// IsWLocked checks whether the mutex is locked by writing lock.
// Note that the result might be changed after it's called,
// so it cannot be the criterion for atomic operations.
func (m *Mutex) IsWLocked() bool {
	return atomic.LoadInt32(&m.wLocked) > 0
}

// IsRLocked checks whether the mutex is locked by reading lock.
// Note that the result might be changed after it's called,
// so it cannot be the criterion for atomic operations.
func (m *Mutex) IsRLocked() bool {
	return atomic.LoadInt64(&m.rLocked) > 0
}

// LockFunc locks the mutex for writing with given callback function `f`.
// If there's a write/reading lock the mutex, it will blocks until the lock is released.
//
// It releases the lock after `f` is executed.
func (m *Mutex) LockFunc(f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

// RLockFunc locks the mutex for reading with given callback function `f`.
// If there's a writing lock the mutex, it will blocks until the lock is released.
//
// It releases the lock after `f` is executed.
func (m *Mutex) RLockFunc(f func()) {
	m.RLock()
	defer m.RUnlock()
	f()
}

// TryLockFunc tries locking the mutex for writing with given callback function `f`.
// it returns true immediately if success, or if there's a write/reading lock on the mutex,
// it returns false immediately.
//
// It releases the lock after `f` is executed.
func (m *Mutex) TryLockFunc(f func()) (result bool) {
	if m.TryLock() {
		result = true
		defer m.Unlock()
		f()
	}
	return
}

// TryRLockFunc tries locking the mutex for reading with given callback function `f`.
// It returns true immediately if success, or if there's a writing lock on the mutex,
// it returns false immediately.
//
// It releases the lock after `f` is executed.
func (m *Mutex) TryRLockFunc(f func()) (result bool) {
	if m.TryRLock() {
		result = true
		defer m.RUnlock()
		f()
	}
	return
}
