// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock

import (
	"github.com/gogf/gf/g/container/gtype"
	"runtime"
	"sync"
)

// The high level RWMutex.
// It wraps the sync.RWMutex to implements more rich features.
type Mutex struct {
    mu      sync.RWMutex
    wid     *gtype.Int64 // Unique id, used for multiple and safe logic Unlock.
    locking *gtype.Bool  // Locking mark for atomic operation for *Lock and Try*Lock functions.
                         // There must be only one locking operation at the same time for concurrent safe purpose.
    state   *gtype.Int32 // Locking state:
                         //   0: writing lock false;
                         //  -1: writing lock true;
                         // >=1: reading lock;
}

// NewMutex creates and returns a new mutex.
func NewMutex() *Mutex {
    return &Mutex{
        wid     : gtype.NewInt64(),
	    state   : gtype.NewInt32(),
		locking : gtype.NewBool(),
    }
}

// Lock locks the mutex for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (m *Mutex) Lock() {
	if m.locking.Cas(false, true) {
		m.mu.Lock()
		// State should be changed after locks.
		m.state.Set(-1)
		m.wid.Add(1)
		m.locking.Set(false)
	} else {
		runtime.Gosched()
		m.Lock()
	}
}

// Unlock unlocks the writing lock.
// It is safe to be called multiple times if there's any locks or not.
func (m *Mutex) Unlock() {
	if m.state.Cas(-1, 0) {
		m.mu.Unlock()
	}
}

// TryLock tries locking the mutex for writing.
// It returns true if success, or if there's a write/reading lock on the mutex,
// it returns false.
func (m *Mutex) TryLock() bool {
	if m.locking.Cas(false, true) {
		m.mu.Lock()
		// State should be changed after locks.
		m.state.Set(-1)
		m.wid.Add(1)
		m.locking.Set(false)
		return true
	}
	return false
}

// RLock locks mutex for reading.
// If the mutex is already locked for writing,
// It blocks until the lock is available.
func (m *Mutex) RLock() {
	if m.locking.Cas(false, true) {
		m.mu.RLock()
		// State should be changed after locks.
		m.state.Add(1)
		m.locking.Set(false)
	} else {
		runtime.Gosched()
		m.RLock()
	}
}

// RUnlock unlocks the reading lock.
// It is safe to be called multiple times if there's any locks or not.
func (m *Mutex) RUnlock() {
    if n := m.state.Val(); n >= 1 {
    	if m.state.Cas(n, n - 1) {
			m.mu.RUnlock()
		} else {
			m.RUnlock()
		}
    }
}

// TryRLock tries locking the mutex for reading.
// It returns true if success, or if there's a writing lock on the mutex, it returns false.
func (m *Mutex) TryRLock() bool {
	if m.locking.Cas(false, true) {
		if m.state.Val() >= 0 {
			m.mu.RLock()
			m.state.Add(1)
			m.locking.Set(false)
			return true
		}
    }
    return false
}

// TryLockFunc tries locking the mutex for writing with given callback function <f>.
// it returns true if success, or if there's a write/reading lock on the mutex,
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
// It returns true if success, or if there's a writing lock on the mutex, it returns false.
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
// If there's a write/reading lock the mutex, it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
func (m *Mutex) LockFunc(f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

// RLockFunc locks the mutex for reading with given callback function <f>.
// If there's a writing lock the mutex, it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
func (m *Mutex) RLockFunc(f func()) {
	m.RLock()
	defer m.RUnlock()
	f()
}
