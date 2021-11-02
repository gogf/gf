// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Load_JSON1(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	// JSON
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1").Int(), 2)
	})
	// JSON
	gtest.C(t, func(t *gtest.T) {
		path := "test.json"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1").Int(), 2)
	})
}

func Test_Load_JSON2(t *testing.T) {
	data := []byte(`{"n":123456789000000000000, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789000000000000")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1").Int(), 2)
	})
}

func Test_Load_XML(t *testing.T) {
	data := []byte(`<doc><a>1</a><a>2</a><a>3</a><m><k>v</k></m><n>123456789</n></doc>`)
	// XML
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("doc.n").String(), "123456789")
		t.Assert(j.Get("doc.m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("doc.m.k").String(), "v")
		t.Assert(j.Get("doc.a").Slice(), g.Slice{"1", "2", "3"})
		t.Assert(j.Get("doc.a.1").Int(), 2)
	})
	// XML
	gtest.C(t, func(t *gtest.T) {
		path := "test.xml"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		t.Assert(err, nil)
		t.Assert(j.Get("doc.n").String(), "123456789")
		t.Assert(j.Get("doc.m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("doc.m.k").String(), "v")
		t.Assert(j.Get("doc.a").Array(), g.Slice{"1", "2", "3"})
		t.Assert(j.Get("doc.a.1").Int(), 2)
	})

	// XML
	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(err, nil)
		t.Assert(j.Get("Output.ipageIndex"), "2")
		t.Assert(j.Get("Output.itotalRecords"), "GF框架")
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
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1").Int(), 2)
	})
	// YAML
	gtest.C(t, func(t *gtest.T) {
		path := "test.yaml"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{1, 2, 3})
		t.Assert(j.Get("a.1").Int(), 2)
	})
}

func Test_Load_YAML2(t *testing.T) {
	data := []byte("i : 123456789")
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("i"), "123456789")
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
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{"1", "2", "3"})
		t.Assert(j.Get("a.1").Int(), 2)
	})
	// TOML
	gtest.C(t, func(t *gtest.T) {
		path := "test.toml"
		gfile.PutBytes(path, data)
		defer gfile.Remove(path)
		j, err := gjson.Load(path)
		t.Assert(err, nil)
		t.Assert(j.Get("n").String(), "123456789")
		t.Assert(j.Get("m").Map(), g.Map{"k": "v"})
		t.Assert(j.Get("m.k").String(), "v")
		t.Assert(j.Get("a").Slice(), g.Slice{"1", "2", "3"})
		t.Assert(j.Get("a.1").Int(), 2)
	})
}

func Test_Load_TOML2(t *testing.T) {
	data := []byte("i=123456789")
	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		t.Assert(err, nil)
		t.Assert(j.Get("i"), "123456789")
	})
}

func Test_Load_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(nil)
		t.Assert(j.Interface(), nil)
		_, err := gjson.Decode(nil)
		t.AssertNE(err, nil)
		_, err = gjson.DecodeToJson(nil)
		t.AssertNE(err, nil)
		j, err = gjson.LoadContent(nil)
		t.Assert(err, nil)
		t.Assert(j.Interface(), nil)

		j, err = gjson.LoadContent(`{"name": "gf"}`)
		t.Assert(err, nil)

		j, err = gjson.LoadContent(`{"name": "gf"""}`)
		t.AssertNE(err, nil)

		j = gjson.New(&g.Map{"name": "gf"})
		t.Assert(j.Get("name").String(), "gf")

	})
}

func Test_Load_Ini(t *testing.T) {
	var data = `

;注释

[addr]
ip = 127.0.0.1
port=9001
enable=true

	[DBINFO]
	type=mysql
	user=root
	password=password

`

	gtest.C(t, func(t *gtest.T) {
		j, err := gjson.LoadContent(data)
		if err != nil {
			gtest.Fatal(err)
		}

		t.Assert(j.Get("addr.ip").String(), "127.0.0.1")
		t.Assert(j.Get("addr.port").String(), "9001")
		t.Assert(j.Get("addr.enable").String(), "true")
		t.Assert(j.Get("DBINFO.type").String(), "mysql")
		t.Assert(j.Get("DBINFO.user").String(), "root")
		t.Assert(j.Get("DBINFO.password").String(), "password")

		_, err = j.ToIni()
		if err != nil {
			gtest.Fatal(err)
		}
	})
}

func Test_Load_YamlWithV3(t *testing.T) {
	content := `
# CLI tool, only in development environment.
# https://goframe.org/pages/viewpage.action?pageId=3673173
gfcli:
  gen:
    dao:
    - path            : "../../pkg/oss/oss/internal"
      group           : "oss"
      stdTime         : true
      descriptionTag  : true
      noJsonTag       : true
      noModelComment  : true
      overwriteDao    : true
      modelFileForDao : "model_dao.go"
      tablesEx        : |
        bpmn_info,
        dlocker,
        dlocker_detail,
        message_table,
        monitor_data,
        resource_param_info,
        version_info,
        version_topology_info,
        work_flow,
        work_flow_step_info,
        work_flow_undo_step_info

    - path            : "../../pkg/oss/workflow/internal"
      group           : "workflow"
      stdTime         : true
      descriptionTag  : true
      noJsonTag       : true
      noModelComment  : true
      overwriteDao    : true
      modelFileForDao : "model_dao.go"
`
	gtest.C(t, func(t *gtest.T) {
		_, err := gjson.LoadContent(content)
		t.AssertNil(err)
	})
}
