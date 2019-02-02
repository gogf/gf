// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
    "github.com/gogf/gf/g/util/gvalid"
    "testing"
)

func Test_CheckStruct(t *testing.T) {
    type Object struct {
        Name string
        Age  int
    }
    rules := map[string]string {
        "Name" : "required|length:6,16",
        "Age"  : "between:18,30",
    }
    msgs  := map[string]interface{} {
        "Name" : map[string]string {
            "required" : "名称不能为空",
            "length"   : "名称长度为:min到:max个字符",
        },
        "Age"  : "年龄为18到30周岁",
    }
    obj := &Object{"john", 16}
    if m := gvalid.CheckStruct(obj, rules, msgs); m == nil {
        t.Error("CheckObject校验失败")
    }
}
