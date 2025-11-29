// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"sync"
	"testing"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_ListKVMap_NewListKVMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string](true)
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_ListKVMap_NewListKVMapFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := map[string]string{"a": "1", "b": "2"}
		m := gmap.NewListKVMapFrom(data)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
	})

	gtest.C(t, func(t *gtest.T) {
		data := map[int]int{1: 10, 2: 20}
		m := gmap.NewListKVMapFrom(data, true)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get(1), 10)
		t.Assert(m.Get(2), 20)
	})
}

func Test_ListKVMap_Set_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("a", "1")
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Size(), 1)

		m.Set("b", "2")
		t.Assert(m.Get("b"), "2")
		t.Assert(m.Size(), 2)

		// Set existing key
		m.Set("a", "10")
		t.Assert(m.Get("a"), "10")
		t.Assert(m.Size(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[int, int]()
		m.Set(1, 100)
		m.Set(2, 200)
		t.Assert(m.Get(1), 100)
		t.Assert(m.Get(2), 200)
	})
}

func Test_ListKVMap_Sets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Sets(map[string]string{"a": "1", "b": "2", "c": "3"})
		t.Assert(m.Size(), 3)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
		t.Assert(m.Get("c"), "3")
	})

	gtest.C(t, func(t *gtest.T) {
		data := map[string]string{"x": "10", "y": "20"}
		m := gmap.NewListKVMapFrom(data)
		m.Sets(map[string]string{"a": "1", "b": "2"})
		t.Assert(m.Size(), 4)
		t.Assert(m.Get("x"), "10")
		t.Assert(m.Get("a"), "1")
	})
}

func Test_ListKVMap_Search(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})

		v, found := m.Search("a")
		t.Assert(found, true)
		t.Assert(v, "1")

		v, found = m.Search("c")
		t.Assert(found, false)
		t.Assert(v, "")
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[int, string]()
		v, found := m.Search(1)
		t.Assert(found, false)
		t.Assert(v, "")
	})
}

func Test_ListKVMap_Contains(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		t.Assert(m.Contains("a"), true)
		t.Assert(m.Contains("b"), true)
		t.Assert(m.Contains("c"), false)
	})
}

func Test_ListKVMap_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		t.Assert(m.Size(), 2)

		v := m.Remove("a")
		t.Assert(v, "1")
		t.Assert(m.Contains("a"), false)
		t.Assert(m.Size(), 1)

		v = m.Remove("c")
		t.Assert(v, "")
		t.Assert(m.Size(), 1)
	})
}

func Test_ListKVMap_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		t.Assert(m.Size(), 3)

		m.Removes([]string{"a", "c"})
		t.Assert(m.Size(), 1)
		t.Assert(m.Contains("a"), false)
		t.Assert(m.Contains("c"), false)
		t.Assert(m.Contains("b"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		m.Removes([]string{"x", "y"})
		t.Assert(m.Size(), 2)
	})
}

func Test_ListKVMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		t.Assert(m.Size(), 2)

		k, v := m.Pop()
		t.AssertIN(k, []string{"a", "b"})
		t.AssertIN(v, []string{"1", "2"})
		t.Assert(m.Size(), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		k, v := m.Pop()
		t.Assert(k, "")
		t.Assert(v, "")
	})
}

func Test_ListKVMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		t.Assert(m.Size(), 3)

		popped := m.Pops(2)
		t.Assert(len(popped), 2)
		t.Assert(m.Size(), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		popped := m.Pops(-1)
		t.Assert(len(popped), 3)
		t.Assert(m.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		popped := m.Pops(10)
		t.Assert(len(popped), 2)
		t.Assert(m.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		popped := m.Pops(1)
		t.AssertNil(popped)
	})
}

func Test_ListKVMap_Keys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		keys := m.Keys()
		t.Assert(len(keys), 3)
		t.AssertIN("a", keys)
		t.AssertIN("b", keys)
		t.AssertIN("c", keys)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[int, string]()
		keys := m.Keys()
		t.Assert(len(keys), 0)
	})
}

