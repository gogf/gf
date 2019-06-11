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

func Test_Array_Basic(t *testing.T) {
    gtest.Case(t, func() {
        expect := []interface{}{0, 1, 2, 3}
        array  := garray.NewArrayFrom(expect)
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
        gtest.Assert(array.Slice(), []interface{}{100, 200, 1, 2, 3, 4})
        array.InsertBefore(5, 300)
        array.InsertAfter(6, 400)
        gtest.Assert(array.Slice(), []interface{}{100, 200, 1, 2, 3, 300, 4, 400})
        gtest.Assert(array.Clear().Len(), 0)
    })
}

func TestArray_Sort(t *testing.T) {
    gtest.Case(t, func() {
        expect1 := []interface{}{0, 1, 2, 3}
        expect2 := []interface{}{3, 2, 1, 0}
        array   := garray.NewArray()
        for i := 3; i >= 0; i-- {
            array.Append(i)
        }
        array.SortFunc(func(v1, v2 interface{}) bool {
            return v1.(int) < v2.(int)
        })
        gtest.Assert(array.Slice(), expect1)
        array.SortFunc(func(v1, v2 interface{}) bool {
            return v1.(int) > v2.(int)
        })
        gtest.Assert(array.Slice(), expect2)
    })
}

func TestArray_Unique(t *testing.T) {
    gtest.Case(t, func() {
        expect := []interface{}{1, 1, 2, 3}
        array  := garray.NewArrayFrom(expect)
        gtest.Assert(array.Unique().Slice(), []interface{}{1, 2, 3})
    })
}

func TestArray_PushAndPop(t *testing.T) {
    gtest.Case(t, func() {
        expect := []interface{}{0, 1, 2, 3}
        array  := garray.NewArrayFrom(expect)
        gtest.Assert(array.Slice(),     expect)
        gtest.Assert(array.PopLeft(),   0)
        gtest.Assert(array.PopRight(),  3)
        gtest.AssertIN(array.PopRand(), []interface{}{1, 2})
        gtest.AssertIN(array.PopRand(), []interface{}{1, 2})
        gtest.Assert(array.Len(), 0)
        array.PushLeft(1).PushRight(2)
        gtest.Assert(array.Slice(),    []interface{}{1, 2})
    })
}

func TestArray_PopRands(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{100, 200, 300, 400, 500, 600}
        array := garray.NewFromCopy(a1)
        gtest.AssertIN(array.PopRands(2), a1)
    })
}

func TestArray_PopLeftsAndPopRights(t *testing.T) {
    gtest.Case(t, func() {
        value1 := []interface{}{0,1,2,3,4,5,6}
        value2 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(value1)
        array2 := garray.NewArrayFrom(value2)
        gtest.Assert(array1.PopLefts(2), []interface{}{0,1})
        gtest.Assert(array1.Slice(), []interface{}{2,3,4,5,6})
        gtest.Assert(array1.PopRights(2), []interface{}{5,6})
        gtest.Assert(array1.Slice(), []interface{}{2,3,4})
        gtest.Assert(array1.PopRights(20), []interface{}{2,3,4})
        gtest.Assert(array1.Slice(), []interface{}{})
        gtest.Assert(array2.PopLefts(20), []interface{}{0,1,2,3,4,5,6})
        gtest.Assert(array2.Slice(), []interface{}{})
    })
}

func TestArray_Range(t *testing.T) {
    gtest.Case(t, func() {
        value1 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(value1)
        gtest.Assert(array1.Range(0, 1), []interface{}{0})
        gtest.Assert(array1.Range(1, 2), []interface{}{1})
        gtest.Assert(array1.Range(0, 2), []interface{}{0, 1})
        gtest.Assert(array1.Range(-1, 10), value1)
    })
}

func TestArray_Merge(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0, 1, 2, 3}
        a2 := []interface{}{4, 5, 6, 7}
        array1 := garray.NewArrayFrom(a1)
        array2 := garray.NewArrayFrom(a2)
        gtest.Assert(array1.Merge(array2).Slice(), []interface{}{0,1,2,3,4,5,6,7})
    })
}

func TestArray_Fill(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0}
        a2 := []interface{}{0}
        array1 := garray.NewArrayFrom(a1)
        array2 := garray.NewArrayFrom(a2)
        gtest.Assert(array1.Fill(1, 2, 100).Slice(), []interface{}{0,100,100})
        gtest.Assert(array2.Fill(0, 2, 100).Slice(), []interface{}{100,100})
    })
}

func TestArray_Chunk(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{1,2,3,4,5}
        array1 := garray.NewArrayFrom(a1)
        chunks := array1.Chunk(2)
        gtest.Assert(len(chunks), 3)
        gtest.Assert(chunks[0], []interface{}{1,2})
        gtest.Assert(chunks[1], []interface{}{3,4})
        gtest.Assert(chunks[2], []interface{}{5})
    })
}

func TestArray_Pad(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0}
        array1 := garray.NewArrayFrom(a1)
        gtest.Assert(array1.Pad(3,  1).Slice(), []interface{}{0,1,1})
        gtest.Assert(array1.Pad(-4, 1).Slice(), []interface{}{1,0,1,1})
        gtest.Assert(array1.Pad(3,  1).Slice(), []interface{}{1,0,1,1})
    })
}

func TestArray_SubSlice(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(a1)
        gtest.Assert(array1.SubSlice(0, 2), []interface{}{0,1})
        gtest.Assert(array1.SubSlice(2, 2), []interface{}{2,3})
        gtest.Assert(array1.SubSlice(5, 8), []interface{}{5,6})
    })
}

func TestArray_Rand(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(a1)
        gtest.Assert(len(array1.Rands(2)),  2)
        gtest.Assert(len(array1.Rands(10)), 7)
        gtest.AssertIN(array1.Rands(1)[0], a1)
    })
}

func TestArray_Shuffle(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(a1)
        gtest.Assert(array1.Shuffle().Len(), 7)
    })
}

func TestArray_Reverse(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(a1)
        gtest.Assert(array1.Reverse().Slice(), []interface{}{6,5,4,3,2,1,0})
    })
}

func TestArray_Join(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{0,1,2,3,4,5,6}
        array1 := garray.NewArrayFrom(a1)
        gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
    })
}