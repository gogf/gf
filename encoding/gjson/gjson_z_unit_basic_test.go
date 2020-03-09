// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"github.com/gogf/gf/container/gmap"
	"testing"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
)

func Test_New(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j := gjson.New(data)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
	})

	gtest.Case(t, func() {
		m := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		})
		j := gjson.New(m)
		gtest.Assert(j.Get("k1"), "v1")
		gtest.Assert(j.Get("k2"), "v2")
		gtest.Assert(j.Get("k3"), nil)
	})
}

func Test_Valid(t *testing.T) {
	data1 := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	data2 := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]`)
	gtest.Case(t, func() {
		gtest.Assert(gjson.Valid(data1), true)
		gtest.Assert(gjson.Valid(data2), false)
	})
}

func Test_Encode(t *testing.T) {
	value := g.Slice{1, 2, 3}
	gtest.Case(t, func() {
		b, err := gjson.Encode(value)
		gtest.Assert(err, nil)
		gtest.Assert(b, []byte(`[1,2,3]`))
	})
}

func Test_Decode(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		v, err := gjson.Decode(data)
		gtest.Assert(err, nil)
		gtest.Assert(v, g.Map{
			"n": 123456789,
			"a": g.Slice{1, 2, 3},
			"m": g.Map{
				"k": "v",
			},
		})
	})
	gtest.Case(t, func() {
		var v interface{}
		err := gjson.DecodeTo(data, &v)
		gtest.Assert(err, nil)
		gtest.Assert(v, g.Map{
			"n": 123456789,
			"a": g.Slice{1, 2, 3},
			"m": g.Map{
				"k": "v",
			},
		})
	})
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
}

func Test_SplitChar(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		j.SetSplitChar(byte('#'))
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m#k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a#1"), 2)
	})
}

func Test_ViolenceCheck(t *testing.T) {
	data := []byte(`{"m":{"a":[1,2,3], "v1.v2":"4"}}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("m.a.2"), 3)
		gtest.Assert(j.Get("m.v1.v2"), nil)
		j.SetViolenceCheck(true)
		gtest.Assert(j.Get("m.v1.v2"), 4)
	})
}

func Test_GetVar(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.GetVar("n").String(), "123456789")
		gtest.Assert(j.GetVar("m").Map(), g.Map{"k": "v"})
		gtest.Assert(j.GetVar("a").Interfaces(), g.Slice{1, 2, 3})
		gtest.Assert(j.GetVar("a").Slice(), g.Slice{1, 2, 3})
		gtest.Assert(j.GetVar("a").Array(), g.Slice{1, 2, 3})
	})
}

func Test_GetMap(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.GetMap("n"), nil)
		gtest.Assert(j.GetMap("m"), g.Map{"k": "v"})
		gtest.Assert(j.GetMap("a"), g.Map{"1": "2", "3": nil})
	})
}

func Test_GetJson(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		j2 := j.GetJson("m")
		gtest.AssertNE(j2, nil)
		gtest.Assert(j2.Get("k"), "v")
		gtest.Assert(j2.Get("a"), nil)
		gtest.Assert(j2.Get("n"), nil)
	})
}

func Test_GetArray(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.GetArray("n"), g.Array{123456789})
		gtest.Assert(j.GetArray("m"), g.Array{g.Map{"k": "v"}})
		gtest.Assert(j.GetArray("a"), g.Array{1, 2, 3})
	})
}

func Test_GetString(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.AssertEQ(j.GetString("n"), "123456789")
		gtest.AssertEQ(j.GetString("m"), `{"k":"v"}`)
		gtest.AssertEQ(j.GetString("a"), `[1,2,3]`)
		gtest.AssertEQ(j.GetString("i"), "")
	})
}

func Test_GetStrings(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.AssertEQ(j.GetStrings("n"), g.SliceStr{"123456789"})
		gtest.AssertEQ(j.GetStrings("m"), g.SliceStr{`{"k":"v"}`})
		gtest.AssertEQ(j.GetStrings("a"), g.SliceStr{"1", "2", "3"})
		gtest.AssertEQ(j.GetStrings("i"), nil)
	})
}

func Test_GetInterfaces(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.DecodeToJson(data)
		gtest.Assert(err, nil)
		gtest.AssertEQ(j.GetInterfaces("n"), g.Array{123456789})
		gtest.AssertEQ(j.GetInterfaces("m"), g.Array{g.Map{"k": "v"}})
		gtest.AssertEQ(j.GetInterfaces("a"), g.Array{1, 2, 3})
	})
}

func Test_Len(t *testing.T) {
	gtest.Case(t, func() {
		p := gjson.New(nil)
		p.Append("a", 1)
		p.Append("a", 2)
		gtest.Assert(p.Len("a"), 2)
	})
	gtest.Case(t, func() {
		p := gjson.New(nil)
		p.Append("a.b", 1)
		p.Append("a.c", 2)
		gtest.Assert(p.Len("a"), 2)
	})
	gtest.Case(t, func() {
		p := gjson.New(nil)
		p.Set("a", 1)
		gtest.Assert(p.Len("a"), -1)
	})
}

