// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gyaml_test

import (
	"github.com/gogf/gf/g/encoding/gparser"
	"github.com/gogf/gf/g/encoding/gyaml"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

var yamlStr string = `
#即表示url属性值；
url: http://www.wolfcode.cn 
#即表示server.host属性的值；
server:
    host: http://www.wolfcode.cn 
#数组，即表示server为[a,b,c]
server:
    - 120.168.117.21
    - 120.168.117.22
    - 120.168.117.23
#常量
pi: 3.14   #定义一个数值3.14
hasChild: true  #定义一个boolean值
name: '你好YAML'   #定义一个字符串
`

var yamlErr string = `
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
	gtest.Case(t, func() {
		m := make(map[string]string)
		m["yaml"] = yamlStr
		res, err := gyaml.Encode(m)
		if err != nil {
			t.Errorf("encode failed. %v", err)
			return
		}

		p, err := gparser.LoadContent(res)
		if err != nil {
			t.Errorf("parser failed. %v", err)
			return
		}

		gtest.Assert(p.GetString("yaml"), yamlStr)
	})

}

func TestDecode(t *testing.T) {
	gtest.Case(t, func() {
		m := make(map[string]string)
		m["yaml"] = yamlStr
		res, err := gyaml.Encode(m)
		if err != nil {
			t.Errorf("encode failed. %v", err)
			return
		}

		decodeStr, err := gyaml.Decode(res)
		if err != nil {
			t.Errorf("decode failed. %v", err)
			return
		}

		gtest.Assert(decodeStr.(map[string]interface{})["yaml"], yamlStr)

		decodeStr1 := make(map[string]interface{})
		err = gyaml.DecodeTo(res, &decodeStr1)
		if err != nil {
			t.Errorf("decodeTo failed. %v", err)
			return
		}
		gtest.Assert(decodeStr1["yaml"], yamlStr)
	})

	gtest.Case(t, func() {
		_, err := gyaml.Decode([]byte(yamlErr))
		if err == nil {
			t.Errorf("decode failed. %v", err)
			return
		}

		decodeStr1 := make(map[string]interface{})
		err = gyaml.DecodeTo([]byte(yamlErr), &decodeStr1)
		if err == nil {
			t.Errorf("decodeTo failed. %v", err)
			return
		}
	})
}

func TestToJson(t *testing.T) {
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
