// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray_test

import (
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)

func Test_IntArray_Basic(t *testing.T) {
    gtest.Case(t, func() {
        expect := []int{0, 1, 2, 3}
        array  := garray.NewIntArrayFrom(expect)
        gtest.Assert(array.Slice(), expect)
        array.Set(0, 100)
        gtest.Assert(array.Get(0),        100)
        gtest.Assert(array.Get(1),        1)
        gtest.Assert(array.Search(100),   0)
        gtest.Assert(array.Contains(100), true)
        gtest.Assert(array.Remove(0),     100)
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
        array   := garray.NewIntArray()
        for i := 3; i >= 0; i-- {
            array.Append(i)
        }
        array.Sort()
        gtest.Assert(array.Slice(), expect1)
        array.Sort(true)
        gtest.Assert(array.Slice(), expect2)
    })
}

func TestIntArray_Unique(t *testing.T) {
    gtest.Case(t, func() {
        expect := []int{1, 1, 2, 3}
        array  := garray.NewIntArrayFrom(expect)
        gtest.Assert(array.Unique().Slice(), []int{1, 2, 3})
    })
}

func TestIntArray_PushAndPop(t *testing.T) {
    gtest.Case(t, func() {
        expect := []int{0, 1, 2, 3}
        array  := garray.NewIntArrayFrom(expect)
        gtest.Assert(array.Slice(),     expect)
        gtest.Assert(array.PopLeft(),   0)
        gtest.Assert(array.PopRight(),  3)
        gtest.AssertIN(array.PopRand(), []int{1, 2})
        gtest.AssertIN(array.PopRand(), []int{1, 2})
        gtest.Assert(array.Len(), 0)
        array.PushLeft(1).PushRight(2)
        gtest.Assert(array.Slice(),    []int{1, 2})
    })
}

func TestIntArray_PopLeftsAndPopRights(t *testing.T) {
    gtest.Case(t, func() {
        value1 := []int{0,1,2,3,4,5,6}
        value2 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(value1)
        array2 := garray.NewIntArrayFrom(value2)
        gtest.Assert(array1.PopLefts(2), []int{0,1})
        gtest.Assert(array1.Slice(), []int{2,3,4,5,6})
        gtest.Assert(array1.PopRights(2), []int{5,6})
        gtest.Assert(array1.Slice(), []int{2,3,4})
        gtest.Assert(array1.PopRights(20), []int{2,3,4})
        gtest.Assert(array1.Slice(), []int{})
        gtest.Assert(array2.PopLefts(20), []int{0,1,2,3,4,5,6})
        gtest.Assert(array2.Slice(), []int{})
    })
}

func TestIntArray_Range(t *testing.T) {
    gtest.Case(t, func() {
        value1 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(value1)
        gtest.Assert(array1.Range(0, 1), []int{0})
        gtest.Assert(array1.Range(1, 2), []int{1})
        gtest.Assert(array1.Range(0, 2), []int{0, 1})
        gtest.Assert(array1.Range(-1, 10), value1)
    })
}

func TestIntArray_Merge(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0, 1, 2, 3}
        a2 := []int{4, 5, 6, 7}
        array1 := garray.NewIntArrayFrom(a1)
        array2 := garray.NewIntArrayFrom(a2)
        gtest.Assert(array1.Merge(array2).Slice(), []int{0,1,2,3,4,5,6,7})
    })
}

func TestIntArray_Fill(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0}
        a2 := []int{0}
        array1 := garray.NewIntArrayFrom(a1)
        array2 := garray.NewIntArrayFrom(a2)
        gtest.Assert(array1.Fill(1, 2, 100).Slice(), []int{0,100,100})
        gtest.Assert(array2.Fill(0, 2, 100).Slice(), []int{100,100})
    })
}

func TestIntArray_Chunk(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{1,2,3,4,5}
        array1 := garray.NewIntArrayFrom(a1)
        chunks := array1.Chunk(2)
        gtest.Assert(len(chunks), 3)
        gtest.Assert(chunks[0], []int{1,2})
        gtest.Assert(chunks[1], []int{3,4})
        gtest.Assert(chunks[2], []int{5})
    })
}

func TestIntArray_Pad(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0}
        array1 := garray.NewIntArrayFrom(a1)
        gtest.Assert(array1.Pad(3,  1).Slice(), []int{0,1,1})
        gtest.Assert(array1.Pad(-4, 1).Slice(), []int{1,0,1,1})
        gtest.Assert(array1.Pad(3,  1).Slice(), []int{1,0,1,1})
    })
}

func TestIntArray_SubSlice(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(a1)
        gtest.Assert(array1.SubSlice(0, 2), []int{0,1})
        gtest.Assert(array1.SubSlice(2, 2), []int{2,3})
        gtest.Assert(array1.SubSlice(5, 8), []int{5,6})
    })
}

func TestIntArray_Rand(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(a1)
        gtest.Assert(len(array1.Rands(2)),  2)
        gtest.Assert(len(array1.Rands(10)), 7)
        gtest.AssertIN(array1.Rands(1)[0], a1)
        gtest.AssertIN(array1.Rand(), a1)
    })
}

func TestIntArray_PopRands(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{100, 200, 300, 400, 500, 600}
        array := garray.NewFromCopy(a1)
        gtest.AssertIN(array.PopRands(2), a1)
    })
}

func TestIntArray_Shuffle(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(a1)
        gtest.Assert(array1.Shuffle().Len(), 7)
    })
}

func TestIntArray_Reverse(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(a1)
        gtest.Assert(array1.Reverse().Slice(), []int{6,5,4,3,2,1,0})
    })
}

func TestIntArray_Join(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []int{0,1,2,3,4,5,6}
        array1 := garray.NewIntArrayFrom(a1)
        gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
    })
}