func Test_ListKVMap_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		values := m.Values()
		t.Assert(len(values), 3)
		t.AssertIN("1", values)
		t.AssertIN("2", values)
		t.AssertIN("3", values)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		values := m.Values()
		t.Assert(len(values), 0)
	})
}

func Test_ListKVMap_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		t.Assert(m.Size(), 0)

		m.Set("a", "1")
		t.Assert(m.Size(), 1)

		m.Set("b", "2")
		t.Assert(m.Size(), 2)

		m.Remove("a")
		t.Assert(m.Size(), 1)

		m.Clear()
		t.Assert(m.Size(), 0)
	})
}

func Test_ListKVMap_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		t.Assert(m.IsEmpty(), true)

		m.Set("a", "1")
		t.Assert(m.IsEmpty(), false)

		m.Remove("a")
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_ListKVMap_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		t.Assert(m.Size(), 3)

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
		t.Assert(m.Get("a"), "")
	})
}

func Test_ListKVMap_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		data := m.Map()
		t.Assert(data["a"], "1")
		t.Assert(data["b"], "2")
		t.Assert(len(data), 2)
	})
}

func Test_ListKVMap_MapStrAny(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]int{"a": 1, "b": 2})
		data := m.MapStrAny()
		t.Assert(len(data), 2)
		t.Assert(data["a"], 1)
		t.Assert(data["b"], 2)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[int]string{1: "a", 2: "b"})
		data := m.MapStrAny()
		t.Assert(len(data), 2)
		t.Assert(data["1"], "a")
		t.Assert(data["2"], "b")
	})
}

func Test_ListKVMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "", "b": "2", "c": "3"})
		t.Assert(m.Size(), 3)

		m.FilterEmpty()
		t.Assert(m.Size(), 2)
		t.Assert(m.Contains("a"), false)
		t.Assert(m.Contains("b"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]int{"a": 0, "b": 1, "c": 2})
		t.Assert(m.Size(), 3)

		m.FilterEmpty()
		t.Assert(m.Size(), 2)
		t.Assert(m.Contains("a"), false)
	})
}

func Test_ListKVMap_GetOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()

		v := m.GetOrSet("a", "1")
		t.Assert(v, "1")
		t.Assert(m.Get("a"), "1")

		v = m.GetOrSet("a", "10")
		t.Assert(v, "1")
		t.Assert(m.Get("a"), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]int{"a": 10})

		v := m.GetOrSet("a", 20)
		t.Assert(v, 10)

		v = m.GetOrSet("b", 30)
		t.Assert(v, 30)
		t.Assert(m.Get("b"), 30)
	})
}

func Test_ListKVMap_GetOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()

		v := m.GetOrSetFunc("a", func() string { return "1" })
		t.Assert(v, "1")

		v = m.GetOrSetFunc("a", func() string { return "10" })
		t.Assert(v, "1")

		v = m.GetOrSetFunc("b", func() string { return "2" })
		t.Assert(v, "2")
	})
}

func Test_ListKVMap_GetOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		counter := 0

		v := m.GetOrSetFuncLock("a", func() int {
			counter++
			return 10
		})
		t.Assert(v, 10)
		t.Assert(counter, 1)

		v = m.GetOrSetFuncLock("a", func() int {
			counter++
			return 20
		})
		t.Assert(v, 10)
		t.Assert(counter, 1)
	})
}

func Test_ListKVMap_SetIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()

		ok := m.SetIfNotExist("a", "1")
		t.Assert(ok, true)
		t.Assert(m.Get("a"), "1")

		ok = m.SetIfNotExist("a", "10")
		t.Assert(ok, false)
		t.Assert(m.Get("a"), "1")
	})
}

func Test_ListKVMap_SetIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()

		ok := m.SetIfNotExistFunc("a", func() int { return 10 })
		t.Assert(ok, true)
		t.Assert(m.Get("a"), 10)

		ok = m.SetIfNotExistFunc("a", func() int { return 20 })
		t.Assert(ok, false)
		t.Assert(m.Get("a"), 10)
	})
}

