// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains unit tests for unused-import cleanup and package-name
// resolving used by GoFmt.

package utils

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// stubResolvePackageNames replaces the package name resolving with the given names and
// returns a function to restore the original one.
func stubResolvePackageNames(names map[string]string) func() {
	var original = resolvePackageNames
	resolvePackageNames = func(_ string, importPaths []string) map[string]string {
		var result = make(map[string]string, len(importPaths))
		for _, importPath := range importPaths {
			result[importPath] = names[importPath]
		}
		return result
	}
	return func() {
		resolvePackageNames = original
	}
}

// Test_GoFmt_ResolvedPackageNames asserts that unused imports are determined by the real
// package name instead of any segment of the import path.
func Test_GoFmt_ResolvedPackageNames(t *testing.T) {
	var restore = stubResolvePackageNames(map[string]string{
		"github.com/minio/minio-go/v7": "minio",
		"k8s.io/api/core/v1":           "v1",
		"go.etcd.io/etcd/client/v3":    "clientv3",
	})
	defer restore()

	gtest.C(t, func(t *gtest.T) {
		// A unused package whose name is not any segment of the import path is removed.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func main() {
	_ = context.Background()
}
`
			expect = `package main

import (
	"context"
)

func main() {
	_ = context.Background()
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// A package whose name is the last path segment keeps working when it is used.
		// See the k8s.io/api/core/v1 case: the package name is "v1" instead of "core".
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"k8s.io/api/core/v1"
)

func main() {
	var pod v1.Pod
	_ = pod
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})

	gtest.C(t, func(t *gtest.T) {
		// A package whose name is neither path segment keeps working when it is used.
		// See the go.etcd.io/etcd/client/v3 case: the package name is "clientv3".
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"go.etcd.io/etcd/client/v3"
)

func main() {
	_ = clientv3.Config{}
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})

	gtest.C(t, func(t *gtest.T) {
		// The import block inside a raw string is not an import block of the file, so it
		// must be kept as is instead of being grouped or sorted.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"fmt"
)

var template = ` + "`" + `import (
	"github.com/not/used"
	"os"
)` + "`" + `

func main() {
	fmt.Println(template)
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})
}

// Test_GoFmt_UnresolvedPackageName asserts that an import whose package name cannot be
// resolved is handled by the name assumed from its import path, as goimports does. The
// generated go files do not compile with an unused import left behind, so the unused ones
// are removed while the used ones are kept.
func Test_GoFmt_UnresolvedPackageName(t *testing.T) {
	var restore = stubResolvePackageNames(nil)
	defer restore()

	gtest.C(t, func(t *gtest.T) {
		// An unused import is removed by its assumed name, even though its real package
		// name is unknown.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"context"

	"github.com/foo/bar"
)

func main() {
	_ = context.Background()
}
`
			expect = `package main

import (
	"context"
)

func main() {
	_ = context.Background()
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// An import whose assumed name is in use is kept, as the package name of
		// "github.com/minio/minio-go/v7" is "minio" instead of any path segment.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
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
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})
}

// Test_GoFmt_ShadowedImport asserts that an unused import shadowed by a declaration of the
// same name inside the file is removed, as golang.org/x/tools/imports.Process does. An
// identifier bound to a declaration never refers to the imported package.
func Test_GoFmt_ShadowedImport(t *testing.T) {
	var restore = stubResolvePackageNames(map[string]string{
		"strings": "strings",
	})
	defer restore()

	gtest.C(t, func(t *gtest.T) {
		// The local variable shadows the import, so the import is unused and removed.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"strings"
)

type T struct {
	Name string
}

func main() {
	strings := T{}
	_ = strings.Name
}
`
			expect = `package main

type T struct {
	Name string
}

func main() {
	strings := T{}
	_ = strings.Name
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// The import is used in `a` while the local variable shadows it in `b` only, so the
		// import must be kept: the scopes of the two functions do not overlap.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"strings"
)

type T struct {
	Name string
}

func a() bool {
	return strings.Contains("a", "b")
}

func b() string {
	strings := T{}
	return strings.Name
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})

	gtest.C(t, func(t *gtest.T) {
		// The function parameter and the type switch variable shadow the import as well.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	"strings"
)

type T struct {
	Name string
}

func f(strings T) string {
	switch strings := any(strings).(type) {
	case T:
		return strings.Name
	}
	return ""
}

func g() bool {
	for strings := range []int{1} {
		_ = strings
	}
	return strings.Contains("a", "b")
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(gstr.Contains(out, `"strings"`), true)
	})
}

