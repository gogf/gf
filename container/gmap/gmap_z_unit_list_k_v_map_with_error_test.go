// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/test/gtest"
)

// --------------------------------------------------------------------------
// GetOrSetFuncWithError
// --------------------------------------------------------------------------

func Test_ListKVMap_GetOrSetFuncWithError(t *testing.T) {
	type MyVal struct {
		Valid bool
		Data  string
	}
	// Case: key not exist, f returns valid value → value stored and returned.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		v, err := m.GetOrSetFuncWithError("k1", func() (string, error) {
			return "val1", nil
		})
		t.AssertNil(err)
		t.Assert(v, "val1")
		t.Assert(m.Get("k1"), "val1")
		t.Assert(m.Size(), 1)
	})

	// Case: key already exists → existing value returned, f is NOT called.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("k1", "existing")
		called := false
		v, err := m.GetOrSetFuncWithError("k1", func() (string, error) {
			called = true
			return "new", nil
		})
		t.AssertNil(err)
		t.Assert(v, "existing")
		t.Assert(called, false)
		t.Assert(m.Get("k1"), "existing")
	})

	// Case: f returns an error → zero value returned, error propagated, key NOT stored.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		testErr := errors.New("load error")
		v, err := m.GetOrSetFuncWithError("k1", func() (string, error) {
			return "", testErr
		})
		t.AssertNE(err, nil)
		t.Assert(err, testErr)
		t.Assert(v, "")
		t.Assert(m.Contains("k1"), false)
		t.Assert(m.Size(), 0)
	})

	// Case: f returns nil pointer → nil not stored in map, nil returned.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, *string]()
		v, err := m.GetOrSetFuncWithError("k1", func() (*string, error) {
			return nil, nil
		})
		t.AssertNil(err)
		t.AssertNil(v)
		t.Assert(m.Contains("k1"), false)
		t.Assert(m.Size(), 0)
	})

	// Case: custom NilChecker — struct treated as "nil" by checker → not stored.
	gtest.C(t, func(t *gtest.T) {
		checker := func(v *MyVal) bool { return !v.Valid }
		m := gmap.NewListKVMapWithChecker[string, *MyVal](checker)

		// "nil" per checker: valid=false
		v, err := m.GetOrSetFuncWithError("k1", func() (*MyVal, error) {
			return &MyVal{Valid: false, Data: "ignored"}, nil
		})
		t.AssertNil(err)
		t.Assert(v.Valid, false)
		t.Assert(m.Contains("k1"), false)

		// valid value → stored
		v, err = m.GetOrSetFuncWithError("k2", func() (*MyVal, error) {
			return &MyVal{Valid: true, Data: "hello"}, nil
		})
		t.AssertNil(err)
		t.Assert(v.Valid, true)
		t.Assert(v.Data, "hello")
		t.Assert(m.Contains("k2"), true)
		t.Assert(m.Get("k2").Data, "hello")
	})

	// Case: insertion order is maintained after GetOrSetFuncWithError adds entries.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)

		_, err := m.GetOrSetFuncWithError("c", func() (int, error) {
			return 3, nil
		})
		t.AssertNil(err)

		keys := m.Keys()
		t.Assert(len(keys), 3)
		t.Assert(keys[0], "a")
		t.Assert(keys[1], "b")
		t.Assert(keys[2], "c")
	})

	// Case: after f returns error, key is absent and a subsequent call can succeed.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		attempts := gtype.NewInt32()

		_, err := m.GetOrSetFuncWithError("k1", func() (int, error) {
			attempts.Add(1)
			return 0, errors.New("temporary error")
		})
		t.AssertNE(err, nil)
		t.Assert(m.Contains("k1"), false)

		v, err := m.GetOrSetFuncWithError("k1", func() (int, error) {
			attempts.Add(1)
			return 42, nil
		})
		t.AssertNil(err)
		t.Assert(v, 42)
		t.Assert(m.Get("k1"), 42)
		t.Assert(attempts.Val(), 2)
	})

	// Case: safe mode (concurrent-safe=true) — basic functionality is correct.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		v, err := m.GetOrSetFuncWithError("k1", func() (int, error) {
			return 99, nil
		})
		t.AssertNil(err)
		t.Assert(v, 99)
		t.Assert(m.Get("k1"), 99)
	})
}