func Test_ListKVMap_SetIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		counter := 0

		ok := m.SetIfNotExistFuncLock("a", func() string {
			counter++
			return "1"
		})
		t.Assert(ok, true)
		t.Assert(counter, 1)

		ok = m.SetIfNotExistFuncLock("a", func() string {
			counter++
			return "2"
		})
		t.Assert(ok, false)
		t.Assert(counter, 1)
	})
}

func Test_ListKVMap_GetVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})

		v := m.GetVar("a")
		t.AssertNE(v, nil)
		t.Assert(v.Val(), "1")

		v = m.GetVar("c")
		t.Assert(v.Val(), nil)
	})
}

func Test_ListKVMap_GetVarOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()

		v := m.GetVarOrSet("a", "1")
		t.AssertNE(v, nil)
		t.Assert(v.Val(), "1")

		v = m.GetVarOrSet("a", "10")
		t.Assert(v.Val(), "1")
	})
}

func Test_ListKVMap_GetVarOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()

		v := m.GetVarOrSetFunc("a", func() int { return 10 })
		t.AssertNE(v, nil)
		t.Assert(v.Val(), 10)

		v = m.GetVarOrSetFunc("a", func() int { return 20 })
		t.Assert(v.Val(), 10)
	})
}

func Test_ListKVMap_GetVarOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()

		v := m.GetVarOrSetFuncLock("a", func() string { return "1" })
		t.AssertNE(v, nil)
		t.Assert(v.Val(), "1")

		v = m.GetVarOrSetFuncLock("a", func() string { return "10" })
		t.Assert(v.Val(), "1")
	})
}

func Test_ListKVMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := map[string]string{"a": "1", "b": "2", "c": "3"}
		m := gmap.NewListKVMapFrom(data)

		count := 0
		m.Iterator(func(k string, v string) bool {
			t.Assert(data[k], v)
			count++
			return true
		})
		t.Assert(count, 3)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})

		count := 0
		m.Iterator(func(k int, v string) bool {
			count++
			return count < 2
		})
		t.Assert(count, 2)
	})
}

func Test_ListKVMap_IteratorAsc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("k1", "v1")
		m.Set("k2", "v2")
		m.Set("k3", "v3")

		var keys []string
		var values []string
		m.IteratorAsc(func(k string, v string) bool {
			keys = append(keys, k)
			values = append(values, v)
			return true
		})
		t.Assert(keys, g.Slice{"k1", "k2", "k3"})
		t.Assert(values, g.Slice{"v1", "v2", "v3"})
	})
}

func Test_ListKVMap_IteratorDesc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("k1", "v1")
		m.Set("k2", "v2")
		m.Set("k3", "v3")

		var keys []string
		var values []string
		m.IteratorDesc(func(k string, v string) bool {
			keys = append(keys, k)
			values = append(values, v)
			return true
		})
		t.Assert(keys, g.Slice{"k3", "k2", "k1"})
		t.Assert(values, g.Slice{"v3", "v2", "v1"})
	})
}

func Test_ListKVMap_Replace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		t.Assert(m.Size(), 2)

		m.Replace(map[string]string{"x": "10", "y": "20", "z": "30"})
		t.Assert(m.Size(), 3)
		t.Assert(m.Get("a"), "")
		t.Assert(m.Get("x"), "10")
	})
}

func Test_ListKVMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		m2 := m.Clone()

		t.Assert(m2.Get("a"), "1")
		t.Assert(m2.Get("b"), "2")
		t.Assert(m2.Size(), 2)

		m.Set("a", "10")
		t.Assert(m2.Get("a"), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]int{"a": 1, "b": 2}, false)
		m2 := m.Clone(true)

		t.Assert(m2.Size(), 2)
	})
}

func Test_ListKVMap_Flip(t *testing.T) {
	// Test with same type for key and value (string -> string)
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		err := m.Flip()
		t.AssertNil(err)

		t.Assert(m.Get("1"), "a")
		t.Assert(m.Get("2"), "b")
		t.Assert(m.Get("3"), "c")
	})

	// Test with same type for key and value (int -> int)
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[int]int{1: 10, 2: 20})
		err := m.Flip()
		t.AssertNil(err)

		t.Assert(m.Get(10), 1)
		t.Assert(m.Get(20), 2)
	})
}

