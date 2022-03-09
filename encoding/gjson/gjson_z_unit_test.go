// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_New(t *testing.T) {
	// New with json map.
	gtest.C(t, func(t *gtest.T) {
		data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
		j := gjson.New(data)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("a").Array(), g.Slice{1, 2, 3})
	})
	// New with json array map.
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(`[{"a":1},{"b":2},{"c":3}]`)
		t.Assert(j.Get(".").String(), `[{"a":1},{"b":2},{"c":3}]`)
		t.Assert(j.Get("2.c").String(), `3`)
	})
	// New with gvar.
	// https://github.com/gogf/gf/issues/1571
	gtest.C(t, func(t *gtest.T) {
		v := gvar.New(`[{"a":1},{"b":2},{"c":3}]`)
		j := gjson.New(v)
		t.Assert(j.Get(".").String(), `[{"a":1},{"b":2},{"c":3}]`)
		t.Assert(j.Get("2.c").String(), `3`)
	})
	// New with gmap.
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		})
		j := gjson.New(m)
		t.Assert(j.Get("k1"), "v1")
		t.Assert(j.Get("k2"), "v2")
		t.Assert(j.Get("k3"), nil)
	})
}

func Test_Valid(t *testing.T) {
	data1 := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	data2 := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]`)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gjson.Valid(data1), true)
		t.Assert(gjson.Valid(data2), false)
	})
}

func Test_Encode(t *testing.T) {
	value := g.Slice{1, 2, 3}
	gtest.C(t, func(t *gtest.T) {
		b, err := gjson.Encode(value)
		t.Assert(err, nil)
		t.Assert(b, []byte(`[1,2,3]`))
	})
}

func Test_Decode(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		v, err := gjson.Decode(data)
		t.Assert(err, nil)
		t.Assert(v, g.Map{
			"n": 123456789,
			"a": g.Slice{1, 2, 3},
			"m": g.Map{
				"k": "v",
			},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		var v interface{}
		err := gjson.DecodeTo(data, &v)
		t.Assert(err, nil)
		t.Assert(v, g.Map{
			"n": 123456789,
			"a": g.Slice{1, 2, 3},
			"m": g.Map{
				"k": "v",
			},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k"), "v")
		t.Assert(j.Get("a").Array(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1").Int(), 2)
	})
}

func Test_SplitChar(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		j.SetSplitChar(byte('#'))
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m#k").String(), "v")
		t.Assert(j.Get("a").Array(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a#1").Int(), 2)
	})
}

func Test_ViolenceCheck(t *testing.T) {
	data := []byte(`{"m":{"a":[1,2,3], "v1.v2":"4"}}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.Assert(j.Get("m.a.2"), 3)
		t.Assert(j.Get("m.v1.v2"), nil)
		j.SetViolenceCheck(true)
		t.Assert(j.Get("m.v1.v2"), 4)
	})
}

func Test_GetVar(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("a").Interfaces(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a").Array(), g.Slice{1, 2, 3})
	})
}

func Test_GetMap(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").Map(), nil)
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("a").Map(), g.Map{"1": "2", "3": nil})
	})
}

func Test_GetJson(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		j2 := j.GetJson("m")
		t.AssertNE(j2, nil)
		t.Assert(j2.Get("k"), "v")
		t.Assert(j2.Get("a"), nil)
		t.Assert(j2.Get("n"), nil)
	})
}

func Test_GetArray(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").Array(), g.Array{123456789})
		t.Assert(j.Get("m").Array(), g.Array{g.Map{"k": "v"}})
		t.Assert(j.Get("a").Array(), g.Array{1, 2, 3})
	})
}

func Test_GetString(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.AssertEQ(j.Get("n").String(), "123456789")
		t.AssertEQ(j.Get("m").String(), `{"k":"v"}`)
		t.AssertEQ(j.Get("a").String(), `[1,2,3]`)
		t.AssertEQ(j.Get("i").String(), "")
	})
}

func Test_GetStrings(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.AssertEQ(j.Get("n").Strings(), g.SliceStr{"123456789"})
		t.AssertEQ(j.Get("m").Strings(), g.SliceStr{`{"k":"v"}`})
		t.AssertEQ(j.Get("a").Strings(), g.SliceStr{"1", "2", "3"})
		t.AssertEQ(j.Get("i").Strings(), nil)
	})
}

