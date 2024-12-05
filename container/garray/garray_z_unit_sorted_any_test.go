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

func TestSortedArray_NewSortedArrayFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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

		t.Assert(array1.Len(), 3)
		t.Assert(array1, []interface{}{"a", "c", "f"})

		t.Assert(array2.Len(), 4)
		t.Assert(array2, []interface{}{"k", "i", "j", "h"})
	})
}

func TestNewSortedArrayFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "f", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 interface{}) int {
			return -1
		}
		array1 := garray.NewSortedArrayFromCopy(a1, func1)
		array2 := garray.NewSortedArrayFromCopy(a1, func2)
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []interface{}{"a", "c", "f"})
		t.Assert(array1.Len(), 3)
		t.Assert(array2, []interface{}{"c", "f", "a"})
	})
}

func TestNewSortedArrayRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 interface{}) int {
			return gconv.Int(v1) - gconv.Int(v2)
		}

		array1 := garray.NewSortedArrayRange(1, 5, 1, func1)
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []interface{}{1, 2, 3, 4, 5})
	})
}

func TestSortedArray_SetArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "f", "c"}
		a2 := []interface{}{"e", "h", "g", "k"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}

		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.SetArray(a2)
		t.Assert(array1.Len(), 4)
		t.Assert(array1, []interface{}{"e", "g", "h", "k"})
	})

}

func TestSortedArray_Sort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.Sort()
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []interface{}{"a", "c", "f"})
	})

}

func TestSortedArray_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
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

func TestSortedArray_At(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		v := array1.At(2)
		t.Assert(v, "f")
	})
}

func TestSortedArray_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
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

func TestSortedArray_PopLeft(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := garray.NewSortedArrayFrom(
			[]interface{}{"a", "d", "c", "b"},
			gutil.ComparatorString,
		)
		i1, ok := array1.PopLeft()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "a")
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []interface{}{"b", "c", "d"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{1, 2, 3}, gutil.ComparatorInt)
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

func TestSortedArray_PopRight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := garray.NewSortedArrayFrom(
			[]interface{}{"a", "d", "c", "b"},
			gutil.ComparatorString,
		)
		i1, ok := array1.PopRight()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "d")
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []interface{}{"a", "b", "c"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{1, 2, 3}, gutil.ComparatorInt)
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

func TestSortedArray_PopRand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1, ok := array1.PopRand()
		t.Assert(ok, true)
		t.AssertIN(i1, []interface{}{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 3)

	})
}

func TestSortedArray_PopRands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRands(2)
		t.Assert(len(i1), 2)
		t.AssertIN(i1, []interface{}{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 2)

		i2 := array1.PopRands(3)
		t.Assert(len(i1), 2)
		t.AssertIN(i2, []interface{}{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArray(gutil.ComparatorInt)
		v, ok := array.PopLeft()
		t.Assert(v, nil)
		t.Assert(ok, false)
		t.Assert(array.PopLefts(10), nil)

		v, ok = array.PopRight()
		t.Assert(v, nil)
		t.Assert(ok, false)
		t.Assert(array.PopRights(10), nil)

		v, ok = array.PopRand()
		t.Assert(v, nil)
		t.Assert(ok, false)
		t.Assert(array.PopRands(10), nil)
	})
}

func TestSortedArray_PopLefts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopLefts(2)
		t.Assert(len(i1), 2)
		t.AssertIN(i1, []interface{}{"a", "d", "c", "b", "e", "f"})
		t.Assert(array1.Len(), 4)

		i2 := array1.PopLefts(5)
		t.Assert(len(i2), 4)
		t.AssertIN(i1, []interface{}{"a", "d", "c", "b", "e", "f"})
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedArray_PopRights(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRights(2)
		t.Assert(len(i1), 2)
		t.Assert(i1, []interface{}{"e", "f"})
		t.Assert(array1.Len(), 4)

		i2 := array1.PopRights(10)
		t.Assert(len(i2), 4)
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedArray_Range(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a1, func1, true)
		i1 := array1.Range(2, 5)
		t.Assert(i1, []interface{}{"c", "d", "e"})
		t.Assert(array1.Len(), 6)

		i2 := array1.Range(7, 5)
		t.Assert(len(i2), 0)
		i2 = array1.Range(-1, 2)
		t.Assert(i2, []interface{}{"a", "b"})

		i2 = array1.Range(4, 10)
		t.Assert(len(i2), 2)
		t.Assert(i2, []interface{}{"e", "f"})

		t.Assert(array2.Range(1, 3), []interface{}{"b", "c"})

	})
}