func Test_ListKVMap_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewListKVMapFrom(map[string]string{"a": "1"})
		m2 := gmap.NewListKVMapFrom(map[string]string{"b": "2", "c": "3"})

		m1.Merge(m2)
		t.Assert(m1.Size(), 3)
		t.Assert(m1.Get("a"), "1")
		t.Assert(m1.Get("b"), "2")
		t.Assert(m1.Get("c"), "3")
	})

	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewListKVMap[string, int]()
		m2 := gmap.NewListKVMapFrom(map[string]int{"a": 10, "b": 20})

		m1.Merge(m2)
		t.Assert(m1.Size(), 2)
		t.Assert(m1.Get("a"), 10)
	})

	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewListKVMapFrom(map[string]string{"a": "1"})
		m2 := gmap.NewListKVMapFrom(map[string]string{"a": "10", "b": "2"})

		m1.Merge(m2)
		t.Assert(m1.Get("a"), "10")
	})
}

func Test_ListKVMap_Merge_Self(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		m.Merge(m)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
	})
}

func Test_ListKVMap_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1"})
		s := m.String()
		t.AssertNE(s, "")
		t.AssertIN("a", s)
	})

	gtest.C(t, func(t *gtest.T) {
		var m *gmap.ListKVMap[string, string]
		s := m.String()
		t.Assert(s, "")
	})
}

func Test_ListKVMap_MarshalJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		b, err := json.Marshal(m)
		t.AssertNil(err)
		t.AssertNE(b, nil)

		var data map[string]int
		err = json.Unmarshal(b, &data)
		t.AssertNil(err)
		t.Assert(data["a"], 1)
		t.Assert(data["b"], 2)
	})

	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		b, err := m.MarshalJSON()
		t.AssertNil(err)
		t.Assert(string(b), "{}")
	})
}

func Test_ListKVMap_UnmarshalJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		data := []byte(`{"a":1,"b":2,"c":3}`)

		err := json.UnmarshalUseNumber(data, m)
		t.AssertNil(err)
		t.Assert(m.Get("a"), 1)
		t.Assert(m.Get("b"), 2)
		t.Assert(m.Get("c"), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		var m gmap.ListKVMap[string, string]
		data := []byte(`{"x":"10","y":"20"}`)

		err := json.UnmarshalUseNumber(data, &m)
		t.AssertNil(err)
		t.Assert(m.Get("x"), "10")
		t.Assert(m.Get("y"), "20")
	})
}

func Test_ListKVMap_UnmarshalValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		err := m.UnmarshalValue(map[string]any{
			"a": "1",
			"b": "2",
		})
		t.AssertNil(err)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
	})
}

func Test_ListKVMap_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string][]string{
			"a": {"1", "2"},
			"b": {"3", "4"},
		})

		n := m.DeepCopy().(*gmap.ListKVMap[string, []string])
		t.Assert(n.Size(), 2)
		t.Assert(n.Get("a"), []string{"1", "2"})

		// Modifying original doesn't affect copy
		m.Get("a")[0] = "10"
		t.Assert(n.Get("a")[0], "1")
	})

	gtest.C(t, func(t *gtest.T) {
		var m *gmap.ListKVMap[string, int]
		n := m.DeepCopy()
		t.AssertNil(n)
	})
}

func Test_ListKVMap_Order(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("k1", "v1")
		m.Set("k2", "v2")
		m.Set("k3", "v3")
		t.Assert(m.Keys(), g.Slice{"k1", "k2", "k3"})
		t.Assert(m.Values(), g.Slice{"v1", "v2", "v3"})
	})
}

