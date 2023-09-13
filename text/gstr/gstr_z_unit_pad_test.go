// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Pad(t *testing.T) {
	str := "abc"
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.PadStart(str, " ", 5), "  abc")
		t.Assert(gstr.PadEnd(str, " ", 5), "abc  ")
	})
}
