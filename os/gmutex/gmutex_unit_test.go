// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmutex_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gmutex"
	"github.com/gogf/gf/test/gtest"
)

func Test_Mutex_RUnlock(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmutex.New()
		for index := 0; index < 1000; index++ {
			go func() {
				mu.RLockFunc(func() {
					time.Sleep(200 * time.Millisecond)
				})
			}()
		}
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(mu.IsRLocked(), true)
		gtest.Assert(mu.IsLocked(), true)
		gtest.Assert(mu.IsWLocked(), false)
		for index := 0; index < 1000; index++ {
			go func() {
				mu.RUnlock()
			}()
		}
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(mu.IsRLocked(), false)

	})

	//RLock before Lock
	gtest.Case(t, func() {
		mu := gmutex.New()
		mu.RLock()
		go func() {
			mu.Lock()
			time.Sleep(300 * time.Millisecond)
			mu.Unlock()
		}()
		time.Sleep(100 * time.Millisecond)
		mu.RUnlock()
		gtest.Assert(mu.IsRLocked(), false)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(mu.IsLocked(), true)
		time.Sleep(400 * time.Millisecond)
		gtest.Assert(mu.IsLocked(), false)
	})
}

func Test_Mutex_IsLocked(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmutex.New()
		go func() {
			mu.LockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(mu.IsLocked(), true)
		gtest.Assert(mu.IsWLocked(), true)
		gtest.Assert(mu.IsRLocked(), false)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(mu.IsLocked(), false)
		gtest.Assert(mu.IsWLocked(), false)

		go func() {
			mu.RLockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(mu.IsRLocked(), true)
		gtest.Assert(mu.IsLocked(), true)
		gtest.Assert(mu.IsWLocked(), false)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(mu.IsRLocked(), false)
	})
}

func Test_Mutex_Unlock(t *testing.T) {
	gtest.Case(t, func() {
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
		gtest.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Mutex_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
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
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Mutex_TryLockFunc(t *testing.T) {
	gtest.Case(t, func() {
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
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Mutex_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
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
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	gtest.Case(t, func() {
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
		gtest.Assert(array.Len(), 0)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Mutex_TryRLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		mu := gmutex.New()
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(500 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		for index := 0; index < 1000; index++ {
			go func() {
				time.Sleep(1000 * time.Millisecond)
				mu.TryRLockFunc(func() {
					array.Append(1)
				})
			}()
		}
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(150 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(600 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(600 * time.Millisecond)
		gtest.Assert(array.Len(), 1001)
	})
}
