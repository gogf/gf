// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmutex_test

import (
	"context"
	"github.com/gogf/gf/os/glog"
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gmutex"
	"github.com/gogf/gf/test/gtest"
)

func Test_Mutex_RUnlock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		for index := 0; index < 1000; index++ {
			go func() {
				mu.RLockFunc(func() {
					time.Sleep(200 * time.Millisecond)
				})
			}()
		}
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsRLocked(), true)
		t.Assert(mu.IsLocked(), true)
		t.Assert(mu.IsWLocked(), false)
		for index := 0; index < 1000; index++ {
			go func() {
				mu.RUnlock()
			}()
		}
		time.Sleep(300 * time.Millisecond)
		t.Assert(mu.IsRLocked(), false)

	})

	//RLock before Lock
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		mu.RLock()
		go func() {
			mu.Lock()
			time.Sleep(300 * time.Millisecond)
			mu.Unlock()
		}()
		time.Sleep(100 * time.Millisecond)
		mu.RUnlock()
		t.Assert(mu.IsRLocked(), false)
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsLocked(), true)
		time.Sleep(400 * time.Millisecond)
		t.Assert(mu.IsLocked(), false)
	})
}

func Test_Mutex_IsLocked(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		go func() {
			mu.LockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsLocked(), true)
		t.Assert(mu.IsWLocked(), true)
		t.Assert(mu.IsRLocked(), false)
		time.Sleep(300 * time.Millisecond)
		t.Assert(mu.IsLocked(), false)
		t.Assert(mu.IsWLocked(), false)

		go func() {
			mu.RLockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsRLocked(), true)
		t.Assert(mu.IsLocked(), true)
		t.Assert(mu.IsWLocked(), false)
		time.Sleep(300 * time.Millisecond)
		t.Assert(mu.IsRLocked(), false)
	})
}

func Test_Mutex_Unlock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()

		go func() {
			time.Sleep(200 * time.Millisecond)
			mu.Unlock()
			mu.Unlock()
			mu.Unlock()
			mu.Unlock()
		}()

		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Mutex_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Mutex_TryLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Mutex_RLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.New()
		array := garray.New(true)
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		t.Assert(array.Len(), 0)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Mutex_TryRLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu    = gmutex.New()
			array = garray.New(true)
		)
		// First writing lock
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				glog.Println(context.TODO(), "lock1 done")
				time.Sleep(2000 * time.Millisecond)
			})
		}()
		// This goroutine never gets the lock.
		go func() {
			time.Sleep(1000 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		for index := 0; index < 1000; index++ {
			go func() {
				time.Sleep(4000 * time.Millisecond)
				mu.TryRLockFunc(func() {
					array.Append(1)
				})
			}()
		}
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 1001)
	})
}
