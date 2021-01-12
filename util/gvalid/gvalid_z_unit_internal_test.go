// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"testing"

	"github.com/gogf/gf/test/gtest"
)

func Test_parseSequenceTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "name@required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致"
		field, rule, msg := parseSequenceTag(s)
		t.Assert(field, "name")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "||密码强度不足|两次密码不一致")
	})
	gtest.C(t, func(t *gtest.T) {
		s := "required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致"
		field, rule, msg := parseSequenceTag(s)
		t.Assert(field, "")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "||密码强度不足|两次密码不一致")
	})
	gtest.C(t, func(t *gtest.T) {
		s := "required|length:2,20|password3|same:password1"
		field, rule, msg := parseSequenceTag(s)
		t.Assert(field, "")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "")
	})
	gtest.C(t, func(t *gtest.T) {
		s := "required"
		field, rule, msg := parseSequenceTag(s)
		t.Assert(field, "")
		t.Assert(rule, "required")
		t.Assert(msg, "")
	})
}