func Test_Append(t *testing.T) {
	gtest.Case(t, func() {
		p := gjson.New(nil)
		p.Append("a", 1)
		p.Append("a", 2)
		gtest.Assert(p.Get("a"), g.Slice{1, 2})
	})
	gtest.Case(t, func() {
		p := gjson.New(nil)
		p.Append("a.b", 1)
		p.Append("a.c", 2)
		gtest.Assert(p.Get("a"), g.Map{
			"b": g.Slice{1},
			"c": g.Slice{2},
		})
	})
	gtest.Case(t, func() {
		p := gjson.New(nil)
		p.Set("a", 1)
		err := p.Append("a", 2)
		gtest.AssertNE(err, nil)
		gtest.Assert(p.Get("a"), 1)
	})
}

func TestJson_ToJson(t *testing.T) {
	gtest.Case(t, func() {
		p := gjson.New("1")
		s, e := p.ToJsonString()
		gtest.Assert(e, nil)
		gtest.Assert(s, "1")
	})
	gtest.Case(t, func() {
		p := gjson.New("a")
		s, e := p.ToJsonString()
		gtest.Assert(e, nil)
		gtest.Assert(s, `"a"`)
	})
}

func TestJson_Default(t *testing.T) {
	gtest.Case(t, func() {
		j := gjson.New(nil)
		gtest.AssertEQ(j.Get("no", 100), 100)
		gtest.AssertEQ(j.GetString("no", 100), "100")
		gtest.AssertEQ(j.GetBool("no", "on"), true)
		gtest.AssertEQ(j.GetInt("no", 100), 100)
		gtest.AssertEQ(j.GetInt8("no", 100), int8(100))
		gtest.AssertEQ(j.GetInt16("no", 100), int16(100))
		gtest.AssertEQ(j.GetInt32("no", 100), int32(100))
		gtest.AssertEQ(j.GetInt64("no", 100), int64(100))
		gtest.AssertEQ(j.GetUint("no", 100), uint(100))
		gtest.AssertEQ(j.GetUint8("no", 100), uint8(100))
		gtest.AssertEQ(j.GetUint16("no", 100), uint16(100))
		gtest.AssertEQ(j.GetUint32("no", 100), uint32(100))
		gtest.AssertEQ(j.GetUint64("no", 100), uint64(100))
		gtest.AssertEQ(j.GetFloat32("no", 123.456), float32(123.456))
		gtest.AssertEQ(j.GetFloat64("no", 123.456), float64(123.456))
		gtest.AssertEQ(j.GetArray("no", g.Slice{1, 2, 3}), g.Slice{1, 2, 3})
		gtest.AssertEQ(j.GetInts("no", g.Slice{1, 2, 3}), g.SliceInt{1, 2, 3})
		gtest.AssertEQ(j.GetFloats("no", g.Slice{1, 2, 3}), []float64{1, 2, 3})
		gtest.AssertEQ(j.GetMap("no", g.Map{"k": "v"}), g.Map{"k": "v"})
		gtest.AssertEQ(j.GetVar("no", 123.456).Float64(), float64(123.456))
		gtest.AssertEQ(j.GetJson("no", g.Map{"k": "v"}).Get("k"), "v")
		gtest.AssertEQ(j.GetJsons("no", g.Slice{
			g.Map{"k1": "v1"},
			g.Map{"k2": "v2"},
			g.Map{"k3": "v3"},
		})[0].Get("k1"), "v1")
		gtest.AssertEQ(j.GetJsonMap("no", g.Map{
			"m1": g.Map{"k1": "v1"},
			"m2": g.Map{"k2": "v2"},
		})["m2"].Get("k2"), "v2")
	})
}

func Test_Convert(t *testing.T) {
	gtest.Case(t, func() {
		j := gjson.New(`{"name":"gf"}`)
		arr, err := j.ToXml()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "<name>gf</name>")
		arr, err = j.ToXmlIndent()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "<name>gf</name>")
		str, err := j.ToXmlString()
		gtest.Assert(err, nil)
		gtest.Assert(str, "<name>gf</name>")
		str, err = j.ToXmlIndentString()
		gtest.Assert(err, nil)
		gtest.Assert(str, "<name>gf</name>")

		arr, err = j.ToJsonIndent()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "{\n\t\"name\": \"gf\"\n}")
		str, err = j.ToJsonIndentString()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "{\n\t\"name\": \"gf\"\n}")

		arr, err = j.ToYaml()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "name: gf\n")
		str, err = j.ToYamlString()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "name: gf\n")

		arr, err = j.ToToml()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "name = \"gf\"\n")
		str, err = j.ToTomlString()
		gtest.Assert(err, nil)
		gtest.Assert(string(arr), "name = \"gf\"\n")
	})
}

