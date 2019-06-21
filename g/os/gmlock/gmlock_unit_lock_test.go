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

func Test_Locker_Lock(t *testing.T) {
	//no expire
	gtest.Case(t, func() {
		key := "testLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(10 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
	})
	//expire
	gtest.Case(t, func() {
		key := "testLockExpire"
		array := garray.New()
		go func() {
			gmlock.Lock(key, 100*time.Millisecond)
			array.Append(1)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.Lock(key)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(150 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(250 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Locker_TryLock(t *testing.T) {
	gtest.Case(t, func() {
		key := "testTryLock"
		array := garray.New()
		go func() {
			if gmlock.TryLock(key, 200*time.Millisecond) {
				array.Append(1)
			}
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if !gmlock.TryLock(key) {
				array.Append(1)
			} else {
				gmlock.Unlock(key)
			}
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			if gmlock.TryLock(key) {
				array.Append(1)
				gmlock.Unlock(key)
			}
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		time.Sleep(350 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Locker_LockFunc(t *testing.T) {
	//no expire
	gtest.Case(t, func() {
		key := "testLockFunc"
		array := garray.New()
		go func() {
			gmlock.LockFunc(key, func() {
				array.Append(1)
				time.Sleep(50 * time.Millisecond)
			}) //
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.LockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(10 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1) //
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	//expire
	gtest.Case(t, func() {
		key := "testLockFuncExpire"
		array := garray.New()
		go func() {
			gmlock.LockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			}, 100*time.Millisecond) //
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			gmlock.LockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 2) //
		time.Sleep(350 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
func Test_Locker_TryLockFunc(t *testing.T) {
	//no expire
	gtest.Case(t, func() {
		key := "testTryLockFunc"
		array := garray.New()
		go func() {
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
				time.Sleep(50 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(70 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
	//expire1
	gtest.Case(t, func() {
		key := "testTryLockFuncExpire1"
		array := garray.New()
		go func() {
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			}, 50*time.Millisecond)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(70 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})

	//expire2
	gtest.Case(t, func() {
		key := "testTryLockFuncExpire2"
		array := garray.New()
		go func() {
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			}, 50*time.Millisecond) //unlock after expire, before func finish.
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(70 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(10 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(70 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
