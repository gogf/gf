// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

func Test_CI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gvalid.CheckValue(ctx, "id", "in:Id,Name", nil)
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		err := gvalid.CheckValue(ctx, "id", "ci|in:Id,Name", nil)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().CaseInsensitive().Rules("in:Id,Name").CheckValue(ctx, "id")
		t.AssertNil(err)
	})
}
