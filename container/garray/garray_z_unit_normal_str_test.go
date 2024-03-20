// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray_test

import (
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/internal/empty"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestStrArrayBasic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"0", "1", "2", "3"}
		array := garray.NewStrArrayFrom(expect)
		array2 := garray.NewStrArrayFrom(expect, true)
		array3 := garray.NewStrArrayFrom([]string{})
		t.Assert(array.Slice(), expect)
		t.Assert(array.Interfaces(), expect)
		array.Set(0, "100")

		v, ok := array.Get(0)
		t.Assert(v, 100)
		t.Assert(ok, true)

		v, ok = array3.Get(0)
		t.Assert(v, "")
		t.Assert(ok, false)

		t.Assert(array.Search("100"), 0)
		t.Assert(array.Contains("100"), true)

		v, ok = array.Remove(0)
		t.Assert(v, 100)
		t.Assert(ok, true)

		v, ok = array.Remove(-1)
		t.Assert(v, "")
		t.Assert(ok, false)

		v, ok = array.Remove(100000)
		t.Assert(v, "")
		t.Assert(ok, false)

		t.Assert(array.Contains("100"), false)
		array.Append("4")
		t.Assert(array.Len(), 4)
		array.InsertBefore(0, "100")
		array.InsertAfter(0, "200")
		t.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "4"})
		array.InsertBefore(5, "300")
		array.InsertAfter(6, "400")
		t.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "300", "4", "400"})
		t.Assert(array.Clear().Len(), 0)
		t.Assert(array2.Slice(), expect)
		t.Assert(array3.Search("100"), -1)
		err := array.InsertBefore(99, "300")
		t.AssertNE(err, nil)
		array.InsertAfter(99, "400")
		t.AssertNE(err, nil)

	})

	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom([]string{"0", "1", "2", "3"})

		copyArray := array.DeepCopy().(*garray.StrArray)
		copyArray.Set(0, "1")
		cval, _ := copyArray.Get(0)
		val, _ := array.Get(0)
		t.AssertNE(cval, val)
	})
}

func TestStrArrayContainsI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := garray.NewStrArray()
		t.Assert(s.Contains("A"), false)
		s.Append("a", "b", "C")
		t.Assert(s.Contains("A"), false)
		t.Assert(s.Contains("a"), true)
		t.Assert(s.ContainsI("A"), true)
	})
}

func TestStrArraySort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect1 := []string{"0", "1", "2", "3"}
		expect2 := []string{"3", "2", "1", "0"}
		array := garray.NewStrArray()
		for i := 3; i >= 0; i-- {
			array.Append(gconv.String(i))
		}
		array.Sort()
		t.Assert(array.Slice(), expect1)
		array.Sort(true)
		t.Assert(array.Slice(), expect2)
	})
}

func TestStrArrayUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"1", "1", "2", "2", "3", "3", "2", "2"}
		array := garray.NewStrArrayFrom(expect)
		t.Assert(array.Unique().Slice(), []string{"1", "2", "3"})
		array1 := garray.NewStrArrayFrom([]string{})
		t.Assert(array1.Unique().Slice(), []string{})
	})
}

func TestStrArrayPushAndPop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"0", "1", "2", "3"}
		array := garray.NewStrArrayFrom(expect)
		t.Assert(array.Slice(), expect)

		v, ok := array.PopLeft()
		t.Assert(v, "0")
		t.Assert(ok, true)

		v, ok = array.PopRight()
		t.Assert(v, "3")
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.AssertIN(v, []string{"1", "2"})
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.AssertIN(v, []string{"1", "2"})
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.Assert(v, "")
		t.Assert(ok, false)

		t.Assert(array.Len(), 0)
		array.PushLeft("1").PushRight("2")
		t.Assert(array.Slice(), []string{"1", "2"})
	})
}

func TestStrArrayPopLeft(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2", "3"})
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

func TestStrArrayPopRight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2", "3"})

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

func TestStrArrayPopLefts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2", "3"})
		t.Assert(array.PopLefts(2), g.Slice{"1", "2"})
		t.Assert(array.Len(), 1)
		t.Assert(array.PopLefts(2), g.Slice{"3"})
		t.Assert(array.Len(), 0)
	})
}

func TestStrArrayPopRights(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2", "3"})
		t.Assert(array.PopRights(2), g.Slice{"2", "3"})
		t.Assert(array.Len(), 1)
		t.Assert(array.PopLefts(2), g.Slice{"1"})
		t.Assert(array.Len(), 0)
	})
}

func TestStrArrayPopLeftsAndPopRights(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArray()
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

	gtest.C(t, func(t *gtest.T) {
		value1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		value2 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(value1)
		array2 := garray.NewStrArrayFrom(value2)
		t.Assert(array1.PopLefts(2), []interface{}{"0", "1"})
		t.Assert(array1.Slice(), []interface{}{"2", "3", "4", "5", "6"})
		t.Assert(array1.PopRights(2), []interface{}{"5", "6"})
		t.Assert(array1.Slice(), []interface{}{"2", "3", "4"})
		t.Assert(array1.PopRights(20), []interface{}{"2", "3", "4"})
		t.Assert(array1.Slice(), []interface{}{})
		t.Assert(array2.PopLefts(20), []interface{}{"0", "1", "2", "3", "4", "5", "6"})
		t.Assert(array2.Slice(), []interface{}{})
	})
}

func TestStringRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(value1)
		array2 := garray.NewStrArrayFrom(value1, true)
		t.Assert(array1.Range(0, 1), []interface{}{"0"})
		t.Assert(array1.Range(1, 2), []interface{}{"1"})
		t.Assert(array1.Range(0, 2), []interface{}{"0", "1"})
		t.Assert(array1.Range(-1, 10), value1)
		t.Assert(array1.Range(10, 1), nil)
		t.Assert(array2.Range(0, 1), []interface{}{"0"})
	})
}

func TestStrArrayMerge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a11 := []string{"0", "1", "2", "3"}
		a21 := []string{"4", "5", "6", "7"}
		array1 := garray.NewStrArrayFrom(a11)
		array2 := garray.NewStrArrayFrom(a21)
		t.Assert(array1.Merge(array2).Slice(), []string{"0", "1", "2", "3", "4", "5", "6", "7"})

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
		a1 := garray.NewStrArrayFrom(s1)

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i1).Len(), 9)
		t.Assert(a1.Merge(i2).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestStrArrayFill(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0"}
		a2 := []string{"0"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := garray.NewStrArrayFrom(a2)
		t.Assert(array1.Fill(1, 2, "100"), nil)
		t.Assert(array1.Slice(), []string{"0", "100", "100"})
		t.Assert(array2.Fill(0, 2, "100"), nil)
		t.Assert(array2.Slice(), []string{"100", "100"})
		t.AssertNE(array2.Fill(-1, 2, "100"), nil)
		t.Assert(array2.Len(), 2)
	})
}

func TestStrArrayChunk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5"}
		array1 := garray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []string{"1", "2"})
		t.Assert(chunks[1], []string{"3", "4"})
		t.Assert(chunks[2], []string{"5"})
		t.Assert(len(array1.Chunk(0)), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5"}
		array1 := garray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []string{"1", "2", "3"})
		t.Assert(chunks[1], []string{"4", "5"})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []string{"1", "2"})
		t.Assert(chunks[1], []string{"3", "4"})
		t.Assert(chunks[2], []string{"5", "6"})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []string{"1", "2", "3"})
		t.Assert(chunks[1], []string{"4", "5", "6"})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestStrArrayPad(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Pad(3, "1").Slice(), []string{"0", "1", "1"})
		t.Assert(array1.Pad(-4, "1").Slice(), []string{"1", "0", "1", "1"})
		t.Assert(array1.Pad(3, "1").Slice(), []string{"1", "0", "1", "1"})
	})
}

func TestStrArraySubSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := garray.NewStrArrayFrom(a1, true)
		t.Assert(array1.SubSlice(0, 2), []string{"0", "1"})
		t.Assert(array1.SubSlice(2, 2), []string{"2", "3"})
		t.Assert(array1.SubSlice(5, 8), []string{"5", "6"})
		t.Assert(array1.SubSlice(8, 2), nil)
		t.Assert(array1.SubSlice(1, -2), nil)
		t.Assert(array1.SubSlice(-5, 2), []string{"2", "3"})
		t.Assert(array1.SubSlice(-10, 1), nil)
		t.Assert(array2.SubSlice(0, 2), []string{"0", "1"})
	})
}

func TestStrArrayRand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(len(array1.Rands(2)), "2")
		t.Assert(len(array1.Rands(10)), 10)
		t.AssertIN(array1.Rands(1)[0], a1)
		v, ok := array1.Rand()
		t.Assert(ok, true)
		t.AssertIN(v, a1)

		array2 := garray.NewStrArrayFrom([]string{})
		v, ok = array2.Rand()
		t.Assert(ok, false)
		t.Assert(v, "")
		strArray := array2.Rands(1)
		t.AssertNil(strArray)
	})
}

func TestStrArrayPopRands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "b", "c", "d", "e", "f", "g"}
		array1 := garray.NewStrArrayFrom(a1)
		t.AssertIN(array1.PopRands(1), []string{"a", "b", "c", "d", "e", "f", "g"})
		t.AssertIN(array1.PopRands(1), []string{"a", "b", "c", "d", "e", "f", "g"})
		t.AssertNI(array1.PopRands(1), array1.Slice())
		t.AssertNI(array1.PopRands(1), array1.Slice())
		t.Assert(len(array1.PopRands(10)), 3)
	})
}

func TestStrArrayShuffle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestStrArrayReverse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Reverse().Slice(), []string{"6", "5", "4", "3", "2", "1", "0"})
	})
}

func TestStrArrayJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Join("."), `0.1.2.3.4.5.6`)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", `"a"`, `\a`}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Join("."), `0.1."a".\a`)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Join("."), "")
	})
}

func TestStrArrayString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.String(), `["0","1","2","3","4","5","6"]`)

		array1 = nil
		t.Assert(array1.String(), "")
	})
}

func TestNewStrArrayFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := garray.NewStrArrayFromCopy(a1)
		a3 := garray.NewStrArrayFromCopy(a1, true)
		t.Assert(a2.Contains("1"), true)
		t.Assert(a2.Len(), 7)
		t.Assert(a2, a3)
	})
}

func TestStrArraySetArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"a", "b", "c", "d"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Contains("2"), true)
		t.Assert(array1.Len(), 7)

		array1 = array1.SetArray(a2)
		t.Assert(array1.Contains("2"), false)
		t.Assert(array1.Contains("c"), true)
		t.Assert(array1.Len(), 4)
	})
}

func TestStrArrayReplace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"a", "b", "c", "d"}
		a3 := []string{"o", "p", "q", "x", "y", "z", "w", "r", "v"}
		array1 := garray.NewStrArrayFrom(a1)
		t.Assert(array1.Contains("2"), true)
		t.Assert(array1.Len(), 7)

		array1 = array1.Replace(a2)
		t.Assert(array1.Contains("2"), false)
		t.Assert(array1.Contains("c"), true)
		t.Assert(array1.Contains("5"), true)
		t.Assert(array1.Len(), 7)

		array1 = array1.Replace(a3)
		t.Assert(array1.Contains("2"), false)
		t.Assert(array1.Contains("c"), false)
		t.Assert(array1.Contains("5"), false)
		t.Assert(array1.Contains("p"), true)
		t.Assert(array1.Contains("r"), false)
		t.Assert(array1.Len(), 7)

	})
}

func TestStrArraySum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"0", "a", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := garray.NewStrArrayFrom(a2)
		t.Assert(array1.Sum(), 21)
		t.Assert(array2.Sum(), 18)
	})
}

func TestStrArrayPopRand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		str1, ok := array1.PopRand()
		t.Assert(strings.Contains("0,1,2,3,4,5,6", str1), true)
		t.Assert(array1.Len(), 6)
		t.Assert(ok, true)
	})
}

func TestStrArrayClone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := array1.Clone()
		t.Assert(array2, array1)
		t.Assert(array2.Len(), 7)
	})
}

func TestStrArrayCountValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "4", "6"}
		array1 := garray.NewStrArrayFrom(a1)

		m1 := array1.CountValues()
		t.Assert(len(m1), 6)
		t.Assert(m1["2"], 1)
		t.Assert(m1["4"], 2)
	})
}

func TestStrArrayRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewStrArrayFrom(a1)
		s1, ok := array1.Remove(1)
		t.Assert(s1, "a")
		t.Assert(ok, true)
		t.Assert(array1.Len(), 4)
		s1, ok = array1.Remove(3)
		t.Assert(s1, "c")
		t.Assert(ok, true)
		t.Assert(array1.Len(), 3)
	})
}

func TestStrArrayRLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewStrArrayFrom(s1, true)

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

func TestStrArraySortFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "d", "c", "b"}
		a1 := garray.NewStrArrayFrom(s1)
		func1 := func(v1, v2 string) bool {
			return v1 < v2
		}
		a11 := a1.SortFunc(func1)
		t.Assert(a11, []string{"a", "b", "c", "d"})
	})
}

func TestStrArrayLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewStrArrayFrom(s1, true)

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

func TestStrArrayJson(t *testing.T) {
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := garray.NewStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewStrArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s1)

		var a3 garray.StrArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := *garray.NewStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewStrArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s1)

		var a3 garray.StrArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
	})
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores *garray.StrArray
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
		t.Assert(user.Scores, data["Scores"])
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores garray.StrArray
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
		t.Assert(user.Scores, data["Scores"])
	})
}

func TestStrArrayIterator(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewStrArrayFrom(slice)
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

func TestStrArrayRemoveValue(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewStrArrayFrom(slice)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestStrArrayRemoveValues(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewStrArrayFrom(slice)
	gtest.C(t, func(t *gtest.T) {
		array.RemoveValues("a", "b", "c")
		t.Assert(array.Slice(), g.SliceStr{"d"})
	})
}

func TestStrArrayUnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *garray.StrArray
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`["1","2","3"]`),
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
			"array": g.SliceStr{"1", "2", "3"},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.SliceStr{"1", "2", "3"})
	})
}
func TestStrArrayFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"", "1", "2", "0"})
		t.Assert(array.Filter(func(index int, value string) bool {
			return empty.IsEmpty(value)
		}), g.SliceStr{"1", "2", "0"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2"})
		t.Assert(array.Filter(func(index int, value string) bool {
			return empty.IsEmpty(value)
		}), g.SliceStr{"1", "2"})
	})
}

func TestStrArrayFilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"", "1", "2", "0"})
		t.Assert(array.FilterEmpty(), g.SliceStr{"1", "2", "0"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2"})
		t.Assert(array.FilterEmpty(), g.SliceStr{"1", "2"})
	})
}

func TestStrArrayWalk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewStrArrayFrom(g.SliceStr{"1", "2"})
		t.Assert(array.Walk(func(value string) string {
			return "key-" + value
		}), g.Slice{"key-1", "key-2"})
	})
}
