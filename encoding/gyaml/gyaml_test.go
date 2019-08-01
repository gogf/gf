// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gyaml_test

import (
	"testing"

	"github.com/gogf/gf/encoding/gparser"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/encoding/gyaml"
	"github.com/gogf/gf/test/gtest"
)

var yamlStr string = `
#即表示url属性值；
url: https://goframe.org

#数组，即表示server为[a,b,c]
server:
    - 120.168.117.21
    - 120.168.117.22
#常量
pi: 3.14   #定义一个数值3.14
hasChild: true  #定义一个boolean值
name: '你好YAML'   #定义一个字符串
`

var yamlErr string = `
[redis]
dd = 11
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`

func Test_Decode(t *testing.T) {
	gtest.Case(t, func() {
		result, err := gyaml.Decode([]byte(yamlStr))
		gtest.Assert(err, nil)

		m, ok := result.(map[string]interface{})
		gtest.Assert(ok, true)
		gtest.Assert(m, map[string]interface{}{
			"url":      "https://goframe.org",
			"server":   g.Slice{"120.168.117.21", "120.168.117.22"},
			"pi":       3.14,
			"hasChild": true,
			"name":     "你好YAML",
		})
	})
}

func Test_DecodeTo(t *testing.T) {
	gtest.Case(t, func() {
		result := make(map[string]interface{})
		err := gyaml.DecodeTo([]byte(yamlStr), &result)
		gtest.Assert(err, nil)
		gtest.Assert(result, map[string]interface{}{
			"url":      "https://goframe.org",
			"server":   g.Slice{"120.168.117.21", "120.168.117.22"},
			"pi":       3.14,
			"hasChild": true,
			"name":     "你好YAML",
		})
	})
}

func Test_DecodeError(t *testing.T) {
	gtest.Case(t, func() {
		_, err := gyaml.Decode([]byte(yamlErr))
		gtest.AssertNE(err, nil)

		result := make(map[string]interface{})
		err = gyaml.DecodeTo([]byte(yamlErr), &result)
		gtest.AssertNE(err, nil)
	})
}

func Test_ToJson(t *testing.T) {
	gtest.Case(t, func() {
		m := make(map[string]string)
		m["yaml"] = yamlStr
		res, err := gyaml.Encode(m)
		if err != nil {
			t.Errorf("encode failed. %v", err)
			return
		}

		jsonyaml, err := gyaml.ToJson(res)
		if err != nil {
			t.Errorf("ToJson failed. %v", err)
			return
		}

		p, err := gparser.LoadContent(res)
		if err != nil {
			t.Errorf("parser failed. %v", err)
			return
		}
		expectJson, err := p.ToJson()
		if err != nil {
			t.Errorf("parser ToJson failed. %v", err)
			return
		}
		gtest.Assert(jsonyaml, expectJson)
	})

	gtest.Case(t, func() {
		_, err := gyaml.ToJson([]byte(yamlErr))
		if err == nil {
			t.Errorf("ToJson failed. %v", err)
			return
		}
	})
}
