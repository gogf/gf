// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gparser_test

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/encoding/gparser"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)

func Test_New(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.Assert(j.Get("n"), "123456789")
        gtest.Assert(j.Get("m"), g.Map{"k" : "v"})
        gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
    })
}

func Test_NewUnsafe(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.NewUnsafe(data)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
}

func Test_Encode(t *testing.T) {
    value := g.Slice{1, 2, 3}
    gtest.Case(t, func() {
        b, err := gparser.VarToJson(value)
        gtest.Assert(err, nil)
        gtest.Assert(b,   []byte(`[1,2,3]`))
    })
}

func Test_Decode(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j,          nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
}

func Test_SplitChar(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        j.SetSplitChar(byte('#'))
        gtest.AssertNE(j, nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m#k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a#1"), 2)
    })
}

func Test_ViolenceCheck(t *testing.T) {
    data := []byte(`{"m":{"a":[1,2,3], "v1.v2":"4"}}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j, nil)
        gtest.Assert(j.Get("m.a.2"),   3)
        gtest.Assert(j.Get("m.v1.v2"), nil)
        j.SetViolenceCheck(true)
        gtest.Assert(j.Get("m.v1.v2"), 4)
    })
}

func Test_GetToVar(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        var m map[string]string
        var n int
        var a []int
        j := gparser.New(data)
        gtest.AssertNE(j, nil)

        j.GetToVar("n", &n)
        j.GetToVar("m", &m)
        j.GetToVar("a", &a)

        gtest.Assert(n, "123456789")
        gtest.Assert(m, g.Map{"k" : "v"})
        gtest.Assert(a, g.Slice{1, 2, 3})
    })
}

func Test_GetMap(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j, nil)
        gtest.Assert(j.GetMap("n"), nil)
        gtest.Assert(j.GetMap("m"), g.Map{"k" : "v"})
        gtest.Assert(j.GetMap("a"), nil)
    })
}

func Test_GetArray(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j, nil)
        gtest.Assert(j.GetArray("n"), g.Array{123456789})
        gtest.Assert(j.GetArray("m"), g.Array{g.Map{"k":"v"}})
        gtest.Assert(j.GetArray("a"), g.Array{1,2,3})
    })
}

func Test_GetString(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j, nil)
        gtest.AssertEQ(j.GetString("n"), "123456789")
        gtest.AssertEQ(j.GetString("m"), `{"k":"v"}`)
        gtest.AssertEQ(j.GetString("a"), `[1,2,3]`)
        gtest.AssertEQ(j.GetString("i"), "")
    })
}

func Test_GetStrings(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j, nil)
        gtest.AssertEQ(j.GetStrings("n"), g.SliceStr{"123456789"})
        gtest.AssertEQ(j.GetStrings("m"), g.SliceStr{`{"k":"v"}`})
        gtest.AssertEQ(j.GetStrings("a"), g.SliceStr{"1", "2", "3"})
        gtest.AssertEQ(j.GetStrings("i"), nil)
    })
}

func Test_GetInterfaces(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    gtest.Case(t, func() {
        j := gparser.New(data)
        gtest.AssertNE(j, nil)
        gtest.AssertEQ(j.GetInterfaces("n"), g.Array{123456789})
        gtest.AssertEQ(j.GetInterfaces("m"), g.Array{g.Map{"k":"v"}})
        gtest.AssertEQ(j.GetInterfaces("a"), g.Array{1,2,3})
    })
}

func Test_Len(t *testing.T) {
    gtest.Case(t, func() {
        p := gparser.New(nil)
        p.Append("a", 1)
        p.Append("a", 2)
        gtest.Assert(p.Len("a"), 2)
    })
    gtest.Case(t, func() {
        p := gparser.New(nil)
        p.Append("a.b", 1)
        p.Append("a.c", 2)
        gtest.Assert(p.Len("a"), 2)
    })
    gtest.Case(t, func() {
        p := gparser.New(nil)
        p.Set("a", 1)
        gtest.Assert(p.Len("a"), -1)
    })
}

func Test_Append(t *testing.T) {
    gtest.Case(t, func() {
        p := gparser.New(nil)
        p.Append("a", 1)
        p.Append("a", 2)
        gtest.Assert(p.Get("a"), g.Slice{1, 2})
    })
    gtest.Case(t, func() {
        p := gparser.New(nil)
        p.Append("a.b", 1)
        p.Append("a.c", 2)
        gtest.Assert(p.Get("a"), g.Map{
            "b" : g.Slice{1},
            "c" : g.Slice{2},
        })
    })
    gtest.Case(t, func() {
        p := gparser.New(nil)
        p.Set("a", 1)
        err := p.Append("a", 2)
        gtest.AssertNE(err, nil)
        gtest.Assert(p.Get("a"), 1)
    })
}