func Test_ListKVMap_Json_Sequence(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int32]()
		for i := 'z'; i >= 'a'; i-- {
			m.Set(string(i), i)
		}
		b, err := json.Marshal(m)
		t.AssertNil(err)
		t.Assert(b, `{"z":122,"y":121,"x":120,"w":119,"v":118,"u":117,"t":116,"s":115,"r":114,"q":113,"p":112,"o":111,"n":110,"m":109,"l":108,"k":107,"j":106,"i":105,"h":104,"g":103,"f":102,"e":101,"d":100,"c":99,"b":98,"a":97}`)
	})
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int32]()
		for i := 'a'; i <= 'z'; i++ {
			m.Set(string(i), i)
		}
		b, err := json.Marshal(m)
		t.AssertNil(err)
		t.Assert(b, `{"a":97,"b":98,"c":99,"d":100,"e":101,"f":102,"g":103,"h":104,"i":105,"j":106,"k":107,"l":108,"m":109,"n":110,"o":111,"p":112,"q":113,"r":114,"s":115,"t":116,"u":117,"v":118,"w":119,"x":120,"y":121,"z":122}`)
	})
}

// Test Set with nil data
func Test_ListKVMap_Set_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		m.Set("a", "1")
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Size(), 1)
	})
}

// Test Sets with nil data
func Test_ListKVMap_Sets_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		m.Sets(map[string]string{"a": "1", "b": "2"})
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
		t.Assert(m.Size(), 2)
	})
}

// Test GetOrSet with nil value (using any type)
func Test_ListKVMap_GetOrSet_NilValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, any]()
		v := m.GetOrSet("a", nil)
		t.Assert(v, nil)
		// nil interface value should not be stored
		t.Assert(m.Contains("a"), false)
	})
}

// Test Search with nil data
func Test_ListKVMap_Search_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		v, found := m.Search("a")
		t.Assert(found, false)
		t.Assert(v, "")
	})
}

// Test Get with nil data
func Test_ListKVMap_Get_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, int](nil)
		v := m.Get("a")
		t.Assert(v, 0)
	})
}

// Test Contains with nil data
func Test_ListKVMap_Contains_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		t.Assert(m.Contains("a"), false)
	})
}

// Test Remove with nil data
func Test_ListKVMap_Remove_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		v := m.Remove("a")
		t.Assert(v, "")
	})
}

// Test Removes with nil data
func Test_ListKVMap_Removes_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		m.Removes([]string{"a", "b"})
		t.Assert(m.Size(), 0)
	})
}

// Test Pops with size 0
func Test_ListKVMap_Pops_ZeroSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		popped := m.Pops(0)
		t.AssertNil(popped)
		t.Assert(m.Size(), 2)
	})
}

// Test Iterator early break
func Test_ListKVMap_Iterator_Break(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2", "c": "3"})
		count := 0
		m.Iterator(func(k string, v string) bool {
			count++
			return false // Break immediately
		})
		t.Assert(count, 1)
	})
}

// Test IteratorAsc with nil list
func Test_ListKVMap_IteratorAsc_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		count := 0
		m.IteratorAsc(func(k string, v string) bool {
			count++
			return true
		})
		t.Assert(count, 0)
	})
}

// Test IteratorDesc with nil list
func Test_ListKVMap_IteratorDesc_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		count := 0
		m.IteratorDesc(func(k string, v string) bool {
			count++
			return true
		})
		t.Assert(count, 0)
	})
}

// Test Map with nil list
func Test_ListKVMap_Map_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		data := m.Map()
		t.Assert(data, "{}")
	})
}

// Test MapStrAny with nil list
func Test_ListKVMap_MapStrAny_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		data := m.MapStrAny()
		t.Assert(data, "{}")
	})
}

// Test FilterEmpty with nil list
func Test_ListKVMap_FilterEmpty_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		m.FilterEmpty()
		t.Assert(m.Size(), 0)
	})
}

// Test Keys with nil list
func Test_ListKVMap_Keys_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		keys := m.Keys()
		t.Assert(len(keys), 0)
	})
}

// Test Values with nil list
func Test_ListKVMap_Values_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		values := m.Values()
		t.Assert(len(values), 0)
	})
}

// Test DeepCopy with nil list
func Test_ListKVMap_DeepCopy_NilList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		n := m.DeepCopy().(*gmap.ListKVMap[string, string])
		t.Assert(n.Size(), 0)
	})
}