// --------------------------------------------------------------------------
// GetOrSetFuncLockWithError
// --------------------------------------------------------------------------

func Test_ListKVMap_GetOrSetFuncLockWithError(t *testing.T) {
	type MyVal struct {
		Valid bool
		Data  string
	}
	// Case: key not exist, f returns valid value → value stored and returned.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		v, err := m.GetOrSetFuncLockWithError("k1", func() (string, error) {
			return "val1", nil
		})
		t.AssertNil(err)
		t.Assert(v, "val1")
		t.Assert(m.Get("k1"), "val1")
		t.Assert(m.Size(), 1)
	})

	// Case: key already exists → existing value returned, f is NOT called.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("k1", "existing")
		called := false
		v, err := m.GetOrSetFuncLockWithError("k1", func() (string, error) {
			called = true
			return "new", nil
		})
		t.AssertNil(err)
		t.Assert(v, "existing")
		t.Assert(called, false)
		t.Assert(m.Get("k1"), "existing")
	})

	// Case: f returns an error → zero value returned, error propagated, key NOT stored.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		testErr := errors.New("lock load error")
		v, err := m.GetOrSetFuncLockWithError("k1", func() (string, error) {
			return "", testErr
		})
		t.AssertNE(err, nil)
		t.Assert(err, testErr)
		t.Assert(v, "")
		t.Assert(m.Contains("k1"), false)
		t.Assert(m.Size(), 0)
	})

	// Case: f returns nil pointer → nil not stored in map, nil returned.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, *string]()
		v, err := m.GetOrSetFuncLockWithError("k1", func() (*string, error) {
			return nil, nil
		})
		t.AssertNil(err)
		t.AssertNil(v)
		t.Assert(m.Contains("k1"), false)
		t.Assert(m.Size(), 0)
	})

	// Case: custom NilChecker — struct treated as "nil" by checker → not stored.
	gtest.C(t, func(t *gtest.T) {
		checker := func(v *MyVal) bool { return !v.Valid }
		m := gmap.NewListKVMapWithChecker[string, *MyVal](checker)

		// "nil" per checker → not stored
		v, err := m.GetOrSetFuncLockWithError("k1", func() (*MyVal, error) {
			return &MyVal{Valid: false, Data: "ignored"}, nil
		})
		t.AssertNil(err)
		t.Assert(v.Valid, false)
		t.Assert(m.Contains("k1"), false)

		// valid value → stored
		v, err = m.GetOrSetFuncLockWithError("k2", func() (*MyVal, error) {
			return &MyVal{Valid: true, Data: "world"}, nil
		})
		t.AssertNil(err)
		t.Assert(v.Data, "world")
		t.Assert(m.Contains("k2"), true)
		t.Assert(m.Get("k2").Data, "world")
	})

	// Case: insertion order is maintained after GetOrSetFuncLockWithError adds entries.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("x", 10)
		m.Set("y", 20)

		_, err := m.GetOrSetFuncLockWithError("z", func() (int, error) {
			return 30, nil
		})
		t.AssertNil(err)

		keys := m.Keys()
		t.Assert(len(keys), 3)
		t.Assert(keys[0], "x")
		t.Assert(keys[1], "y")
		t.Assert(keys[2], "z")
	})

	// Case: after f returns error, key is absent and a subsequent call can succeed.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()

		_, err := m.GetOrSetFuncLockWithError("k1", func() (int, error) {
			return 0, errors.New("temporary")
		})
		t.AssertNE(err, nil)
		t.Assert(m.Contains("k1"), false)

		v, err := m.GetOrSetFuncLockWithError("k1", func() (int, error) {
			return 77, nil
		})
		t.AssertNil(err)
		t.Assert(v, 77)
		t.Assert(m.Get("k1"), 77)
	})
}

