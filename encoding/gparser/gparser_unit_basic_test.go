// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gparser_test

import (
	"testing"

	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
)

func Test_New(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.Assert(j.Get("n"), "123456789")
		t.Assert(j.Get("m"), g.Map{"k": "v"})
		t.Assert(j.Get("a"), g.Slice{1, 2, 3})
		v := j.Value().(g.Map)
		t.Assert(v["n"], 123456789)
	})
}

func Test_NewUnsafe(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.Assert(j.Get("n"), "123456789")
		t.Assert(j.Get("m"), g.Map{"k": "v"})
		t.Assert(j.Get("m.k"), "v")
		t.Assert(j.Get("a"), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1"), 2)
	})
}

func Test_Encode(t *testing.T) {
	value := g.Slice{1, 2, 3}
	gtest.C(t, func(t *gtest.T) {
		b, err := gparser.VarToJson(value)
		t.Assert(err, nil)
		t.Assert(b, []byte(`[1,2,3]`))
	})
}

func Test_Decode(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.Assert(j.Get("n"), "123456789")
		t.Assert(j.Get("m"), g.Map{"k": "v"})
		t.Assert(j.Get("m.k"), "v")
		t.Assert(j.Get("a"), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1"), 2)
	})
}

func Test_SplitChar(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		j.SetSplitChar(byte('#'))
		t.AssertNE(j, nil)
		t.Assert(j.Get("n"), "123456789")
		t.Assert(j.Get("m"), g.Map{"k": "v"})
		t.Assert(j.Get("m#k"), "v")
		t.Assert(j.Get("a"), g.Slice{1, 2, 3})
		t.Assert(j.Get("a#1"), 2)
	})
}

