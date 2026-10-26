// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains unit tests for the public GoFmt helper.

package utils_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func Test_GoFmt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "example.go")
			source   = "package main\n\nimport (\n\"fmt\"\n\"os\"\n)\n\nfunc main() {\nfmt.Println()\n}\n"
			expect   = "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println()\n}\n"
		)
		t.AssertNil(gfile.PutContents(filePath, source))
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
		t.AssertNil(gfile.PutContents(filePath, source))
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// Import path with a version suffix whose package name is the second-to-last
		// segment must be recognized as used. See issue 3835.
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "example.go")
			source   = `package main

import (
	"context"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	_ = context.Background()
	var d mysql.Driver
	_ = d
}
`
			expect = `package main

import (
	"context"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	_ = context.Background()
	var d mysql.Driver
	_ = d
}
`
		)
		t.AssertNil(gfile.PutContents(filePath, source))
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// Import whose package name cannot be determined from the path is kept as is,
		// because removing it mistakenly would break the code.
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "example.go")
			source   = `package main

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func main() {
	_ = context.Background()
	var c *minio.Client
	_ = c
}
`
			expect = `package main

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func main() {
	_ = context.Background()
	var c *minio.Client
	_ = c
}
`
		)
		t.AssertNil(gfile.PutContents(filePath, source))
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// Imports are grouped into the standard library section and the others, and
		// sorted by import path within each section, as imports.Process did.
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "example.go")
			source   = `package main

import (
	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"go/ast"
	"context"
)

func main() {
	_ = context.Background()
	_ = ast.NewIdent("x")
	var d mysql.Driver
	_ = d
}
`
			expect = `package main

import (
	"context"
	"go/ast"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	_ = context.Background()
	_ = ast.NewIdent("x")
	var d mysql.Driver
	_ = d
}
`
		)
		t.AssertNil(gfile.PutContents(filePath, source))
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})
}