// Test_ListKVMap_GetOrSetFuncLockWithError_Race verifies that f is called exactly once
// under high concurrency because GetOrSetFuncLockWithError holds the mutex while calling f.
// This differs from GetOrSetFuncWithError, which calls f outside the lock and may invoke
// f multiple times when multiple goroutines all see the key as absent simultaneously.
func Test_ListKVMap_GetOrSetFuncLockWithError_Race(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "shared"
		callCount := gtype.NewInt32()
		goroutines := 100

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for range goroutines {
			go func() {
				defer wg.Done()
				v, err := m.GetOrSetFuncLockWithError(key, func() (int, error) {
					callCount.Add(1)
					time.Sleep(time.Microsecond)
					return 999, nil
				})
				t.AssertNil(err)
				t.Assert(v, 999)
			}()
		}
		wg.Wait()

		// f must be called exactly once because it executes inside the write lock.
		t.Assert(callCount.Val(), 1)
		t.Assert(m.Get(key), 999)
		t.Assert(m.Size(), 1)
	})
}

// Test_ListKVMap_GetOrSetFuncLockWithError_Race_ErrorCase verifies that when f returns
// an error concurrently, the key is never stored and a later successful call works.
func Test_ListKVMap_GetOrSetFuncLockWithError_Race_ErrorCase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "retry"
		attempts := gtype.NewInt32()

		goroutines := 20
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for range goroutines {
			go func() {
				defer wg.Done()
				_, err := m.GetOrSetFuncLockWithError(key, func() (int, error) {
					attempts.Add(1)
					return 0, errors.New("transient")
				})
				t.AssertNE(err, nil)
			}()
		}
		wg.Wait()

		t.Assert(m.Contains(key), false)

		v, err := m.GetOrSetFuncLockWithError(key, func() (int, error) {
			return 55, nil
		})
		t.AssertNil(err)
		t.Assert(v, 55)
		t.Assert(m.Get(key), 55)
	})
}

// --------------------------------------------------------------------------
// SetIfNotExistFuncWithError
// --------------------------------------------------------------------------

func Test_ListKVMap_SetIfNotExistFuncWithError(t *testing.T) {
	type MyVal struct {
		Valid bool
		Data  string
	}
	// Case: key not exist, f returns valid value → true, nil, value stored.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		ok, err := m.SetIfNotExistFuncWithError("k1", func() (int, error) {
			return 100, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Get("k1"), 100)
		t.Assert(m.Size(), 1)
	})

	// Case: key already exists → f NOT called, returns (false, nil), original value unchanged.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("k1", 42)
		called := false
		ok, err := m.SetIfNotExistFuncWithError("k1", func() (int, error) {
			called = true
			return 999, nil
		})
		t.AssertNil(err)
		t.Assert(ok, false)
		t.Assert(called, false)
		t.Assert(m.Get("k1"), 42)
	})

	// Case: f returns an error → (false, error), key NOT stored.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		testErr := errors.New("set error")
		ok, err := m.SetIfNotExistFuncWithError("k1", func() (int, error) {
			return 0, testErr
		})
		t.Assert(err, testErr)
		t.Assert(ok, false)
		t.Assert(m.Contains("k1"), false)
		t.Assert(m.Size(), 0)
	})

	// Case: f returns nil pointer → returns (true, nil) but value NOT stored in map.
	// This is the special nil-value behavior: the operation reports "intent to set"
	// but skips storage when the value is nil.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, *int]()
		ok, err := m.SetIfNotExistFuncWithError("k1", func() (*int, error) {
			return nil, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)                // returns true
		t.Assert(m.Contains("k1"), false) // but NOT stored
		t.Assert(m.Size(), 0)

		// A subsequent call can still set the value (key was never stored).
		n := 7
		ok, err = m.SetIfNotExistFuncWithError("k1", func() (*int, error) {
			return &n, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Contains("k1"), true)
		t.Assert(*m.Get("k1"), 7)
	})

	// Case: custom NilChecker — struct treated as "nil" by checker → not stored, returns true.
	gtest.C(t, func(t *gtest.T) {
		checker := func(v *MyVal) bool { return !v.Valid }
		m := gmap.NewListKVMapWithChecker[string, *MyVal](checker)

		// "nil" per checker → not stored, but returns true
		ok, err := m.SetIfNotExistFuncWithError("k1", func() (*MyVal, error) {
			return &MyVal{Valid: false, Data: "irrelevant"}, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Contains("k1"), false)

		// valid value → stored
		ok, err = m.SetIfNotExistFuncWithError("k2", func() (*MyVal, error) {
			return &MyVal{Valid: true, Data: "stored"}, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Contains("k2"), true)
		t.Assert(m.Get("k2").Data, "stored")
	})

	// Case: insertion order is maintained after SetIfNotExistFuncWithError adds entries.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)

		ok, err := m.SetIfNotExistFuncWithError("c", func() (int, error) {
			return 3, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)

		keys := m.Keys()
		t.Assert(len(keys), 3)
		t.Assert(keys[0], "a")
		t.Assert(keys[1], "b")
		t.Assert(keys[2], "c")
	})

	// Case: safe mode (concurrent-safe=true) — basic functionality is correct.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		ok, err := m.SetIfNotExistFuncWithError("k1", func() (int, error) {
			return 55, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Get("k1"), 55)
	})
}

