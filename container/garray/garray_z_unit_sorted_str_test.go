// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestNewSortedStrArrayComparator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayComparator(func(a, b string) int {
			return gstr.Compare(a, b)
		})
		s1.Add(a1...)
		t.Assert(s1.Len(), 4)
		t.Assert(s1, []string{"a", "b", "c", "d"})
	})
}

func TestNewSortedStrArrayFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayFrom(a1, true)
		t.Assert(s1, []string{"a", "b", "c", "d"})
		s2 := garray.NewSortedStrArrayFrom(a1, false)
		t.Assert(s2, []string{"a", "b", "c", "d"})
	})
}

func TestNewSortedStrArrayFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayFromCopy(a1, true)
		t.Assert(s1.Len(), 4)
		t.Assert(s1, []string{"a", "b", "c", "d"})
	})
}

func TestSortedStrArray_SetArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		a2 := []string{"f", "g", "h"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.SetArray(a2)
		t.Assert(array1.Len(), 3)
		t.Assert(array1.Contains("d"), false)
		t.Assert(array1.Contains("b"), false)
		t.Assert(array1.Contains("g"), true)
	})
}

func TestSortedStrArray_ContainsI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := garray.NewSortedStrArray()
		s.Append("a", "b", "C")
		t.Assert(s.Contains("A"), false)
		t.Assert(s.Contains("a"), true)
		t.Assert(s.ContainsI("A"), true)

		s = garray.NewSortedStrArray()
		t.Assert(s.Contains("A"), false)
	})
}

func TestSortedStrArray_Sort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)

		t.Assert(array1, []string{"a", "b", "c", "d"})
		array1.Sort()
		t.Assert(array1.Len(), 4)
		t.Assert(array1.Contains("c"), true)
		t.Assert(array1, []string{"a", "b", "c", "d"})
	})
}

func TestSortedStrArray_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		v, ok := array1.Get(2)
		t.Assert(v, "c")
		t.Assert(ok, true)

		v, ok = array1.Get(0)
		t.Assert(v, "a")
		t.Assert(ok, true)

		v, ok = array1.Get(99)
		t.Assert(v, "")
		t.Assert(ok, false)
	})
}

func TestSortedStrArray_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)

		v, ok := array1.Remove(-1)
		t.Assert(v, "")
		t.Assert(ok, false)

		v, ok = array1.Remove(100000)
		t.Assert(v, "")
		t.Assert(ok, false)

		v, ok = array1.Remove(2)
		t.Assert(v, "c")
		t.Assert(ok, true)

		v, ok = array1.Get(2)
		t.Assert(v, "d")
		t.Assert(ok, true)

		t.Assert(array1.Len(), 3)
		t.Assert(array1.Contains("c"), false)

		v, ok = array1.Remove(0)
		t.Assert(v, "a")
		t.Assert(ok, true)

		t.Assert(array1.Len(), 2)
		t.Assert(array1.Contains("a"), false)

		v, ok = array1.Remove(1)
		t.Assert(v, "d")
		t.Assert(ok, true)

		t.Assert(array1.Len(), 1)
	})
}

func TestSortedStrArray_PopLeft(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		v, ok := array1.PopLeft()
		t.Assert(v, "a")
		t.Assert(ok, true)
		t.Assert(array1.Len(), 4)
		t.Assert(array1.Contains("a"), false)
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"1", "2", "3"})
		v, ok := array.PopLeft()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array.Len(), 2)
		v, ok = array.PopLeft()
		t.Assert(v, 2)
		t.Assert(ok, true)
		t.Assert(array.Len(), 1)
		v, ok = array.PopLeft()
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array.Len(), 0)
	})
}

func TestSortedStrArray_PopRight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		v, ok := array1.PopRight()
		t.Assert(v, "e")
		t.Assert(ok, ok)
		t.Assert(array1.Len(), 4)
		t.Assert(array1.Contains("e"), false)
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"1", "2", "3"})
		v, ok := array.PopRight()
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array.Len(), 2)

		v, ok = array.PopRight()
		t.Assert(v, 2)
		t.Assert(ok, true)
		t.Assert(array.Len(), 1)

		v, ok = array.PopRight()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array.Len(), 0)
	})
}

func TestSortedStrArray_PopRand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1, ok := array1.PopRand()
		t.Assert(ok, true)
		t.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
		t.Assert(array1.Len(), 4)
		t.Assert(array1.Contains(s1), false)
	})
}

func TestSortedStrArray_PopRands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRands(2)
		t.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
		t.Assert(array1.Len(), 3)
		t.Assert(len(s1), 2)

		s1 = array1.PopRands(4)
		t.Assert(len(s1), 3)
		t.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
	})
}

