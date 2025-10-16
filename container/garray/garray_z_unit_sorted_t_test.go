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

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

func TestSortedTArray_NewSortedTArrayFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		a2 := []string{"h", "j", "i", "k"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(v1, v2)
		}
		func2 := func(v1, v2 string) int {
			return -1
		}
		array1 := garray.NewSortedTArrayFrom(a1, func1)
		array2 := garray.NewSortedTArrayFrom(a2, func2)

		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"a", "c", "f"})

		t.Assert(array2.Len(), 4)
		t.Assert(array2, []string{"k", "i", "j", "h"})
	})
}

func TestNewSortedTArrayFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 string) int {
			return -1
		}
		array1 := garray.NewSortedTArrayFromCopy(a1, func1)
		array2 := garray.NewSortedTArrayFromCopy(a1, func2)
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"a", "c", "f"})
		t.Assert(array1.Len(), 3)
		t.Assert(array2, []string{"c", "f", "a"})
	})
}

func TestSortedTArray_SetArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		a2 := []string{"e", "h", "g", "k"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(v1, v2)
		}

		array1 := garray.NewSortedTArrayFrom(a1, func1)
		array1.SetArray(a2)
		t.Assert(array1.Len(), 4)
		t.Assert(array1, []string{"e", "g", "h", "k"})
	})

}

func TestSortedTArray_Sort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		array1.Sort()
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []any{"a", "c", "f"})
	})

}

func TestSortedTArray_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		v, ok := array1.Get(2)
		t.Assert(v, "f")
		t.Assert(ok, true)

		v, ok = array1.Get(1)
		t.Assert(v, "c")
		t.Assert(ok, true)

		v, ok = array1.Get(99)
		t.Assert(v, nil)
		t.Assert(ok, false)
	})

}

func TestSortedTArray_At(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		v := array1.At(2)
		t.Assert(v, "f")
	})
}

func TestSortedTArray_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		i1, ok := array1.Remove(1)
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "b")
		t.Assert(array1.Len(), 3)
		t.Assert(array1.Contains("b"), false)

		v, ok := array1.Remove(-1)
		t.Assert(v, nil)
		t.Assert(ok, false)

		v, ok = array1.Remove(100000)
		t.Assert(v, nil)
		t.Assert(ok, false)

		i2, ok := array1.Remove(0)
		t.Assert(ok, true)
		t.Assert(gconv.String(i2), "a")
		t.Assert(array1.Len(), 2)
		t.Assert(array1.Contains("a"), false)

		i3, ok := array1.Remove(1)
		t.Assert(ok, true)
		t.Assert(gconv.String(i3), "d")
		t.Assert(array1.Len(), 1)
		t.Assert(array1.Contains("d"), false)
	})

}

func TestSortedTArray_PopLeft(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := garray.NewSortedTArrayFrom(
			[]string{"a", "d", "c", "b"},
			gutil.ComparatorT,
		)
		i1, ok := array1.PopLeft()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "a")
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []any{"b", "c", "d"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceInt{1, 2, 3}, gutil.ComparatorT)
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

func TestSortedTArray_PopRight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := garray.NewSortedTArrayFrom(
			[]string{"a", "d", "c", "b"},
			gutil.ComparatorT,
		)
		i1, ok := array1.PopRight()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "d")
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []any{"a", "b", "c"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceInt{1, 2, 3}, gutil.ComparatorT)
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

func TestSortedTArray_PopRand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		i1, ok := array1.PopRand()
		t.Assert(ok, true)
		t.AssertIN(i1, []string{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 3)

	})
}

func TestSortedTArray_PopRands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		i1 := array1.PopRands(2)
		t.Assert(len(i1), 2)
		t.AssertIN(i1, []string{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 2)

		i2 := array1.PopRands(3)
		t.Assert(len(i2), 2)
		t.AssertIN(i2, []string{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 0)

	})
}

