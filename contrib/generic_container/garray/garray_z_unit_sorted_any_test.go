// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray

import (
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/generic_container/v2/comparator"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/empty"

	"github.com/gogf/gf/contrib/generic_container/v2/internal/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestSortedArray_NewSortedArrayFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		a2 := []string{"h", "j", "i", "k"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 string) int {
			return -1
		}
		array1 := NewSortedArrayFrom(a1, func1)
		array2 := NewSortedArrayFrom(a2, func2)

		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"a", "c", "f"})

		t.Assert(array2.Len(), 4)
		t.Assert(array2, []string{"k", "i", "j", "h"})
	})
}

func TestNewSortedArrayFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 string) int {
			return -1
		}
		array1 := NewSortedArrayFromCopy(a1, func1)
		array2 := NewSortedArrayFromCopy(a1, func2)
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"a", "c", "f"})
		t.Assert(array1.Len(), 3)
		t.Assert(array2, []string{"c", "f", "a"})
	})
}

func TestNewSortedArrayRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 int) int {
			return v1 - v2
		}

		array1 := NewSortedArrayRange(1, 5, 1, func1)
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})
	})
}

func TestSortedArray_SetArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		a2 := []string{"e", "h", "g", "k"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}

		array1 := NewSortedArrayFrom(a1, func1)
		array1.SetArray(a2)
		t.Assert(array1.Len(), 4)
		t.Assert(array1, []string{"e", "g", "h", "k"})
	})

}

func TestSortedArray_Sort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		array1.Sort()
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"a", "c", "f"})
	})

}

func TestSortedArray_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
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
		a1 := []string{"a", "f", "c"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		v := array1.At(2)
		t.Assert(v, "f")
	})
}

func TestSortedArray_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
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
		array1 := NewSortedArrayFrom(
			[]string{"a", "d", "c", "b"},
			comparator.ComparatorString,
		)
		i1, ok := array1.PopLeft()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "a")
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"b", "c", "d"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := NewSortedArrayFrom[int]([]int{1, 2, 3}, comparator.ComparatorInt)
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
		array1 := NewSortedArrayFrom(
			[]string{"a", "d", "c", "b"},
			comparator.ComparatorString,
		)
		i1, ok := array1.PopRight()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "d")
		t.Assert(array1.Len(), 3)
		t.Assert(array1, []string{"a", "b", "c"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := NewSortedArrayFrom[int]([]int{1, 2, 3}, comparator.ComparatorInt)
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
		a1 := []string{"a", "d", "c", "b"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1, ok := array1.PopRand()
		t.Assert(ok, true)
		t.AssertIN(i1, []string{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 3)

	})
}

func TestSortedArray_PopRands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRands(2)
		t.Assert(len(i1), 2)
		t.AssertIN(i1, []string{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 2)

		i2 := array1.PopRands(3)
		t.Assert(len(i1), 2)
		t.AssertIN(i2, []string{"a", "d", "c", "b"})
		t.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := NewSortedArray[int](comparator.ComparatorInt)
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

func TestSortedArray_PopLefts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1 := array1.PopLefts(2)
		t.Assert(len(i1), 2)
		t.AssertIN(i1, []string{"a", "d", "c", "b", "e", "f"})
		t.Assert(array1.Len(), 4)

		i2 := array1.PopLefts(5)
		t.Assert(len(i2), 4)
		t.AssertIN(i1, []string{"a", "d", "c", "b", "e", "f"})
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedArray_PopRights(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRights(2)
		t.Assert(len(i1), 2)
		t.Assert(i1, []string{"e", "f"})
		t.Assert(array1.Len(), 4)

		i2 := array1.PopRights(10)
		t.Assert(len(i2), 4)
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedArray_Range(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		array2 := NewSortedArrayFrom(a1, func1, true)
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

func TestSortedArray_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		a2 := []string{"1", "2", "3", "b", "e", "f"}
		a3 := []string{"4", "5", "6"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		array2 := NewSortedArrayFrom(a2, func1)
		array3 := NewSortedArrayFrom(a3, func1)
		t.Assert(array1.Sum(), 0)
		t.Assert(array2.Sum(), 6)
		t.Assert(array3.Sum(), 15)

	})
}

func TestSortedArray_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		array2 := array1.Clone()
		t.Assert(array1, array2)
		array1.Remove(1)
		t.AssertNE(array1, array2)

	})
}