// Test_EnsureImports asserts that the required imports are added, that an existing import is
// left untouched, and that the import alias is respected.
func Test_EnsureImports(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// A file without any import declaration gets an import block.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

func main() {
	_ = grpcx.GrpcServer{}
	_ = context.Background()
}
`
			expect = `package main

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

func main() {
	_ = grpcx.GrpcServer{}
	_ = context.Background()
}
`
		)
		t.AssertNil(gfile.PutContents(filePath, source))
		changed, err := EnsureImports(
			filePath,
			ImportSpec{Path: "context"},
			ImportSpec{Path: "github.com/gogf/gf/contrib/rpc/grpcx/v2"},
		)
		t.AssertNil(err)
		t.Assert(changed, true)
		GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)

		// The imports already exist, so nothing is changed and the file is not written.
		t.AssertNil(gfile.PutContents(filePath, expect))
		changed, err = EnsureImports(
			filePath,
			ImportSpec{Path: "context"},
			ImportSpec{Path: "github.com/gogf/gf/contrib/rpc/grpcx/v2"},
		)
		t.AssertNil(err)
		t.Assert(changed, false)
		t.Assert(gfile.GetContents(filePath), expect)
	})

	gtest.C(t, func(t *gtest.T) {
		// The alias is used for the api package, as the package name is not the version.
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

type Controller struct {
	v1.UnimplementedAccountServer
}
`
			// A single import is rendered without parentheses, which is what gofmt does.
			expect = `package main

import v1 "github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3953/account/v1"

type Controller struct {
	v1.UnimplementedAccountServer
}
`
		)
		t.AssertNil(gfile.PutContents(filePath, source))
		changed, err := EnsureImports(filePath, ImportSpec{
			Name: "v1",
			Path: "github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3953/account/v1",
		})
		t.AssertNil(err)
		t.Assert(changed, true)
		t.Assert(gfile.GetContents(filePath), expect)
		// The added import is used, so GoFmt keeps it.
		GoFmt(filePath)
		t.Assert(gfile.GetContents(filePath), expect)
	})
}

// Test_EnsureImports_Error asserts the invalid arguments of EnsureImports are reported.
func Test_EnsureImports_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// No spec, nothing to do.
		changed, err := EnsureImports(gfile.Join(t.TempDir(), "example.go"))
		t.AssertNil(err)
		t.Assert(changed, false)

		// The file does not exist.
		_, err = EnsureImports(
			gfile.Join(t.TempDir(), "example.go"),
			ImportSpec{Path: "context"},
		)
		t.AssertNE(err, nil)

		// The import path is empty.
		var filePath = gfile.Join(t.TempDir(), "example.go")
		t.AssertNil(gfile.PutContents(filePath, "package main\n"))
		_, err = EnsureImports(filePath, ImportSpec{})
		t.AssertNE(err, nil)

		// The content is not valid go code.
		t.AssertNil(gfile.PutContents(filePath, "package main\n\nfunc {"))
		_, err = EnsureImports(filePath, ImportSpec{Path: "context"})
		t.AssertNE(err, nil)
	})
}

// Test_AssumedPackageName asserts the package names assumed from the import paths.
func Test_AssumedPackageName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(assumedPackageName("github.com/minio/minio-go/v7"), "minio")
		t.Assert(assumedPackageName("github.com/gogf/gf/v2"), "gf")
		t.Assert(assumedPackageName("go.etcd.io/etcd/client/v3"), "client")
		t.Assert(assumedPackageName("github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"), "service")
		t.Assert(assumedPackageName("k8s.io/api/core/v1"), "core")
		t.Assert(assumedPackageName("example.com/go-redis"), "redis")
		t.Assert(assumedPackageName("example.com/foo.bar"), "foo")
	})
}

// Test_GoFmt_DotImport asserts that dot imports are kept as is, as their usages cannot
// be determined. It is the same behavior as golang.org/x/tools/imports.Process.
func Test_GoFmt_DotImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = `package main

import (
	. "fmt"
	"os"
)

func main() {
	_ = os.Args
}
`
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})
}

// Test_GoFmt_RealPackageNames asserts the unused imports are removed by the package names
// resolved from the module, which is the case of gf gen service.
func Test_GoFmt_RealPackageNames(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The directory is created inside the module so that the package names are
		// resolvable by `go list`.
		var dir, err = os.MkdirTemp(".", "tmp-resolve-")
		t.AssertNil(err)

		var (
			filePath = gfile.Join(dir, "example.go")
			source   = `package main

import (
	"context"

	"golang.org/x/tools/go/packages"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	_ = context.Background()
	var d mysql.Driver
	_ = d
}
`
			// The unused `golang.org/x/tools/go/packages` is removed by its real name
			// "packages", and the used mysql driver is kept by its real name "mysql".
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
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, expect)
		t.AssertNil(os.RemoveAll(dir))
	})
}
