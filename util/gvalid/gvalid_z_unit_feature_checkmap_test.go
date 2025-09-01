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
		data := map[string]any{
			"id":   "0",
			"name": "john",
		}
		rules := map[string]string{
			"id":   "required|between:1,100",
			"name": "required|length:6,16",
		}
		if m := g.Validator().Data(data).Rules(rules).Run(context.TODO()); m == nil {
			t.Error("CheckMap校验失败")
		} else {
			t.Assert(len(m.Maps()), 2)
			t.Assert(m.Maps()["id"]["between"], "The id value `0` must be between 1 and 100")
			t.Assert(m.Maps()["name"]["length"], "The name value `john` length must be between 6 and 16")
		}
	})
}

func Test_CheckMap2(t *testing.T) {
	var params any
	gtest.C(t, func(t *gtest.T) {
		if err := g.Validator().Data(params).Run(context.TODO()); err == nil {
			t.AssertNil(err)
		}
	})

	kvmap := map[string]any{
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
	if m := g.Validator().Data(kvmap).Rules(rules).Messages(msgs).Run(context.TODO()); m == nil {
		t.Error("CheckMap校验失败")
	}

	kvmap = map[string]any{
		"id":   "1",
		"name": "john",
	}
	rules = map[string]string{
		"id":   "required|between:1,100",
		"name": "required|length:4,16",
	}
	msgs = map[string]any{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]any{
		"id":   "1",
		"name": "john",
	}
	rules = map[string]string{
		"id":   "",
		"name": "",
	}
	msgs = map[string]any{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]any{
		"id":   "1",
		"name": "john",
	}
	rules2 := []string{
		"@required|between:1,100",
		"@required|length:4,16",
	}
	msgs = map[string]any{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules2).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]any{
		"id":   "1",
		"name": "john",
	}
	rules2 = []string{
		"id@required|between:1,100",
		"name@required|length:4,16#名称不能为空|",
	}
	msgs = map[string]any{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules2).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]any{
		"id":   "1",
		"name": "john",
	}
	rules2 = []string{
		"id@required|between:1,100",
		"name@required|length:4,16#名称不能为空",
	}
	msgs = map[string]any{
		"id": "ID不能为空|ID范围应当为{min}到{max}",
		"name": map[string]string{
			"required": "名称不能为空",
			"length":   "名称长度为{min}到{max}个字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules2).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}
}

func Test_CheckMapWithNilAndNotRequiredField(t *testing.T) {
	data := map[string]any{
		"id": "1",
	}
	rules := map[string]string{
		"id":   "required",
		"name": "length:4,16",
	}
	if m := g.Validator().Data(data).Rules(rules).Run(context.TODO()); m != nil {
		t.Error(m)
	}
}

func Test_Sequence(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		params := map[string]any{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@required|length:6,16#账号不能为空|账号长度应当在{min}到{max}之间",
			"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Data(params).Rules(rules).Run(context.TODO())
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
		params := map[string]any{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@required|length:6,16#账号不能为空|账号长度应当在{min}到{max}之间",
			"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Bail().Rules(rules).Data(params).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.String(), "账号不能为空")
	})
	// global bail with rule bail
	gtest.C(t, func(t *gtest.T) {
		params := map[string]any{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@bail|required|length:6,16#|账号不能为空|账号长度应当在{min}到{max}之间",
			"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Bail().Rules(rules).Data(params).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.String(), "账号不能为空")
	})
}
