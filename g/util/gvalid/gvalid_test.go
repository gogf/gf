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

func Test_Required(t *testing.T) {
    if m := gvalid.Check("1", "required", nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("", "required", nil);  m == nil {
        t.Error(m)
    }
    if m := gvalid.Check("", "required-if:id,1,age,18", nil, map[string]string{"id" : "1", "age" : "19"});  m == nil {
        t.Error("required校验失败")
    }
    if m := gvalid.Check("", "required-if:id,1,age,18", nil, map[string]string{"id" : "2", "age" : "19"});  m != nil {
        t.Error("required校验失败")
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
    if m := gvalid.Check("10", rule, nil);  m != nil {
        t.Error(m)
    }
    if m := gvalid.Check("10.02", rule, nil);  m == nil {
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