// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gini_test

import (
	"testing"

	"github.com/gogf/gf/v2/encoding/gini"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/test/gtest"
)

var iniContent = `

;注释
aa=bb
[addr] 
#注释
ip = 127.0.0.1
port=9001
enable=true
command=/bin/echo "gf=GoFrame"

	[DBINFO]
	type=mysql
	user=root
	password=password
[键]
呵呵=值

`

func TestDecode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		res, err := gini.Decode([]byte(iniContent))
		if err != nil {
			gtest.Fatal(err)
		}
		t.Assert(res["addr"].(map[string]any)["ip"], "127.0.0.1")
		t.Assert(res["addr"].(map[string]any)["port"], "9001")
		t.Assert(res["addr"].(map[string]any)["command"], `/bin/echo "gf=GoFrame"`)
		t.Assert(res["DBINFO"].(map[string]any)["user"], "root")
		t.Assert(res["DBINFO"].(map[string]any)["type"], "mysql")
		t.Assert(res["键"].(map[string]any)["呵呵"], "值")
	})

	gtest.C(t, func(t *gtest.T) {
		errContent := `
		a = b
`
		_, err := gini.Decode([]byte(errContent))
		if err == nil {
			gtest.Fatal(err)
		}
	})
}

func TestEncode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		iniMap, err := gini.Decode([]byte(iniContent))
		if err != nil {
			gtest.Fatal(err)
		}

		iniStr, err := gini.Encode(iniMap)
		if err != nil {
			gtest.Fatal(err)
		}

		res, err := gini.Decode(iniStr)
		if err != nil {
			gtest.Fatal(err)
		}

		t.Assert(res["addr"].(map[string]any)["ip"], "127.0.0.1")
		t.Assert(res["addr"].(map[string]any)["port"], "9001")
		t.Assert(res["DBINFO"].(map[string]any)["user"], "root")
		t.Assert(res["DBINFO"].(map[string]any)["type"], "mysql")

	})
}

func TestToJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		jsonStr, err := gini.ToJson([]byte(iniContent))
		if err != nil {
			gtest.Fatal(err)
		}

		json, err := gjson.LoadContent(jsonStr)
		if err != nil {
			gtest.Fatal(err)
		}

		iniMap, err := gini.Decode([]byte(iniContent))
		t.AssertNil(err)

		t.Assert(iniMap["addr"].(map[string]any)["ip"], json.Get("addr.ip").String())
		t.Assert(iniMap["addr"].(map[string]any)["port"], json.Get("addr.port").String())
		t.Assert(iniMap["DBINFO"].(map[string]any)["user"], json.Get("DBINFO.user").String())
		t.Assert(iniMap["DBINFO"].(map[string]any)["type"], json.Get("DBINFO.type").String())
	})
}
