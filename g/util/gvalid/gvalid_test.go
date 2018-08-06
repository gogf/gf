// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 单元测试
// go test *.go -bench=".*"

package gvalid_test

import (
    "testing"
    "gitee.com/johng/gf/g/util/gvalid"
    "strings"
)

func Test_Regex(t *testing.T) {
    rule := `regex:\d{6}|\D{6}|length:6,16`
    if m := gvalid.Check("123456", rule, nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("abcde6", rule, nil);  m == nil {
        t.Error("校验失败")
    }
}

func Test_CheckMap(t *testing.T) {
    kvmap := map[string]interface{} {
        "id"   : "0",
        "name" : "john",
    }
    rules := map[string]string {
        "id"   : "required|between:1,100",
        "name" : "required|length:6,16",
    }
    msgs  := map[string]interface{} {
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

func Test_CheckObject(t *testing.T) {
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

func Test_Required(t *testing.T) {
    if m := gvalid.Check("1", "required", nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("", "required", nil);  m == nil {
        t.Error(m)
    }
    if m := gvalid.Check("", "required-if:id,1,age,18", nil, map[string]interface{}{"id" : 1, "age" : 19});  m == nil {
        t.Error("Required校验失败")
    }
    if m := gvalid.Check("", "required-if:id,1,age,18", nil, map[string]interface{}{"id" : 2, "age" : 19});  m != nil {
        t.Error("Required校验失败")
    }
}

func Test_Ip(t *testing.T) {
    if m := gvalid.Check("10.0.0.1", "ipv4", nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("0.0.0.0", "ipv4", nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("1920.0.0.0", "ipv4", nil);  m == nil {
        t.Error("ipv4校验失败")
    }
    if m := gvalid.Check("fe80::5484:7aff:fefe:9799", "ipv6", nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("fe80::5484:7aff:fefe:9799123", "ipv6", nil);  m == nil {
        t.Error(m)
    }
}

func Test_Length(t *testing.T) {
    rule := "length:6,16"
    if m := gvalid.Check("123456", rule, nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("12345", rule, nil);  m == nil {
        t.Error("长度校验失败")
    }
}

func Test_MinLength(t *testing.T) {
    rule := "min-length:6"
    if m := gvalid.Check("123456", rule, nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("12345", rule, nil);  m == nil {
        t.Error("长度校验失败")
    }
}

func Test_MaxLength(t *testing.T) {
    rule := "max-length:6"
    if m := gvalid.Check("12345", rule, nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("1234567", rule, nil);  m == nil {
        t.Error("长度校验失败")
    }
}

func Test_Between(t *testing.T) {
    rule := "between:6.01, 10.01"
    if m := gvalid.Check(10, rule, nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check(10.02, rule, nil);  m == nil {
        t.Error("大小范围校验失败")
    }
}

func Test_SetDefaultErrorMsgs(t *testing.T) {
    rule := "integer|length:6,16"
    msgs := map[string]string {
        "integer" : "请输入一个整数",
        "length"  : "参数长度不对啊老铁",
    }
    gvalid.SetDefaultErrorMsgs(msgs)
    m := gvalid.Check("6.66", rule, nil)
    if len(m) != 2 {
        t.Error("规则校验失败")
    } else {
        if v, ok := m["integer"]; ok {
            if strings.Compare(v, msgs["integer"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
        if v, ok := m["length"]; ok {
            if strings.Compare(v, msgs["length"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
    }
}

func Test_CustomError1(t *testing.T) {
    rule := "integer|length:6,16"
    msgs := map[string]string {
        "integer" : "请输入一个整数",
        "length"  : "参数长度不对啊老铁",
    }
    m := gvalid.Check("6.66", rule, msgs)
    if len(m) != 2 {
        t.Error("规则校验失败")
    } else {
        if v, ok := m["integer"]; ok {
            if strings.Compare(v, msgs["integer"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
        if v, ok := m["length"]; ok {
            if strings.Compare(v, msgs["length"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
    }
}

func Test_CustomError2(t *testing.T) {
    rule := "integer|length:6,16"
    msgs := "请输入一个整数|参数长度不对啊老铁"
    m := gvalid.Check("6.66", rule, msgs)
    if len(m) != 2 {
        t.Error("规则校验失败")
    } else {
        if v, ok := m["integer"]; ok {
            if strings.Compare(v, "请输入一个整数") != 0 {
                t.Error("错误信息不匹配")
            }
        }
        if v, ok := m["length"]; ok {
            if strings.Compare(v, "参数长度不对啊老铁") != 0 {
                t.Error("错误信息不匹配")
            }
        }
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