// --------------------------------------------------------------------------
// SetIfNotExistFuncLockWithError
// --------------------------------------------------------------------------

func Test_ListKVMap_SetIfNotExistFuncLockWithError(t *testing.T) {
	type MyVal struct {
		Valid bool
		Data  string
	}
	// Case: key not exist, f returns valid value → true, nil, value stored.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		ok, err := m.SetIfNotExistFuncLockWithError("k1", func() (int, error) {
			return 200, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Get("k1"), 200)
		t.Assert(m.Size(), 1)
	})

	// Case: key already exists → f NOT called, returns (false, nil), original value unchanged.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("k1", 99)
		called := false
		ok, err := m.SetIfNotExistFuncLockWithError("k1", func() (int, error) {
			called = true
			return 999, nil
		})
		t.AssertNil(err)
		t.Assert(ok, false)
		t.Assert(called, false)
		t.Assert(m.Get("k1"), 99)
	})

	// Case: f returns an error → (false, error), key NOT stored.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		testErr := errors.New("lock set error")
		ok, err := m.SetIfNotExistFuncLockWithError("k1", func() (int, error) {
			return 0, testErr
		})
		t.Assert(err, testErr)
		t.Assert(ok, false)
		t.Assert(m.Contains("k1"), false)
		t.Assert(m.Size(), 0)
	})

	// Case: f returns nil pointer → returns (true, nil) but value NOT stored in map.
	// Special behavior: the method signals "key was absent and no error" via true,
	// but skips the actual insertion because the value is nil.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, *int]()
		ok, err := m.SetIfNotExistFuncLockWithError("k1", func() (*int, error) {
			return nil, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)                // true: key was absent and f had no error
		t.Assert(m.Contains("k1"), false) // but nil value was NOT stored
		t.Assert(m.Size(), 0)

		// A subsequent call can still store a real value (key remains absent).
		n := 9
		ok, err = m.SetIfNotExistFuncLockWithError("k1", func() (*int, error) {
			return &n, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Contains("k1"), true)
		t.Assert(*m.Get("k1"), 9)
	})

	// Case: custom NilChecker — struct treated as "nil" by checker → not stored, returns true.
	gtest.C(t, func(t *gtest.T) {
		checker := func(v *MyVal) bool { return !v.Valid }
		m := gmap.NewListKVMapWithChecker[string, *MyVal](checker)

		// "nil" per checker → not stored, but returns true
		ok, err := m.SetIfNotExistFuncLockWithError("k1", func() (*MyVal, error) {
			return &MyVal{Valid: false, Data: "irrelevant"}, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Contains("k1"), false)

		// valid value → stored
		ok, err = m.SetIfNotExistFuncLockWithError("k2", func() (*MyVal, error) {
			return &MyVal{Valid: true, Data: "hello"}, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Contains("k2"), true)
		t.Assert(m.Get("k2").Data, "hello")
	})

	// Case: insertion order is maintained after SetIfNotExistFuncLockWithError adds entries.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)

		ok, err := m.SetIfNotExistFuncLockWithError("c", func() (int, error) {
			return 3, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)

		keys := m.Keys()
		t.Assert(len(keys), 3)
		t.Assert(keys[0], "a")
		t.Assert(keys[1], "b")
		t.Assert(keys[2], "c")
	})
}

