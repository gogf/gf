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
    "github.com/gogf/gf/g/util/gconv"
    "testing"
)

func Test_StringArray_Basic(t *testing.T) {
    gtest.Case(t, func() {
        expect := []string{"0", "1", "2", "3"}
        array  := garray.NewStringArrayFrom(expect)
        gtest.Assert(array.Slice(), expect)
        array.Set(0, "100")
        gtest.Assert(array.Get(0),        100)
        gtest.Assert(array.Get(1),        1)
        gtest.Assert(array.Search("100"),   0)
        gtest.Assert(array.Contains("100"), true)
        gtest.Assert(array.Remove(0),     100)
        gtest.Assert(array.Contains("100"), false)
        array.Append("4")
        gtest.Assert(array.Len(), 4)
        array.InsertBefore(0, "100")
        array.InsertAfter(0, "200")
        gtest.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "4"})
        array.InsertBefore(5, "300")
        array.InsertAfter(6, "400")
        gtest.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "300", "4", "400"})
        gtest.Assert(array.Clear().Len(), 0)
    })
}

func TestStringArray_Sort(t *testing.T) {
    gtest.Case(t, func() {
        expect1 := []string{"0", "1", "2", "3"}
        expect2 := []string{"3", "2", "1", "0"}
        array   := garray.NewStringArray()
        for i := 3; i >= 0; i-- {
            array.Append(gconv.String(i))
        }
        array.Sort()
        gtest.Assert(array.Slice(), expect1)
        array.Sort(true)
        gtest.Assert(array.Slice(), expect2)
    })
}

func TestStringArray_Unique(t *testing.T) {
    gtest.Case(t, func() {
        expect := []string{"1", "1", "2", "3"}
        array  := garray.NewStringArrayFrom(expect)
        gtest.Assert(array.Unique().Slice(), []string{"1", "2", "3"})
    })
}

func TestStringArray_PushAndPop(t *testing.T) {
    gtest.Case(t, func() {
        expect := []string{"0", "1", "2", "3"}
        array  := garray.NewStringArrayFrom(expect)
        gtest.Assert(array.Slice(),     expect)
        gtest.Assert(array.PopLeft(),   "0")
        gtest.Assert(array.PopRight(),  "3")
        gtest.AssertIN(array.PopRand(), []string{"1", "2"})
        gtest.AssertIN(array.PopRand(), []string{"1", "2"})
        gtest.Assert(array.Len(), 0)
        array.PushLeft("1").PushRight("2")
        gtest.Assert(array.Slice(),    []string{"1", "2"})
    })
}

func TestStringArray_PopLeftsAndPopRights(t *testing.T) {
    gtest.Case(t, func() {
        value1 := []string{"0","1","2","3","4","5","6"}
        value2 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(value1)
        array2 := garray.NewStringArrayFrom(value2)
        gtest.Assert(array1.PopLefts(2), []interface{}{"0","1"})
        gtest.Assert(array1.Slice(), []interface{}{"2","3","4","5","6"})
        gtest.Assert(array1.PopRights(2), []interface{}{"5","6"})
        gtest.Assert(array1.Slice(), []interface{}{"2","3","4"})
        gtest.Assert(array1.PopRights(20), []interface{}{"2","3","4"})
        gtest.Assert(array1.Slice(), []interface{}{})
        gtest.Assert(array2.PopLefts(20), []interface{}{"0","1","2","3","4","5","6"})
        gtest.Assert(array2.Slice(), []interface{}{})
    })
}

func TestString_Range(t *testing.T) {
    gtest.Case(t, func() {
        value1 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(value1)
        gtest.Assert(array1.Range(0, 1), []interface{}{"0"})
        gtest.Assert(array1.Range(1, 2), []interface{}{"1"})
        gtest.Assert(array1.Range(0, 2), []interface{}{"0", "1"})
        gtest.Assert(array1.Range(-1, 10), value1)
    })
}

func TestStringArray_Merge(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0", "1", "2", "3"}
        a2 := []string{"4", "5", "6", "7"}
        array1 := garray.NewStringArrayFrom(a1)
        array2 := garray.NewStringArrayFrom(a2)
        gtest.Assert(array1.Merge(array2).Slice(), []string{"0","1","2","3","4","5","6","7"})
    })
}

func TestStringArray_Fill(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0"}
        a2 := []string{"0"}
        array1 := garray.NewStringArrayFrom(a1)
        array2 := garray.NewStringArrayFrom(a2)
        gtest.Assert(array1.Fill(1, 2, "100").Slice(), []string{"0","100","100"})
        gtest.Assert(array2.Fill(0, 2, "100").Slice(), []string{"100","100"})
    })
}

func TestStringArray_Chunk(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"1","2","3","4","5"}
        array1 := garray.NewStringArrayFrom(a1)
        chunks := array1.Chunk(2)
        gtest.Assert(len(chunks), 3)
        gtest.Assert(chunks[0], []string{"1","2"})
        gtest.Assert(chunks[1], []string{"3","4"})
        gtest.Assert(chunks[2], []string{"5"})
    })
}

func TestStringArray_Pad(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0"}
        array1 := garray.NewStringArrayFrom(a1)
        gtest.Assert(array1.Pad(3,  "1").Slice(), []string{"0","1","1"})
        gtest.Assert(array1.Pad(-4, "1").Slice(), []string{"1","0","1","1"})
        gtest.Assert(array1.Pad(3,  "1").Slice(), []string{"1","0","1","1"})
    })
}

func TestStringArray_SubSlice(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(a1)
        gtest.Assert(array1.SubSlice(0, 2), []string{"0","1"})
        gtest.Assert(array1.SubSlice(2, 2), []string{"2","3"})
        gtest.Assert(array1.SubSlice(5, 8), []string{"5","6"})
    })
}

func TestStringArray_Rand(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(a1)
        gtest.Assert(len(array1.Rands(2)),  "2")
        gtest.Assert(len(array1.Rands(10)), "7")
        gtest.AssertIN(array1.Rands(1)[0], a1)
    })
}

func TestStringArray_PopRands(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []interface{}{"100", "200", "300", "400", "500", "600"}
        array := garray.NewFromCopy(a1)
        gtest.AssertIN(array.PopRands(2), a1)
    })
}

func TestStringArray_Shuffle(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(a1)
        gtest.Assert(array1.Shuffle().Len(), 7)
    })
}

func TestStringArray_Reverse(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(a1)
        gtest.Assert(array1.Reverse().Slice(), []string{"6","5","4","3","2","1","0"})
    })
}

func TestStringArray_Join(t *testing.T) {
    gtest.Case(t, func() {
        a1 := []string{"0","1","2","3","4","5","6"}
        array1 := garray.NewStringArrayFrom(a1)
        gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
    })
}