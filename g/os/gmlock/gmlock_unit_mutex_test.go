// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gmlock"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Mutex_Unlock(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmlock.NewMutex()
		array := garray.New()
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()

		go func() {
			time.Sleep(60 * time.Millisecond)
			mu.Unlock()
			mu.Unlock()
			mu.Unlock()
		}()

		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Mutex_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmlock.NewMutex()
		array := garray.New()
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Mutex_TryLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmlock.NewMutex()
		array := garray.New()
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(110 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Mutex_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmlock.NewMutex()
		array := garray.New()
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	gtest.Case(t, func() {
		mu := gmlock.NewMutex()
		array := garray.New()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		gtest.Assert(array.Len(), 0)
		time.Sleep(80 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Mutex_TryRLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmlock.NewMutex()
		array := garray.New()
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(110 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(110 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}
