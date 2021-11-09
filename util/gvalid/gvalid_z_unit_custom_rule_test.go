// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

func Test_CustomRule1(t *testing.T) {
	rule := "custom"
	err := gvalid.RegisterRule(
		rule,
		func(ctx context.Context, in gvalid.RuleFuncInput) error {
			pass := in.Value.String()
			if len(pass) != 6 {
				return errors.New(in.Message)
			}
			m := in.Data.Map()
			if m["data"] != pass {
				return errors.New(in.Message)
			}
			return nil
		},
	)
	gtest.Assert(err, nil)
	gtest.C(t, func(t *gtest.T) {
		err := gvalid.CheckValue(context.TODO(), "123456", rule, "custom message")
		t.Assert(err.String(), "custom message")
		err = gvalid.CheckValue(context.TODO(), "123456", rule, "custom message", g.Map{"data": "123456"})
		t.Assert(err, nil)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := gvalid.CheckStruct(context.TODO(), st, nil)
		t.Assert(err.String(), "自定义错误")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := gvalid.CheckStruct(context.TODO(), st, nil)
		t.Assert(err, nil)
	})
}

func Test_CustomRule2(t *testing.T) {
	rule := "required-map"
	err := gvalid.RegisterRule(rule, func(ctx context.Context, in gvalid.RuleFuncInput) error {
		m := in.Value.Map()
		if len(m) == 0 {
			return errors.New(in.Message)
		}
		return nil
	})
	gtest.Assert(err, nil)
	// Check.
	gtest.C(t, func(t *gtest.T) {
		errStr := "data map should not be empty"
		t.Assert(gvalid.CheckValue(context.TODO(), g.Map{}, rule, errStr).String(), errStr)
		t.Assert(gvalid.CheckValue(context.TODO(), g.Map{"k": "v"}, rule, errStr), nil)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value map[string]string `v:"uid@required-map#自定义错误"`
			Data  string            `p:"data"`
		}
		st := &T{
			Value: map[string]string{},
			Data:  "123456",
		}
		err := gvalid.CheckStruct(context.TODO(), st, nil)
		t.Assert(err.String(), "自定义错误")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value map[string]string `v:"uid@required-map#自定义错误"`
			Data  string            `p:"data"`
		}
		st := &T{
			Value: map[string]string{"k": "v"},
			Data:  "123456",
		}
		err := gvalid.CheckStruct(context.TODO(), st, nil)
		t.Assert(err, nil)
	})
}

func Test_CustomRule_AllowEmpty(t *testing.T) {
	rule := "allow-empty-str"
	err := gvalid.RegisterRule(rule, func(ctx context.Context, in gvalid.RuleFuncInput) error {
		s := in.Value.String()
		if len(s) == 0 || s == "gf" {
			return nil
		}
		return errors.New(in.Message)
	})
	gtest.Assert(err, nil)
	// Check.
	gtest.C(t, func(t *gtest.T) {
		errStr := "error"
		t.Assert(gvalid.CheckValue(context.TODO(), "", rule, errStr), nil)
		t.Assert(gvalid.CheckValue(context.TODO(), "gf", rule, errStr), nil)
		t.Assert(gvalid.CheckValue(context.TODO(), "gf2", rule, errStr).String(), errStr)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@allow-empty-str#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "",
			Data:  "123456",
		}
		err := gvalid.CheckStruct(context.TODO(), st, nil)
		t.Assert(err, nil)
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@allow-empty-str#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "john",
			Data:  "123456",
		}
		err := gvalid.CheckStruct(context.TODO(), st, nil)
		t.Assert(err.String(), "自定义错误")
	})
}

func TestValidator_RuleFunc(t *testing.T) {
	ruleName := "custom_1"
	ruleFunc := func(ctx context.Context, in gvalid.RuleFuncInput) error {
		pass := in.Value.String()
		if len(pass) != 6 {
			return errors.New(in.Message)
		}
		if m := in.Data.Map(); m["data"] != pass {
			return errors.New(in.Message)
		}
		return nil
	}
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Rules(ruleName).
			Messages("custom message").
			RuleFunc(ruleName, ruleFunc).
			CheckValue(ctx, "123456")
		t.Assert(err.String(), "custom message")
		err = g.Validator().
			Rules(ruleName).
			Messages("custom message").
			Data(g.Map{"data": "123456"}).
			RuleFunc(ruleName, ruleFunc).
			CheckValue(ctx, "123456")
		t.AssertNil(err)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := g.Validator().RuleFunc(ruleName, ruleFunc).CheckStruct(ctx, st)
		t.Assert(err.String(), "自定义错误")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := g.Validator().RuleFunc(ruleName, ruleFunc).CheckStruct(ctx, st)
		t.AssertNil(err)
	})
}

func TestValidator_RuleFuncMap(t *testing.T) {
	ruleName := "custom_1"
	ruleFunc := func(ctx context.Context, in gvalid.RuleFuncInput) error {
		pass := in.Value.String()
		if len(pass) != 6 {
			return errors.New(in.Message)
		}
		if m := in.Data.Map(); m["data"] != pass {
			return errors.New(in.Message)
		}
		return nil
	}
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().
			Rules(ruleName).
			Messages("custom message").
			RuleFuncMap(map[string]gvalid.RuleFunc{
				ruleName: ruleFunc,
			}).CheckValue(ctx, "123456")
		t.Assert(err.String(), "custom message")
		err = g.Validator().
			Rules(ruleName).
			Messages("custom message").
			Data(g.Map{"data": "123456"}).
			RuleFuncMap(map[string]gvalid.RuleFunc{
				ruleName: ruleFunc,
			}).
			CheckValue(ctx, "123456")
		t.AssertNil(err)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := g.Validator().
			RuleFuncMap(map[string]gvalid.RuleFunc{
				ruleName: ruleFunc,
			}).CheckStruct(ctx, st)
		t.Assert(err.String(), "自定义错误")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := g.Validator().
			RuleFuncMap(map[string]gvalid.RuleFunc{
				ruleName: ruleFunc,
			}).CheckStruct(ctx, st)
		t.AssertNil(err)
	})
}
