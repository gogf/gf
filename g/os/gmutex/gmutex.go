// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmutex implements graceful concurrent-safe mutex with more rich features.
package gmutex

import (
	"math"
	"runtime"

	"github.com/gogf/gf/g/container/gtype"
)

// The high level Mutex, which implements more rich features for mutex.
type Mutex struct {
	state   *gtype.Int32  // Indicates the state of mutex.
	writer  *gtype.Int32  // Pending writer count.
	reader  *gtype.Int32  // Pending reader count.
	writing chan struct{} // Channel for writer blocking.
	reading chan struct{} // Channel for reader blocking.
}

// New creates and returns a new mutex.
func New() *Mutex {
	return &Mutex{
		state:   gtype.NewInt32(),
		writer:  gtype.NewInt32(),
		reader:  gtype.NewInt32(),
		writing: make(chan struct{}, 1),
		reading: make(chan struct{}, math.MaxInt32),
	}
}

// Lock locks the mutex for writing.
// If the lock is already locked for reading or writing,
// it blocks until the lock is available.
func (m *Mutex) Lock() {
	for {
		// Using CAS operation to get the writing lock atomically.
		if m.state.Cas(0, -1) {
			return
		}
		// Or else blocks to wait for the next chance.
		m.writer.Add(1)
		<-m.writing
		m.writer.Add(-1)
	}
}

// Unlock unlocks the writing lock.
// It is safe to be called multiple times if there's any locks or not.
func (m *Mutex) Unlock() {
	if m.state.Cas(-1, 0) {
		// Writing lock unlocks, then first check the blocked readers.
		// If there're readers blocked, unlock them with preemption.
		for n := m.reader.Val(); n > 0; n-- {
			m.reading <- struct{}{}
		}
		// It then also kindly feeds the pending writers with one chance.
		if m.writer.Val() > 0 {
			m.writing <- struct{}{}
		}
	}
}

// TryLock tries locking the mutex for writing.
// It returns true immediately if success, or if there's a write/reading lock on the mutex,
// it returns false immediately.
func (m *Mutex) TryLock() bool {
	if m.state.Cas(0, -1) {
		return true
	}
	return false
}

// RLock locks mutex for reading purpose.
// If the mutex is already locked for writing,
// it blocks until the lock is available.
func (m *Mutex) RLock() {
	var n int32
	for {
		// If there're no writing lock currently, then do the reading lock checks.
		if n = m.state.Val(); n >= 0 {
			if m.state.Cas(n, n+1) {
				return
			} else {
				runtime.Gosched()
				continue
			}
		}

		// Or else pending the reader.
		m.reader.Add(1)
		<-m.reading
		m.reader.Add(-1)
	}
}

// RUnlock unlocks the reading lock.
// It is safe to be called multiple times if there's any locks or not.
func (m *Mutex) RUnlock() {
	var n int32
	for {
		if n = m.state.Val(); n >= 1 {
			if m.state.Cas(n, n-1) {
				break
			} else {
				runtime.Gosched()
			}
		} else {
			break
		}
	}
	// Reading lock unlocks, then only check the blocked writers.
	// Note that it is not necessary to check the pending readers here.
	if n == 1 && m.writer.Val() > 0 {
		m.writing <- struct{}{}
	}
}

// TryRLock tries locking the mutex for reading.
// It returns true immediately if success, or if there's a writing lock on the mutex,
// it returns false immediately.
func (m *Mutex) TryRLock() bool {
	var n int32
	for {
		if n = m.state.Val(); n >= 0 {
			if m.state.Cas(n, n+1) {
				return true
			} else {
				runtime.Gosched()
			}
		} else {
			return false
		}
	}
}

// IsLocked checks whether the mutex is locked with writing or reading lock.
// Note that the result might be changed after it's called,
// so it cannot be the criterion for atomic operations.
func (m *Mutex) IsLocked() bool {
	return m.state.Val() != 0
}

// IsWLocked checks whether the mutex is locked by writing lock.
// Note that the result might be changed after it's called,
// so it cannot be the criterion for atomic operations.
func (m *Mutex) IsWLocked() bool {
	return m.state.Val() < 0
}

// IsRLocked checks whether the mutex is locked by reading lock.
// Note that the result might be changed after it's called,
// so it cannot be the criterion for atomic operations.
func (m *Mutex) IsRLocked() bool {
	return m.state.Val() > 0
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

// TryLockFunc tries locking the mutex for writing with given callback function <f>.
// it returns true immediately if success, or if there's a write/reading lock on the mutex,
// it returns false immediately.
//
// It releases the lock after <f> is executed.
func (m *Mutex) TryLockFunc(f func()) (result bool) {
	if m.TryLock() {
		result = true
		defer m.Unlock()
		f()
	}
	return
}

// TryRLockFunc tries locking the mutex for reading with given callback function <f>.
// It returns true immediately if success, or if there's a writing lock on the mutex,
// it returns false immediately.
//
// It releases the lock after <f> is executed.
func (m *Mutex) TryRLockFunc(f func()) (result bool) {
	if m.TryRLock() {
		result = true
		defer m.RUnlock()
		f()
	}
	return
}
