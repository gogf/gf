// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Fix_doFixV25Content(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			content = gtest.DataContent(`fix`, `fix25_content.go`)
			f       = cFix{}
		)
		_, err := f.doFixV25Content(content)
		t.AssertNil(err)
	})
}
