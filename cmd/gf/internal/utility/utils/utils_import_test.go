// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
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
			source   = "package main\n\nimport (\n\t\"context\"\n\n\t\"github.com/minio/minio-go/v7\"\n)\n\nfunc main() {\n\t_ = context.Background()\n}\n"
			expect   = "package main\n\nimport (\n\t\"context\"\n)\n\nfunc main() {\n\t_ = context.Background()\n}\n"
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
			source   = "package main\n\nimport (\n\t\"k8s.io/api/core/v1\"\n)\n\nfunc main() {\n\tvar pod v1.Pod\n\t_ = pod\n}\n"
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
			source   = "package main\n\nimport (\n\t\"go.etcd.io/etcd/client/v3\"\n)\n\nfunc main() {\n\t_ = clientv3.Config{}\n}\n"
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
			source   = "package main\n\nimport (\n\t\"fmt\"\n)\n\nvar template = `import (\n\t\"github.com/not/used\"\n\t\"os\"\n)`\n\nfunc main() {\n\tfmt.Println(template)\n}\n"
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})
}

// Test_GoFmt_UnresolvedPackageName asserts that an import is kept when its package name
// cannot be resolved, even though it looks unused. The package name is not always a
// segment of the import path, so removing the import by a guessed name would break the
// generated code that references the package.
func Test_GoFmt_UnresolvedPackageName(t *testing.T) {
	var restore = stubResolvePackageNames(nil)
	defer restore()

	gtest.C(t, func(t *gtest.T) {
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			// The import path "github.com/foo/bar" would be guessed as package "bar",
			// but the real package name could be anything, for example "baz".
			source = "package main\n\nimport (\n\t\"context\"\n\n\t\"github.com/foo/bar\"\n)\n\nfunc main() {\n\t_ = context.Background()\n}\n"
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, source)
	})
}

// Test_GoFmt_DotImport asserts that dot imports are kept as is, as their usages cannot
// be determined. It is the same behavior as golang.org/x/tools/imports.Process.
func Test_GoFmt_DotImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filePath = gfile.Join(t.TempDir(), "example.go")
			source   = "package main\n\nimport (\n\t. \"fmt\"\n\t\"os\"\n)\n\nfunc main() {\n\t_ = os.Args\n}\n"
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
			source   = "package main\n\nimport (\n\t\"context\"\n\n\t\"golang.org/x/tools/go/packages\"\n\n\t\"github.com/gogf/gf/contrib/drivers/mysql/v2\"\n)\n\nfunc main() {\n\t_ = context.Background()\n\tvar d mysql.Driver\n\t_ = d\n}\n"
			// The unused `golang.org/x/tools/go/packages` is removed by its real name
			// "packages", and the used mysql driver is kept by its real name "mysql".
			expect = "package main\n\nimport (\n\t\"context\"\n\n\t\"github.com/gogf/gf/contrib/drivers/mysql/v2\"\n)\n\nfunc main() {\n\t_ = context.Background()\n\tvar d mysql.Driver\n\t_ = d\n}\n"
		)
		out, err := goFormatContent(filePath, source)
		t.AssertNil(err)
		t.Assert(out, expect)
		t.AssertNil(os.RemoveAll(dir))
	})
}
