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
)

func Test_CI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("id").Rules("in:Id,Name").Run(ctx)
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("id").Rules("ci|in:Id,Name").Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Ci().Rules("in:Id,Name").Data("id").Run(ctx)
		t.AssertNil(err)
	})
}
