// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/encoding/gjson"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)


func Test_Load_JSON(t *testing.T) {
    data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
    // JSON
    gtest.Case(t, func() {
        j, err := gjson.LoadContent(data)
        gtest.Assert(err,        nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
    // JSON
    gtest.Case(t, func() {
        path := "test.json"
        gfile.PutBinContents(path, data)
        defer gfile.Remove(path)
        j, err := gjson.Load(path)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
}

func Test_Load_XML(t *testing.T) {
    data := []byte(`<doc><a>1</a><a>2</a><a>3</a><m><k>v</k></m><n>123456789</n></doc>`)
    // XML
    gtest.Case(t, func() {
        j, err := gjson.LoadContent(data)
        gtest.Assert(err, nil)
        gtest.Assert(j.Get("doc.n"),   "123456789")
        gtest.Assert(j.Get("doc.m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("doc.m.k"), "v")
        gtest.Assert(j.Get("doc.a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("doc.a.1"), 2)
    })
    // XML
    gtest.Case(t, func() {
        path := "test.xml"
        gfile.PutBinContents(path, data)
        defer gfile.Remove(path)
        j, err := gjson.Load(path)
        gtest.Assert(err, nil)
        gtest.Assert(j.Get("doc.n"),   "123456789")
        gtest.Assert(j.Get("doc.m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("doc.m.k"), "v")
        gtest.Assert(j.Get("doc.a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("doc.a.1"), 2)
    })
}

func Test_Load_YAML1(t *testing.T) {
    data := []byte(`
a:
- 1
- 2
- 3
m:
 k: v
"n": 123456789
    `)
    // YAML
    gtest.Case(t, func() {
        j, err := gjson.LoadContent(data)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
    // YAML
    gtest.Case(t, func() {
        path := "test.yaml"
        gfile.PutBinContents(path, data)
        defer gfile.Remove(path)
        j, err := gjson.Load(path)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
}

func Test_Load_YAML2(t *testing.T) {
    data := []byte("i : 123456789")
    gtest.Case(t, func() {
        j, err := gjson.LoadContent(data)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("i"),   "123456789")
    })
}

func Test_Load_TOML1(t *testing.T) {
    data := []byte(`
a = ["1", "2", "3"]
n = "123456789"

[m]
  k = "v"
`)
    // TOML
    gtest.Case(t, func() {
        j, err := gjson.LoadContent(data)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
    // TOML
    gtest.Case(t, func() {
        path := "test.toml"
        gfile.PutBinContents(path, data)
        defer gfile.Remove(path)
        j, err := gjson.Load(path)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("n"),   "123456789")
        gtest.Assert(j.Get("m"),   g.Map{"k" : "v"})
        gtest.Assert(j.Get("m.k"), "v")
        gtest.Assert(j.Get("a"),   g.Slice{1, 2, 3})
        gtest.Assert(j.Get("a.1"), 2)
    })
}

func Test_Load_TOML2(t *testing.T) {
    data := []byte("i=123456789")
    gtest.Case(t, func() {
        j, err := gjson.LoadContent(data)
        gtest.Assert(err,          nil)
        gtest.Assert(j.Get("i"),   "123456789")
    })
}
