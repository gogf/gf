// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gini_test

import (
	"github.com/jin502437344/gf/encoding/gini"
	"github.com/jin502437344/gf/encoding/gjson"
	"github.com/jin502437344/gf/test/gtest"
	"testing"
)

var iniContent = `

;注释
aa=bb
[addr] 
#注释
ip = 127.0.0.1
port=9001
enable=true

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
		t.Assert(res["addr"].(map[string]interface{})["ip"], "127.0.0.1")
		t.Assert(res["addr"].(map[string]interface{})["port"], "9001")
		t.Assert(res["DBINFO"].(map[string]interface{})["user"], "root")
		t.Assert(res["DBINFO"].(map[string]interface{})["type"], "mysql")
		t.Assert(res["键"].(map[string]interface{})["呵呵"], "值")
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

		t.Assert(res["addr"].(map[string]interface{})["ip"], "127.0.0.1")
		t.Assert(res["addr"].(map[string]interface{})["port"], "9001")
		t.Assert(res["DBINFO"].(map[string]interface{})["user"], "root")
		t.Assert(res["DBINFO"].(map[string]interface{})["type"], "mysql")

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
		t.Assert(err, nil)

		t.Assert(iniMap["addr"].(map[string]interface{})["ip"], json.GetString("addr.ip"))
		t.Assert(iniMap["addr"].(map[string]interface{})["port"], json.GetString("addr.port"))
		t.Assert(iniMap["DBINFO"].(map[string]interface{})["user"], json.GetString("DBINFO.user"))
		t.Assert(iniMap["DBINFO"].(map[string]interface{})["type"], json.GetString("DBINFO.type"))
	})
}