func TestSortedTArray_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArray[int](gutil.ComparatorT)
		v, ok := array.PopLeft()
		t.Assert(v, 0)
		t.Assert(ok, false)
		t.Assert(array.PopLefts(10), nil)

		v, ok = array.PopRight()
		t.Assert(v, 0)
		t.Assert(ok, false)
		t.Assert(array.PopRights(10), nil)

		v, ok = array.PopRand()
		t.Assert(v, 0)
		t.Assert(ok, false)
		t.Assert(array.PopRands(10), nil)
	})
}

func TestSortedTArray_PopLefts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		i1 := array1.PopLefts(2)
		t.Assert(len(i1), 2)
		t.AssertIN(i1, []string{"a", "d", "c", "b", "e", "f"})
		t.Assert(array1.Len(), 4)

		i2 := array1.PopLefts(5)
		t.Assert(len(i2), 4)
		t.AssertIN(i2, []string{"a", "d", "c", "b", "e", "f"})
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedTArray_PopRights(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		i1 := array1.PopRights(2)
		t.Assert(len(i1), 2)
		t.Assert(i1, []string{"e", "f"})
		t.Assert(array1.Len(), 4)

		i2 := array1.PopRights(10)
		t.Assert(len(i2), 4)
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedTArray_Range(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		array2 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT, true)
		i1 := array1.Range(2, 5)
		t.Assert(i1, []string{"c", "d", "e"})
		t.Assert(array1.Len(), 6)

		i2 := array1.Range(7, 5)
		t.Assert(len(i2), 0)
		i2 = array1.Range(-1, 2)
		t.Assert(i2, []string{"a", "b"})

		i2 = array1.Range(4, 10)
		t.Assert(len(i2), 2)
		t.Assert(i2, []string{"e", "f"})

		t.Assert(array2.Range(1, 3), []string{"b", "c"})

	})
}

func TestSortedTArray_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		a2 := []string{"1", "2", "3", "b", "e", "f"}
		a3 := []string{"4", "5", "6"}
		array1 := garray.NewSortedTArrayFrom(a1, gutil.ComparatorT)
		array2 := garray.NewSortedTArrayFrom(a2, gutil.ComparatorT)
		array3 := garray.NewSortedTArrayFrom(a3, gutil.ComparatorT)
		t.Assert(array1.Sum(), 0)
		t.Assert(array2.Sum(), 6)
		t.Assert(array3.Sum(), 15)

	})
}

func TestSortedTArray_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		array2 := array1.Clone()
		t.Assert(array1, array2)
		array1.Remove(1)
		t.AssertNE(array1, array2)

	})
}

func TestSortedTArray_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		t.Assert(array1.Len(), 6)
		array1.Clear()
		t.Assert(array1.Len(), 0)

	})
}

func TestSortedTArray_Chunk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		i1 := array1.Chunk(2)
		t.Assert(len(i1), 3)
		t.Assert(i1[0], []any{"a", "b"})
		t.Assert(i1[2], []any{"e"})

		i1 = array1.Chunk(0)
		t.Assert(len(i1), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []int32{1, 2, 3, 4, 5}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []int32{1, 2, 3})
		t.Assert(chunks[1], []int32{4, 5})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []int{1, 2})
		t.Assert(chunks[1], []int{3, 4})
		t.Assert(chunks[2], []int{5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []int{1, 2, 3})
		t.Assert(chunks[1], []int{4, 5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedTArray_SubSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		array2 := garray.NewSortedTArrayFrom(a1, nil, true)
		i1 := array1.SubSlice(2, 3)
		t.Assert(len(i1), 3)
		t.Assert(i1, []string{"c", "d", "e"})

		i1 = array1.SubSlice(2, 6)
		t.Assert(len(i1), 3)
		t.Assert(i1, []string{"c", "d", "e"})

		i1 = array1.SubSlice(7, 2)
		t.Assert(len(i1), 0)

		s1 := array1.SubSlice(1, -2)
		t.Assert(s1, nil)

		s1 = array1.SubSlice(-9, 2)
		t.Assert(s1, nil)
		t.Assert(array2.SubSlice(1, 3), []string{"b", "c", "d"})

	})
}

func TestSortedTArray_Rand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		i1, ok := array1.Rand()
		t.Assert(ok, true)
		t.AssertIN(i1, []string{"a", "d", "c"})
		t.Assert(array1.Len(), 3)

		array2 := garray.NewSortedTArrayFrom([]string{}, nil)
		v, ok := array2.Rand()
		t.Assert(ok, false)
		t.Assert(v, nil)
	})
}

