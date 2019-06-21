// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmlock implements a concurrent-safe memory-based locker.
package gmlock

import "time"

var (
	// Default locker.
	locker = New()
)

// TryLock tries locking the <key> with writing lock,
// it returns true if success, or if there's a write/reading lock the <key>,
// it returns false. The parameter <expire> specifies the max duration it locks.
func TryLock(key string, expire ...time.Duration) bool {
	return locker.TryLock(key, expire...)
}

// Lock locks the <key> with writing lock.
// If there's a write/reading lock the <key>,
// it will blocks until the lock is released.
// The parameter <expire> specifies the max duration it locks.
func Lock(key string, expire ...time.Duration) {
	locker.Lock(key, expire...)
}

// Unlock unlocks the writing lock of the <key>.
func Unlock(key string) {
	locker.Unlock(key)
}

// TryRLock tries locking the <key> with reading lock.
// It returns true if success, or if there's a writing lock on <key>, it returns false.
func TryRLock(key string) bool {
	return locker.TryRLock(key)
}

// RLock locks the <key> with reading lock.
// If there's a writing lock on <key>,
// it will blocks until the writing lock is released.
func RLock(key string) {
	locker.RLock(key)
}

// RUnlock unlocks the reading lock of the <key>.
func RUnlock(key string) {
	locker.RUnlock(key)
}

// TryLockFunc locks the <key> with writing lock and callback function <f>.
// It returns true if success, or else if there's a write/reading lock the <key>, it return false.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func TryLockFunc(key string, f func(), expire ...time.Duration) bool {
	return locker.TryLockFunc(key, f, expire...)
}

// TryRLockFunc locks the <key> with reading lock and callback function <f>.
// It returns true if success, or else if there's a writing lock the <key>, it returns false.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func TryRLockFunc(key string, f func()) bool {
	return locker.TryRLockFunc(key, f)
}

// LockFunc locks the <key> with writing lock and callback function <f>.
// If there's a write/reading lock the <key>,
// it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func LockFunc(key string, f func(), expire ...time.Duration) {
	locker.LockFunc(key, f, expire...)
}

// RLockFunc locks the <key> with reading lock and callback function <f>.
// If there's a writing lock the <key>,
// it will blocks until the lock is released.
//
// It releases the lock after <f> is executed.
//
// The parameter <expire> specifies the max duration it locks.
func RLockFunc(key string, f func()) {
	locker.RLockFunc(key, f)
}
