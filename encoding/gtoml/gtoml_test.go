// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gtoml_test

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gtoml"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

var tomlStr string = `
# 模板引擎目录
viewpath = "/home/www/templates/"
# MySQL数据库配置
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`

var tomlErr string = `
# 模板引擎目录
viewpath = "/home/www/templates/"
# MySQL数据库配置
[redis]
dd = 11
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`

func TestEncode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string]string)
		m["toml"] = tomlStr
		res, err := gtoml.Encode(m)
		if err != nil {
			t.Errorf("encode failed. %v", err)
			return
		}

		t.Assert(gjson.New(res).Get("toml").String(), tomlStr)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := gtoml.Encode(tomlErr)
		if err == nil {
			t.Errorf("encode should be failed. %v", err)
			return
		}
	})
}

func TestDecode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string]string)
		m["toml"] = tomlStr
		res, err := gtoml.Encode(m)
		if err != nil {
			t.Errorf("encode failed. %v", err)
			return
		}

		decodeStr, err := gtoml.Decode(res)
		if err != nil {
			t.Errorf("decode failed. %v", err)
			return
		}

		t.Assert(decodeStr.(map[string]interface{})["toml"], tomlStr)

		decodeStr1 := make(map[string]interface{})
		err = gtoml.DecodeTo(res, &decodeStr1)
		if err != nil {
			t.Errorf("decodeTo failed. %v", err)
			return
		}
		t.Assert(decodeStr1["toml"], tomlStr)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := gtoml.Decode([]byte(tomlErr))
		if err == nil {
			t.Errorf("decode failed. %v", err)
			return
		}

		decodeStr1 := make(map[string]interface{})
		err = gtoml.DecodeTo([]byte(tomlErr), &decodeStr1)
		if err == nil {
			t.Errorf("decodeTo failed. %v", err)
			return
		}
	})
}

func TestToJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string]string)
		m["toml"] = tomlStr
		res, err := gtoml.Encode(m)
		if err != nil {
			t.Errorf("encode failed. %v", err)
			return
		}

		jsonToml, err := gtoml.ToJson(res)
		if err != nil {
			t.Errorf("ToJson failed. %v", err)
			return
		}

		p := gjson.New(res)
		expectJson, err := p.ToJson()
		if err != nil {
			t.Errorf("parser ToJson failed. %v", err)
			return
		}
		t.Assert(jsonToml, expectJson)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := gtoml.ToJson([]byte(tomlErr))
		if err == nil {
			t.Errorf("ToJson failed. %v", err)
			return
		}
	})
}
