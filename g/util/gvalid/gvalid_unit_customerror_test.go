// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gvalid_test

import (
    "gitee.com/johng/gf/g/util/gvalid"
    "strings"
    "testing"
)

func Test_SetDefaultErrorMsgs(t *testing.T) {
    rule := "integer|length:6,16"
    msgs := map[string]string {
        "integer" : "请输入一个整数",
        "length"  : "参数长度不对啊老铁",
    }
    gvalid.SetDefaultErrorMsgs(msgs)
    e := gvalid.Check("6.66", rule, nil)
    if e == nil || len(e.Map()) != 2 {
        t.Error("规则校验失败")
    } else {
        if v, ok := e.Map()["integer"]; ok {
            if strings.Compare(v, msgs["integer"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
        if v, ok := e.Map()["length"]; ok {
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
    e := gvalid.Check("6.66", rule, msgs)
    if e == nil || len(e.Map()) != 2 {
        t.Error("规则校验失败")
    } else {
        if v, ok := e.Map()["integer"]; ok {
            if strings.Compare(v, msgs["integer"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
        if v, ok := e.Map()["length"]; ok {
            if strings.Compare(v, msgs["length"]) != 0 {
                t.Error("错误信息不匹配")
            }
        }
    }
}

func Test_CustomError2(t *testing.T) {
    rule := "integer|length:6,16"
    msgs := "请输入一个整数|参数长度不对啊老铁"
    e := gvalid.Check("6.66", rule, msgs)
    if e == nil || len(e.Map()) != 2 {
        t.Error("规则校验失败")
    } else {
        if v, ok := e.Map()["integer"]; ok {
            if strings.Compare(v, "请输入一个整数") != 0 {
                t.Error("错误信息不匹配")
            }
        }
        if v, ok := e.Map()["length"]; ok {
            if strings.Compare(v, "参数长度不对啊老铁") != 0 {
                t.Error("错误信息不匹配")
            }
        }
    }
}