// Copyright 2015 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

// +build !windows

package flock

import (
	"syscall"
)

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
func (f *Flock) Lock() error {
	return f.lock(&f.l, syscall.LOCK_EX)
}

// RLock is a blocking call to try and take a ahred file lock. It will wait
// until it is able to obtain the shared file lock. It's recommended that
// TryRLock() be used over this function. This function may block the ability to
// query the current Locked() or RLocked() status due to a RW-mutex lock.
//
// If we are already shared-locked, this function short-circuits and returns
// immediately assuming it can take the mutex lock.
func (f *Flock) RLock() error {
	return f.lock(&f.r, syscall.LOCK_SH)
}

func (f *Flock) lock(locked *bool, flag int) error {
	f.m.Lock()
	defer f.m.Unlock()

	if *locked {
		return nil
	}

	if f.fh == nil {
		if err := f.setFh(); err != nil {
			return err
		}
	}

	if err := syscall.Flock(int(f.fh.Fd()), flag); err != nil {
		return err
	}

	*locked = true
	return nil
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
func (f *Flock) Unlock() error {
	f.m.Lock()
	defer f.m.Unlock()

	// if we aren't locked or if the lockfile instance is nil
	// just return a nil error because we are unlocked
	if (!f.l && !f.r) || f.fh == nil {
		return nil
	}

	// mark the file as unlocked
	if err := syscall.Flock(int(f.fh.Fd()), syscall.LOCK_UN); err != nil {
		return err
	}

	f.fh.Close()

	f.l = false
	f.r = false
	f.fh = nil

	return nil
}

// TryLock is the preferred function for taking an exclusive file lock. This
// function takes an RW-mutex lock before it tries to lock the file, so there is
// the possibility that this function may block for a short time if another
// goroutine is trying to take any action.
//
// The actual file lock is non-blocking. If we are unable to get the exclusive
// file lock, the function will return false instead of waiting for the lock. If
// we get the lock, we also set the *Flock instance as being exclusive-locked.
func (f *Flock) TryLock() (bool, error) {
	return f.try(&f.l, syscall.LOCK_EX)
}

// TryRLock is the preferred function for taking a shared file lock. This
// function takes an RW-mutex lock before it tries to lock the file, so there is
// the possibility that this function may block for a short time if another
// goroutine is trying to take any action.
//
// The actual file lock is non-blocking. If we are unable to get the shared file
// lock, the function will return false instead of waiting for the lock. If we
// get the lock, we also set the *Flock instance as being share-locked.
func (f *Flock) TryRLock() (bool, error) {
	return f.try(&f.r, syscall.LOCK_SH)
}

func (f *Flock) try(locked *bool, flag int) (bool, error) {
	f.m.Lock()
	defer f.m.Unlock()

	if *locked {
		return true, nil
	}

	if f.fh == nil {
		if err := f.setFh(); err != nil {
			return false, err
		}
	}

	err := syscall.Flock(int(f.fh.Fd()), flag|syscall.LOCK_NB)

	switch err {
	case syscall.EWOULDBLOCK:
		return false, nil
	case nil:
		*locked = true
		return true, nil
	}

	return false, err
}
