// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gmode_test

import (
	"github.com/gogf/gf/v2/util/gmode"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_AutoCheckSourceCodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gmode.IsDevelop(), true)
	})
}

func Test_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		oldMode := gmode.Mode()
		defer gmode.Set(oldMode)
		gmode.SetDevelop()
		t.Assert(gmode.IsDevelop(), true)
		t.Assert(gmode.IsTesting(), false)
		t.Assert(gmode.IsStaging(), false)
		t.Assert(gmode.IsProduct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		oldMode := gmode.Mode()
		defer gmode.Set(oldMode)
		gmode.SetTesting()
		t.Assert(gmode.IsDevelop(), false)
		t.Assert(gmode.IsTesting(), true)
		t.Assert(gmode.IsStaging(), false)
		t.Assert(gmode.IsProduct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		oldMode := gmode.Mode()
		defer gmode.Set(oldMode)
		gmode.SetStaging()
		t.Assert(gmode.IsDevelop(), false)
		t.Assert(gmode.IsTesting(), false)
		t.Assert(gmode.IsStaging(), true)
		t.Assert(gmode.IsProduct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		oldMode := gmode.Mode()
		defer gmode.Set(oldMode)
		gmode.SetProduct()
		t.Assert(gmode.IsDevelop(), false)
		t.Assert(gmode.IsTesting(), false)
		t.Assert(gmode.IsStaging(), false)
		t.Assert(gmode.IsProduct(), true)
	})
}
