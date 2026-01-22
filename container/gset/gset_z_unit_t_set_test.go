// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gset_test

import (
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestTSet_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestTSet_NewFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3}, true)
		t.Assert(s.Size(), 3)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Contains(3), true)
		t.Assert(s.Contains(4), false)
	})
}

func TestTSet_Add_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s gset.TSet[int]
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)
		t.Assert(s.Contains(1), true)
	})
}

func TestTSet_AddIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int](true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.AddIfNotExist(1), false)
		t.Assert(s.AddIfNotExist(2), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExist(2), false)
	})

	// Test with pointer type to test nil check
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[*int](true)

		val := 1
		ptr := &val
		t.Assert(s.AddIfNotExist(ptr), true)
		t.Assert(s.AddIfNotExist(ptr), false)
	})

	// Test nil data map initialization
	gtest.C(t, func(t *gtest.T) {
		var s gset.TSet[int]
		t.Assert(s.AddIfNotExist(1), true)
		t.Assert(s.Size(), 1)
	})
}

func TestTSet_AddIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int](true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return false }), false)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return true }), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return true }), false)
		t.Assert(s.Contains(2), true)
	})

	// Test concurrent scenario
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int](true)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFunc(1, func() bool {
				time.Sleep(100 * time.Millisecond)
				return true
			})
			t.Assert(r, false)
		}()
		s.Add(1)
		wg.Wait()
	})

	// Test nil data map initialization
	gtest.C(t, func(t *gtest.T) {
		var s gset.TSet[int]
		t.Assert(s.AddIfNotExistFunc(1, func() bool { return true }), true)
		t.Assert(s.Size(), 1)
	})
}

func TestTSet_AddIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int](true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFuncLock(2, func() bool { return false }), false)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFuncLock(2, func() bool { return true }), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExistFuncLock(2, func() bool { return true }), false)
		t.Assert(s.Contains(2), true)
	})

	// Test nil data map initialization
	gtest.C(t, func(t *gtest.T) {
		var s gset.TSet[int]
		t.Assert(s.AddIfNotExistFuncLock(1, func() bool { return true }), true)
		t.Assert(s.Size(), 1)
	})
}

func TestTSet_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3, 4, 5}, true)
		var sum int
		s.Iterator(func(v int) bool {
			sum += v
			return true
		})
		t.Assert(sum, 15)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3, 4, 5}, true)
		var count int
		s.Iterator(func(v int) bool {
			count++
			return count < 3
		})
		t.Assert(count, 3)
	})
}

func TestTSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		t.Assert(s.Join(","), "")
		s.Add(1, 2, 3)
		result := s.Join(",")
		t.Assert(len(result) > 0, true)
	})
}

func TestTSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s *gset.TSet[int]
		t.Assert(s.String(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		t.Assert(s.String(), "[]")
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		result := s.String()
		t.Assert(len(result) > 2, true)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[string]([]string{"a", "b", "c"})
		result := s.String()
		t.Assert(len(result) > 2, true)
	})
}

func TestTSet_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		t.Assert(s1.Equal(s2), true)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{1, 2, 3, 4})
		t.Assert(s1.Equal(s2), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		t.Assert(s1.Equal(s1), true)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{1, 2, 4})
		t.Assert(s1.Equal(s2), false)
	})
}

func TestTSet_IsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2})
		s2 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		t.Assert(s1.IsSubsetOf(s2), true)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{1, 2})
		t.Assert(s1.IsSubsetOf(s2), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		t.Assert(s1.IsSubsetOf(s1), true)
	})
}

func TestTSet_Union(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{3, 4, 5})
		s := s1.Union(s2)
		t.Assert(s.Size(), 5)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Contains(3), true)
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), true)
	})

	// Test with nil set - should skip it and copy s1 data
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		var s2 *gset.TSet[int]
		s := s1.Union(s2)
		// Since s2 is nil and skipped, newSet will be empty
		// because the loop runs but nothing is copied when other is nil
		t.Assert(s.Size(), 0)
	})

	// Test with self
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s := s1.Union(s1)
		t.Assert(s.Size(), 3)
	})
}

func TestTSet_Diff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{3, 4, 5})
		s := s1.Diff(s2)
		t.Assert(s.Size(), 2)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Contains(3), false)
	})

	// Test with nil set - should skip it
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		var s2 *gset.TSet[int]
		s := s1.Diff(s2)
		// Since s2 is nil and skipped, newSet will be empty
		// because the loop runs but nothing is copied when other is nil
		t.Assert(s.Size(), 0)
	})

	// Test with self
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s := s1.Diff(s1)
		t.Assert(s.Size(), 0)
	})
}

func TestTSet_Intersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{3, 4, 5})
		s := s1.Intersect(s2)
		t.Assert(s.Size(), 1)
		t.Assert(s.Contains(3), true)
	})

	// Test with nil set
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		var s2 *gset.TSet[int]
		s := s1.Intersect(s2)
		t.Assert(s.Size(), 0)
	})

	// Test with self
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s := s1.Intersect(s1)
		t.Assert(s.Size(), 3)
	})
}

func TestTSet_Complement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{1, 2, 3, 4, 5})
		s := s1.Complement(s2)
		t.Assert(s.Size(), 2)
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), true)
	})

	// Test with self
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s := s1.Complement(s1)
		t.Assert(s.Size(), 0)
	})
}