func TestSortedArray_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		a2 := []interface{}{"1", "2", "3", "b", "e", "f"}
		a3 := []interface{}{"4", "5", "6"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a2, func1)
		array3 := garray.NewSortedArrayFrom(a3, func1)
		t.Assert(array1.Sum(), 0)
		t.Assert(array2.Sum(), 6)
		t.Assert(array3.Sum(), 15)

	})
}

func TestSortedArray_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := array1.Clone()
		t.Assert(array1, array2)
		array1.Remove(1)
		t.AssertNE(array1, array2)

	})
}

func TestSortedArray_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		t.Assert(array1.Len(), 6)
		array1.Clear()
		t.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_Chunk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Chunk(2)
		t.Assert(len(i1), 3)
		t.Assert(i1[0], []interface{}{"a", "b"})
		t.Assert(i1[2], []interface{}{"e"})

		i1 = array1.Chunk(0)
		t.Assert(len(i1), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []interface{}{1, 2, 3})
		t.Assert(chunks[1], []interface{}{4, 5})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5, 6}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []interface{}{1, 2})
		t.Assert(chunks[1], []interface{}{3, 4})
		t.Assert(chunks[2], []interface{}{5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5, 6}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []interface{}{1, 2, 3})
		t.Assert(chunks[1], []interface{}{4, 5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedArray_SubSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a1, func1, true)
		i1 := array1.SubSlice(2, 3)
		t.Assert(len(i1), 3)
		t.Assert(i1, []interface{}{"c", "d", "e"})

		i1 = array1.SubSlice(2, 6)
		t.Assert(len(i1), 3)
		t.Assert(i1, []interface{}{"c", "d", "e"})

		i1 = array1.SubSlice(7, 2)
		t.Assert(len(i1), 0)

		s1 := array1.SubSlice(1, -2)
		t.Assert(s1, nil)

		s1 = array1.SubSlice(-9, 2)
		t.Assert(s1, nil)
		t.Assert(array2.SubSlice(1, 3), []interface{}{"b", "c", "d"})

	})
}

func TestSortedArray_Rand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1, ok := array1.Rand()
		t.Assert(ok, true)
		t.AssertIN(i1, []interface{}{"a", "d", "c"})
		t.Assert(array1.Len(), 3)

		array2 := garray.NewSortedArrayFrom([]interface{}{}, func1)
		v, ok := array2.Rand()
		t.Assert(ok, false)
		t.Assert(v, nil)
	})
}

func TestSortedArray_Rands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Rands(2)
		t.AssertIN(i1, []interface{}{"a", "d", "c"})
		t.Assert(len(i1), 2)
		t.Assert(array1.Len(), 3)

		i1 = array1.Rands(4)
		t.Assert(len(i1), 4)

		array2 := garray.NewSortedArrayFrom([]interface{}{}, func1)
		v := array2.Rands(1)
		t.Assert(v, nil)
	})
}

func TestSortedArray_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		t.Assert(array1.Join(","), `a,c,d`)
		t.Assert(array1.Join("."), `a.c.d`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{0, 1, `"a"`, `\a`}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorString)
		t.Assert(array1.Join("."), `"a".0.1.\a`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorString)
		t.Assert(array1.Join("."), "")
	})
}

func TestSortedArray_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{0, 1, "a", "b"}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorString)
		t.Assert(array1.String(), `[0,1,"a","b"]`)

		array1 = nil
		t.Assert(array1.String(), "")
	})
}

func TestSortedArray_CountValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"a", "d", "c", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		m1 := array1.CountValues()
		t.Assert(len(m1), 3)
		t.Assert(m1["c"], 2)
		t.Assert(m1["a"], 1)

	})
}

func TestSortedArray_SetUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		array1.SetUnique(true)
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []interface{}{1, 2, 3, 4, 5})
	})
}

func TestSortedArray_Unique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := garray.NewSortedArrayFrom(a1, gutil.ComparatorInt)
		array1.Unique()
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []interface{}{1, 2, 3, 4, 5})

		array2 := garray.NewSortedArrayFrom([]interface{}{}, gutil.ComparatorInt)
		array2.Unique()
		t.Assert(array2.Len(), 0)
		t.Assert(array2, []interface{}{})
	})
}

