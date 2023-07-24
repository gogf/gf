// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	json2 "encoding/json"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
)

var (
	jsonStr1 = `{"name":"john","slice":[1,2,3]}`
	jsonStr2 = `{"CallbackCommand":"Group.CallbackAfterSendMsg","From_Account":"61934946","GroupId":"@TGS#2FLGX67FD","MsgBody":[{"MsgContent":{"Text":"是的"},"MsgType":"TIMTextElem"}],"MsgSeq":23,"MsgTime":1567032819,"Operator_Account":"61934946","Random":2804799576,"Type":"Public"}`
	jsonObj1 = gjson.New(jsonStr1)
	jsonObj2 = gjson.New(jsonStr2)
)

func Benchmark_Validate_Simple_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Valid(jsonStr1)
	}
}

func Benchmark_Validate_Complicated_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Valid(jsonStr2)
	}
}

func Benchmark_Get_Simple_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonObj1.Get("name")
	}
}

func Benchmark_Get_Complicated_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonObj2.Get("GroupId")
	}
}

func Benchmark_Stdlib_Json_Unmarshal_Simple_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		json2.Unmarshal([]byte(jsonStr1), &m)
	}
}

func Benchmark_Stdlib_Json_Unmarshal_Complicated_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		json2.Unmarshal([]byte(jsonStr2), &m)
	}
}

func Benchmark_New_Simple_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.New(jsonStr1)
	}
}

func Benchmark_New_Complicated_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.New(jsonStr2)
	}
}

func Benchmark_Remove_Simple_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonObj1.Remove("name")
	}
}

func Benchmark_Remove_Complicated_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonObj2.Remove("GroupId")
	}
}

func Benchmark_New_Nil_And_Set_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := gjson.New(nil)
		p.Set("k", "v")
	}
}

func Benchmark_New_Nil_And_Set_Multiple_Level(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := gjson.New(nil)
		p.Set("0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0", []int{1, 2, 3})
	}
}
