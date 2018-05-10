// Copyright 2015 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package flock_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/theckman/go-flock"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	path  string
	flock *flock.Flock
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) SetUpTest(c *C) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "go-flock-")
	c.Assert(err, IsNil)
	c.Assert(tmpFile, Not(IsNil))

	t.path = tmpFile.Name()

	defer os.Remove(t.path)
	tmpFile.Close()

	t.flock = flock.NewFlock(t.path)
}

func (t *TestSuite) TearDownTest(c *C) {
	t.flock.Unlock()
	os.Remove(t.path)
}

func (t *TestSuite) TestNewFlock(c *C) {
	var f *flock.Flock

	f = flock.NewFlock(t.path)
	c.Assert(f, Not(IsNil))
	c.Check(f.Path(), Equals, t.path)
	c.Check(f.Locked(), Equals, false)
	c.Check(f.RLocked(), Equals, false)
}

func (t *TestSuite) TestFlock_Path(c *C) {
	var path string
	path = t.flock.Path()
	c.Check(path, Equals, t.path)
}

func (t *TestSuite) TestFlock_Locked(c *C) {
	var locked bool
	locked = t.flock.Locked()
	c.Check(locked, Equals, false)
}

func (t *TestSuite) TestFlock_RLocked(c *C) {
	var locked bool
	locked = t.flock.RLocked()
	c.Check(locked, Equals, false)
}

func (t *TestSuite) TestFlock_String(c *C) {
	var str string
	str = t.flock.String()
	c.Assert(str, Equals, t.path)
}

func (t *TestSuite) TestFlock_TryLock(c *C) {
	c.Assert(t.flock.Locked(), Equals, false)
	c.Assert(t.flock.RLocked(), Equals, false)

	var locked bool
	var err error

	locked, err = t.flock.TryLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)
	c.Check(t.flock.Locked(), Equals, true)
	c.Check(t.flock.RLocked(), Equals, false)

	locked, err = t.flock.TryLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)

	// make sure we just return false with no error in cases
	// where we would have been blocked
	locked, err = flock.NewFlock(t.path).TryLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, false)
}

func (t *TestSuite) TestFlock_TryRLock(c *C) {
	c.Assert(t.flock.Locked(), Equals, false)
	c.Assert(t.flock.RLocked(), Equals, false)

	var locked bool
	var err error

	locked, err = t.flock.TryRLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)
	c.Check(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, true)

	locked, err = t.flock.TryRLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)

	// shared lock should not block.
	flock2 := flock.NewFlock(t.path)
	locked, err = flock2.TryRLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)

	// make sure we just return false with no error in cases
	// where we would have been blocked
	t.flock.Unlock()
	flock2.Unlock()
	t.flock.Lock()
	locked, err = flock.NewFlock(t.path).TryRLock()
	c.Assert(err, IsNil)
	c.Check(locked, Equals, false)
}

func (t *TestSuite) TestFlock_TryLockContext(c *C) {
	// happy path
	ctx, cancel := context.WithCancel(context.Background())
	locked, err := t.flock.TryLockContext(ctx, time.Second)
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)

	// context already canceled
	cancel()
	locked, err = flock.NewFlock(t.path).TryLockContext(ctx, time.Second)
	c.Assert(err, Equals, context.Canceled)
	c.Check(locked, Equals, false)

	// timeout
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	locked, err = flock.NewFlock(t.path).TryLockContext(ctx, time.Second)
	c.Assert(err, Equals, context.DeadlineExceeded)
	c.Check(locked, Equals, false)
}

func (t *TestSuite) TestFlock_TryRLockContext(c *C) {
	// happy path
	ctx, cancel := context.WithCancel(context.Background())
	locked, err := t.flock.TryRLockContext(ctx, time.Second)
	c.Assert(err, IsNil)
	c.Check(locked, Equals, true)

	// context already canceled
	cancel()
	locked, err = flock.NewFlock(t.path).TryRLockContext(ctx, time.Second)
	c.Assert(err, Equals, context.Canceled)
	c.Check(locked, Equals, false)

	// timeout
	t.flock.Unlock()
	t.flock.Lock()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	locked, err = flock.NewFlock(t.path).TryRLockContext(ctx, time.Second)
	c.Assert(err, Equals, context.DeadlineExceeded)
	c.Check(locked, Equals, false)
}

func (t *TestSuite) TestFlock_Unlock(c *C) {
	var err error

	err = t.flock.Unlock()
	c.Assert(err, IsNil)

	// get a lock for us to unlock
	locked, err := t.flock.TryLock()
	c.Assert(err, IsNil)
	c.Assert(locked, Equals, true)
	c.Assert(t.flock.Locked(), Equals, true)
	c.Check(t.flock.RLocked(), Equals, false)

	_, err = os.Stat(t.path)
	c.Assert(os.IsNotExist(err), Equals, false)

	err = t.flock.Unlock()
	c.Assert(err, IsNil)
	c.Check(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, false)
}

func (t *TestSuite) TestFlock_Lock(c *C) {
	c.Assert(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, false)

	var err error

	err = t.flock.Lock()
	c.Assert(err, IsNil)
	c.Check(t.flock.Locked(), Equals, true)
	c.Check(t.flock.RLocked(), Equals, false)

	// test that the short-circuit works
	err = t.flock.Lock()
	c.Assert(err, IsNil)

	//
	// Test that Lock() is a blocking call
	//
	ch := make(chan error, 2)
	gf := flock.NewFlock(t.path)
	defer gf.Unlock()

	go func(ch chan<- error) {
		ch <- nil
		ch <- gf.Lock()
		close(ch)
	}(ch)

	errCh, ok := <-ch
	c.Assert(ok, Equals, true)
	c.Assert(errCh, IsNil)

	err = t.flock.Unlock()
	c.Assert(err, IsNil)

	errCh, ok = <-ch
	c.Assert(ok, Equals, true)
	c.Assert(errCh, IsNil)
	c.Check(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, false)
	c.Check(gf.Locked(), Equals, true)
	c.Check(gf.RLocked(), Equals, false)
}

func (t *TestSuite) TestFlock_RLock(c *C) {
	c.Assert(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, false)

	var err error

	err = t.flock.RLock()
	c.Assert(err, IsNil)
	c.Check(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, true)

	// test that the short-circuit works
	err = t.flock.RLock()
	c.Assert(err, IsNil)

	//
	// Test that RLock() is a blocking call
	//
	ch := make(chan error, 2)
	gf := flock.NewFlock(t.path)
	defer gf.Unlock()

	go func(ch chan<- error) {
		ch <- nil
		ch <- gf.RLock()
		close(ch)
	}(ch)

	errCh, ok := <-ch
	c.Assert(ok, Equals, true)
	c.Assert(errCh, IsNil)

	err = t.flock.Unlock()
	c.Assert(err, IsNil)

	errCh, ok = <-ch
	c.Assert(ok, Equals, true)
	c.Assert(errCh, IsNil)
	c.Check(t.flock.Locked(), Equals, false)
	c.Check(t.flock.RLocked(), Equals, false)
	c.Check(gf.Locked(), Equals, false)
	c.Check(gf.RLocked(), Equals, true)
}