func TestSortedTArray_Rands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		i1 := array1.Rands(2)
		t.AssertIN(i1, []string{"a", "d", "c"})
		t.Assert(len(i1), 2)
		t.Assert(array1.Len(), 3)

		i1 = array1.Rands(4)
		t.Assert(len(i1), 4)

		array2 := garray.NewSortedTArrayFrom([]string{}, nil)
		v := array2.Rands(1)
		t.Assert(v, nil)
	})
}

func TestSortedTArray_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		t.Assert(array1.Join(","), `a,c,d`)
		t.Assert(array1.Join("."), `a.c.d`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", `"a"`, `\a`}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		t.Assert(array1.Join("."), `"a".0.1.\a`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []string{}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		t.Assert(array1.Join("."), "")
	})
}

func TestSortedTArray_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"0", "1", "a", "b"}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		t.Assert(array1.String(), `[0,1,"a","b"]`)

		array1 = nil
		t.Assert(array1.String(), "")
	})
}

func TestSortedTArray_CountValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "c"}

		array1 := garray.NewSortedTArrayFrom(a1, nil)
		m1 := array1.CountValues()
		t.Assert(len(m1), 3)
		t.Assert(m1["c"], 2)
		t.Assert(m1["a"], 1)

	})
}

func TestSortedTArray_SetUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		array1.SetUnique(true)
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})
	})
}

func TestSortedTArray_Unique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := garray.NewSortedTArrayFrom(a1, nil)
		array1.Unique()
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})

		array2 := garray.NewSortedTArrayFrom([]int{}, nil)
		array2.Unique()
		t.Assert(array2.Len(), 0)
		t.Assert(array2, []int{})
	})
}

func TestSortedTArray_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedTArrayFrom(s1, nil, true)

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

func TestSortedTArray_RLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedTArrayFrom(s1, nil, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		// go1
		go a1.RLockFunc(func(n1 []string) { // 读写锁
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
		t.AssertLT(t2-t1, 20) // go1加的读锁，所go2读的时候不会被阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestSortedTArray_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]any{3})
		s3 := garray.NewStrArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedTArrayFrom([]int{4, 5}, nil)
		s5 := garray.NewSortedStrArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})

		a1 := garray.NewSortedTArrayFrom(s1, nil)

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i1).Len(), 9)
		t.Assert(a1.Merge(i2).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedTArray_Json(t *testing.T) {
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedTArrayFrom(s1, nil)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewSortedTArray[string](nil)
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s2)

		var a3 garray.SortedTArray[string]
		a3.SetComparator(nil)
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := *garray.NewSortedTArrayFrom(s1, nil)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewSortedTArray[string](nil)
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s2)

		var a3 garray.SortedTArray[string]
		a3.SetComparator(nil)
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores *garray.SortedTArray[int]
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []int{99, 100, 98},
		}
		b, err := json.Marshal(data)
		t.AssertNil(err)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.AssertNil(err)
		t.Assert(user.Name, data["Name"])
		t.AssertNE(user.Scores, nil)
		t.Assert(user.Scores.Len(), 3)

		v, ok := user.Scores.PopLeft()
		t.AssertIN(v, data["Scores"])
		t.Assert(ok, true)

		v, ok = user.Scores.PopLeft()
		t.AssertIN(v, data["Scores"])
		t.Assert(ok, true)

		v, ok = user.Scores.PopLeft()
		t.AssertIN(v, data["Scores"])
		t.Assert(ok, true)

		v, ok = user.Scores.PopLeft()
		t.Assert(v, 0)
		t.Assert(ok, false)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores garray.SortedTArray[int]
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []int{99, 100, 98},
		}
		b, err := json.Marshal(data)
		t.AssertNil(err)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.AssertNil(err)
		t.Assert(user.Name, data["Name"])
		t.AssertNE(user.Scores, nil)
		t.Assert(user.Scores.Len(), 3)

		v, ok := user.Scores.PopLeft()
		t.AssertIN(v, data["Scores"])
		t.Assert(ok, true)

		v, ok = user.Scores.PopLeft()
		t.AssertIN(v, data["Scores"])
		t.Assert(ok, true)

		v, ok = user.Scores.PopLeft()
		t.AssertIN(v, data["Scores"])
		t.Assert(ok, true)

		v, ok = user.Scores.PopLeft()
		t.Assert(v, 0)
		t.Assert(ok, false)
	})
}

