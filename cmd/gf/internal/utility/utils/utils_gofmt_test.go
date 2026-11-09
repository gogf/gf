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
			source   = `package main

import (
"fmt"
"os"
)

func main() {
fmt.Println()
}
`
			expect = `package main

import (
	"fmt"
)

func main() {
	fmt.Println()
}
`
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
			source   = `package main

import (
fm "fmt"
"os"
)

func main() {
fm.Println()
}
`
			expect = `package main

import (
	fm "fmt"
)

func main() {
	fm.Println()
}
`
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

	gtest.C(t, func(t *gtest.T) {
		// The blank-line-separated groups written by hand are kept, and only the imports
		// inside every group are sorted. gf gen service formats the main.go of the user,
		// whose import block usually contains three groups: the blank imports, the gf
		// packages and the project packages. Merging them into one sorted group as this
		// tool used to do is unwanted, as it rewrites the style of the user's code.
		var (
			dir      = t.TempDir()
			filePath = gfile.Join(dir, "main.go")
			source   = `package main

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/frame/g"

	"demo/internal/logic"
	"demo/internal/cmd"
)

func main(ctx context.Context) {
	g.Log().Print(ctx, cmd.Name, logic.Name, gctx.CtxId(ctx))
}
`
			// The specs inside every group are sorted, the groups are kept as is: the
			// standard library, the blank import, the gf packages and the project packages
			// are never merged into a single group.
			expect = `package main

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"demo/internal/cmd"
	"demo/internal/logic"
)

func main(ctx context.Context) {
	g.Log().Print(ctx, cmd.Name, logic.Name, gctx.CtxId(ctx))
}
`
		)
		t.AssertNil(gfile.PutContents(filePath, source))
		utils.GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})
}