func Test_GetInterfaces(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.DecodeToJson(data)
		t.Assert(err, nil)
		t.AssertEQ(j.Get("n").Interfaces(), g.Array{123456789})
		t.AssertEQ(j.Get("m").Interfaces(), g.Array{g.Map{"k": "v"}})
		t.AssertEQ(j.Get("a").Interfaces(), g.Array{1, 2, 3})
	})
}

func Test_Len(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(nil)
		p.Append("a", 1)
		p.Append("a", 2)
		t.Assert(p.Len("a"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(nil)
		p.Append("a.b", 1)
		p.Append("a.c", 2)
		t.Assert(p.Len("a"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(nil)
		p.Set("a", 1)
		t.Assert(p.Len("a"), -1)
	})
}

func Test_Append(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(nil)
		p.Append("a", 1)
		p.Append("a", 2)
		t.Assert(p.Get("a"), g.Slice{1, 2})
	})
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(nil)
		p.Append("a.b", 1)
		p.Append("a.c", 2)
		t.Assert(p.Get("a").Map(), g.Map{
			"b": g.Slice{1},
			"c": g.Slice{2},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(nil)
		p.Set("a", 1)
		err := p.Append("a", 2)
		t.AssertNE(err, nil)
		t.Assert(p.Get("a"), 1)
	})
}

func Test_RawArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(nil)
		t.AssertNil(j.Set("0", 1))
		t.AssertNil(j.Set("1", 2))
		t.Assert(j.MustToJsonString(), `[1,2]`)
	})

	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(nil)
		t.AssertNil(j.Append(".", 1))
		t.AssertNil(j.Append(".", 2))
		t.Assert(j.MustToJsonString(), `[1,2]`)
	})
}

func TestJson_ToJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New(1)
		s, e := p.ToJsonString()
		t.Assert(e, nil)
		t.Assert(s, "1")
	})
	gtest.C(t, func(t *gtest.T) {
		p := gjson.New("a")
		s, e := p.ToJsonString()
		t.Assert(e, nil)
		t.Assert(s, `"a"`)
	})
}

func TestJson_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(nil)
		t.AssertEQ(j.Get("no", 100).Int(), 100)
		t.AssertEQ(j.Get("no", 100).String(), "100")
		t.AssertEQ(j.Get("no", "on").Bool(), true)
		t.AssertEQ(j.Get("no", 100).Int(), 100)
		t.AssertEQ(j.Get("no", 100).Int8(), int8(100))
		t.AssertEQ(j.Get("no", 100).Int16(), int16(100))
		t.AssertEQ(j.Get("no", 100).Int32(), int32(100))
		t.AssertEQ(j.Get("no", 100).Int64(), int64(100))
		t.AssertEQ(j.Get("no", 100).Uint(), uint(100))
		t.AssertEQ(j.Get("no", 100).Uint8(), uint8(100))
		t.AssertEQ(j.Get("no", 100).Uint16(), uint16(100))
		t.AssertEQ(j.Get("no", 100).Uint32(), uint32(100))
		t.AssertEQ(j.Get("no", 100).Uint64(), uint64(100))
		t.AssertEQ(j.Get("no", 123.456).Float32(), float32(123.456))
		t.AssertEQ(j.Get("no", 123.456).Float64(), float64(123.456))
		t.AssertEQ(j.Get("no", g.Slice{1, 2, 3}).Array(), g.Slice{1, 2, 3})
		t.AssertEQ(j.Get("no", g.Slice{1, 2, 3}).Ints(), g.SliceInt{1, 2, 3})
		t.AssertEQ(j.Get("no", g.Slice{1, 2, 3}).Floats(), []float64{1, 2, 3})
		t.AssertEQ(j.Get("no", g.Map{"k": "v"}).Map(), g.Map{"k": "v"})
		t.AssertEQ(j.Get("no", 123.456).Float64(), float64(123.456))
		t.AssertEQ(j.GetJson("no", g.Map{"k": "v"}).Get("k").String(), "v")
		t.AssertEQ(j.GetJsons("no", g.Slice{
			g.Map{"k1": "v1"},
			g.Map{"k2": "v2"},
			g.Map{"k3": "v3"},
		})[0].Get("k1").String(), "v1")
		t.AssertEQ(j.GetJsonMap("no", g.Map{
			"m1": g.Map{"k1": "v1"},
			"m2": g.Map{"k2": "v2"},
		})["m2"].Get("k2").String(), "v2")
	})
}

