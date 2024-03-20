// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package gset_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestStrSetVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s gset.StrSet
		s.Add("1", "1", "2")
		s.Add([]string{"3", "4"}...)
		t.Assert(s.Size(), 4)
		t.AssertIN("1", s.Slice())
		t.AssertIN("2", s.Slice())
		t.AssertIN("3", s.Slice())
		t.AssertIN("4", s.Slice())
		t.AssertNI("0", s.Slice())
		t.Assert(s.Contains("4"), true)
		t.Assert(s.Contains("5"), false)
		s.Remove("1")
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestStrSetBasic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet()
		s.Add("1", "1", "2")
		s.Add([]string{"3", "4"}...)
		t.Assert(s.Size(), 4)
		t.AssertIN("1", s.Slice())
		t.AssertIN("2", s.Slice())
		t.AssertIN("3", s.Slice())
		t.AssertIN("4", s.Slice())
		t.AssertNI("0", s.Slice())
		t.Assert(s.Contains("4"), true)
		t.Assert(s.Contains("5"), false)
		s.Remove("1")
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestStrSetContainsI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet()
		s.Add("a", "b", "C")
		t.Assert(s.Contains("A"), false)
		t.Assert(s.Contains("a"), true)
		t.Assert(s.ContainsI("A"), true)
		t.Assert(s.ContainsI("d"), false)
	})
}

func TestStrSetIterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet()
		s.Add("1", "2", "3")
		t.Assert(s.Size(), 3)

		a1 := garray.New(true)
		a2 := garray.New(true)
		s.Iterator(func(v string) bool {
			a1.Append("1")
			return false
		})
		s.Iterator(func(v string) bool {
			a2.Append("1")
			return true
		})
		t.Assert(a1.Len(), 1)
		t.Assert(a2.Len(), 3)
	})
}

func TestStrSetLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet()
		s.Add("1", "2", "3")
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[string]struct{}) {
			delete(m, "1")
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[string]struct{}) {
			t.Assert(m, map[string]struct{}{
				"3": {},
				"2": {},
			})
		})
	})
}

func TestStrSetEqual(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s3 := gset.NewStrSet()
		s4 := gset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("1", "2", "3")
		s3.Add("1", "2", "3", "4")
		s4.Add("4", "5", "6")
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
		t.Assert(s1.Equal(s4), false)
		s5 := s1
		t.Assert(s1.Equal(s5), true)
	})
}

func TestStrSetIsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s3 := gset.NewStrSet()
		s1.Add("1", "2")
		s2.Add("1", "2", "3")
		s3.Add("1", "2", "3", "4")
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)

		s4 := s1
		t.Assert(s1.IsSubsetOf(s4), true)
	})
}

func TestStrSetUnion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1", "2")
		s2.Add("3", "4")
		s3 := s1.Union(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("2"), true)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s3.Contains("4"), true)
	})
}

func TestStrSetDiff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("2"), true)
		t.Assert(s3.Contains("3"), false)
		t.Assert(s3.Contains("4"), false)

		s4 := s1
		s5 := s1.Diff(s2, s4)
		t.Assert(s5.Contains("1"), true)
		t.Assert(s5.Contains("2"), true)
		t.Assert(s5.Contains("3"), false)
		t.Assert(s5.Contains("4"), false)
	})
}

func TestStrSetIntersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains("1"), false)
		t.Assert(s3.Contains("2"), false)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s3.Contains("4"), false)
	})
}

func TestStrSetComplement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains("1"), false)
		t.Assert(s3.Contains("2"), false)
		t.Assert(s3.Contains("4"), true)
		t.Assert(s3.Contains("5"), true)
	})
}

func TestNewIntSetFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSetFrom([]int{1, 2, 3, 4})
		s2 := gset.NewIntSetFrom([]int{5, 6, 7, 8})
		t.Assert(s1.Contains(3), true)
		t.Assert(s1.Contains(5), false)
		t.Assert(s2.Contains(3), false)
		t.Assert(s2.Contains(5), true)
	})
}

func TestStrSetMerge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Merge(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("6"), false)
		t.Assert(s3.Contains("4"), true)
		t.Assert(s3.Contains("5"), true)
	})
}

func TestNewStrSetFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		t.Assert(s1.Contains("b"), true)
		t.Assert(s1.Contains("d"), false)
	})
}

func TestStrSetJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, "b"), true)
		t.Assert(strings.Contains(str1, "d"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSet()
		t.Assert(s1.Join(","), "")
		s1.Add("a", `"b"`, `\c`)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, `"b"`), true)
		t.Assert(strings.Contains(str1, `\c`), true)
		t.Assert(strings.Contains(str1, `a`), true)
	})
}

func TestStrSetString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "b"), true)
		t.Assert(strings.Contains(str1, "d"), false)
		s1 = nil
		t.Assert(s1.String(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New(true)
		s1.Add("a", "a2", "b", "c")
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "["), true)
		t.Assert(strings.Contains(str1, "]"), true)
		t.Assert(strings.Contains(str1, "a2"), true)
	})
}

