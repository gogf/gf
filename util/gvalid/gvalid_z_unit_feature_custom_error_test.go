// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			rule = "ipv4"
			val  = "0.0.0"
			err  = g.Validator().Data(val).Rules(rule).Run(context.TODO())
			msg  = map[string]string{
				"ipv4": "The value `0.0.0` is not a valid IPv4 address",
			}
		)
		t.Assert(err.Map(), msg)
	})
}

func Test_FirstString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			rule = "ipv4"
			val  = "0.0.0"
			err  = g.Validator().Data(val).Rules(rule).Run(context.TODO())
		)
		t.Assert(err.FirstError(), "The value `0.0.0` is not a valid IPv4 address")
	})
}

func Test_CustomError1(t *testing.T) {
	rule := "integer|length:6,16"
	msgs := map[string]string{
		"integer": "请输入一个整数",
		"length":  "参数长度不对啊老铁",
	}
	e := g.Validator().Data("6.66").Rules(rule).Messages(msgs).Run(context.TODO())
	if e == nil || len(e.Map()) != 2 {
		t.Error("规则校验失败")
	} else {
		if v, ok := e.Map()["integer"]; ok {
			if strings.Compare(v.Error(), msgs["integer"]) != 0 {
				t.Error("错误信息不匹配")
			}
		}
		if v, ok := e.Map()["length"]; ok {
			if strings.Compare(v.Error(), msgs["length"]) != 0 {
				t.Error("错误信息不匹配")
			}
		}
	}
}

func Test_CustomError2(t *testing.T) {
	rule := "integer|length:6,16"
	msgs := "请输入一个整数|参数长度不对啊老铁"
	e := g.Validator().Data("6.66").Rules(rule).Messages(msgs).Run(context.TODO())
	if e == nil || len(e.Map()) != 2 {
		t.Error("规则校验失败")
	} else {
		if v, ok := e.Map()["integer"]; ok {
			if strings.Compare(v.Error(), "请输入一个整数") != 0 {
				t.Error("错误信息不匹配")
			}
		}
		if v, ok := e.Map()["length"]; ok {
			if strings.Compare(v.Error(), "参数长度不对啊老铁") != 0 {
				t.Error("错误信息不匹配")
			}
		}
	}
}
