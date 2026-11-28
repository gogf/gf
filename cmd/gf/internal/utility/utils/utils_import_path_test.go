// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies that CLI import paths follow go.mod module directives.

package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

// Test_GetImportPath_ModuleDirective covers valid module layouts and missing declarations.
func Test_GetImportPath_ModuleDirective(t *testing.T) {
	for _, test := range []struct {
		name    string
		content string
		want    string
	}{
		{"plain", "module example.com/project\n\ngo 1.26.0\n", "example.com/project/internal/model"},
		{"leading comment", "// module example.com/old\nmodule example.com/project\n", "example.com/project/internal/model"},
		{"leading blank line", "\n\nmodule example.com/project\n", "example.com/project/internal/model"},
		{"leading whitespace", "\tmodule example.com/project\n", "example.com/project/internal/model"},
		{"quoted module", "module \"example.com/project\"\n", "example.com/project/internal/model"},
		{"inline comment", "module example.com/project // project module\n", "example.com/project/internal/model"},
		{"crlf", "// project\r\nmodule example.com/project\r\n", "example.com/project/internal/model"},
		{"module after go directive", "go 1.26.0\nmodule example.com/project\n", "example.com/project/internal/model"},
		{"no module directive", "go 1.26.0\n", ""},
		{"missing module path", "module\n", ""},
		{"unterminated quoted path", "module \"example.com/project\n", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					root = t.TempDir()
					dir  = filepath.Join(root, "internal", "model")
				)
				t.AssertNil(gfile.Mkdir(dir))
				t.AssertNil(gfile.PutContents(filepath.Join(root, "go.mod"), test.content))
				t.Assert(utils.GetImportPath(dir), test.want)
			})
		})
	}
}
