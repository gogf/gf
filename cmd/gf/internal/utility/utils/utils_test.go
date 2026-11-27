// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file provides tests for the general utilities of the gf CLI, including the source
// formatting behavior that the code generators rely on.

package utils_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func Test_GetModPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		goModPath := utils.GetModPath()
		fmt.Println(goModPath)
	})
}

// Test_GoFmt_RemoveUnusedImport asserts that GoFmt drops the imports the file does not use,
// which is what the generators rely on for the content they render from a template.
func Test_GoFmt_RemoveUnusedImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path    = gfile.Join(t.TempDir(), "example.go")
			content = `package example

import (
	"fmt"
	"strings"
)

// Hello returns a greeting.
func Hello() string {
	return fmt.Sprint("hello")
}
`
		)
		t.AssertNil(gfile.PutContents(path, content))
		utils.GoFmt(path)
		var formatted = gfile.GetContents(path)
		t.Assert(gstr.Contains(formatted, `"fmt"`), true)
		t.Assert(gstr.Contains(formatted, `"strings"`), false)
	})
}

// Test_GoFmt_AddMissingImport asserts that GoFmt adds the imports the file references but does
// not import. The generators write their content without an import block, so GoFmt has to add
// the imports back: this is the capability a "remove only" implementation cannot provide.
func Test_GoFmt_AddMissingImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path    = gfile.Join(t.TempDir(), "example.go")
			content = `package example

// Hello returns a formatted greeting.
func Hello() string {
	return fmt.Sprint(strings.TrimSpace(" hello "))
}
`
		)
		t.AssertNil(gfile.PutContents(path, content))
		utils.GoFmt(path)
		var formatted = gfile.GetContents(path)
		t.Assert(gstr.Contains(formatted, `"fmt"`), true)
		t.Assert(gstr.Contains(formatted, `"strings"`), true)
	})
}

// Test_GoFmt_AddModuleImport asserts that GoFmt resolves the package name of a module import
// and adds it with the real package name, for example the `gerror` package of an import path
// whose last segment is not the package name of the top level.
func Test_GoFmt_AddModuleImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The file has to live inside the module, as the imports of a module package are
		// only resolved with a go.mod context.
		var (
			dir  = gfile.Join(gfile.Pwd(), "tmp_gofmt_"+guid.S())
			path = gfile.Join(dir, "example.go")
		)
		defer gfile.Remove(dir)
		t.AssertNil(gfile.Mkdir(dir))
		t.AssertNil(gfile.PutContents(path, `package example

// NewError returns a wrapped error.
func NewError(err error) error {
	return gerror.Wrap(err, "wrapped")
}
`))
		utils.GoFmt(path)
		var formatted = gfile.GetContents(path)
		t.Assert(gstr.Contains(formatted, `"github.com/gogf/gf/v2/errors/gerror"`), true)
	})
}

// Test_GoFmt_Folder asserts that GoFmt formats every go file of a folder recursively, and that
// the files with another extension are left untouched.
func Test_GoFmt_Folder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dir        = t.TempDir()
			nestedDir  = filepath.Join(dir, "nested")
			goPath     = filepath.Join(dir, "root.go")
			nestedPath = filepath.Join(nestedDir, "nested.go")
			textPath   = filepath.Join(dir, "readme.txt")
			content    = `package example

import (
	"fmt"
	"strings"
)

func Hello() string {
	return fmt.Sprint("hello")
}
`
		)
		t.AssertNil(gfile.Mkdir(nestedDir))
		t.AssertNil(gfile.PutContents(goPath, content))
		t.AssertNil(gfile.PutContents(nestedPath, gstr.Replace(content, "package example", "package nested")))
		t.AssertNil(gfile.PutContents(textPath, content))

		utils.GoFmt(dir)

		for _, path := range []string{goPath, nestedPath} {
			var formatted = gfile.GetContents(path)
			t.Assert(gstr.Contains(formatted, `"strings"`), false)
		}
		// A non-go file is never formatted.
		t.Assert(gfile.GetContents(textPath), content)
	})
}

// Test_GoFmt_InvalidSource asserts that an unparsable file is left untouched instead of being
// written back in a broken form.
func Test_GoFmt_InvalidSource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path    = gfile.Join(t.TempDir(), "example.go")
			content = "package example\n\nfunc Broken( {\n"
		)
		t.AssertNil(gfile.PutContents(path, content))
		utils.GoFmt(path)
		t.Assert(gfile.GetContents(path), content)
	})
}

// Test_GoFmt_NotGoFile asserts that a file with an extension other than "go" is skipped.
func Test_GoFmt_NotGoFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path    = gfile.Join(t.TempDir(), "example.txt")
			content = "package example\nimport \"fmt\"\n"
		)
		t.AssertNil(gfile.PutContents(path, content))
		utils.GoFmt(path)
		t.Assert(gfile.GetContents(path), content)
	})
}

// Test_GoFmt_Idempotent asserts that formatting an already formatted file twice keeps the
// content unchanged, so that re-running a generator does not produce a diff.
func Test_GoFmt_Idempotent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var path = gfile.Join(t.TempDir(), "example.go")
		t.AssertNil(gfile.PutContents(path, `package example

import (
	"fmt"
	"strings"
)

func Hello() string {
	return fmt.Sprint(strings.TrimSpace(" hello "))
}
`))
		utils.GoFmt(path)
		var once = gfile.GetContents(path)
		utils.GoFmt(path)
		t.Assert(gfile.GetContents(path), once)
	})
}