func TestSortedArray_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, func1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		// go1
		go a1.LockFunc(func(n1 []interface{}) { // 读写锁
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

func TestSortedArray_RLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, func1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		// go1
		go a1.RLockFunc(func(n1 []interface{}) { // 读写锁
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

func TestSortedArray_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i1).Len(), 9)
		t.Assert(a1.Merge(i2).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedArray_Json(t *testing.T) {
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		s1 := []interface{}{"a", "b", "d", "c"}
		s2 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, gutil.ComparatorString)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewSortedArray(gutil.ComparatorString)
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s2)

		var a3 garray.SortedArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		s1 := []interface{}{"a", "b", "d", "c"}
		s2 := []interface{}{"a", "b", "c", "d"}
		a1 := *garray.NewSortedArrayFrom(s1, gutil.ComparatorString)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := garray.NewSortedArray(gutil.ComparatorString)
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s2)

		var a3 garray.SortedArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores *garray.SortedArray
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
		t.Assert(v, nil)
		t.Assert(ok, false)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores garray.SortedArray
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
		t.Assert(v, nil)
		t.Assert(ok, false)
	})
}

func TestSortedArray_Iterator(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := garray.NewSortedArrayFrom(slice, gutil.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		array.Iterator(func(k int, v interface{}) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	gtest.C(t, func(t *gtest.T) {
		array.IteratorAsc(func(k int, v interface{}) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	gtest.C(t, func(t *gtest.T) {
		array.IteratorDesc(func(k int, v interface{}) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		array.Iterator(func(k int, v interface{}) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		array.IteratorAsc(func(k int, v interface{}) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		array.IteratorDesc(func(k int, v interface{}) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestSortedArray_RemoveValue(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := garray.NewSortedArrayFrom(slice, gutil.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestSortedArray_RemoveValues(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := garray.NewSortedArrayFrom(slice, gutil.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		array.RemoveValues("a", "b", "c")
		t.Assert(array.Slice(), g.SliceStr{"d"})
	})
}

func TestSortedArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *garray.SortedArray
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
			"array": g.Slice{2, 3, 1},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
}
func TestSortedArray_Filter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}
		array := garray.NewSortedArrayFromCopy(values, gutil.ComparatorInt)
		t.Assert(array.Filter(func(index int, value interface{}) bool {
			return empty.IsNil(value)
		}).Slice(), g.Slice{0, "", g.Slice{}, 1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFromCopy(g.Slice{nil, 1, 2, 3, 4, nil}, gutil.ComparatorInt)
		t.Assert(array.Filter(func(index int, value interface{}) bool {
			return empty.IsNil(value)
		}), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}, gutil.ComparatorInt)
		t.Assert(array.Filter(func(index int, value interface{}) bool {
			return empty.IsEmpty(value)
		}), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{1, 2, 3, 4}, gutil.ComparatorInt)
		t.Assert(array.Filter(func(index int, value interface{}) bool {
			return empty.IsEmpty(value)
		}), g.Slice{1, 2, 3, 4})
	})
}

func TestSortedArray_FilterNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}
		array := garray.NewSortedArrayFromCopy(values, gutil.ComparatorInt)
		t.Assert(array.FilterNil().Slice(), g.Slice{0, "", g.Slice{}, 1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFromCopy(g.Slice{nil, 1, 2, 3, 4, nil}, gutil.ComparatorInt)
		t.Assert(array.FilterNil(), g.Slice{1, 2, 3, 4})
	})
}

func TestSortedArray_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}, gutil.ComparatorInt)
		t.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{1, 2, 3, 4}, gutil.ComparatorInt)
		t.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
}

func TestSortedArray_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom(g.Slice{"1", "2"}, gutil.ComparatorString)
		t.Assert(array.Walk(func(value interface{}) interface{} {
			return "key-" + gconv.String(value)
		}), g.Slice{"key-1", "key-2"})
	})
}

func TestSortedArray_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom([]interface{}{}, gutil.ComparatorString)
		t.Assert(array.IsEmpty(), true)
	})
}

func TestSortedArray_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewSortedArrayFrom([]interface{}{1, 2, 3, 4, 5}, gutil.ComparatorString)
		copyArray := array.DeepCopy().(*garray.SortedArray)
		array.Add(6)
		copyArray.Add(7)
		cval, _ := copyArray.Get(5)
		val, _ := array.Get(5)
		t.AssertNE(cval, val)
	})
}