func Test_ViolenceCheck(t *testing.T) {
	data := []byte(`{"m":{"a":[1,2,3], "v1.v2":"4"}}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.Assert(j.Get("m.a.2"), 3)
		t.Assert(j.Get("m.v1.v2"), nil)
		j.SetViolenceCheck(true)
		t.Assert(j.Get("m.v1.v2"), 4)
	})
}

func Test_GetVar(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.Assert(j.GetVar("n").String(), "123456789")
		t.Assert(j.GetVar("m").Map(), g.Map{"k": "v"})
		t.Assert(j.GetVar("a").Interfaces(), g.Slice{1, 2, 3})
		t.Assert(j.GetVar("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.GetMap("a"), g.Map{"1": "2", "3": nil})
	})
}

func Test_GetMap(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.Assert(j.GetMap("n"), nil)
		t.Assert(j.GetMap("m"), g.Map{"k": "v"})
		t.Assert(j.GetMap("a"), g.Map{"1": "2", "3": nil})
	})
}

func Test_GetArray(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.Assert(j.GetArray("n"), g.Array{123456789})
		t.Assert(j.GetArray("m"), g.Array{"k", "v"})
		t.Assert(j.GetArray("a"), g.Array{1, 2, 3})
	})
}

func Test_GetString(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.AssertEQ(j.GetString("n"), "123456789")
		t.AssertEQ(j.GetString("m"), `{"k":"v"}`)
		t.AssertEQ(j.GetString("a"), `[1,2,3]`)
		t.AssertEQ(j.GetString("i"), "")
	})
}

func Test_GetStrings(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.AssertEQ(j.GetStrings("n"), g.SliceStr{"123456789"})
		t.AssertEQ(j.GetStrings("m"), g.SliceStr{`{"k":"v"}`})
		t.AssertEQ(j.GetStrings("a"), g.SliceStr{"1", "2", "3"})
		t.AssertEQ(j.GetStrings("i"), nil)
	})
}

func Test_GetInterfaces(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.AssertEQ(j.GetInterfaces("n"), g.Array{123456789})
		t.AssertEQ(j.GetInterfaces("m"), g.Array{"k", "v"})
		t.AssertEQ(j.GetInterfaces("a"), g.Array{1, 2, 3})
	})
}

func Test_Len(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(nil)
		p.Append("a", 1)
		p.Append("a", 2)
		t.Assert(p.Len("a"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(nil)
		p.Append("a.b", 1)
		p.Append("a.c", 2)
		t.Assert(p.Len("a"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(nil)
		p.Set("a", 1)
		t.Assert(p.Len("a"), -1)
	})
}

func Test_Append(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(nil)
		p.Append("a", 1)
		p.Append("a", 2)
		t.Assert(p.Get("a"), g.Slice{1, 2})
	})
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(nil)
		p.Append("a.b", 1)
		p.Append("a.c", 2)
		t.Assert(p.Get("a"), g.Map{
			"b": g.Slice{1},
			"c": g.Slice{2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(nil)
		p.Set("a", 1)
		err := p.Append("a", 2)
		t.AssertNE(err, nil)
		t.Assert(p.Get("a"), 1)
	})
}

func Test_Convert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := gparser.New(`{"name":"gf","bool":true,"int":1,"float":1,"ints":[1,2],"floats":[1,2],"time":"2019-06-12","person": {"name": "gf"}}`)
		t.Assert(p.GetVar("name").String(), "gf")
		t.Assert(p.GetString("name"), "gf")
		t.Assert(p.GetBool("bool"), true)
		t.Assert(p.GetInt("int"), 1)
		t.Assert(p.GetInt8("int"), 1)
		t.Assert(p.GetInt16("int"), 1)
		t.Assert(p.GetInt32("int"), 1)
		t.Assert(p.GetInt64("int"), 1)
		t.Assert(p.GetUint("int"), 1)
		t.Assert(p.GetUint8("int"), 1)
		t.Assert(p.GetUint16("int"), 1)
		t.Assert(p.GetUint32("int"), 1)
		t.Assert(p.GetUint64("int"), 1)
		t.Assert(p.GetInts("ints")[0], 1)
		t.Assert(p.GetFloat32("float"), 1)
		t.Assert(p.GetFloat64("float"), 1)
		t.Assert(p.GetFloats("floats")[0], 1)
		t.Assert(p.GetTime("time").Format("2006-01-02"), "2019-06-12")
		t.Assert(p.GetGTime("time").Format("Y-m-d"), "2019-06-12")
		t.Assert(p.GetDuration("time").String(), "0s")
		name := struct {
			Name string
		}{}
		err := p.GetStruct("person", &name)
		t.Assert(err, nil)
		t.Assert(name.Name, "gf")
		t.Assert(p.ToMap()["name"], "gf")
		err = p.ToStruct(&name)
		t.Assert(err, nil)
		t.Assert(name.Name, "gf")
		//p.Dump()

		p = gparser.New(`[0,1,2]`)
		t.Assert(p.ToArray()[0], 0)
	})
}

func Test_Convert2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		xmlArr := []byte{60, 114, 111, 111, 116, 47, 62}
		p := gparser.New(`<root></root>`)
		arr, err := p.ToXml("root")
		t.Assert(err, nil)
		t.Assert(arr, xmlArr)
		arr, err = gparser.VarToXml(`<root></root>`, "root")
		t.Assert(err, nil)
		t.Assert(arr, xmlArr)

		arr, err = p.ToXmlIndent("root")
		t.Assert(err, nil)
		t.Assert(arr, xmlArr)
		arr, err = gparser.VarToXmlIndent(`<root></root>`, "root")
		t.Assert(err, nil)
		t.Assert(arr, xmlArr)

		p = gparser.New(`{"name":"gf"}`)
		str, err := p.ToJsonString()
		t.Assert(err, nil)
		t.Assert(str, `{"name":"gf"}`)
		str, err = gparser.VarToJsonString(`{"name":"gf"}`)
		t.Assert(err, nil)
		t.Assert(str, `{"name":"gf"}`)

		jsonIndentArr := []byte{123, 10, 9, 34, 110, 97, 109, 101, 34, 58, 32, 34, 103, 102, 34, 10, 125}
		arr, err = p.ToJsonIndent()
		t.Assert(err, nil)
		t.Assert(arr, jsonIndentArr)
		arr, err = gparser.VarToJsonIndent(`{"name":"gf"}`)
		t.Assert(err, nil)
		t.Assert(arr, jsonIndentArr)

		str, err = p.ToJsonIndentString()
		t.Assert(err, nil)
		t.Assert(str, "{\n\t\"name\": \"gf\"\n}")
		str, err = gparser.VarToJsonIndentString(`{"name":"gf"}`)
		t.Assert(err, nil)
		t.Assert(str, "{\n\t\"name\": \"gf\"\n}")

		p = gparser.New(g.Map{"name": "gf"})
		arr, err = p.ToYaml()
		t.Assert(err, nil)
		t.Assert(arr, "name: gf\n")
		arr, err = gparser.VarToYaml(g.Map{"name": "gf"})
		t.Assert(err, nil)
		t.Assert(arr, "name: gf\n")

		tomlArr := []byte{110, 97, 109, 101, 32, 61, 32, 34, 103, 102, 34, 10}
		p = gparser.New(`
name= "gf"
`)
		arr, err = p.ToToml()
		t.Assert(err, nil)
		t.Assert(arr, tomlArr)
		arr, err = gparser.VarToToml(`
name= "gf"
`)
		t.Assert(err, nil)
		t.Assert(arr, tomlArr)
	})
}