func Test_Convert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(`{"name":"gf"}`)
		arr, err := j.ToXml()
		t.Assert(err, nil)
		t.Assert(string(arr), "<name>gf</name>")
		arr, err = j.ToXmlIndent()
		t.Assert(err, nil)
		t.Assert(string(arr), "<name>gf</name>")
		str, err := j.ToXmlString()
		t.Assert(err, nil)
		t.Assert(str, "<name>gf</name>")
		str, err = j.ToXmlIndentString()
		t.Assert(err, nil)
		t.Assert(str, "<name>gf</name>")

		arr, err = j.ToJsonIndent()
		t.Assert(err, nil)
		t.Assert(string(arr), "{\n\t\"name\": \"gf\"\n}")
		str, err = j.ToJsonIndentString()
		t.Assert(err, nil)
		t.Assert(string(arr), "{\n\t\"name\": \"gf\"\n}")

		arr, err = j.ToYaml()
		t.Assert(err, nil)
		t.Assert(string(arr), "name: gf\n")
		str, err = j.ToYamlString()
		t.Assert(err, nil)
		t.Assert(string(arr), "name: gf\n")

		arr, err = j.ToToml()
		t.Assert(err, nil)
		t.Assert(string(arr), "name = \"gf\"\n")
		str, err = j.ToTomlString()
		t.Assert(err, nil)
		t.Assert(string(arr), "name = \"gf\"\n")
	})
}

func Test_Convert2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		name := struct {
			Name string
		}{}
		j := gjson.New(`{"name":"gf","time":"2019-06-12"}`)
		t.Assert(j.Interface().(g.Map)["name"], "gf")
		t.Assert(j.Get("name1").Map(), nil)
		t.AssertNE(j.GetJson("name1"), nil)
		t.Assert(j.GetJsons("name1"), nil)
		t.Assert(j.GetJsonMap("name1"), nil)
		t.Assert(j.Contains("name1"), false)
		t.Assert(j.Get("name1").IsNil(), true)
		t.Assert(j.Get("name").IsNil(), false)
		t.Assert(j.Len("name1"), -1)
		t.Assert(j.Get("time").Time().Format("2006-01-02"), "2019-06-12")
		t.Assert(j.Get("time").GTime().Format("Y-m-d"), "2019-06-12")
		t.Assert(j.Get("time").Duration().String(), "0s")

		err := j.Var().Scan(&name)
		t.Assert(err, nil)
		t.Assert(name.Name, "gf")
		// j.Dump()
		t.Assert(err, nil)

		j = gjson.New(`{"person":{"name":"gf"}}`)
		err = j.Get("person").Scan(&name)
		t.Assert(err, nil)
		t.Assert(name.Name, "gf")

		j = gjson.New(`{"name":"gf""}`)
		// j.Dump()
		t.Assert(err, nil)

		j = gjson.New(`[1,2,3]`)
		t.Assert(len(j.Var().Array()), 3)
	})
}

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(`{"name":"gf","time":"2019-06-12"}`)
		j.SetViolenceCheck(true)
		t.Assert(j.Get(""), nil)
		t.Assert(j.Get(".").Interface().(g.Map)["name"], "gf")
		t.Assert(j.Get(".").Interface().(g.Map)["name1"], nil)
		j.SetViolenceCheck(false)
		t.Assert(j.Get(".").Interface().(g.Map)["name"], "gf")

		err := j.Set("name", "gf1")
		t.Assert(err, nil)
		t.Assert(j.Get("name"), "gf1")

		j = gjson.New(`[1,2,3]`)
		err = j.Set("\"0\".1", 11)
		t.Assert(err, nil)
		t.Assert(j.Get("1"), 11)

		j = gjson.New(`[1,2,3]`)
		err = j.Set("11111111111111111111111", 11)
		t.AssertNE(err, nil)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("1")
		t.Assert(err, nil)
		t.Assert(j.Get("0"), 1)
		t.Assert(len(j.Var().Array()), 2)

		j = gjson.New(`[1,2,3]`)
		// If index 0 is delete, its next item will be at index 0.
		t.Assert(j.Remove("0"), nil)
		t.Assert(j.Remove("0"), nil)
		t.Assert(j.Remove("0"), nil)
		t.Assert(j.Get("0"), nil)
		t.Assert(len(j.Var().Array()), 0)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("3")
		t.Assert(err, nil)
		t.Assert(j.Get("0"), 1)
		t.Assert(len(j.Var().Array()), 3)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("0.3")
		t.Assert(err, nil)
		t.Assert(j.Get("0"), 1)

		j = gjson.New(`[1,2,3]`)
		err = j.Remove("0.a")
		t.Assert(err, nil)
		t.Assert(j.Get("0"), 1)

		name := struct {
			Name string
		}{Name: "gf"}
		j = gjson.New(name)
		t.Assert(j.Get("Name"), "gf")
		err = j.Remove("Name")
		t.Assert(err, nil)
		t.Assert(j.Get("Name"), nil)

		err = j.Set("Name", "gf1")
		t.Assert(err, nil)
		t.Assert(j.Get("Name"), "gf1")

		j = gjson.New(nil)
		err = j.Remove("Name")
		t.Assert(err, nil)
		t.Assert(j.Get("Name"), nil)

		j = gjson.New(name)
		t.Assert(j.Get("Name"), "gf")
		err = j.Set("Name1", g.Map{"Name": "gf1"})
		t.Assert(err, nil)
		t.Assert(j.Get("Name1").Interface().(g.Map)["Name"], "gf1")
		err = j.Set("Name2", g.Slice{1, 2, 3})
		t.Assert(err, nil)
		t.Assert(j.Get("Name2").Interface().(g.Slice)[0], 1)
		err = j.Set("Name3", name)
		t.Assert(err, nil)
		t.Assert(j.Get("Name3").Interface().(g.Map)["Name"], "gf")
		err = j.Set("Name4", &name)
		t.Assert(err, nil)
		t.Assert(j.Get("Name4").Interface().(g.Map)["Name"], "gf")
		arr := [3]int{1, 2, 3}
		err = j.Set("Name5", arr)
		t.Assert(err, nil)
		t.Assert(j.Get("Name5").Interface().(g.Array)[0], 1)

	})
}

