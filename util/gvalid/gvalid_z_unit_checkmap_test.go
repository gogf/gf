// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

func Test_CheckMap1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := map[string]interface{}{
			"id":   "0",
			"name": "john",
		}
		rules := map[string]string{
			"id":   "required|between:1,100",
			"name": "required|length:6,16",
		}
		if m := gvalid.CheckMap(context.TODO(), data, rules); m == nil {
			t.Error("CheckMap校验失败")
		} else {
			t.Assert(len(m.Maps()), 2)
			t.Assert(m.Maps()["id"]["between"], "The id value `0` must be between 1 and 100")
			t.Assert(m.Maps()["name"]["length"], "The name value `john` length must be between 6 and 16")
		}
	})
}

func Test_CheckMap2(t *testing.T) {
	var params interface{}
	gtest.C(t, func(t *gtest.T) {
		if err := gvalid.CheckMap(context.TODO(), params, nil, nil); err == nil {
			t.Assert(err, nil)
		}
	})

	kvmap := map[string]interface{}{
		"id":   "0",
		"name": "john",
	}
	rules := map[string]string{
		"id":   "required|between:1,100",
		"name": "required|length:6,16",
	}
	msgs := gvalid.CustomMsg{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := gvalid.CheckMap(context.TODO(), kvmap, rules, msgs); m == nil {
		t.Error("CheckMap校验失败")
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules = map[string]string{
		"id":   "required|between:1,100",
		"name": "required|length:4,16",
	}
	msgs = map[string]interface{}{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := gvalid.CheckMap(context.TODO(), kvmap, rules, msgs); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules = map[string]string{
		"id":   "",
		"name": "",
	}
	msgs = map[string]interface{}{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := gvalid.CheckMap(context.TODO(), kvmap, rules, msgs); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules2 := []string{
		"@required|between:1,100",
		"@required|length:4,16",
	}
	msgs = map[string]interface{}{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := gvalid.CheckMap(context.TODO(), kvmap, rules2, msgs); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules2 = []string{
		"id@required|between:1,100",
		"name@required|length:4,16#名称不能为空|",
	}
	msgs = map[string]interface{}{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := gvalid.CheckMap(context.TODO(), kvmap, rules2, msgs); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules2 = []string{
		"id@required|between:1,100",
		"name@required|length:4,16#名称不能为空",
	}
	msgs = map[string]interface{}{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := gvalid.CheckMap(context.TODO(), kvmap, rules2, msgs); m != nil {
		t.Error(m)
	}
}

// 如果值为nil，并且不需要require*验证时，其他验证失效
func Test_CheckMapWithNilAndNotRequiredField(t *testing.T) {
	data := map[string]interface{}{
		"id": "1",
	}
	rules := map[string]string{
		"id":   "required",
		"name": "length:4,16",
	}
	if m := gvalid.CheckMap(context.TODO(), data, rules); m != nil {
		t.Error(m)
	}
}

func Test_Sequence(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		params := map[string]interface{}{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@required|length:6,16#账号不能为空|账号长度应当在{min}到{max}之间",
			"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
			"password2@required|length:6,16#",
		}
		err := gvalid.CheckMap(context.TODO(), params, rules)
		t.AssertNE(err, nil)
		t.Assert(len(err.Map()), 2)
		t.Assert(err.Map()["required"], "账号不能为空")
		t.Assert(err.Map()["length"], "账号长度应当在6到16之间")
		t.Assert(len(err.Maps()), 2)

		t.Assert(len(err.Items()), 2)
		t.Assert(err.Items()[0]["passport"]["length"], "账号长度应当在6到16之间")
		t.Assert(err.Items()[0]["passport"]["required"], "账号不能为空")
		t.Assert(err.Items()[1]["password"]["same"], "两次密码输入不相等")

		t.Assert(err.String(), "账号不能为空; 账号长度应当在6到16之间; 两次密码输入不相等")
		t.Assert(err.Strings(), []string{"账号不能为空", "账号长度应当在6到16之间", "两次密码输入不相等"})

		k, m := err.FirstItem()
		t.Assert(k, "passport")
		t.Assert(m, err.Map())

		r, s := err.FirstRule()
		t.Assert(r, "required")
		t.Assert(s, "账号不能为空")

		t.Assert(gerror.Current(err), "账号不能为空")
	})
}

func Test_Map_Bail(t *testing.T) {
	// global bail
	gtest.C(t, func(t *gtest.T) {
		params := map[string]interface{}{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@required|length:6,16#账号不能为空|账号长度应当在{min}到{max}之间",
			"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Bail().Rules(rules).CheckMap(ctx, params)
		t.AssertNE(err, nil)
		t.Assert(err.String(), "账号不能为空; 账号长度应当在6到16之间")
	})
	// global bail with rule bail
	gtest.C(t, func(t *gtest.T) {
		params := map[string]interface{}{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@bail|required|length:6,16#|账号不能为空|账号长度应当在{min}到{max}之间",
			"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Bail().Rules(rules).CheckMap(ctx, params)
		t.AssertNE(err, nil)
		t.Assert(err.String(), "账号不能为空")
	})
}