func TestSortedArray_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		t.Assert(array1.Len(), 6)
		array1.Clear()
		t.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_Chunk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1 := array1.Chunk(2)
		t.Assert(len(i1), 3)
		t.Assert(i1[0], []string{"a", "b"})
		t.Assert(i1[2], []string{"e"})

		i1 = array1.Chunk(0)
		t.Assert(len(i1), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := NewSortedArrayFrom[int](a1, comparator.ComparatorInt)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []int{1, 2, 3})
		t.Assert(chunks[1], []int{4, 5})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorInt)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []int{1, 2})
		t.Assert(chunks[1], []int{3, 4})
		t.Assert(chunks[2], []int{5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorInt)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []int{1, 2, 3})
		t.Assert(chunks[1], []int{4, 5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedArray_SubSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		array2 := NewSortedArrayFrom(a1, func1, true)
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

func TestSortedArray_Rand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1, ok := array1.Rand()
		t.Assert(ok, true)
		t.AssertIN(i1, []string{"a", "d", "c"})
		t.Assert(array1.Len(), 3)

		array2 := NewSortedArrayFrom([]string{}, func1)
		v, ok := array2.Rand()
		t.Assert(ok, false)
		t.Assert(v, nil)
	})
}

func TestSortedArray_Rands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		i1 := array1.Rands(2)
		t.AssertIN(i1, []string{"a", "d", "c"})
		t.Assert(len(i1), 2)
		t.Assert(array1.Len(), 3)

		i1 = array1.Rands(4)
		t.Assert(len(i1), 4)

		array2 := NewSortedArrayFrom([]string{}, func1)
		v := array2.Rands(1)
		t.Assert(v, nil)
	})
}

func TestSortedArray_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		t.Assert(array1.Join(","), `a,c,d`)
		t.Assert(array1.Join("."), `a.c.d`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []int{0, 1, 2, 3}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorInt)
		t.Assert(array1.Join("."), `0.1.2.3`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []string{}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorString)
		t.Assert(array1.Join("."), "")
	})
}

func TestSortedArray_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{0, 1, 2, 3}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorInt)
		t.Assert(array1.String(), `[0,1,2,3]`)

		array1 = nil
		t.Assert(array1.String(), "")
	})
}

func TestSortedArray_CountValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "c"}

		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := NewSortedArrayFrom(a1, func1)
		m1 := array1.CountValues()
		t.Assert(len(m1), 3)
		t.Assert(m1["c"], 2)
		t.Assert(m1["a"], 1)

	})
}

func TestSortedArray_SetUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorInt)
		array1.SetUnique(true)
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})
	})
}

func TestSortedArray_Unique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := NewSortedArrayFrom(a1, comparator.ComparatorInt)
		array1.Unique()
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})

		array2 := NewSortedArrayFrom([]int{}, comparator.ComparatorInt)
		array2.Unique()
		t.Assert(array2.Len(), 0)
		t.Assert(array2, []int{})
	})
}

func TestSortedArray_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []string{"a", "b", "c", "d"}
		a1 := NewSortedArrayFrom(s1, func1, true)

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

func TestSortedArray_RLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		func1 := func(v1, v2 string) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []string{"a", "b", "c", "d"}
		a1 := NewSortedArrayFrom(s1, func1, true)

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

func TestSortedArray_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i2 := NewArrayFrom([]string{"3"})
		s3 := NewArrayFrom([]string{"g", "h"})
		s4 := NewSortedArrayFrom([]string{"4", "5"}, comparator.ComparatorString)
		s5 := NewSortedArrayFrom(s2, comparator.ComparatorString)
		s6 := NewSortedArrayFrom([]string{"1", "2", "3"}, comparator.ComparatorString)

		a1 := NewSortedArrayFrom(s1, comparator.ComparatorString)

		t.Assert(a1.MergeSlice(s2).Len(), 6)
		t.Assert(a1.Merge(s3).Len(), 8)
		t.Assert(a1.Merge(i2).Len(), 9)
		t.Assert(a1.Merge(s3).Len(), 11)
		t.Assert(a1.Merge(s4).Len(), 13)
		t.Assert(a1.Merge(s5).Len(), 15)
		t.Assert(a1.Merge(s6).Len(), 18)
	})
}

func TestSortedArray_Json(t *testing.T) {
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := NewSortedArrayFrom(s1, comparator.ComparatorString)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := NewSortedArray(comparator.ComparatorString)
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s2)

		var a3 SortedArray[string]
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array value
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		s2 := []string{"a", "b", "c", "d"}
		a1 := *NewSortedArrayFrom(s1, comparator.ComparatorString)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := NewSortedArray(comparator.ComparatorString)
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.AssertNil(err1)
		t.Assert(a2.Slice(), s2)

		var a3 SortedArray[string]
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Slice(), s1)
		t.Assert(a3.Interfaces(), s1)
	})
	// array pointer
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name   string
			Scores *SortedArray[int]
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
			Scores SortedArray[int]
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