func Test_Convert2(t *testing.T) {
	gtest.Case(t, func() {
		name := struct {
			Name string
		}{}
		j := gjson.New(`{"name":"gf","time":"2019-06-12"}`)
		gtest.Assert(j.Value().(g.Map)["name"], "gf")
		gtest.Assert(j.GetMap("name1"), nil)
		gtest.AssertNE(j.GetJson("name1"), nil)
		gtest.Assert(j.GetJsons("name1"), nil)
		gtest.Assert(j.GetJsonMap("name1"), nil)
		gtest.Assert(j.Contains("name1"), false)
		gtest.Assert(j.GetVar("name1").IsNil(), true)
		gtest.Assert(j.GetVar("name").IsNil(), false)
		gtest.Assert(j.Len("name1"), -1)
		gtest.Assert(j.GetTime("time").Format("2006-01-02"), "2019-06-12")
		gtest.Assert(j.GetGTime("time").Format("Y-m-d"), "2019-06-12")
		gtest.Assert(j.GetDuration("time").String(), "0s")

		err := j.ToStruct(&name)
		gtest.Assert(err, nil)
		gtest.Assert(name.Name, "gf")
		//j.Dump()
		gtest.Assert(err, nil)

		j = gjson.New(`{"person":{"name":"gf"}}`)
		err = j.GetStruct("person", &name)
		gtest.Assert(err, nil)
		gtest.Assert(name.Name, "gf")

		j = gjson.New(`{"name":"gf""}`)
		//j.Dump()
		gtest.Assert(err, nil)

		j = gjson.New(`[1,2,3]`)
		gtest.Assert(len(j.ToArray()), 3)
	})
}

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		j := gjson.New(`{"name":"gf","time":"2019-06-12"}`)
		j.SetViolenceCheck(true)
		gtest.Assert(j.Get(""), nil)
		gtest.Assert(j.Get(".").(g.Map)["name"], "gf")
		gtest.Assert(j.Get(".").(g.Map)["name1"], nil)
		j.SetViolenceCheck(false)
		gtest.Assert(j.Get(".").(g.Map)["name"], "gf")

		err := j.Set("name", "gf1")
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("name"), "gf1")

		j = gjson.New(`[1,2,3]`)
		err = j.Set("\"0\".1", 11)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("1"), 11)

		j = gjson.New(`[1,2,3]`)
		err = j.Set("11111111111111111111111", 11)
		gtest.AssertNE(err, nil)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("1")
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("0"), 1)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("3")
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("0"), 1)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("0.3")
		gtest.Assert(err, nil)
		gtest.Assert(len(j.Get("0").([]interface{})), 3)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("0.a")
		gtest.Assert(err, nil)
		gtest.Assert(len(j.Get("0").(g.Map)), 0)

		name := struct {
			Name string
		}{Name: "gf"}
		j = gjson.New(name)
		gtest.Assert(j.Get("Name"), "gf")
		err = j.Remove("Name")
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name"), nil)

		err = j.Set("Name", "gf1")
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name"), "gf1")

		j = gjson.New(nil)
		err = j.Remove("Name")
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name"), nil)

		j = gjson.New(name)
		gtest.Assert(j.Get("Name"), "gf")
		err = j.Set("Name1", g.Map{"Name": "gf1"})
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name1").(g.Map)["Name"], "gf1")
		err = j.Set("Name2", g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name2").(g.Slice)[0], 1)
		err = j.Set("Name3", name)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name3").(g.Map)["Name"], "gf")
		err = j.Set("Name4", &name)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name4").(g.Map)["Name"], "gf")
		arr := [3]int{1, 2, 3}
		err = j.Set("Name5", arr)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Name5").(g.Array)[0], 1)

	})
}

func Test_IsNil(t *testing.T) {
	gtest.Case(t, func() {
		j := gjson.New(nil)
		gtest.Assert(j.IsNil(), true)
	})
}

func Test_ToStructDeep(t *testing.T) {
	gtest.Case(t, func() {
		type Item struct {
			Title string `json:"title"`
			Key   string `json:"key"`
		}

		type M struct {
			Id    string                 `json:"id"`
			Me    map[string]interface{} `json:"me"`
			Txt   string                 `json:"txt"`
			Items []*Item                `json:"items"`
		}

		txt := `{
		  "id":"88888",
		  "me":{"name":"mikey","day":"20009"},
		  "txt":"hello",
		  "items":null
		 }`

		j, err := gjson.LoadContent(txt)
		gtest.Assert(err, nil)
		gtest.Assert(j.GetString("me.name"), "mikey")
		gtest.Assert(j.GetString("items"), "")
		gtest.Assert(j.GetBool("items"), false)
		gtest.Assert(j.GetArray("items"), nil)
		m := new(M)
		err = j.ToStructDeep(m)
		gtest.Assert(err, nil)
		gtest.AssertNE(m.Me, nil)
		gtest.Assert(m.Me["day"], "20009")
		gtest.Assert(m.Items, nil)
	})
}