func TestJson_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := []byte("[9223372036854775807, 9223372036854775806]")
		array := gjson.New(data).Var().Array()
		t.Assert(array, []uint64{9223372036854776000, 9223372036854776000})
	})
	gtest.C(t, func(t *gtest.T) {
		data := []byte("[9223372036854775807, 9223372036854775806]")
		array := gjson.NewWithOptions(data, gjson.Options{StrNumber: true}).Var().Array()
		t.Assert(array, []uint64{9223372036854775807, 9223372036854775806})
	})
}

func TestJson_IsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(nil)
		t.Assert(j.IsNil(), true)
	})
}

func TestJson_Set_With_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gjson.New(g.Map{
			"user1": g.Map{"name": "user1"},
			"user2": g.Map{"name": "user2"},
			"user3": g.Map{"name": "user3"},
		})
		user1 := v.GetJson("user1")
		t.AssertNil(user1.Set("id", 111))
		t.AssertNil(v.Set("user1", user1))
		t.Assert(v.Get("user1.id"), 111)
	})
}

func TestJson_Options(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S struct {
			Id int64
		}
		s := S{
			Id: 53687091200,
		}
		m := make(map[string]interface{})
		t.AssertNil(gjson.DecodeTo(gjson.MustEncode(s), &m, gjson.Options{
			StrNumber: false,
		}))
		t.Assert(fmt.Sprintf(`%v`, m["Id"]), `5.36870912e+10`)
		t.AssertNil(gjson.DecodeTo(gjson.MustEncode(s), &m, gjson.Options{
			StrNumber: true,
		}))
		t.Assert(fmt.Sprintf(`%v`, m["Id"]), `53687091200`)
	})
}

// https://github.com/gogf/gf/issues/1617
func Test_Issue1617(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type MyJsonName struct {
			F中文   int64 `json:"F中文"`
			F英文   int64 `json:"F英文"`
			F法文   int64 `json:"F法文"`
			F西班牙语 int64 `json:"F西班牙语"`
		}
		jso := `{"F中文":1,"F英文":2,"F法文":3,"F西班牙语":4}`
		var a MyJsonName
		json, err := gjson.DecodeToJson(jso)
		t.AssertNil(err)
		err = json.Scan(&a)
		t.AssertNil(err)
		t.Assert(a, MyJsonName{
			F中文:   1,
			F英文:   2,
			F法文:   3,
			F西班牙语: 4,
		})
	})
}