func TestSortedStrArray_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArray()
		v, ok := array.PopLeft()
		t.Assert(v, "")
		t.Assert(ok, false)
		t.Assert(array.PopLefts(10), nil)

		v, ok = array.PopRight()
		t.Assert(v, "")
		t.Assert(ok, false)
		t.Assert(array.PopRights(10), nil)

		v, ok = array.PopRand()
		t.Assert(v, "")
		t.Assert(ok, false)
		t.Assert(array.PopRands(10), nil)
	})
}

func TestSortedStrArray_PopLefts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopLefts(2)
		t.Assert(s1, []string{"a", "b"})
		t.Assert(array1.Len(), 3)
		t.Assert(len(s1), 2)

		s1 = array1.PopLefts(4)
		t.Assert(len(s1), 3)
		t.Assert(s1, []string{"c", "d", "e"})
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_PopRights(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRights(2)
		t.Assert(s1, []string{"f", "g"})
		t.Assert(array1.Len(), 5)
		t.Assert(len(s1), 2)
		s1 = array1.PopRights(6)
		t.Assert(len(s1), 5)
		t.Assert(s1, []string{"a", "b", "c", "d", "e"})
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_Range(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a1, true)
		s1 := array1.Range(2, 4)
		t.Assert(len(s1), 2)
		t.Assert(s1, []string{"c", "d"})

		s1 = array1.Range(-1, 2)
		t.Assert(len(s1), 2)
		t.Assert(s1, []string{"a", "b"})

		s1 = array1.Range(4, 8)
		t.Assert(len(s1), 3)
		t.Assert(s1, []string{"e", "f", "g"})
		t.Assert(array1.Range(10, 2), nil)

		s2 := array2.Range(2, 4)
		t.Assert(s2, []string{"c", "d"})

	})
}

func TestSortedStrArray_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		a2 := []string{"1", "2", "3", "4", "a"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a2)
		t.Assert(array1.Sum(), 0)
		t.Assert(array2.Sum(), 10)
	})
}

func TestSortedStrArray_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.Clone()
		t.Assert(array1, array2)
		array1.Remove(1)
		t.Assert(array2.Len(), 7)
	})
}

func TestSortedStrArray_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.Clear()
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_SubSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a1, true)
		s1 := array1.SubSlice(1, 3)
		t.Assert(len(s1), 3)
		t.Assert(s1, []string{"b", "c", "d"})
		t.Assert(array1.Len(), 7)

		s2 := array1.SubSlice(1, 10)
		t.Assert(len(s2), 6)

		s3 := array1.SubSlice(10, 2)
		t.Assert(len(s3), 0)

		s3 = array1.SubSlice(-5, 2)
		t.Assert(s3, []string{"c", "d"})

		s3 = array1.SubSlice(-10, 2)
		t.Assert(s3, nil)

		s3 = array1.SubSlice(1, -2)
		t.Assert(s3, nil)

		t.Assert(array2.SubSlice(1, 3), []string{"b", "c", "d"})
	})
}

func TestSortedStrArray_Len(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		t.Assert(array1.Len(), 7)

	})
}

func TestSortedStrArray_Rand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		v, ok := array1.Rand()
		t.AssertIN(v, []string{"e", "a", "d"})
		t.Assert(ok, true)

		array2 := garray.NewSortedStrArrayFrom([]string{})
		v, ok = array2.Rand()
		t.Assert(v, "")
		t.Assert(ok, false)
	})
}

func TestSortedStrArray_Rands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.Rands(2)

		t.AssertIN(s1, []string{"e", "a", "d"})
		t.Assert(len(s1), 2)

		s1 = array1.Rands(4)
		t.Assert(len(s1), 4)

		array2 := garray.NewSortedStrArrayFrom([]string{})
		val := array2.Rands(1)
		t.Assert(val, nil)
	})
}

func TestSortedStrArray_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		t.Assert(array1.Join(","), `a,d,e`)
		t.Assert(array1.Join("."), `a.d.e`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", `"b"`, `\c`}
		array1 := garray.NewSortedStrArrayFrom(a1)
		t.Assert(array1.Join("."), `"b".\c.a`)
	})

	gtest.C(t, func(t *gtest.T) {
		array1 := garray.NewSortedStrArrayFrom([]string{})
		t.Assert(array1.Join("."), "")
	})
}

func TestSortedStrArray_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		t.Assert(array1.String(), `["a","d","e"]`)

		array1 = nil
		t.Assert(array1.String(), "")
	})
}

func TestSortedStrArray_CountValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		m1 := array1.CountValues()
		t.Assert(m1["a"], 2)
		t.Assert(m1["d"], 1)

	})
}