func TestStrSetSum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		s2 := gset.NewIntSetFrom([]int{2, 3, 4}, true)
		t.Assert(s1.Sum(), 0)
		t.Assert(s2.Sum(), 9)
	})
}

func TestStrSetSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		t.Assert(s1.Size(), 3)

	})
}

func TestStrSetRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		s1.Remove("b")
		t.Assert(s1.Contains("b"), false)
		t.Assert(s1.Contains("c"), true)
	})
}

func TestStrSetPop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := []string{"a", "b", "c", "d"}
		s := gset.NewStrSetFrom(a, true)
		t.Assert(s.Size(), 4)
		t.AssertIN(s.Pop(), a)
		t.Assert(s.Size(), 3)
		t.AssertIN(s.Pop(), a)
		t.Assert(s.Size(), 2)

		s1 := gset.StrSet{}
		t.Assert(s1.Pop(), "")
	})
}

func TestStrSetPops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := []string{"a", "b", "c", "d"}
		s := gset.NewStrSetFrom(a, true)
		array := s.Pops(2)
		t.Assert(len(array), 2)
		t.Assert(s.Size(), 2)
		t.AssertIN(array, a)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(2), a)
		t.Assert(s.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet(true)
		a := []string{"1", "2", "3", "4"}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestStrSetAddIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet(true)
		s.Add("1")
		t.Assert(s.Contains("1"), true)
		t.Assert(s.AddIfNotExist("1"), false)
		t.Assert(s.AddIfNotExist("2"), true)
		t.Assert(s.Contains("2"), true)
		t.Assert(s.AddIfNotExist("2"), false)
		t.Assert(s.Contains("2"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.StrSet{}
		t.Assert(s.AddIfNotExist("1"), true)
	})
}

func TestStrSetAddIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet(true)
		s.Add("1")
		t.Assert(s.Contains("1"), true)
		t.Assert(s.Contains("2"), false)
		t.Assert(s.AddIfNotExistFunc("2", func() bool { return false }), false)
		t.Assert(s.Contains("2"), false)
		t.Assert(s.AddIfNotExistFunc("2", func() bool { return true }), true)
		t.Assert(s.Contains("2"), true)
		t.Assert(s.AddIfNotExistFunc("2", func() bool { return true }), false)
		t.Assert(s.Contains("2"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet(true)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFunc("1", func() bool {
				time.Sleep(100 * time.Millisecond)
				return true
			})
			t.Assert(r, false)
		}()
		s.Add("1")
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.StrSet{}
		t.Assert(s.AddIfNotExistFunc("1", func() bool { return true }), true)
	})
}

func TestStrSetAddIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewStrSet(true)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFuncLock("1", func() bool {
				time.Sleep(500 * time.Millisecond)
				return true
			})
			t.Assert(r, true)
		}()
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFuncLock("1", func() bool {
				return true
			})
			t.Assert(r, false)
		}()
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.StrSet{}
		t.Assert(s.AddIfNotExistFuncLock("1", func() bool { return true }), true)
	})
}

func TestStrSetJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := gset.NewStrSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := gset.NewStrSet()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains("a"), true)
		t.Assert(a2.Contains("b"), true)
		t.Assert(a2.Contains("c"), true)
		t.Assert(a2.Contains("d"), true)
		t.Assert(a2.Contains("e"), false)

		var a3 gset.StrSet
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Contains("a"), true)
		t.Assert(a3.Contains("b"), true)
		t.Assert(a3.Contains("c"), true)
		t.Assert(a3.Contains("d"), true)
		t.Assert(a3.Contains("e"), false)
	})
}

func TestStrSetWalk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			set    gset.StrSet
			names  = g.SliceStr{"user", "user_detail"}
			prefix = "gf_"
		)
		set.Add(names...)
		// Add prefix for given table names.
		set.Walk(func(item string) string {
			return prefix + item
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains("gf_user"), true)
		t.Assert(set.Contains("gf_user_detail"), true)
	})
}

func TestStrSetUnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *gset.StrSet
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  []byte(`["1","2","3"]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("1"), true)
		t.Assert(v.Set.Contains("2"), true)
		t.Assert(v.Set.Contains("3"), true)
		t.Assert(v.Set.Contains("4"), false)
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  g.SliceStr{"1", "2", "3"},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("1"), true)
		t.Assert(v.Set.Contains("2"), true)
		t.Assert(v.Set.Contains("3"), true)
		t.Assert(v.Set.Contains("4"), false)
	})
}

func TestStrSetDeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := gset.NewStrSet()
		set.Add("1", "2", "3")

		copySet := set.DeepCopy().(*gset.StrSet)
		copySet.Add("4")
		t.AssertNE(set.Size(), copySet.Size())
		t.AssertNE(set.String(), copySet.String())

		set = nil
		t.AssertNil(set.DeepCopy())
	})
}
