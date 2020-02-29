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
	"github.com/gogf/gf/util/gutil"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func TestSortedArray_NewSortedArrayFrom(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		a2 := []interface{}{"h", "j", "i", "k"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 interface{}) int {
			return -1
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a2, func2)

		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "f"})

		gtest.Assert(array2.Len(), 4)
		gtest.Assert(array2, []interface{}{"k", "i", "j", "h"})
	})
}

func TestNewSortedArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 interface{}) int {
			return -1
		}
		array1 := garray.NewSortedArrayFromCopy(a1, func1)
		array2 := garray.NewSortedArrayFromCopy(a1, func2)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "f"})
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array2, []interface{}{"c", "f", "a"})
	})
}

func TestSortedArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		a2 := []interface{}{"e", "h", "g", "k"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}

		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.SetArray(a2)
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1, []interface{}{"e", "g", "h", "k"})
	})

}

func TestSortedArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.Sort()
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "f"})
	})

}

func TestSortedArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		gtest.Assert(array1.Get(2), "f")
		gtest.Assert(array1.Get(1), "c")
	})

}

func TestSortedArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Remove(1)
		gtest.Assert(gconv.String(i1), "b")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("b"), false)

		gtest.Assert(array1.Remove(-1), nil)
		gtest.Assert(array1.Remove(100000), nil)

		i2 := array1.Remove(0)
		gtest.Assert(gconv.String(i2), "a")
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(array1.Contains("a"), false)

		i3 := array1.Remove(1)
		gtest.Assert(gconv.String(i3), "d")
		gtest.Assert(array1.Len(), 1)
		gtest.Assert(array1.Contains("d"), false)
	})

}

func TestSortedArray_PopLeft(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopLeft()
		gtest.Assert(gconv.String(i1), "a")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"b", "c", "d"})
	})

}

func TestSortedArray_PopRight(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRight()
		gtest.Assert(gconv.String(i1), "d")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "b", "c"})
	})

}

func TestSortedArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRand()
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 3)

	})
}

func TestSortedArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRands(2)
		gtest.Assert(len(i1), 2)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 2)

		i2 := array1.PopRands(3)
		gtest.Assert(len(i1), 2)
		gtest.AssertIN(i2, []interface{}{"a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_PopLefts(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopLefts(2)
		gtest.Assert(len(i1), 2)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b", "e", "f"})
		gtest.Assert(array1.Len(), 4)

		i2 := array1.PopLefts(5)
		gtest.Assert(len(i2), 4)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b", "e", "f"})
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_PopRights(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRights(2)
		gtest.Assert(len(i1), 2)
		gtest.Assert(i1, []interface{}{"e", "f"})
		gtest.Assert(array1.Len(), 4)

		i2 := array1.PopRights(10)
		gtest.Assert(len(i2), 4)

	})
}

func TestSortedArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a1, func1, true)
		i1 := array1.Range(2, 5)
		gtest.Assert(i1, []interface{}{"c", "d", "e"})
		gtest.Assert(array1.Len(), 6)

		i2 := array1.Range(7, 5)
		gtest.Assert(len(i2), 0)
		i2 = array1.Range(-1, 2)
		gtest.Assert(i2, []interface{}{"a", "b"})

		i2 = array1.Range(4, 10)
		gtest.Assert(len(i2), 2)
		gtest.Assert(i2, []interface{}{"e", "f"})

		gtest.Assert(array2.Range(1, 3), []interface{}{"b", "c"})

	})
}

func TestSortedArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		a2 := []interface{}{"1", "2", "3", "b", "e", "f"}
		a3 := []interface{}{"4", "5", "6"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a2, func1)
		array3 := garray.NewSortedArrayFrom(a3, func1)
		gtest.Assert(array1.Sum(), 0)
		gtest.Assert(array2.Sum(), 6)
		gtest.Assert(array3.Sum(), 15)

	})
}

func TestSortedArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := array1.Clone()
		gtest.Assert(array1, array2)
		array1.Remove(1)
		gtest.AssertNE(array1, array2)

	})
}

func TestSortedArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		gtest.Assert(array1.Len(), 6)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Chunk(2)
		gtest.Assert(len(i1), 3)
		gtest.Assert(i1[0], []interface{}{"a", "b"})
		gtest.Assert(i1[2], []interface{}{"e"})

		i1 = array1.Chunk(0)
		gtest.Assert(len(i1), 0)
	})
	gtest.Case(t, func() {
		a1 := []interface{}{1, 2, 3, 4, 5}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		chunks := array1.Chunk(3)
		gtest.Assert(len(chunks), 2)
		gtest.Assert(chunks[0], []interface{}{1, 2, 3})
		gtest.Assert(chunks[1], []interface{}{4, 5})
		gtest.Assert(array1.Chunk(0), nil)
	})
	gtest.Case(t, func() {
		a1 := []interface{}{1, 2, 3, 4, 5, 6}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		chunks := array1.Chunk(2)
		gtest.Assert(len(chunks), 3)
		gtest.Assert(chunks[0], []interface{}{1, 2})
		gtest.Assert(chunks[1], []interface{}{3, 4})
		gtest.Assert(chunks[2], []interface{}{5, 6})
		gtest.Assert(array1.Chunk(0), nil)
	})
	gtest.Case(t, func() {
		a1 := []interface{}{1, 2, 3, 4, 5, 6}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		chunks := array1.Chunk(3)
		gtest.Assert(len(chunks), 2)
		gtest.Assert(chunks[0], []interface{}{1, 2, 3})
		gtest.Assert(chunks[1], []interface{}{4, 5, 6})
		gtest.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a1, func1, true)
		i1 := array1.SubSlice(2, 3)
		gtest.Assert(len(i1), 3)
		gtest.Assert(i1, []interface{}{"c", "d", "e"})

		i1 = array1.SubSlice(2, 6)
		gtest.Assert(len(i1), 3)
		gtest.Assert(i1, []interface{}{"c", "d", "e"})

		i1 = array1.SubSlice(7, 2)
		gtest.Assert(len(i1), 0)

		s1 := array1.SubSlice(1, -2)
		gtest.Assert(s1, nil)

		s1 = array1.SubSlice(-9, 2)
		gtest.Assert(s1, nil)
		gtest.Assert(array2.SubSlice(1, 3), []interface{}{"b", "c", "d"})

	})
}

func TestSortedArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Rand()
		gtest.AssertIN(i1, []interface{}{"a", "d", "c"})
		gtest.Assert(array1.Len(), 3)
	})
}

func TestSortedArray_Rands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Rands(2)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c"})
		gtest.Assert(len(i1), 2)
		gtest.Assert(array1.Len(), 3)

		i1 = array1.Rands(4)
		gtest.Assert(len(i1), 3)
	})
}

func TestSortedArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		gtest.Assert(array1.Join(","), `a,c,d`)
		gtest.Assert(array1.Join("."), `a.c.d`)
	})

	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, `"a"`, `\a`}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorString)
		gtest.Assert(array1.Join("."), `"a".0.1.\a`)
	})
}

func TestSortedArray_String(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, "a", "b"}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorString)
		gtest.Assert(array1.String(), `[0,1,"a","b"]`)
	})
}

func TestSortedArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		m1 := array1.CountValues()
		gtest.Assert(len(m1), 3)
		gtest.Assert(m1["c"], 2)
		gtest.Assert(m1["a"], 1)

	})
}

func TestSortedArray_SetUnique(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.SetUnique(true)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "d"})
	})
}

func TestSortedArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, func1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []interface{}) { //读写锁
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

func TestSortedArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, func1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.RLockFunc(func(n1 []interface{}) { //读写锁
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
		gtest.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候不会被阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		s1 := []interface{}{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStrArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStrArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})

		a1 := garray.NewSortedArrayFrom(s1, func1)

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i1).Len(), 9)
		gtest.Assert(a1.Merge(i2).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedArray_Json(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []interface{}{"a", "b", "d", "c"}
		s2 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, gutil.ComparatorString)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		gtest.Assert(b1, b2)
		gtest.Assert(err1, err2)

		a2 := garray.NewSortedArray(gutil.ComparatorString)
		err1 = json.Unmarshal(b2, &a2)
		gtest.Assert(a2.Slice(), s2)

		var a3 garray.SortedArray
		err := json.Unmarshal(b2, &a3)
		gtest.Assert(err, nil)
		gtest.Assert(a3.Slice(), s1)
		gtest.Assert(a3.Interfaces(), s1)
	})

	gtest.Case(t, func() {
		type User struct {
			Name   string
			Scores *garray.SortedArray
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
		gtest.AssertNE(user.Scores, nil)
		gtest.Assert(user.Scores.Len(), 3)
		gtest.AssertIN(user.Scores.PopLeft(), data["Scores"])
		gtest.AssertIN(user.Scores.PopLeft(), data["Scores"])
		gtest.AssertIN(user.Scores.PopLeft(), data["Scores"])
	})
}

func TestSortedArray_Iterator(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := garray.NewSortedArrayFrom(slice, gutil.ComparatorString)
	gtest.Case(t, func() {
		array.Iterator(func(k int, v interface{}) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		array.IteratorAsc(func(k int, v interface{}) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		array.IteratorDesc(func(k int, v interface{}) bool {
			gtest.Assert(v, slice[k])
			return true
		})
	})
	gtest.Case(t, func() {
		index := 0
		array.Iterator(func(k int, v interface{}) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
	gtest.Case(t, func() {
		index := 0
		array.IteratorAsc(func(k int, v interface{}) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
	gtest.Case(t, func() {
		index := 0
		array.IteratorDesc(func(k int, v interface{}) bool {
			index++
			return false
		})
		gtest.Assert(index, 1)
	})
}

func TestSortedArray_RemoveValue(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := garray.NewSortedArrayFrom(slice, gutil.ComparatorString)
	gtest.Case(t, func() {
		gtest.Assert(array.RemoveValue("e"), false)
		gtest.Assert(array.RemoveValue("b"), true)
		gtest.Assert(array.RemoveValue("a"), true)
		gtest.Assert(array.RemoveValue("c"), true)
		gtest.Assert(array.RemoveValue("f"), false)
	})
}

func TestSortedArray_UnmarshalValue(t *testing.T) {
	type T struct {
		Name  string
		Array *garray.SortedArray
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

func TestSortedArray_FilterNil(t *testing.T) {
	gtest.Case(t, func() {
		values := g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}
		array := garray.NewSortedArrayFromCopy(values, gutil.ComparatorInt)
		gtest.Assert(array.FilterNil().Slice(), g.Slice{0, "", g.Slice{}, 1, 2, 3, 4})
	})
	gtest.Case(t, func() {
		array := garray.NewSortedArrayFromCopy(g.Slice{nil, 1, 2, 3, 4, nil}, gutil.ComparatorInt)
		gtest.Assert(array.FilterNil(), g.Slice{1, 2, 3, 4})
	})
}

func TestSortedArray_FilterEmpty(t *testing.T) {
	gtest.Case(t, func() {
		array := garray.NewSortedArrayFrom(g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}, gutil.ComparatorInt)
		gtest.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
	gtest.Case(t, func() {
		array := garray.NewSortedArrayFrom(g.Slice{1, 2, 3, 4}, gutil.ComparatorInt)
		gtest.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
}
