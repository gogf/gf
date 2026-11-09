// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_parseSequenceTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "name@required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "name")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "||密码强度不足|两次密码不一致")
	})
	gtest.C(t, func(t *gtest.T) {
		s := "required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "||密码强度不足|两次密码不一致")
	})
	gtest.C(t, func(t *gtest.T) {
		s := "required|length:2,20|password3|same:password1"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "")
	})
	gtest.C(t, func(t *gtest.T) {
		s := "required"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "")
		t.Assert(rule, "required")
		t.Assert(msg, "")
	})
}

func Test_GetTags(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(structTagPriority, GetTags())
	})
}

// It tests that the parsed rule value cache is disabled for the default
// Validator, and is enabled if it is created with function New(true).
func Test_New_Cache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx         = context.Background()
			uncachedTag = `zztmp_uncached_field@required#uncached`
			cachedTag   = `zztmp_cached_field@required#cached`
		)
		// The default Validator does not cache the parsed rule values.
		err := New().Data(map[string]any{"zztmp_uncached_field": ""}).Rules([]string{uncachedTag}).Run(ctx)
		t.AssertNE(err, nil)
		_, ok := parsedTagValueCache.Load(uncachedTag)
		t.Assert(ok, false)

		// The Validator created with New(true) caches the parsed rule values.
		err = New(true).Data(map[string]any{"zztmp_cached_field": ""}).Rules([]string{cachedTag}).Run(ctx)
		t.AssertNE(err, nil)
		_, ok = parsedTagValueCache.Load(cachedTag)
		t.Assert(ok, true)
	})
}
