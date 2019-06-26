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

func Test_Locker_RLock(t *testing.T) {
	//RLock before Lock
	gtest.Case(t, func() {
		key := "testRLockBeforeLock"
		array := garray.New()
		go func() {
			gmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	//Lock before RLock
	gtest.Case(t, func() {
		key := "testLockBeforeRLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLock(key)
			array.Append(1)
			gmlock.RUnlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	//Lock before RLocks
	gtest.Case(t, func() {
		key := "testLockBeforeRLocks"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.RUnlock(key)
		}()
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLock(t *testing.T) {
	//Lock before TryRLock
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if gmlock.TryRLock(key) {
				array.Append(1)
				gmlock.RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
	})

	//Lock before TryRLocks
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLocks"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if gmlock.TryRLock(key) {
				array.Append(1)
				gmlock.RUnlock(key)
			}
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			if gmlock.TryRLock(key) {
				array.Append(1)
				gmlock.RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Locker_RLockFunc(t *testing.T) {
	//RLockFunc before Lock
	gtest.Case(t, func() {
		key := "testRLockFuncBeforeLock"
		array := garray.New()
		go func() {
			gmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(150 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	//Lock before RLockFunc
	gtest.Case(t, func() {
		key := "testLockBeforeRLockFunc"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	//Lock before RLockFuncs
	gtest.Case(t, func() {
		key := "testLockBeforeRLockFuncs"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLockFunc(t *testing.T) {
	//Lock before TryRLockFunc
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLockFunc"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
	})

	//Lock before TryRLockFuncs
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLockFuncs"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			gmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
