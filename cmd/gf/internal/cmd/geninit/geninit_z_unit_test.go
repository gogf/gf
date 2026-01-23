// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package geninit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_ParseGitURL_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test basic github URL
		info, err := ParseGitURL("github.com/gogf/gf")
		t.AssertNil(err)
		t.Assert(info.Host, "github.com")
		t.Assert(info.Owner, "gogf")
		t.Assert(info.Repo, "gf")
		t.Assert(info.SubPath, "")
		t.Assert(info.Branch, "main")
		t.Assert(info.CloneURL, "https://github.com/gogf/gf.git")
	})
}

func Test_ParseGitURL_WithHTTPS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test URL with https prefix
		info, err := ParseGitURL("https://github.com/gogf/gf")
		t.AssertNil(err)
		t.Assert(info.Host, "github.com")
		t.Assert(info.Owner, "gogf")
		t.Assert(info.Repo, "gf")
	})
}

func Test_ParseGitURL_WithGitSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test URL with .git suffix
		info, err := ParseGitURL("github.com/gogf/gf.git")
		t.AssertNil(err)
		t.Assert(info.Host, "github.com")
		t.Assert(info.Owner, "gogf")
		t.Assert(info.Repo, "gf")
	})
}

func Test_ParseGitURL_WithSubPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test URL with subdirectory
		info, err := ParseGitURL("github.com/gogf/examples/httpserver/jwt")
		t.AssertNil(err)
		t.Assert(info.Host, "github.com")
		t.Assert(info.Owner, "gogf")
		t.Assert(info.Repo, "examples")
		t.Assert(info.SubPath, "httpserver/jwt")
		t.Assert(info.CloneURL, "https://github.com/gogf/examples.git")
	})
}

func Test_ParseGitURL_WithTreeBranch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test GitHub web URL with /tree/branch/
		info, err := ParseGitURL("github.com/gogf/examples/tree/develop/httpserver/jwt")
		t.AssertNil(err)
		t.Assert(info.Host, "github.com")
		t.Assert(info.Owner, "gogf")
		t.Assert(info.Repo, "examples")
		t.Assert(info.Branch, "develop")
		t.Assert(info.SubPath, "httpserver/jwt")
	})
}

func Test_ParseGitURL_WithVersion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test URL with version suffix
		info, err := ParseGitURL("github.com/gogf/gf/cmd/gf/v2@v2.9.7")
		t.AssertNil(err)
		t.Assert(info.Host, "github.com")
		t.Assert(info.Owner, "gogf")
		t.Assert(info.Repo, "gf")
		t.Assert(info.SubPath, "cmd/gf/v2")
	})
}

func Test_ParseGitURL_Invalid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test invalid URL (too short)
		_, err := ParseGitURL("github.com/gogf")
		t.AssertNE(err, nil)
	})
}

func Test_IsSubdirRepo_NotSubdir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Standard Go module paths should not be detected as subdirectory
		t.Assert(IsSubdirRepo("github.com/gogf/gf"), false)
		t.Assert(IsSubdirRepo("github.com/gogf/gf/v2"), false)
	})
}

func Test_IsSubdirRepo_GoModuleWithCmd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Go module paths with common patterns should not be detected as subdirectory
		t.Assert(IsSubdirRepo("github.com/gogf/gf/cmd/gf/v2"), false)
		t.Assert(IsSubdirRepo("github.com/gogf/gf/contrib/drivers/mysql/v2"), false)
	})
}

func Test_IsSubdirRepo_ActualSubdir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Actual subdirectories should be detected
		t.Assert(IsSubdirRepo("github.com/gogf/examples/httpserver/jwt"), true)
		t.Assert(IsSubdirRepo("github.com/gogf/examples/grpc/basic"), true)
	})
}

func Test_GetModuleNameFromGoMod_Valid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create temp directory with go.mod
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Write go.mod file
		goModContent := `module github.com/test/myproject

go 1.21

require (
	github.com/gogf/gf/v2 v2.9.0
)
`
		err = gfile.PutContents(filepath.Join(tempDir, "go.mod"), goModContent)
		t.AssertNil(err)

		// Test extraction
		moduleName := GetModuleNameFromGoMod(tempDir)
		t.Assert(moduleName, "github.com/test/myproject")
	})
}