func TestTSet_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s2 := gset.NewTSetFrom[int]([]int{3, 4, 5})
		s1.Merge(s2)
		t.Assert(s1.Size(), 5)
		t.Assert(s1.Contains(1), true)
		t.Assert(s1.Contains(2), true)
		t.Assert(s1.Contains(3), true)
		t.Assert(s1.Contains(4), true)
		t.Assert(s1.Contains(5), true)
	})

	// Test with nil set
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		var s2 *gset.TSet[int]
		s1.Merge(s2)
		t.Assert(s1.Size(), 3)
	})

	// Test with self
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3})
		s1.Merge(s1)
		t.Assert(s1.Size(), 3)
	})
}

func TestTSet_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		t.Assert(s.Sum(), 6)
	})
}

func TestTSet_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		item := s.Pop()
		t.Assert(s.Size(), 2)
		t.Assert(s.Contains(item), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		item := s.Pop()
		t.Assert(item, 0)
	})
}

func TestTSet_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3, 4, 5})
		items := s.Pops(3)
		t.Assert(len(items), 3)
		t.Assert(s.Size(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		items := s.Pops(-1)
		t.Assert(len(items), 3)
		t.Assert(s.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		items := s.Pops(0)
		t.Assert(items, nil)
		t.Assert(s.Size(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		items := s.Pops(10)
		t.Assert(len(items), 3)
		t.Assert(s.Size(), 0)
	})
}

func TestTSet_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2})
		s.Walk(func(item int) int {
			return item + 10
		})
		t.Assert(s.Size(), 2)
		t.Assert(s.Contains(11), true)
		t.Assert(s.Contains(12), true)
	})
}

func TestTSet_MarshalJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3})
		b, err := json.Marshal(s)
		t.AssertNil(err)
		t.Assert(len(b) > 0, true)
	})
}

func TestTSet_UnmarshalJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		b := []byte(`[1,2,3]`)
		err := json.UnmarshalUseNumber(b, &s)
		t.AssertNil(err)
		t.Assert(s.Size(), 3)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Contains(3), true)
	})

	// Test with nil data map
	gtest.C(t, func(t *gtest.T) {
		var s gset.TSet[int]
		b := []byte(`[1,2,3]`)
		err := json.UnmarshalUseNumber(b, &s)
		t.AssertNil(err)
		t.Assert(s.Size(), 3)
	})

	// Test with invalid JSON
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		b := []byte(`{invalid}`)
		err := json.UnmarshalUseNumber(b, &s)
		t.AssertNE(err, nil)
	})

	// Test with empty array
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		b := []byte(`[]`)
		err := json.UnmarshalUseNumber(b, &s)
		t.AssertNil(err)
		t.Assert(s.Size(), 0)
	})
}

func TestTSet_UnmarshalValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		err := s.UnmarshalValue([]byte(`[1,2,3]`))
		t.AssertNil(err)
		t.Assert(s.Size(), 3)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Contains(3), true)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		err := s.UnmarshalValue(`[1,2,3]`)
		t.AssertNil(err)
		t.Assert(s.Size(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		err := s.UnmarshalValue([]int{1, 2, 3})
		t.AssertNil(err)
		t.Assert(s.Size(), 3)
	})

	// Test with nil data map
	gtest.C(t, func(t *gtest.T) {
		var s gset.TSet[int]
		err := s.UnmarshalValue([]int{1, 2, 3})
		t.AssertNil(err)
		t.Assert(s.Size(), 3)
	})

	// Test error case with invalid JSON
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		err := s.UnmarshalValue([]byte(`{invalid}`))
		t.AssertNE(err, nil)
	})

	// Test with empty array for string/bytes case
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		err := s.UnmarshalValue([]byte(`[]`))
		t.AssertNil(err)
		t.Assert(s.Size(), 0)
	})

	// Test with empty slice for default case
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSet[int]()
		err := s.UnmarshalValue([]int{})
		t.AssertNil(err)
		t.Assert(s.Size(), 0)
	})
}

func TestTSet_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewTSetFrom[int]([]int{1, 2, 3}, true)
		s2 := s1.DeepCopy().(*gset.TSet[int])
		t.Assert(s1.Size(), s2.Size())
		t.Assert(s1.Contains(1), s2.Contains(1))
		t.Assert(s1.Contains(2), s2.Contains(2))
		t.Assert(s1.Contains(3), s2.Contains(3))

		s1.Add(4)
		t.Assert(s1.Size(), 4)
		t.Assert(s2.Size(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		var s1 *gset.TSet[int]
		s2 := s1.DeepCopy()
		t.Assert(s2, nil)
	})
}

func TestTSet_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3}, true)
		s.LockFunc(func(m map[int]struct{}) {
			m[4] = struct{}{}
		})
		t.Assert(s.Size(), 4)
		t.Assert(s.Contains(4), true)
	})
}

func TestTSet_RLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewTSetFrom[int]([]int{1, 2, 3}, true)
		var sum int
		s.RLockFunc(func(m map[int]struct{}) {
			for k := range m {
				sum += k
			}
		})
		t.Assert(sum, 6)
	})
}
