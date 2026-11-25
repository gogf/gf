// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_FormatDirHref(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, testCase := range []struct {
			path     string
			expected string
		}{
			{"/", "/"},
			{"/dir", "/dir"},
			{"/a b", "/a%20b"},
			{"/a#b", "/a%23b"},
			{"/a?b", "/a%3Fb"},
			{`/a"b`, "/a%22b"},
			{"/a&b", "/a&amp;b"},
			{"/<script>", "/%3Cscript%3E"},
			{"/中文", "/%E4%B8%AD%E6%96%87"},
			{`/<img src=x onerror=alert(1)>/x"y.txt`, "/%3Cimg%20src=x%20onerror=alert%281%29%3E/x%22y.txt"},
		} {
			t.Assert(formatDirHref(testCase.path), testCase.expected)
		}
	})
}