func TestSortedStrArray_Chunk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.Chunk(2)
		t.Assert(len(array2), 3)
		t.Assert(len(array2[0]), 2)
		t.Assert(array2[1], []string{"c", "d"})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []string{"1", "2", "3"})
		t.Assert(chunks[1], []string{"4", "5"})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5", "6"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []string{"1", "2"})
		t.Assert(chunks[1], []string{"3", "4"})
		t.Assert(chunks[2], []string{"5", "6"})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5", "6"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []string{"1", "2", "3"})
		t.Assert(chunks[1], []string{"4", "5", "6"})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedStrArray_SetUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "1", "2", "2", "3", "3", "2", "2"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.SetUnique(true)
		t.Assert(array2.Len(), 3)
		t.Assert(array2, []string{"1", "2", "3"})
	})
}

func TestSortedStrArray_Unique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "1", "2", "2", "3", "3", "2", "2"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.Unique()
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"1", "2", "3"})

		array2 := garray.NewSortedStrArrayFrom([]string{})
		array2.Unique()
		t.Assert(array2.Len(), 0)
		t.Assert(array2, []string{})
	})
}

func TestSortedStrArray_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		// go1
		go a1.LockFunc(func(n1 []string) { // 读写锁
			time.Sleep(2 * time.Second) // 暂停2秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		// go2
		go func() {
			time.Sleep(100 * time.Millisecond) // 故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 // 等待go1完成

		// 防止ci抖动,以豪秒为单位
		t.AssertGT(t2-t1, 20) // go1加的读写互斥锁，所go2读的时候被阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_RLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		// go1
		go a1.RLockFunc(func(n1 []string) { // 读锁
			time.Sleep(2 * time.Second) // 暂停1秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		// go2
		go func() {
			time.Sleep(100 * time.Millisecond) // 故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 // 等待go1完成

		// 防止ci抖动,以豪秒为单位
		t.AssertLT(t2-t1, 20) // go1加的读锁，所go2读的时候，并没有阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStrArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStrArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := garray.NewSortedStrArrayFrom(s1)

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i1).Len(), 9)
		t.Assert(a1.Merge(i2).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedStrArray_Json(t *testing.T) {
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewSortedStrArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(a2.Slice(), s2)
		t.Assert(a2.Interfaces(), s2)

		var a3 garray.SortedStrArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := *garray.NewSortedStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewSortedStrArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(a2.Slice(), s2)
		t.Assert(a2.Interfaces(), s2)

		var a3 garray.SortedStrArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores *garray.SortedStrArray
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []string{"A+", "A", "A"},
		}
		b, err := json.Marshal(data)
		t.AssertNil(err)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.AssertNil(err)
		t.Assert(user.Name, data["Name"])
		t.Assert(user.Scores, []string{"A", "A", "A+"})
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores garray.SortedStrArray
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []string{"A+", "A", "A"},
		}
		b, err := json.Marshal(data)
		t.AssertNil(err)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.AssertNil(err)
		t.Assert(user.Name, data["Name"])
		t.Assert(user.Scores, []string{"A", "A", "A+"})
	})
}

func TestSortedStrArray_Iterator(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedStrArrayFrom(slice)
	gtest.C(t, func(t *gtest.T) {
		array.Iterator(func(k int, v string) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	gtest.C(t, func(t *gtest.T) {
		array.IteratorAsc(func(k int, v string) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	gtest.C(t, func(t *gtest.T) {
		array.IteratorDesc(func(k int, v string) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		array.Iterator(func(k int, v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		array.IteratorAsc(func(k int, v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		array.IteratorDesc(func(k int, v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestSortedStrArray_RemoveValue(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedStrArrayFrom(slice)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestSortedStrArray_RemoveValues(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedStrArrayFrom(slice)
	gtest.C(t, func(t *gtest.T) {
		array.RemoveValues("a", "b", "c")
		t.Assert(array.Slice(), g.SliceStr{"d"})
	})
}

func TestSortedStrArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *garray.SortedStrArray
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`["1","3","2"]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.SliceStr{"1", "2", "3"})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": g.SliceStr{"1", "3", "2"},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.SliceStr{"1", "2", "3"})
	})
}

func TestSortedStrArray_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"", "1", "2", "0"})
		t.Assert(array.FilterEmpty(), g.SliceStr{"0", "1", "2"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"1", "2"})
		t.Assert(array.FilterEmpty(), g.SliceStr{"1", "2"})
	})
}

func TestSortedStrArray_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"1", "2"})
		t.Assert(array.Walk(func(value string) string {
			return "key-" + value
		}), g.Slice{"key-1", "key-2"})
	})
}

func TestSortedStrArray_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedStrArrayFrom([]string{"a", "b", "c", "d"})
		copyArray := array.DeepCopy().(*garray.SortedStrArray)
		array.Add("e")
		copyArray.Add("f")
		cval, _ := copyArray.Get(4)
		val, _ := array.Get(4)
		t.AssertNE(cval, val)
	})
}