func TestSortedArray_Iterator(t *testing.T) {
	slice := []string{"a", "b", "d", "c"}
	array := NewSortedArrayFrom(slice, comparator.ComparatorString)
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

func TestSortedArray_RemoveValue(t *testing.T) {
	slice := []string{"a", "b", "d", "c"}
	array := NewSortedArrayFrom(slice, comparator.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestSortedArray_RemoveValues(t *testing.T) {
	slice := []string{"a", "b", "d", "c"}
	array := NewSortedArrayFrom(slice, comparator.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		array.RemoveValues("a", "b", "c")
		t.Assert(array.Slice(), g.SliceStr{"d"})
	})
}

func TestSortedArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *SortedArray[byte]
	}
	type VInt struct {
		Name  string
		Array *SortedArray[int]
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
		t.Assert(v.Array.Slice(), []byte{1, 2, 3})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *VInt
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []int{2, 3, 1},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), []int{1, 2, 3})
	})
}

func comparatorExampleElement(a, b *exampleElement) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil && b != nil {
		return -1
	}
	if a != nil && b == nil {
		return 1
	}
	return a.code - b.code
}

func TestSortedArray_Filter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			{code: 2},
			{code: 0},
			{code: 1},
		}
		array := NewSortedArrayFromCopy[*exampleElement](values, comparatorExampleElement)
		t.Assert(array.Filter(func(index int, value *exampleElement) bool {
			return empty.IsNil(value)
		}).Slice(), []*exampleElement{
			{code: 0},
			{code: 1},
			{code: 2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			nil,
			{code: 2},
			{code: 0},
			{code: 1},
			nil,
		}
		array := NewSortedArrayFromCopy[*exampleElement](values, comparatorExampleElement)
		t.Assert(array.Filter(func(index int, value *exampleElement) bool {
			return empty.IsNil(value)
		}).Slice(), []*exampleElement{
			{code: 0},
			{code: 1},
			{code: 2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			{code: 2},
			{},
			{code: 0},
			{code: 1},
		}
		array := NewArrayFromCopy(values)
		t.Assert(array.Filter(func(index int, value *exampleElement) bool {
			return empty.IsEmpty(value)
		}).Slice(), []*exampleElement{
			{code: 1},
			{code: 2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			{code: 2},
			{code: 3},
			{code: 1},
		}
		array := NewSortedArrayFromCopy[*exampleElement](values, comparatorExampleElement)
		t.Assert(array.Filter(func(index int, value *exampleElement) bool {
			return empty.IsEmpty(value)
		}).Slice(), []*exampleElement{
			{code: 1},
			{code: 2},
			{code: 3},
		})
	})
}

func TestSortedArray_FilterNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			{code: 1},
			{code: 0},
			{code: 2},
		}
		array := NewSortedArrayFromCopy[*exampleElement](values, comparatorExampleElement)
		t.Assert(array.FilterNil().Slice(), []*exampleElement{
			{code: 0},
			{code: 1},
			{code: 2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			nil,
			{code: 1},
			{code: 0},
			{code: 2},
			nil,
		}
		array := NewSortedArrayFromCopy[*exampleElement](values, comparatorExampleElement)
		t.Assert(array.FilterNil().Slice(), []*exampleElement{
			{code: 0},
			{code: 1},
			{code: 2},
		})
	})
}

func TestSortedArray_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			{code: 2},
			{},
			{code: 0},
			{code: 1},
		}
		array := NewArrayFromCopy(values)
		t.Assert(array.FilterEmpty().Slice(), []*exampleElement{
			{code: 1},
			{code: 2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		values := []*exampleElement{
			{code: 2},
			{code: 3},
			{code: 1},
		}
		array := NewArrayFromCopy(values)
		t.Assert(array.FilterEmpty().Slice(), []*exampleElement{
			{code: 1},
			{code: 2},
			{code: 3},
		})
	})
}

func TestSortedArray_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := NewSortedArrayFrom([]string{"1", "2"}, comparator.ComparatorString)
		t.Assert(array.Walk(func(value string) string {
			return "key-" + gconv.String(value)
		}), g.Slice{"key-1", "key-2"})
	})
}

func TestSortedArray_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := NewSortedArrayFrom([]string{}, comparator.ComparatorString)
		t.Assert(array.IsEmpty(), true)
	})
}

func TestSortedArray_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := NewSortedArrayFrom([]int{1, 2, 3, 4, 5}, comparator.ComparatorInt)
		copyArray := array.DeepCopy().(*SortedArray[int])
		array.Add(6)
		copyArray.Add(7)
		cval, _ := copyArray.Get(5)
		val, _ := array.Get(5)
		t.AssertNE(cval, val)
	})
}
