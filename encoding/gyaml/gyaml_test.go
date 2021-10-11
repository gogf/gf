// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gyaml_test

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/internal/json"
	"testing"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/test/gtest"
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
	gtest.C(t, func(t *gtest.T) {
		result, err := gyaml.Decode([]byte(yamlStr))
		t.Assert(err, nil)

		m, ok := result.(map[string]interface{})
		t.Assert(ok, true)
		t.Assert(m, map[string]interface{}{
			"url":      "https://goframe.org",
			"server":   g.Slice{"120.168.117.21", "120.168.117.22"},
			"pi":       3.14,
			"hasChild": true,
			"name":     "你好YAML",
		})
	})
}

func Test_DecodeTo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := make(map[string]interface{})
		err := gyaml.DecodeTo([]byte(yamlStr), &result)
		t.Assert(err, nil)
		t.Assert(result, map[string]interface{}{
			"url":      "https://goframe.org",
			"server":   g.Slice{"120.168.117.21", "120.168.117.22"},
			"pi":       3.14,
			"hasChild": true,
			"name":     "你好YAML",
		})
	})
}

func Test_DecodeError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := gyaml.Decode([]byte(yamlErr))
		t.AssertNE(err, nil)

		result := make(map[string]interface{})
		err = gyaml.DecodeTo([]byte(yamlErr), &result)
		t.AssertNE(err, nil)
	})
}

func Test_DecodeMapToJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := []byte(`
m:
 k: v
    `)
		v, err := gyaml.Decode(data)
		t.Assert(err, nil)
		b, err := json.Marshal(v)
		t.Assert(err, nil)
		t.Assert(b, `{"m":{"k":"v"}}`)
	})
}

func Test_ToJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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

		p := gjson.New(res)
		if err != nil {
			t.Errorf("parser failed. %v", err)
			return
		}
		expectJson, err := p.ToJson()
		if err != nil {
			t.Errorf("parser ToJson failed. %v", err)
			return
		}
		t.Assert(jsonyaml, expectJson)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := gyaml.ToJson([]byte(yamlErr))
		if err == nil {
			t.Errorf("ToJson failed. %v", err)
			return
		}
	})
}