func Test_GetModuleNameFromGoMod_NoFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create temp directory without go.mod
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Test extraction - should return empty
		moduleName := GetModuleNameFromGoMod(tempDir)
		t.Assert(moduleName, "")
	})
}

func Test_GetModuleNameFromGoMod_SimpleModule(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create temp directory with simple go.mod
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Write simple go.mod file
		goModContent := `module main

go 1.21
`
		err = gfile.PutContents(filepath.Join(tempDir, "go.mod"), goModContent)
		t.AssertNil(err)

		// Test extraction
		moduleName := GetModuleNameFromGoMod(tempDir)
		t.Assert(moduleName, "main")
	})
}

func Test_ASTReplacer_ReplaceInFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create temp directory
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Create a Go file with imports
		goFileContent := `package main

import (
	"fmt"

	"github.com/old/module/internal/service"
	"github.com/old/module/pkg/utils"
	"github.com/other/package"
)

func main() {
	fmt.Println("Hello")
}
`
		goFilePath := filepath.Join(tempDir, "main.go")
		err = gfile.PutContents(goFilePath, goFileContent)
		t.AssertNil(err)

		// Replace imports
		replacer := NewASTReplacer("github.com/old/module", "github.com/new/project")
		err = replacer.ReplaceInFile(context.Background(), goFilePath)
		t.AssertNil(err)

		// Verify replacement
		content := gfile.GetContents(goFilePath)
		t.Assert(gfile.Exists(goFilePath), true)

		// Check that old imports are replaced
		t.AssertNE(content, "")
		t.Assert(contains(content, `"github.com/new/project/internal/service"`), true)
		t.Assert(contains(content, `"github.com/new/project/pkg/utils"`), true)

		// Check that other imports are not affected
		t.Assert(contains(content, `"github.com/other/package"`), true)
		t.Assert(contains(content, `"fmt"`), true)
	})
}

func Test_ASTReplacer_ReplaceInDir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create temp directory structure
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Create subdirectory
		subDir := filepath.Join(tempDir, "sub")
		err = gfile.Mkdir(subDir)
		t.AssertNil(err)

		// Create main.go
		mainContent := `package main

import "github.com/old/module/sub"

func main() {
	sub.Hello()
}
`
		err = gfile.PutContents(filepath.Join(tempDir, "main.go"), mainContent)
		t.AssertNil(err)

		// Create sub/sub.go
		subContent := `package sub

import "github.com/old/module/pkg"

func Hello() {
	pkg.Do()
}
`
		err = gfile.PutContents(filepath.Join(subDir, "sub.go"), subContent)
		t.AssertNil(err)

		// Replace imports in directory
		replacer := NewASTReplacer("github.com/old/module", "github.com/new/project")
		err = replacer.ReplaceInDir(context.Background(), tempDir)
		t.AssertNil(err)

		// Verify main.go replacement
		mainResult := gfile.GetContents(filepath.Join(tempDir, "main.go"))
		t.Assert(contains(mainResult, `"github.com/new/project/sub"`), true)

		// Verify sub/sub.go replacement
		subResult := gfile.GetContents(filepath.Join(subDir, "sub.go"))
		t.Assert(contains(subResult, `"github.com/new/project/pkg"`), true)
	})
}

func Test_findGoFiles(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create temp directory structure
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Create subdirectories
		subDir := filepath.Join(tempDir, "sub")
		err = gfile.Mkdir(subDir)
		t.AssertNil(err)

		// Create various files
		err = gfile.PutContents(filepath.Join(tempDir, "main.go"), "package main")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(tempDir, "readme.md"), "# README")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(subDir, "sub.go"), "package sub")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(subDir, "data.json"), "{}")
		t.AssertNil(err)

		// Find Go files
		files, err := findGoFiles(tempDir)
		t.AssertNil(err)

		// Should find exactly 2 Go files
		t.Assert(len(files), 2)

		// Verify file names
		hasMain := false
		hasSub := false
		for _, f := range files {
			if filepath.Base(f) == "main.go" {
				hasMain = true
			}
			if filepath.Base(f) == "sub.go" {
				hasSub = true
			}
		}
		t.Assert(hasMain, true)
		t.Assert(hasSub, true)
	})
}

func Test_findGoFiles_EmptyDir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create empty temp directory
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Find Go files
		files, err := findGoFiles(tempDir)
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
