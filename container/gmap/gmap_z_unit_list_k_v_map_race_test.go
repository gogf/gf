// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_ListKVMap_GetOrSetFuncLock_Race tests the atomicity of GetOrSetFuncLock.
// This test ensures that the callback function is only executed once even under
// high concurrency, which verifies that the function holds the lock during the
// entire check-and-set operation.
func Test_ListKVMap_GetOrSetFuncLock_Race(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "counter"
		callCount := int32(0)
		goroutines := 100

		var wg sync.WaitGroup
		wg.Add(goroutines)

		// Start multiple goroutines trying to set the same key
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				m.GetOrSetFuncLock(key, func() int {
					// Increment call count atomically
					atomic.AddInt32(&callCount, 1)
					// Simulate some work
					time.Sleep(time.Microsecond)
					return 100
				})
			}()
		}

		wg.Wait()

		// The callback should only be called once because of proper locking
		t.Assert(atomic.LoadInt32(&callCount), 1)
		t.Assert(m.Get(key), 100)
		t.Assert(m.Size(), 1)
	})
}

// Test_ListKVMap_SetIfNotExistFuncLock_Race tests the atomicity of SetIfNotExistFuncLock.
// This test ensures that only one goroutine can successfully set the value and
// execute the callback function, even under high concurrency.
func Test_ListKVMap_SetIfNotExistFuncLock_Race(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "counter"
		callCount := int32(0)
		successCount := int32(0)
		goroutines := 100

		var wg sync.WaitGroup
		wg.Add(goroutines)

		// Start multiple goroutines trying to set the same key
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				success := m.SetIfNotExistFuncLock(key, func() int {
					// Increment call count atomically
					atomic.AddInt32(&callCount, 1)
					// Simulate some work
					time.Sleep(time.Microsecond)
					return 200
				})
				if success {
					atomic.AddInt32(&successCount, 1)
				}
			}()
		}

		wg.Wait()

		// The callback should only be called once
		t.Assert(atomic.LoadInt32(&callCount), 1)
		// Only one goroutine should succeed
		t.Assert(atomic.LoadInt32(&successCount), 1)
		t.Assert(m.Get(key), 200)
		t.Assert(m.Size(), 1)
	})
}

// Test_ListKVMap_GetOrSetFuncLock_MultipleKeys tests GetOrSetFuncLock with different keys.
// This ensures that operations on different keys don't interfere with each other.
func Test_ListKVMap_GetOrSetFuncLock_MultipleKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		keys := []string{"key1", "key2", "key3", "key4", "key5"}
		callCounts := make([]int32, len(keys))
		goroutines := 20

		var wg sync.WaitGroup

		// For each key, start multiple goroutines
		for i, key := range keys {
			keyIndex := i
			for j := 0; j < goroutines; j++ {
				wg.Add(1)
				go func(idx int, k string) {
					defer wg.Done()
					m.GetOrSetFuncLock(k, func() int {
						atomic.AddInt32(&callCounts[idx], 1)
						time.Sleep(time.Microsecond)
						return (idx + 1) * 100
					})
				}(keyIndex, key)
			}
		}

		wg.Wait()

		// Each key's callback should only be called once
		for _, count := range callCounts {
			t.Assert(atomic.LoadInt32(&count), 1)
		}

		// Verify all keys are set correctly
		for i, key := range keys {
			t.Assert(m.Get(key), (i+1)*100)
		}
		t.Assert(m.Size(), len(keys))
	})
}

