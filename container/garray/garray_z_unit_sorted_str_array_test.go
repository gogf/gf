// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray_test

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func TestNewSortedStrArrayFrom(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayFrom(a1, true)
		gtest.Assert(s1, []string{"a", "b", "c", "d"})
		s2 := garray.NewSortedStrArrayFrom(a1, false)
		gtest.Assert(s2, []string{"a", "b", "c", "d"})
	})
}

func TestNewSortedStrArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayFromCopy(a1, true)
		gtest.Assert(s1.Len(), 4)
		gtest.Assert(s1, []string{"a", "b", "c", "d"})
	})
}

func TestSortedStrArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		a2 := []string{"f", "g", "h"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.SetArray(a2)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("d"), false)
		gtest.Assert(array1.Contains("b"), false)
		gtest.Assert(array1.Contains("g"), true)
	})
}

func TestSortedStrArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)

		gtest.Assert(array1, []string{"a", "b", "c", "d"})
		array1.Sort()
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains("c"), true)
		gtest.Assert(array1, []string{"a", "b", "c", "d"})
	})
}

func TestSortedStrArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Get(2), "c")
		gtest.Assert(array1.Get(0), "a")
	})
}

func TestSortedStrArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)

		gtest.Assert(array1.Remove(-1), "")
		gtest.Assert(array1.Remove(100000), "")

		gtest.Assert(array1.Remove(2), "c")
		gtest.Assert(array1.Get(2), "d")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("c"), false)

		gtest.Assert(array1.Remove(0), "a")
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(array1.Contains("a"), false)

		// 此时array1里的元素只剩下2个
		gtest.Assert(array1.Remove(1), "d")
		gtest.Assert(array1.Len(), 1)
	})
}

func TestSortedStrArray_PopLeft(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopLeft()
		gtest.Assert(s1, "a")
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains("a"), false)
	})
}

func TestSortedStrArray_PopRight(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRight()
		gtest.Assert(s1, "e")
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains("e"), false)
	})
}

func TestSortedStrArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRand()
		gtest.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains(s1), false)
	})
}

func TestSortedStrArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRands(2)
		gtest.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(len(s1), 2)

		s1 = array1.PopRands(4)
		gtest.Assert(len(s1), 3)
		gtest.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
	})
}

func TestSortedStrArray_PopLefts(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopLefts(2)
		gtest.Assert(s1, []string{"a", "b"})
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(len(s1), 2)

		s1 = array1.PopLefts(4)
		gtest.Assert(len(s1), 3)
		gtest.Assert(s1, []string{"c", "d", "e"})
	})
}

func TestSortedStrArray_PopRights(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRights(2)
		gtest.Assert(s1, []string{"f", "g"})
		gtest.Assert(array1.Len(), 5)
		gtest.Assert(len(s1), 2)
		s1 = array1.PopRights(6)
		gtest.Assert(len(s1), 5)
		gtest.Assert(s1, []string{"a", "b", "c", "d", "e"})
		gtest.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a1, true)
		s1 := array1.Range(2, 4)
		gtest.Assert(len(s1), 2)
		gtest.Assert(s1, []string{"c", "d"})

		s1 = array1.Range(-1, 2)
		gtest.Assert(len(s1), 2)
		gtest.Assert(s1, []string{"a", "b"})

		s1 = array1.Range(4, 8)
		gtest.Assert(len(s1), 3)
		gtest.Assert(s1, []string{"e", "f", "g"})
		gtest.Assert(array1.Range(10, 2), nil)

		s2 := array2.Range(2, 4)
		gtest.Assert(s2, []string{"c", "d"})

	})
}

func TestSortedStrArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		a2 := []string{"1", "2", "3", "4", "a"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a2)
		gtest.Assert(array1.Sum(), 0)
		gtest.Assert(array2.Sum(), 10)
	})
}

func TestSortedStrArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array1, array2)
		array1.Remove(1)
		gtest.Assert(array2.Len(), 7)
	})
}

func TestSortedStrArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a1, true)
		s1 := array1.SubSlice(1, 3)
		gtest.Assert(len(s1), 3)
		gtest.Assert(s1, []string{"b", "c", "d"})
		gtest.Assert(array1.Len(), 7)

		s2 := array1.SubSlice(1, 10)
		gtest.Assert(len(s2), 6)

		s3 := array1.SubSlice(10, 2)
		gtest.Assert(len(s3), 0)

		s3 = array1.SubSlice(-5, 2)
		gtest.Assert(s3, []string{"c", "d"})

		s3 = array1.SubSlice(-10, 2)
		gtest.Assert(s3, nil)

		s3 = array1.SubSlice(1, -2)
		gtest.Assert(s3, nil)

		gtest.Assert(array2.SubSlice(1, 3), []string{"b", "c", "d"})
	})
}

