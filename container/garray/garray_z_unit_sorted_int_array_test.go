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

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/test/gtest"
)

func TestNewSortedIntArrayFrom(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 3, 2, 1, 4, 5, 6}
		array1 := garray.NewSortedIntArrayFrom(a1, true)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
		gtest.Assert(array1.Slice(), a1)
		gtest.Assert(array1.Interfaces(), a1)
	})
}

func TestNewSortedIntArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 5, 2, 1, 4, 3, 6}
		array1 := garray.NewSortedIntArrayFromCopy(a1, false)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestSortedIntArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3}
		a2 := []int{4, 5, 6}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := array1.SetArray(a2)

		gtest.Assert(array2.Len(), 3)
		gtest.Assert(array2.Search(3), -1)
		gtest.Assert(array2.Search(5), 1)
		gtest.Assert(array2.Search(6), 2)
	})
}

func TestSortedIntArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 3, 2, 1}

		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := array1.Sort()

		gtest.Assert(array2.Len(), 4)
		gtest.Assert(array2, []int{0, 1, 2, 3})
	})
}

func TestSortedIntArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 0}
		array1 := garray.NewSortedIntArrayFrom(a1)
		gtest.Assert(array1.Get(0), 0)
		gtest.Assert(array1.Get(1), 1)
		gtest.Assert(array1.Get(3), 5)
	})
}

func TestSortedIntArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 0}
		array1 := garray.NewSortedIntArrayFrom(a1)

		gtest.Assert(array1.Remove(-1), 0)
		gtest.Assert(array1.Remove(100000), 0)

		i1 := array1.Remove(2)
		gtest.Assert(i1, 3)
		gtest.Assert(array1.Search(5), 2)

		// 再次删除剩下的数组中的第一个
		i2 := array1.Remove(0)
		gtest.Assert(i2, 0)
		gtest.Assert(array1.Search(5), 1)

		a2 := []int{1, 3, 4}
		array2 := garray.NewSortedIntArrayFrom(a2)
		i3 := array2.Remove(1)
		gtest.Assert(array2.Search(1), 0)
		gtest.Assert(i3, 3)
		i3 = array2.Remove(1)
		gtest.Assert(array2.Search(4), -1)
		gtest.Assert(i3, 4)
	})
}

func TestSortedIntArray_PopLeft(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.PopLeft()
		gtest.Assert(i1, 1)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Search(1), -1)
	})
}

func TestSortedIntArray_PopRight(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.PopRight()
		gtest.Assert(i1, 5)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Search(5), -1)
	})
}

func TestSortedIntArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.PopRand()
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Search(i1), -1)
		gtest.AssertIN(i1, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopRands(2)
		gtest.Assert(array1.Len(), 2)
		gtest.AssertIN(ns1, []int{1, 3, 5, 2})

		a2 := []int{1, 3, 5, 2}
		array2 := garray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopRands(5)
		gtest.Assert(array2.Len(), 0)
		gtest.Assert(len(ns2), 4)
		gtest.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_PopLefts(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopLefts(2)
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(ns1, []int{1, 2})

		a2 := []int{1, 3, 5, 2}
		array2 := garray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopLefts(5)
		gtest.Assert(array2.Len(), 0)
		gtest.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_PopRights(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopRights(2)
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(ns1, []int{3, 5})

		a2 := []int{1, 3, 5, 2}
		array2 := garray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopRights(5)
		gtest.Assert(array2.Len(), 0)
		gtest.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2, 6, 7}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := garray.NewSortedIntArrayFrom(a1, true)
		ns1 := array1.Range(1, 4)
		gtest.Assert(len(ns1), 3)
		gtest.Assert(ns1, []int{2, 3, 5})

		ns2 := array1.Range(5, 4)
		gtest.Assert(len(ns2), 0)

		ns3 := array1.Range(-1, 4)
		gtest.Assert(len(ns3), 4)

		nsl := array1.Range(5, 8)
		gtest.Assert(len(nsl), 1)
		gtest.Assert(array2.Range(1, 2), []int{2})
	})
}

func TestSortedIntArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		n1 := array1.Sum()
		gtest.Assert(n1, 9)
	})
}

func TestSortedIntArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		gtest.Assert(array1.Join("."), `1.3.5`)
	})
}

func TestSortedIntArray_String(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		gtest.Assert(array1.String(), `[1,3,5]`)
	})
}

func TestSortedIntArray_Contains(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		gtest.Assert(array1.Contains(4), false)
	})
}

func TestSortedIntArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array2.Len(), 3)
		gtest.Assert(array2, array1)
	})
}

func TestSortedIntArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)
	})
}

func TestSortedIntArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Chunk(2) //按每几个元素切成一个数组
		ns2 := array1.Chunk(-1)
		gtest.Assert(len(ns1), 3)
		gtest.Assert(ns1[0], []int{1, 2})
		gtest.Assert(ns1[2], []int{5})
		gtest.Assert(len(ns2), 0)
	})
}

func TestSortedIntArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := garray.NewSortedIntArrayFrom(a1, true)
		ns1 := array1.SubSlice(1, 2)
		gtest.Assert(len(ns1), 2)
		gtest.Assert(ns1, []int{2, 3})

		ns2 := array1.SubSlice(7, 2)
		gtest.Assert(len(ns2), 0)

		ns3 := array1.SubSlice(3, 5)
		gtest.Assert(len(ns3), 2)
		gtest.Assert(ns3, []int{4, 5})

		ns4 := array1.SubSlice(3, 1)
		gtest.Assert(len(ns4), 1)
		gtest.Assert(ns4, []int{4})
		gtest.Assert(array1.SubSlice(-1, 1), []int{5})
		gtest.Assert(array1.SubSlice(-9, 1), nil)
		gtest.Assert(array1.SubSlice(1, -9), nil)
		gtest.Assert(array2.SubSlice(1, 2), []int{2, 3})
	})
}

func TestSortedIntArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Rand() //按每几个元素切成一个数组
		gtest.AssertIN(ns1, a1)
	})
}

func TestSortedIntArray_Rands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Rands(2) //按每几个元素切成一个数组
		gtest.AssertIN(ns1, a1)
		gtest.Assert(len(ns1), 2)

		ns2 := array1.Rands(6) //按每几个元素切成一个数组
		gtest.AssertIN(ns2, a1)
		gtest.Assert(len(ns2), 5)
	})
}

func TestSortedIntArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5, 3}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.CountValues() //按每几个元素切成一个数组
		gtest.Assert(len(ns1), 5)
		gtest.Assert(ns1[2], 1)
		gtest.Assert(ns1[3], 2)
	})
}

func TestSortedIntArray_SetUnique(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5, 3}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array1.SetUnique(true)
		gtest.Assert(array1.Len(), 5)
		gtest.Assert(array1, []int{1, 2, 3, 4, 5})
	})
}

func TestSortedIntArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 2, 3, 4}
		a1 := garray.NewSortedIntArrayFrom(s1, true)
		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []int) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = 6
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
		gtest.Assert(a1.Contains(6), true)
	})
}

func TestSortedIntArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 2, 3, 4}
		a1 := garray.NewSortedIntArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []int) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = 6
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
		gtest.Assert(a1.Contains(6), true)
	})
}

func TestSortedIntArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}
		i0 := []int{1, 2, 3, 4}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStrArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStrArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := garray.NewSortedIntArrayFrom(i0)

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i1).Len(), 9)
		gtest.Assert(a1.Merge(i2).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedIntArray_Json(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 4, 3, 2}
		s2 := []int{1, 2, 3, 4}
		a1 := garray.NewSortedIntArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		gtest.Assert(b1, b2)
		gtest.Assert(err1, err2)

		a2 := garray.NewSortedIntArray()
		err1 = json.Unmarshal(b2, &a2)
		gtest.Assert(a2.Slice(), s2)

		var a3 garray.SortedIntArray
		err := json.Unmarshal(b2, &a3)
		gtest.Assert(err, nil)
		gtest.Assert(a3.Slice(), s1)
	})

	gtest.Case(t, func() {
		type User struct {
			Name   string
			Scores *garray.SortedIntArray
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []int{99, 100, 98},
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		user := new(User)
		err = json.Unmarshal(b, user)
		gtest.Assert(err, nil)
		gtest.Assert(user.Name, data["Name"])
		gtest.Assert(user.Scores, []int{98, 99, 100})
	})
}

func TestSortedIntArray_Iterator(t *testing.T) {
	slice := g.SliceInt{10, 20, 30, 40}
	array := garray.NewSortedIntArrayFrom(slice)
	gtest.Case(t, func() {
		array.Iterator(func(k int, v int) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		array.IteratorAsc(func(k int, v int) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		array.IteratorDesc(func(k int, v int) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		index := 0
		array.Iterator(func(k int, v int) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
	gtest.Case(t, func() {
		index := 0
		array.IteratorAsc(func(k int, v int) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
	gtest.Case(t, func() {
		index := 0
		array.IteratorDesc(func(k int, v int) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
}

func TestSortedIntArray_RemoveValue(t *testing.T) {
	slice := g.SliceInt{10, 20, 30, 40}
	array := garray.NewSortedIntArrayFrom(slice)
	gtest.Case(t, func() {
		gtest.Assert(array.RemoveValue(99), false)
		gtest.Assert(array.RemoveValue(20), true)
		gtest.Assert(array.RemoveValue(10), true)
		gtest.Assert(array.RemoveValue(20), false)
		gtest.Assert(array.RemoveValue(88), false)
		gtest.Assert(array.Len(), 2)
	})
}

func TestSortedIntArray_UnmarshalValue(t *testing.T) {
	type T struct {
		Name  string
		Array *garray.SortedIntArray
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`[2,3,1]`),
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Array.Slice(), g.Slice{1, 2, 3})
	})
	// Map
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": g.Slice{2, 3, 1},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Array.Slice(), g.Slice{1, 2, 3})
	})
}

func TestSortedIntArray_FilterEmpty(t *testing.T) {
	gtest.Case(t, func() {
		array := garray.NewSortedIntArrayFrom(g.SliceInt{0, 1, 2, 3, 4, 0})
		gtest.Assert(array.FilterEmpty(), g.SliceInt{1, 2, 3, 4})
	})
	gtest.Case(t, func() {
		array := garray.NewSortedIntArrayFrom(g.SliceInt{1, 2, 3, 4})
		gtest.Assert(array.FilterEmpty(), g.SliceInt{1, 2, 3, 4})
	})
}
