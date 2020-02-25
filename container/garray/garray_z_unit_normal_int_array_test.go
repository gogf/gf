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

func Test_IntArray_Basic(t *testing.T) {
	gtest.Case(t, func() {
		expect := []int{0, 1, 2, 3}
		expect2 := []int{}
		array := garray.NewIntArrayFrom(expect)
		array2 := garray.NewIntArrayFrom(expect2)
		gtest.Assert(array.Slice(), expect)
		gtest.Assert(array.Interfaces(), expect)
		array.Set(0, 100)
		gtest.Assert(array.Get(0), 100)
		gtest.Assert(array.Get(1), 1)
		gtest.Assert(array.Search(100), 0)
		gtest.Assert(array2.Search(100), -1)
		gtest.Assert(array.Contains(100), true)
		gtest.Assert(array.Remove(0), 100)
		gtest.Assert(array.Remove(-1), 0)
		gtest.Assert(array.Remove(100000), 0)
		gtest.Assert(array.Contains(100), false)
		array.Append(4)
		gtest.Assert(array.Len(), 4)
		array.InsertBefore(0, 100)
		array.InsertAfter(0, 200)
		gtest.Assert(array.Slice(), []int{100, 200, 1, 2, 3, 4})
		array.InsertBefore(5, 300)
		array.InsertAfter(6, 400)
		gtest.Assert(array.Slice(), []int{100, 200, 1, 2, 3, 300, 4, 400})
		gtest.Assert(array.Clear().Len(), 0)
	})
}

func TestIntArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		expect1 := []int{0, 1, 2, 3}
		expect2 := []int{3, 2, 1, 0}
		array := garray.NewIntArray()
		array2 := garray.NewIntArray(true)
		for i := 3; i >= 0; i-- {
			array.Append(i)
			array2.Append(i)
		}
		array.Sort()
		gtest.Assert(array.Slice(), expect1)
		array.Sort(true)
		gtest.Assert(array.Slice(), expect2)
		gtest.Assert(array2.Slice(), expect2)
	})
}

func TestIntArray_Unique(t *testing.T) {
	gtest.Case(t, func() {
		expect := []int{1, 1, 2, 3}
		array := garray.NewIntArrayFrom(expect)
		gtest.Assert(array.Unique().Slice(), []int{1, 2, 3})
	})
}

func TestIntArray_PushAndPop(t *testing.T) {
	gtest.Case(t, func() {
		expect := []int{0, 1, 2, 3}
		array := garray.NewIntArrayFrom(expect)
		gtest.Assert(array.Slice(), expect)
		gtest.Assert(array.PopLeft(), 0)
		gtest.Assert(array.PopRight(), 3)
		gtest.AssertIN(array.PopRand(), []int{1, 2})
		gtest.AssertIN(array.PopRand(), []int{1, 2})
		gtest.Assert(array.Len(), 0)
		array.PushLeft(1).PushRight(2)
		gtest.Assert(array.Slice(), []int{1, 2})
	})
}

func TestIntArray_PopLeftsAndPopRights(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []int{0, 1, 2, 3, 4, 5, 6}
		value2 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(value1)
		array2 := garray.NewIntArrayFrom(value2)
		gtest.Assert(array1.PopLefts(2), []int{0, 1})
		gtest.Assert(array1.Slice(), []int{2, 3, 4, 5, 6})
		gtest.Assert(array1.PopRights(2), []int{5, 6})
		gtest.Assert(array1.Slice(), []int{2, 3, 4})
		gtest.Assert(array1.PopRights(20), []int{2, 3, 4})
		gtest.Assert(array1.Slice(), []int{})
		gtest.Assert(array2.PopLefts(20), []int{0, 1, 2, 3, 4, 5, 6})
		gtest.Assert(array2.Slice(), []int{})
	})
}

func TestIntArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(value1)
		array2 := garray.NewIntArrayFrom(value1, true)
		gtest.Assert(array1.Range(0, 1), []int{0})
		gtest.Assert(array1.Range(1, 2), []int{1})
		gtest.Assert(array1.Range(0, 2), []int{0, 1})
		gtest.Assert(array1.Range(10, 2), nil)
		gtest.Assert(array1.Range(-1, 10), value1)
		gtest.Assert(array2.Range(1, 2), []int{1})
	})
}

func TestIntArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		n1 := []int{0, 1, 2, 3}
		n2 := []int{4, 5, 6, 7}
		i1 := []interface{}{"1", "2"}
		s1 := []string{"a", "b", "c"}
		s2 := []string{"e", "f"}
		a1 := garray.NewIntArrayFrom(n1)
		a2 := garray.NewIntArrayFrom(n2)
		a3 := garray.NewArrayFrom(i1)
		a4 := garray.NewStrArrayFrom(s1)

		a5 := garray.NewSortedStrArrayFrom(s2)
		a6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})

		a7 := garray.NewSortedStrArrayFrom(s1)
		a8 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)

		gtest.Assert(a1.Merge(a2).Slice(), []int{0, 1, 2, 3, 4, 5, 6, 7})
		gtest.Assert(a1.Merge(a3).Len(), 10)
		gtest.Assert(a1.Merge(a4).Len(), 13)
		gtest.Assert(a1.Merge(a5).Len(), 15)
		gtest.Assert(a1.Merge(a6).Len(), 18)
		gtest.Assert(a1.Merge(a7).Len(), 21)
		gtest.Assert(a1.Merge(a8).Len(), 23)
	})
}

func TestIntArray_Fill(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0}
		a2 := []int{0}
		array1 := garray.NewIntArrayFrom(a1)
		array2 := garray.NewIntArrayFrom(a2)
		gtest.Assert(array1.Fill(1, 2, 100).Slice(), []int{0, 100, 100})
		gtest.Assert(array2.Fill(0, 2, 100).Slice(), []int{100, 100})
		gtest.Assert(array2.Fill(-1, 2, 100).Slice(), []int{100, 100})
	})
}

func TestIntArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewIntArrayFrom(a1)
		chunks := array1.Chunk(2)
		gtest.Assert(len(chunks), 3)
		gtest.Assert(chunks[0], []int{1, 2})
		gtest.Assert(chunks[1], []int{3, 4})
		gtest.Assert(chunks[2], []int{5})
		gtest.Assert(array1.Chunk(0), nil)
	})
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewIntArrayFrom(a1)
		chunks := array1.Chunk(3)
		gtest.Assert(len(chunks), 2)
		gtest.Assert(chunks[0], []int{1, 2, 3})
		gtest.Assert(chunks[1], []int{4, 5})
		gtest.Assert(array1.Chunk(0), nil)
	})
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		chunks := array1.Chunk(2)
		gtest.Assert(len(chunks), 3)
		gtest.Assert(chunks[0], []int{1, 2})
		gtest.Assert(chunks[1], []int{3, 4})
		gtest.Assert(chunks[2], []int{5, 6})
		gtest.Assert(array1.Chunk(0), nil)
	})
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		chunks := array1.Chunk(3)
		gtest.Assert(len(chunks), 2)
		gtest.Assert(chunks[0], []int{1, 2, 3})
		gtest.Assert(chunks[1], []int{4, 5, 6})
		gtest.Assert(array1.Chunk(0), nil)
	})
}

func TestIntArray_Pad(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Pad(3, 1).Slice(), []int{0, 1, 1})
		gtest.Assert(array1.Pad(-4, 1).Slice(), []int{1, 0, 1, 1})
		gtest.Assert(array1.Pad(3, 1).Slice(), []int{1, 0, 1, 1})
	})
}

func TestIntArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		array2 := garray.NewIntArrayFrom(a1, true)
		gtest.Assert(array1.SubSlice(6), []int{6})
		gtest.Assert(array1.SubSlice(5), []int{5, 6})
		gtest.Assert(array1.SubSlice(8), nil)
		gtest.Assert(array1.SubSlice(0, 2), []int{0, 1})
		gtest.Assert(array1.SubSlice(2, 2), []int{2, 3})
		gtest.Assert(array1.SubSlice(5, 8), []int{5, 6})
		gtest.Assert(array1.SubSlice(-1, 1), []int{6})
		gtest.Assert(array1.SubSlice(-1, 9), []int{6})
		gtest.Assert(array1.SubSlice(-2, 3), []int{5, 6})
		gtest.Assert(array1.SubSlice(-7, 3), []int{0, 1, 2})
		gtest.Assert(array1.SubSlice(-8, 3), nil)
		gtest.Assert(array1.SubSlice(-1, -3), []int{3, 4, 5})
		gtest.Assert(array1.SubSlice(-9, 3), nil)
		gtest.Assert(array1.SubSlice(1, -1), []int{0})
		gtest.Assert(array1.SubSlice(1, -3), nil)
		gtest.Assert(array2.SubSlice(0, 2), []int{0, 1})
	})
}

func TestIntArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(len(array1.Rands(2)), 2)
		gtest.Assert(len(array1.Rands(10)), 7)
		gtest.AssertIN(array1.Rands(1)[0], a1)
		gtest.AssertIN(array1.Rand(), a1)
	})
}

func TestIntArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{100, 200, 300, 400, 500, 600}
		array := garray.NewIntArrayFrom(a1)
		ns1 := array.PopRands(2)
		gtest.AssertIN(ns1, []int{100, 200, 300, 400, 500, 600})
		gtest.Assert(len(ns1), 2)

		ns2 := array.PopRands(7)
		gtest.Assert(len(ns2), 4)
		gtest.AssertIN(ns2, []int{100, 200, 300, 400, 500, 600})
	})
}