// Test_ListKVMap_SetIfNotExistFuncLockWithError_Race verifies that f is called exactly once
// and only one goroutine succeeds under high concurrency, because
// SetIfNotExistFuncLockWithError holds the mutex for the entire check-and-set operation.
func Test_ListKVMap_SetIfNotExistFuncLockWithError_Race(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		key := "race_key"
		callCount := gtype.NewInt32()
		successCount := gtype.NewInt32()
		goroutines := 100

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for range goroutines {
			go func() {
				defer wg.Done()
				ok, err := m.SetIfNotExistFuncLockWithError(key, func() (int, error) {
					callCount.Add(1)
					time.Sleep(time.Microsecond)
					return 42, nil
				})
				t.AssertNil(err)
				if ok {
					successCount.Add(1)
				}
			}()
		}
		wg.Wait()

		// f must be called exactly once (lock held during f execution).
		t.Assert(callCount.Val(), 1)
		// Exactly one goroutine reports success.
		t.Assert(successCount.Val(), 1)
		t.Assert(m.Get(key), 42)
		t.Assert(m.Size(), 1)
	})
}

// Test_ListKVMap_SetIfNotExistFuncLockWithError_Race_MultipleKeys verifies correctness when
// multiple goroutines compete over different keys simultaneously.
func Test_ListKVMap_SetIfNotExistFuncLockWithError_Race_MultipleKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		keys := []string{"alpha", "beta", "gamma", "delta"}
		callCounts := make([]*gtype.Int32, len(keys))
		successCounts := make([]*gtype.Int32, len(keys))
		for i := range callCounts {
			callCounts[i] = gtype.NewInt32()
			successCounts[i] = gtype.NewInt32()
		}
		goroutines := 30

		var wg sync.WaitGroup
		for i, key := range keys {
			keyIdx := i
			for range goroutines {
				wg.Add(1)
				go func(idx int, k string) {
					defer wg.Done()
					ok, err := m.SetIfNotExistFuncLockWithError(k, func() (int, error) {
						callCounts[idx].Add(1)
						time.Sleep(time.Microsecond)
						return (idx + 1) * 10, nil
					})
					t.AssertNil(err)
					if ok {
						successCounts[idx].Add(1)
					}
				}(keyIdx, key)
			}
		}
		wg.Wait()

		for i, key := range keys {
			// f called exactly once per key
			t.Assert(callCounts[i].Val(), 1)
			// exactly one goroutine succeeded per key
			t.Assert(successCounts[i].Val(), 1)
			t.Assert(m.Get(key), (i+1)*10)
		}
		t.Assert(m.Size(), len(keys))
	})
}

// Test_ListKVMap_SetIfNotExistFuncLockWithError_ErrorRetry verifies that after f returns an error
// the key remains absent and a subsequent successful call stores the value correctly.
func Test_ListKVMap_SetIfNotExistFuncLockWithError_ErrorRetry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		key := "retry"

		ok, err := m.SetIfNotExistFuncLockWithError(key, func() (int, error) {
			return 0, errors.New("transient error")
		})
		t.AssertNE(err, nil)
		t.Assert(ok, false)
		t.Assert(m.Contains(key), false)

		// After the error the key is still absent; a new call succeeds.
		ok, err = m.SetIfNotExistFuncLockWithError(key, func() (int, error) {
			return 123, nil
		})
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(m.Get(key), 123)
	})
}
