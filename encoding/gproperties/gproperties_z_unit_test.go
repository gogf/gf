// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gproperties provides accessing and converting for properties content.
package gproperties_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gproperties"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

var pStr string = `
# template
data = "/home/www/templates/"
# MySQL 
sql.disk.0  = 127.0.0.1:6379,0
sql.cache.0 = 127.0.0.1:6379,1=
sql.cache.1=0
sql.disk.a = 10
`

func TestDecode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		decodeStr, err := gproperties.Decode(([]byte)(pStr))
		if err != nil {
			t.Errorf("decode failed. %v", err)
			return
		}
		fmt.Printf("%v\n", decodeStr)
		v, _ := json.Marshal(decodeStr)
		fmt.Printf("%v\n", string(v))
	})
}
func TestEncode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		encStr, err := gproperties.Encode(map[string]interface{}{
			"sql": g.Map{
				"userName": "admin",
				"password": "123456",
			},
			"user": "admin",
			"no":   123,
		})
		if err != nil {
			t.Errorf("decode failed. %v", err)
			return
		}
		fmt.Printf("%v\n", string(encStr))
	})
}

func TestToJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		res, err := gproperties.Encode(map[string]interface{}{
			"sql": g.Map{
				"userName": "admin",
				"password": "123456",
			},
			"user": "admin",
			"no":   123,
		})
		fmt.Print(string(res))
		jsonPr, err := gproperties.ToJson(res)
		if err != nil {
			t.Errorf("ToJson failed. %v", err)
			return
		}
		fmt.Print(string(jsonPr))

		p := gjson.New(res)
		expectJson, err := p.ToJson()
		if err != nil {
			t.Errorf("parser ToJson failed. %v", err)
			return
		}
		t.Assert(jsonPr, expectJson)
	})
}
