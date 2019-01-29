// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gvalid_test

import (
    "gitee.com/johng/gf/g/util/gvalid"
    "testing"
)

func Test_CheckMap(t *testing.T) {
    kvmap := map[string]interface{} {
        "id"   : "0",
        "name" : "john",
    }
    rules := map[string]string {
        "id"   : "required|between:1,100",
        "name" : "required|length:6,16",
    }
    msgs  := gvalid.CustomMsg {
        "id"   : "ID不能为空|ID范围应当为:min到:max",
        "name" : map[string]string {
            "required" : "名称不能为空",
            "length"   : "名称长度为:min到:max个字符",
        },
    }
    if m := gvalid.CheckMap(kvmap, rules, msgs); m == nil {
        t.Error("CheckMap校验失败")
    }

    kvmap = map[string]interface{} {
        "id"   : "1",
        "name" : "john",
    }
    rules = map[string]string {
        "id"   : "required|between:1,100",
        "name" : "required|length:4,16",
    }
    msgs  = map[string]interface{} {
        "id"   : "ID不能为空|ID范围应当为:min到:max",
        "name" : map[string]string {
            "required" : "名称不能为空",
            "length"   : "名称长度为:min到:max个字符",
        },
    }
    if m := gvalid.CheckMap(kvmap, rules, msgs); m != nil {
        t.Error(m)
    }
}

// 如果值为nil，并且不需要require*验证时，其他验证失效
func Test_CheckMapWithNilAndNotRequiredField(t *testing.T) {
    data  := map[string]interface{} {
        "id"   : "1",
    }
    rules := map[string]string {
        "id"   : "required",
        "name" : "length:4,16",
    }
    if m := gvalid.CheckMap(data, rules); m != nil {
        t.Error(m)
    }
}