// Concurrent safety tests
func Test_ListKVMap_Concurrent_Safe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		ch := make(chan int, 10)

		// Concurrent writes
		for i := 0; i < 10; i++ {
			go func(idx int) {
				m.Set(gconv.String(idx), idx)
				ch <- 1
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-ch
		}

		t.Assert(m.Size(), 10)
	})
}

func Test_ListKVMap_Concurrent_RW(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		m.Sets(map[string]int{"a": 1, "b": 2, "c": 3})

		ch := make(chan int, 20)

		// Concurrent reads and writes
		for i := 0; i < 10; i++ {
			go func() {
				_ = m.Get("a")
				ch <- 1
			}()
		}

		for i := 0; i < 10; i++ {
			go func(idx int) {
				m.Set(gconv.String(idx), idx)
				ch <- 1
			}(i)
		}

		for i := 0; i < 20; i++ {
			<-ch
		}

		t.Assert(m.Size(), 13)
	})
}

// Test concurrent GetOrSet
func Test_ListKVMap_Concurrent_GetOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		ch := make(chan int, 100)

		for i := 0; i < 100; i++ {
			go func(idx int) {
				m.GetOrSet("key", idx)
				ch <- 1
			}(i)
		}

		for i := 0; i < 100; i++ {
			<-ch
		}

		// Only one value should be set
		t.Assert(m.Size(), 1)
		t.Assert(m.Contains("key"), true)
	})
}

// Test concurrent SetIfNotExist
func Test_ListKVMap_Concurrent_SetIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		successCount := 0
		ch := make(chan bool, 100)

		for i := 0; i < 100; i++ {
			go func(idx int) {
				ok := m.SetIfNotExist("key", idx)
				ch <- ok
			}(i)
		}

		for i := 0; i < 100; i++ {
			if <-ch {
				successCount++
			}
		}

		// Only one goroutine should succeed
		t.Assert(successCount, 1)
		t.Assert(m.Size(), 1)
	})
}

// Test concurrent GetOrSetFunc
func Test_ListKVMap_Concurrent_GetOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		ch := make(chan int, 100)

		for i := 0; i < 100; i++ {
			go func(idx int) {
				m.GetOrSetFunc("key", func() int {
					return idx
				})
				ch <- 1
			}(i)
		}

		for i := 0; i < 100; i++ {
			<-ch
		}

		t.Assert(m.Size(), 1)
	})
}

// Test concurrent GetOrSetFuncLock
func Test_ListKVMap_Concurrent_GetOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		ch := make(chan int, 100)

		for i := 0; i < 100; i++ {
			go func(idx int) {
				m.GetOrSetFuncLock("key", func() int {
					return idx
				})
				ch <- 1
			}(i)
		}

		for i := 0; i < 100; i++ {
			<-ch
		}

		t.Assert(m.Size(), 1)
	})
}

// Test concurrent access to doSetWithLockCheck
func Test_ListKVMap_DoSetWithLockCheck_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int](true)
		var wg sync.WaitGroup
		results := make([]int, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				v := m.GetOrSet("key", idx)
				results[idx] = v
			}(i)
		}
		wg.Wait()

		// All results should be the same (the first value that was set)
		firstValue := results[0]
		for _, v := range results {
			t.Assert(v, firstValue)
		}
	})
}

// Test UnmarshalJSON with invalid JSON
func Test_ListKVMap_UnmarshalJSON_InvalidJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		err := m.UnmarshalJSON([]byte(`{invalid json}`))
		t.AssertNE(err, nil)
	})
}

// Test MarshalJSON error handling
func Test_ListKVMap_MarshalJSON_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, string]()
		m.Set("a", "1")
		b, err := m.MarshalJSON()
		t.AssertNil(err)
		t.Assert(string(b), `{"a":"1"}`)
	})
}

// Test empty map operations
func Test_ListKVMap_EmptyMapOperations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()

		// Test Keys on empty map
		keys := m.Keys()
		t.Assert(len(keys), 0)

		// Test Values on empty map
		values := m.Values()
		t.Assert(len(values), 0)

		// Test MapStrAny on empty map
		strAny := m.MapStrAny()
		t.Assert(len(strAny), 0)
	})
}

