// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func Test_GetModPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		goModPath := utils.GetModPath()
		fmt.Println(goModPath)
	})
}

func Test_GoFmt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "example.go")
			source   = "package main\n\nimport (\n\"fmt\"\n\"os\"\n)\n\nfunc main() {\nfmt.Println()\n}\n"
			expect   = "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println()\n}\n"
		)
		gfile.PutContents(filePath, source)
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)

		// The formatted file should keep unchanged on repeated calls.
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// Named import in use and unused import removed.
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "example.go")
			source   = "package main\n\nimport (\nfm \"fmt\"\n\"os\"\n)\n\nfunc main() {\nfm.Println()\n}\n"
			expect   = "package main\n\nimport (\n\tfm \"fmt\"\n)\n\nfunc main() {\n\tfm.Println()\n}\n"
		)
		gfile.PutContents(filePath, source)
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})
}
