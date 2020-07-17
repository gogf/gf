// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gvalid_test

import (
	"errors"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gconv"
	"testing"

	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gvalid"
)

func Test_CustomRule(t *testing.T) {
	rule := "custom"
	err := gvalid.RegisterRule(rule, func(value interface{}, message string, params map[string]interface{}) error {
		pass := gconv.String(value)
		if len(pass) != 6 {
			return errors.New(message)
		}
		if params["data"] != pass {
			return errors.New(message)
		}
		return nil
	})
	gtest.Assert(err, nil)
	gtest.C(t, func(t *gtest.T) {
		err := gvalid.Check("123456", rule, "custom message")
		t.Assert(err.String(), "custom message")
		err = gvalid.Check("123456", rule, "custom message", g.Map{"data": "123456"})
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
		err := gvalid.CheckStruct(st, nil)
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
		err := gvalid.CheckStruct(st, nil)
		t.Assert(err, nil)
	})
}