func TestSortedStrArray_Len(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Len(), 7)

	})
}

func TestSortedStrArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.AssertIN(array1.Rand(), []string{"e", "a", "d"})
	})
}

func TestSortedStrArray_Rands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.Rands(2)

		gtest.AssertIN(s1, []string{"e", "a", "d"})
		gtest.Assert(len(s1), 2)

		s1 = array1.Rands(4)
		gtest.AssertIN(s1, []string{"e", "a", "d"})
		gtest.Assert(len(s1), 3)
	})
}

func TestSortedStrArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Join(","), `a,d,e`)
		gtest.Assert(array1.Join("."), `a.d.e`)
	})

	gtest.Case(t, func() {
		a1 := []string{"a", `"b"`, `\c`}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Join("."), `"b".\c.a`)
	})
}

func TestSortedStrArray_String(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.String(), `["a","d","e"]`)
	})
}

func TestSortedStrArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		m1 := array1.CountValues()
		gtest.Assert(m1["a"], 2)
		gtest.Assert(m1["d"], 1)

	})
}

func TestSortedStrArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.Chunk(2)
		gtest.Assert(len(array2), 3)
		gtest.Assert(len(array2[0]), 2)
		gtest.Assert(array2[1], []string{"c", "d"})
		gtest.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedStrArray_SetUnique(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.SetUnique(true)
		gtest.Assert(array2.Len(), 4)
		gtest.Assert(array2, []string{"a", "c", "d", "e"})
	})
}

func TestSortedStrArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []string) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		gtest.AssertGT(t2-t1, 20) //go1加的读写互斥锁，所go2读的时候被阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []string) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		gtest.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候，并没有阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
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

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i1).Len(), 9)
		gtest.Assert(a1.Merge(i2).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedStrArray_Json(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		gtest.Assert(b1, b2)
		gtest.Assert(err1, err2)

		a2 := garray.NewSortedStrArray()
		err1 = json.Unmarshal(b2, &a2)
		gtest.Assert(a2.Slice(), s2)
		gtest.Assert(a2.Interfaces(), s2)

		var a3 garray.SortedStrArray
		err := json.Unmarshal(b2, &a3)
		gtest.Assert(err, nil)
		gtest.Assert(a3.Slice(), s1)
		gtest.Assert(a3.Interfaces(), s1)
	})

	gtest.Case(t, func() {
		type User struct {
			Name   string
			Scores *garray.SortedStrArray
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []string{"A+", "A", "A"},
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		user := new(User)
		err = json.Unmarshal(b, user)
		gtest.Assert(err, nil)
		gtest.Assert(user.Name, data["Name"])
		gtest.Assert(user.Scores, []string{"A", "A", "A+"})
	})
}

func TestSortedStrArray_Iterator(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedStrArrayFrom(slice)
	gtest.Case(t, func() {
		array.Iterator(func(k int, v string) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		array.IteratorAsc(func(k int, v string) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		array.IteratorDesc(func(k int, v string) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		index := 0
		array.Iterator(func(k int, v string) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
	gtest.Case(t, func() {
		index := 0
		array.IteratorAsc(func(k int, v string) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
	gtest.Case(t, func() {
		index := 0
		array.IteratorDesc(func(k int, v string) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
}

func TestSortedStrArray_RemoveValue(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedStrArrayFrom(slice)
	gtest.Case(t, func() {
		gtest.Assert(array.RemoveValue("e"), false)
		gtest.Assert(array.RemoveValue("b"), true)
		gtest.Assert(array.RemoveValue("a"), true)
		gtest.Assert(array.RemoveValue("c"), true)
		gtest.Assert(array.RemoveValue("f"), false)
	})
}

func TestSortedStrArray_UnmarshalValue(t *testing.T) {
	type T struct {
		Name  string
		Array *garray.SortedStrArray
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`["1","3","2"]`),
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Array.Slice(), g.SliceStr{"1", "2", "3"})
	})
	// Map
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": g.SliceStr{"1", "3", "2"},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Array.Slice(), g.SliceStr{"1", "2", "3"})
	})
}

func TestSortedStrArray_FilterEmpty(t *testing.T) {
	gtest.Case(t, func() {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"", "1", "2", "0"})
		gtest.Assert(array.FilterEmpty(), g.SliceStr{"0", "1", "2"})
	})
	gtest.Case(t, func() {
		array := garray.NewSortedStrArrayFrom(g.SliceStr{"1", "2"})
		gtest.Assert(array.FilterEmpty(), g.SliceStr{"1", "2"})
	})
}
