// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/os/gtimer"
	"time"
)

// Memory locker.
type Locker struct {
	m *gmap.StrAnyMap
}

// New creates and returns a new memory locker.
// A memory locker can lock/unlock with dynamic string key.
func New() *Locker {
	return &Locker{
		m: gmap.NewStrAnyMap(),
	}
}

// TryLock tries locking the <key> with writing lock,
// it returns true if success, or it returns false if there's a writing/reading lock the <key>.
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) TryLock(key string, expire ...time.Duration) bool {
	return l.doLock(key, l.getExpire(expire...), true)
}

// Lock locks the <key> with writing lock.
// If there's a write/reading lock the <key>,
// it will blocks until the lock is released.
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) Lock(key string, expire ...time.Duration) {
	l.doLock(key, l.getExpire(expire...), false)
}

// Unlock unlocks the writing lock of the <key>.
func (l *Locker) Unlock(key string) {
	if v := l.m.Get(key); v != nil {
		v.(*Mutex).Unlock()
	}
}

// TryRLock tries locking the <key> with reading lock.
// It returns true if success, or if there's a writing lock on <key>, it returns false.
func (l *Locker) TryRLock(key string) bool {
	return l.doRLock(key, true)
}

// RLock locks the <key> with reading lock.
// If there's a writing lock on <key>,
// it will blocks until the writing lock is released.
func (l *Locker) RLock(key string) {
	l.doRLock(key, false)
}

// RUnlock unlocks the reading lock of the <key>.
func (l *Locker) RUnlock(key string) {
	if v := l.m.Get(key); v != nil {
		v.(*Mutex).RUnlock()
	}
}

// TryLockFunc locks the <key> with writing lock and callback function <f>.
// It returns true if success, or else if there's a write/reading lock the <key>, it return false.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) TryLockFunc(key string, f func(), expire ...time.Duration) bool {
	if l.TryLock(key, expire...) {
		defer l.Unlock(key)
		f()
		return true
	}
	return false
}

// TryRLockFunc locks the <key> with reading lock and callback function <f>.
// It returns true if success, or else if there's a writing lock the <key>, it returns false.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) TryRLockFunc(key string, f func()) bool {
	if l.TryRLock(key) {
		defer l.RUnlock(key)
		f()
		return true
	}
	return false
}

// LockFunc locks the <key> with writing lock and callback function <f>.
// If there's a write/reading lock the <key>,
// it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) LockFunc(key string, f func(), expire ...time.Duration) {
	l.Lock(key, expire...)
	defer l.Unlock(key)
	f()
}

// RLockFunc locks the <key> with reading lock and callback function <f>.
// If there's a writing lock the <key>,
// it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) RLockFunc(key string, f func()) {
	l.RLock(key)
	defer l.RUnlock(key)
	f()
}

// getExpire returns the duration object passed.
// If <expire> is not passed, it returns a default duration object.
func (l *Locker) getExpire(expire ...time.Duration) time.Duration {
	e := time.Duration(0)
	if len(expire) > 0 {
		e = expire[0]
	}
	return e
}

// doLock locks writing on <key>.
// It returns true if success, or else returns false.
//
// The parameter <try> is true,
// it returns false immediately if it fails getting the writing lock.
// If <true> is false, it blocks until it gets the writing lock.
//
// The parameter <expire> specifies the max duration it locks.
func (l *Locker) doLock(key string, expire time.Duration, try bool) bool {
	mu := l.getOrNewMutex(key)
	ok := true
	if try {
		ok = mu.TryLock()
	} else {
		mu.Lock()
	}
	if ok && expire > 0 {
		wid := mu.wid.Val()
		gtimer.AddOnce(expire, func() {
			if wid == mu.wid.Val() {
				mu.Unlock()
			}
		})
	}
	return ok
}

// doRLock locks reading on <key>.
// It returns true if success, or else returns false.
//
// The parameter <try> is true,
// it returns false immediately if it fails getting the reading lock.
// If <true> is false, it blocks until it gets the reading lock.
func (l *Locker) doRLock(key string, try bool) bool {
	mu := l.getOrNewMutex(key)
	ok := true
	if try {
		ok = mu.TryRLock()
	} else {
		mu.RLock()
	}
	return ok
}

// getOrNewMutex returns the mutex of given <key> if it exists,
// or else creates and returns a new one.
func (l *Locker) getOrNewMutex(key string) *Mutex {
	return l.m.GetOrSetFuncLock(key, func() interface{} {
		return NewMutex()
	}).(*Mutex)
}