// Test_ListKVMap_SetIfNotExistFuncLock_MultipleKeys tests SetIfNotExistFuncLock with different keys.
func Test_ListKVMap_SetIfNotExistFuncLock_MultipleKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[int, string](true)
		keys := []int{1, 2, 3, 4, 5}
		callCounts := make([]int32, len(keys))
		successCounts := make([]int32, len(keys))
		goroutines := 20

		var wg sync.WaitGroup

		// For each key, start multiple goroutines
		for i, key := range keys {
			keyIndex := i
			for j := 0; j < goroutines; j++ {
				wg.Add(1)
				go func(idx int, k int) {
					defer wg.Done()
					success := m.SetIfNotExistFuncLock(k, func() string {
						atomic.AddInt32(&callCounts[idx], 1)
						time.Sleep(time.Microsecond)
						return gtest.DataContent()
					})
					if success {
						atomic.AddInt32(&successCounts[idx], 1)
					}
				}(keyIndex, key)
			}
		}

		wg.Wait()

		// Each key's callback should only be called once
		for _, count := range callCounts {
			t.Assert(atomic.LoadInt32(&count), 1)
		}

		// Each key should have exactly one successful set
		for _, count := range successCounts {
			t.Assert(atomic.LoadInt32(&count), 1)
		}

		t.Assert(m.Size(), len(keys))
	})
}

// Test_ListKVMap_GetOrSetFuncLock_NilValue tests that nil values are handled correctly.
func Test_ListKVMap_GetOrSetFuncLock_NilValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, *int](true)
		key := "nilKey"
		callCount := int32(0)

		var wg sync.WaitGroup
		goroutines := 50
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				m.GetOrSetFuncLock(key, func() *int {
					atomic.AddInt32(&callCount, 1)
					return nil
				})
			}()
		}

		wg.Wait()

		// Callback should be called once
		t.Assert(atomic.LoadInt32(&callCount), 1)
		// Typed nil pointer (*int)(nil) is stored because any(value) != nil for typed nil
		// This is a Go language feature: typed nil is not the same as interface nil
		t.Assert(m.Contains(key), true)
		t.Assert(m.Get(key), (*int)(nil))
		t.Assert(m.Size(), 1)
	})
}

// Test_ListKVMap_SetIfNotExistFuncLock_NilValue tests that nil values are handled correctly.
func Test_ListKVMap_SetIfNotExistFuncLock_NilValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, *string](true)
		key := "nilKey"
		callCount := int32(0)
		successCount := int32(0)

		var wg sync.WaitGroup
		goroutines := 50
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				success := m.SetIfNotExistFuncLock(key, func() *string {
					atomic.AddInt32(&callCount, 1)
					return nil
				})
				if success {
					atomic.AddInt32(&successCount, 1)
				}
			}()
		}

		wg.Wait()

		// Callback should be called once
		t.Assert(atomic.LoadInt32(&callCount), 1)
		// Should report success once
		t.Assert(atomic.LoadInt32(&successCount), 1)
		// Typed nil pointer (*string)(nil) is stored because any(value) != nil for typed nil
		t.Assert(m.Contains(key), true)
		t.Assert(m.Get(key), (*string)(nil))
		t.Assert(m.Size(), 1)
	})
}

// Test_ListKVMap_GetOrSetFuncLock_ExistingKey tests behavior when key already exists.
func Test_ListKVMap_GetOrSetFuncLock_ExistingKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "existing"
		m.Set(key, 999)

		callCount := int32(0)
		goroutines := 50

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				val := m.GetOrSetFuncLock(key, func() int {
					atomic.AddInt32(&callCount, 1)
					return 123
				})
				// Should always get the existing value
				t.Assert(val, 999)
			}()
		}

		wg.Wait()

		// Callback should never be called since key exists
		t.Assert(atomic.LoadInt32(&callCount), 0)
		t.Assert(m.Get(key), 999)
	})
}

// Test_ListKVMap_SetIfNotExistFuncLock_ExistingKey tests behavior when key already exists.
func Test_ListKVMap_SetIfNotExistFuncLock_ExistingKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "existing"
		m.Set(key, 888)

		callCount := int32(0)
		successCount := int32(0)
		goroutines := 50

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				success := m.SetIfNotExistFuncLock(key, func() int {
					atomic.AddInt32(&callCount, 1)
					return 456
				})
				if success {
					atomic.AddInt32(&successCount, 1)
				}
			}()
		}

		wg.Wait()

		// Callback should never be called since key exists
		t.Assert(atomic.LoadInt32(&callCount), 0)
		// No goroutine should succeed
		t.Assert(atomic.LoadInt32(&successCount), 0)
		// Original value should remain
		t.Assert(m.Get(key), 888)
	})
}
