// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gflock implements a concurrent-safe sync.Locker interface for file locking.
package gflock

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/third/github.com/theckman/go-flock"
)

// File locker.
type Locker struct {
	flock *flock.Flock // Underlying file locker.
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
		flock: lock,
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

// IsRLocked returns whether the locker is rlocked.
func (l *Locker) IsRLocked() bool {
	return l.flock.RLocked()
}

// TryLock tries get the writing lock of the locker.
// It returns true if success, or else returns false immediately.
func (l *Locker) TryLock() bool {
	ok, _ := l.flock.TryLock()
	return ok
}

// TryRLock tries get the reading lock of the locker.
// It returns true if success, or else returns false immediately.
func (l *Locker) TryRLock() bool {
	ok, _ := l.flock.TryRLock()
	return ok
}

// Lock is a blocking call to try and take an exclusive file lock. It will wait
// until it is able to obtain the exclusive file lock. It's recommended that
// TryLock() be used over this function. This function may block the ability to
// query the current Locked() or RLocked() status due to a RW-mutex lock.
//
// If we are already exclusive-locked, this function short-circuits and returns
// immediately assuming it can take the mutex lock.
//
// If the *Flock has a shared lock (RLock), this may transparently replace the
// shared lock with an exclusive lock on some UNIX-like operating systems. Be
// careful when using exclusive locks in conjunction with shared locks
// (RLock()), because calling Unlock() may accidentally release the exclusive
// lock that was once a shared lock.
func (l *Locker) Lock() (err error) {
	return l.flock.Lock()
}

// Unlock is a function to unlock the file. This file takes a RW-mutex lock, so
// while it is running the Locked() and RLocked() functions will be blocked.
//
// This function short-circuits if we are unlocked already. If not, it calls
// syscall.LOCK_UN on the file and closes the file descriptor. It does not
// remove the file from disk. It's up to your application to do.
//
// Please note, if your shared lock became an exclusive lock this may
// unintentionally drop the exclusive lock if called by the consumer that
// believes they have a shared lock. Please see Lock() for more details.
func (l *Locker) Unlock() (err error) {
	return l.flock.Unlock()
}

// RLock is a blocking call to try and take a ahred file lock. It will wait
// until it is able to obtain the shared file lock. It's recommended that
// TryRLock() be used over this function. This function may block the ability to
// query the current Locked() or RLocked() status due to a RW-mutex lock.
//
// If we are already shared-locked, this function short-circuits and returns
// immediately assuming it can take the mutex lock.
func (l *Locker) RLock() (err error) {
	return l.flock.RLock()
}

// Unlock is a function to unlock the file. This file takes a RW-mutex lock, so
// while it is running the Locked() and RLocked() functions will be blocked.
//
// This function short-circuits if we are unlocked already. If not, it calls
// syscall.LOCK_UN on the file and closes the file descriptor. It does not
// remove the file from disk. It's up to your application to do.
//
// Please note, if your shared lock became an exclusive lock this may
// unintentionally drop the exclusive lock if called by the consumer that
// believes they have a shared lock. Please see Lock() for more details.
func (l *Locker) RUnlock() (err error) {
	return l.flock.Unlock()
}