func TestSortedTArray_Iterator(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedTArrayFrom(slice, nil)
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

func TestSortedTArray_RemoveValue(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedTArrayFrom(slice, nil)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestSortedTArray_RemoveValues(t *testing.T) {
	slice := g.SliceStr{"a", "b", "d", "c"}
	array := garray.NewSortedTArrayFrom(slice, nil)
	gtest.C(t, func(t *gtest.T) {
		array.RemoveValues("a", "b", "c")
		t.Assert(array.Slice(), g.SliceStr{"d"})
	})
}

func TestSortedTArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *garray.SortedTArray[int]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`[2,3,1]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": g.SliceInt{2, 3, 1},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
}
func TestSortedTArray_Filter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := g.SliceInt{0, 1, 2, 3, 4, -1, -2}
		array := garray.NewSortedTArrayFromCopy(values, nil)
		t.Assert(array.Filter(func(index int, value int) bool {
			return value < 0
		}).Slice(), g.Slice{0, 1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFromCopy(g.SliceInt{-1, 1, 2, 3, 4, -2}, nil)
		t.Assert(array.Filter(func(index int, value int) bool {
			return value < 0
		}), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceInt{0, 1, 2, 3, 4, 0, 0}, nil)
		t.Assert(array.Filter(func(index int, value int) bool {
			return empty.IsEmpty(value)
		}), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceInt{1, 2, 3, 4}, nil)
		t.Assert(array.Filter(func(index int, value int) bool {
			return empty.IsEmpty(value)
		}), g.Slice{1, 2, 3, 4})
	})
}

func TestSortedTArray_FilterNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := g.SliceInt{0, 1, 2, 3, 4, -1, -2}
		array := garray.NewSortedTArrayFromCopy(values, gutil.ComparatorT)
		t.Assert(array.FilterNil().Slice(), g.SliceInt{-2, -1, 0, 1, 2, 3, 4})
	})

	gtest.C(t, func(t *gtest.T) {
		values := g.Slice{0, 1, 2, 3, 4, -1, -2, nil, []any{}, ""}
		array := garray.NewSortedTArrayFromCopy(values, nil)
		t.Assert(array.FilterNil().Slice(), g.Slice{"", -1, -2, 0, 1, 2, 3, 4, []any{}})
	})
}

func TestSortedTArray_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceInt{0, 1, 2, 3, 4, 0, 0}, nil)
		t.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceInt{1, 2, 3, 4}, nil)
		t.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceStr{"a", "", "b", "c", ""}, nil)
		t.Assert(array.FilterEmpty(), g.Slice{"a", "b", "c"})
	})
	gtest.C(t, func(t *gtest.T) {
		values := g.Slice{0, 1, 2, 3, 4, -1, -2, nil, []any{}, ""}
		array := garray.NewSortedTArrayFromCopy(values, nil)
		t.Assert(array.FilterEmpty().Slice(), g.Slice{-1, -2, 1, 2, 3, 4})
	})
}

func TestSortedTArray_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom(g.SliceStr{"1", "2"}, nil)
		t.Assert(array.Walk(func(value string) string {
			return "key-" + value
		}), g.Slice{"key-1", "key-2"})
	})
}

func TestSortedTArray_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom([]string{}, nil)
		t.Assert(array.IsEmpty(), true)
	})
}

func TestSortedTArray_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedTArrayFrom([]int{1, 2, 3, 4, 5}, nil)
		copyArray := array.DeepCopy().(*garray.SortedTArray[int])
		array.Add(6)
		copyArray.Add(7)
		cval, _ := copyArray.Get(5)
		val, _ := array.Get(5)
		t.AssertNE(cval, val)
	})
}
