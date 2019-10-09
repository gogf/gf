// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"testing"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_Load_JSON1(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	// JSON
	gtest.Case(t, func() {
		j, err := gjson.LoadContent(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
	// JSON
	gtest.Case(t, func() {
		path := "test.json"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
}

func Test_Load_JSON2(t *testing.T) {
	data := []byte(`{"n":123456789000000000000, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j, err := gjson.LoadContent(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789000000000000")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
}

func Test_Load_XML(t *testing.T) {
	data := []byte(`<doc><a>1</a><a>2</a><a>3</a><m><k>v</k></m><n>123456789</n></doc>`)
	// XML
	gtest.Case(t, func() {
		j, err := gjson.LoadContent(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("doc.n"), "123456789")
		gtest.Assert(j.Get("doc.m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("doc.m.k"), "v")
		gtest.Assert(j.Get("doc.a"), g.Slice{"1", "2", "3"})
		gtest.Assert(j.Get("doc.a.1"), 2)
	})
	// XML
	gtest.Case(t, func() {
		path := "test.xml"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("doc.n"), "123456789")
		gtest.Assert(j.Get("doc.m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("doc.m.k"), "v")
		gtest.Assert(j.Get("doc.a"), g.Slice{"1", "2", "3"})
		gtest.Assert(j.Get("doc.a.1"), 2)
	})

	// XML
	gtest.Case(t, func() {
		xml := `<?xml version="1.0"?>

	<Output type="o">
	<itotalSize>0</itotalSize>
	<ipageSize>1</ipageSize>
	<ipageIndex>2</ipageIndex>
	<itotalRecords>GF框架</itotalRecords>
	<nworkOrderDtos/>
	<nworkOrderFrontXML/>
	</Output>`
		j, err := gjson.LoadContent(xml)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("Output.ipageIndex"), "2")
		gtest.Assert(j.Get("Output.itotalRecords"), "GF框架")
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
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
	// YAML
	gtest.Case(t, func() {
		path := "test.yaml"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
}

func Test_Load_YAML2(t *testing.T) {
	data := []byte("i : 123456789")
	gtest.Case(t, func() {
		j, err := gjson.LoadContent(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("i"), "123456789")
	})
}

func Test_Load_TOML1(t *testing.T) {
	data := []byte(`
a = ["1", "2", "3"]
n = 123456789

[m]
  k = "v"
`)
	// TOML
	gtest.Case(t, func() {
		j, err := gjson.LoadContent(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{"1", "2", "3"})
		gtest.Assert(j.Get("a.1"), 2)
	})
	// TOML
	gtest.Case(t, func() {
		path := "test.toml"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{"1", "2", "3"})
		gtest.Assert(j.Get("a.1"), 2)
	})
}

func Test_Load_TOML2(t *testing.T) {
	data := []byte("i=123456789")
	gtest.Case(t, func() {
		j, err := gjson.LoadContent(data)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("i"), "123456789")
	})
}

func Test_Load_Basic(t *testing.T) {
	gtest.Case(t, func() {
		j := gjson.New(nil)
		gtest.Assert(j.Value(), nil)
		_, err := gjson.Decode(nil)
		gtest.AssertNE(err, nil)
		_, err = gjson.DecodeToJson(nil)
		gtest.AssertNE(err, nil)
		j, err = gjson.LoadContent(nil)
		gtest.Assert(err, nil)
		gtest.Assert(j.Value(), nil)

		j, err = gjson.LoadContent(`{"name": "gf"}`)
		gtest.Assert(err, nil)

		j, err = gjson.LoadContent(`{"name": "gf"""}`)
		gtest.AssertNE(err, nil)

		j = gjson.New(&g.Map{"name": "gf"})
		gtest.Assert(j.GetString("name"), "gf")

	})
}

func Test_Load_Ini(t *testing.T) {
	var data = `

;注释

[addr] 
#注释
ip = 127.0.0.1
port=9001
enable=true

	[DBINFO]
	type=mysql
	user=root
	password=password

`

	gtest.Case(t, func() {
		json, err := gjson.LoadContent(data)
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(json.GetString("addr.ip"), "127.0.0.1")
		gtest.Assert(json.GetString("addr.port"), "9001")
		gtest.Assert(json.GetString("addr.enable"), "true")
		gtest.Assert(json.GetString("DBINFO.type"), "mysql")
		gtest.Assert(json.GetString("DBINFO.user"), "root")
		gtest.Assert(json.GetString("DBINFO.password"), "password")

		_, err = json.ToIni()
		if err != nil {
			gtest.Fatal(err)
		}
	})
}