func TestIntArray_Shuffle(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestIntArray_Reverse(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Reverse().Slice(), []int{6, 5, 4, 3, 2, 1, 0})
	})
}

func TestIntArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestIntArray_String(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.String(), "[0,1,2,3,4,5,6]")
	})
}

func TestIntArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		a2 := []int{6, 7}
		array1 := garray.NewIntArrayFrom(a1)
		array1.SetArray(a2)
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(array1, []int{6, 7})
	})
}

func TestIntArray_Replace(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		a2 := []int{6, 7}
		a3 := []int{9, 10, 11, 12, 13}
		array1 := garray.NewIntArrayFrom(a1)
		array1.Replace(a2)
		gtest.Assert(array1, []int{6, 7, 3, 5})

		array1.Replace(a3)
		gtest.Assert(array1, []int{9, 10, 11, 12})
	})
}

func TestIntArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)
	})
}

func TestIntArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array1, array2)
	})
}

func TestArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Get(2), 3)
		gtest.Assert(array1.Len(), 4)
	})
}

func TestIntArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Sum(), 11)
	})
}

func TestIntArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5, 3}
		array1 := garray.NewIntArrayFrom(a1)
		m1 := array1.CountValues()
		gtest.Assert(len(m1), 4)
		gtest.Assert(m1[1], 1)
		gtest.Assert(m1[3], 2)
	})
}

func TestNewIntArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5, 3}
		array1 := garray.NewIntArrayFromCopy(a1)
		gtest.Assert(array1.Len(), 5)
		gtest.Assert(array1, a1)
	})
}

func TestIntArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5, 4}
		array1 := garray.NewIntArrayFrom(a1)
		n1 := array1.Remove(1)
		gtest.Assert(n1, 2)
		gtest.Assert(array1.Len(), 4)

		n1 = array1.Remove(0)
		gtest.Assert(n1, 1)
		gtest.Assert(array1.Len(), 3)

		n1 = array1.Remove(2)
		gtest.Assert(n1, 4)
		gtest.Assert(array1.Len(), 2)
	})
}

func TestIntArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 2, 3, 4}
		a1 := garray.NewIntArrayFrom(s1, true)

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

func TestIntArray_SortFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 4, 3, 2}
		a1 := garray.NewIntArrayFrom(s1)
		func1 := func(v1, v2 int) bool {
			return v1 < v2
		}
		a11 := a1.SortFunc(func1)
		gtest.Assert(a11, []int{1, 2, 3, 4})

	})
}

func TestIntArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 2, 3, 4}
		a1 := garray.NewIntArrayFrom(s1, true)

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

func TestIntArray_Json(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []int{1, 4, 3, 2}
		a1 := garray.NewIntArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		gtest.Assert(b1, b2)
		gtest.Assert(err1, err2)

		a2 := garray.NewIntArray()
		err1 = json.Unmarshal(b2, &a2)
		gtest.Assert(a2.Slice(), s1)

		var a3 garray.IntArray
		err := json.Unmarshal(b2, &a3)
		gtest.Assert(err, nil)
		gtest.Assert(a3.Slice(), s1)
	})

	gtest.Case(t, func() {
		type User struct {
			Name   string
			Scores *garray.IntArray
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
		gtest.Assert(user.Scores, data["Scores"])
	})
}

func TestIntArray_Iterator(t *testing.T) {
	slice := g.SliceInt{10, 20, 30, 40}
	array := garray.NewIntArrayFrom(slice)
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

func TestIntArray_RemoveValue(t *testing.T) {
	slice := g.SliceInt{10, 20, 30, 40}
	array := garray.NewIntArrayFrom(slice)
	gtest.Case(t, func() {
		gtest.Assert(array.RemoveValue(99), false)
		gtest.Assert(array.RemoveValue(20), true)
		gtest.Assert(array.RemoveValue(10), true)
		gtest.Assert(array.RemoveValue(20), false)
		gtest.Assert(array.RemoveValue(88), false)
		gtest.Assert(array.Len(), 2)
	})
}

func TestIntArray_UnmarshalValue(t *testing.T) {
	type T struct {
		Name  string
		Array *garray.IntArray
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`[1,2,3]`),
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
			"array": g.Slice{1, 2, 3},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Array.Slice(), g.Slice{1, 2, 3})
	})
}

func TestIntArray_FilterEmpty(t *testing.T) {
	gtest.Case(t, func() {
		array := garray.NewIntArrayFrom(g.SliceInt{0, 1, 2, 3, 4, 0})
		gtest.Assert(array.FilterEmpty(), g.SliceInt{1, 2, 3, 4})
	})
	gtest.Case(t, func() {
		array := garray.NewIntArrayFrom(g.SliceInt{1, 2, 3, 4})
		gtest.Assert(array.FilterEmpty(), g.SliceInt{1, 2, 3, 4})
	})
}