// Test FilterEmpty with various empty values
func Test_ListKVMap_FilterEmpty_Various(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, any]()
		m.Set("nil", nil)
		m.Set("zero", 0)
		m.Set("empty_string", "")
		m.Set("false", false)
		m.Set("valid", "value")
		m.Set("empty_slice", []int{})
		m.Set("empty_map", map[string]int{})

		t.Assert(m.Size(), 7)
		m.FilterEmpty()
		t.Assert(m.Size(), 1)
		t.Assert(m.Contains("valid"), true)
	})
}

// Test Clone with different safe modes
func Test_ListKVMap_Clone_SafeMode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Clone unsafe map to safe
		m := gmap.NewListKVMapFrom(map[string]int{"a": 1}, false)
		m2 := m.Clone(true)
		t.Assert(m2.Get("a"), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		// Clone safe map to unsafe
		m := gmap.NewListKVMapFrom(map[string]int{"a": 1}, true)
		m2 := m.Clone(false)
		t.Assert(m2.Get("a"), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		// Clone with inherited safe mode
		m := gmap.NewListKVMapFrom(map[string]int{"a": 1}, true)
		m2 := m.Clone()
		t.Assert(m2.Get("a"), 1)
	})
}

// Test SetIfNotExistFunc returning false when key exists
func Test_ListKVMap_SetIfNotExistFunc_KeyExists(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, int]()
		m.Set("a", 1)

		called := false
		ok := m.SetIfNotExistFunc("a", func() int {
			called = true
			return 2
		})
		t.Assert(ok, false)
		t.Assert(called, false) // Function should not be called if key exists
		t.Assert(m.Get("a"), 1)
	})
}

// Test struct with ListKVMap for UnmarshalValue
func Test_ListKVMap_UnmarshalValue_Struct(t *testing.T) {
	type V struct {
		Name string
		Map  *gmap.ListKVMap[string, string]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]any{
			"name": "john",
			"map":  []byte(`{"1":"v1","2":"v2"}`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("1"), "v1")
		t.Assert(v.Map.Get("2"), "v2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]any{
			"name": "john",
			"map": g.MapStrStr{
				"1": "v1",
				"2": "v2",
			},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("1"), "v1")
		t.Assert(v.Map.Get("2"), "v2")
	})
}

// Test GetOrSetFuncLock with nil data
func Test_ListKVMap_GetOrSetFuncLock_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		v := m.GetOrSetFuncLock("a", func() string { return "1" })
		t.Assert(v, "1")
		t.Assert(m.Get("a"), "1")
	})

	// Test with nil value (using any type)
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMap[string, any]()
		v := m.GetOrSetFuncLock("a", func() any { return nil })
		t.Assert(v, nil)
		// nil interface value should not be stored
		t.Assert(m.Contains("a"), false)
	})
}

// Test SetIfNotExistFuncLock with nil data
func Test_ListKVMap_SetIfNotExistFuncLock_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		ok := m.SetIfNotExistFuncLock("a", func() string { return "1" })
		t.Assert(ok, true)
		t.Assert(m.Get("a"), "1")
	})
}

// Test Merge with nil data
func Test_ListKVMap_Merge_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		m2 := gmap.NewListKVMapFrom(map[string]string{"a": "1", "b": "2"})
		m.Merge(m2)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
	})
}

// Test UnmarshalJSON with nil data
func Test_ListKVMap_UnmarshalJSON_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		err := m.UnmarshalJSON([]byte(`{"a":"1","b":"2"}`))
		t.AssertNil(err)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
	})
}

// Test UnmarshalValue with nil data
func Test_ListKVMap_UnmarshalValue_NilData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListKVMapFrom[string, string](nil)
		err := m.UnmarshalValue(map[string]any{"a": "1", "b": "2"})
		t.AssertNil(err)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get("a"), "1")
		t.Assert(m.Get("b"), "2")
	})
}
