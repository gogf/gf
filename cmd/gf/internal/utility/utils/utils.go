// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"context"
	"fmt"
	"golang.org/x/tools/imports"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// GoFmt formats the source file and adds or removes import statements as necessary.
func GoFmt(path string) {
	replaceFunc := func(path, content string) string {
		res, err := imports.Process(path, []byte(content), nil)
		if err != nil {
			mlog.Printf(`error format "%s" go files: %v`, path, err)
			return content
		}
		return string(res)
	}

	var err error
	if gfile.IsFile(path) {
		// File format.
		if gfile.ExtName(path) != "go" {
			return
		}
		err = gfile.ReplaceFileFunc(replaceFunc, path)
	} else {
		// Folder format.
		err = gfile.ReplaceDirFunc(replaceFunc, path, "*.go", true)
	}
	if err != nil {
		mlog.Printf(`error format "%s" go files: %v`, path, err)
	}
}

// GoModTidy executes `go mod tidy` at specified directory `dirPath`.
func GoModTidy(ctx context.Context, dirPath string) error {
	command := fmt.Sprintf(`cd %s && go mod tidy`, dirPath)
	err := gproc.ShellRun(ctx, command)
	return err
}

// IsFileDoNotEdit checks and returns whether file contains `do not edit` key.
func IsFileDoNotEdit(filePath string) bool {
	if !gfile.Exists(filePath) {
		return true
	}
	return gstr.Contains(gfile.GetContents(filePath), consts.DoNotEditKey)
}

// ReplaceGeneratedContentGFV2 replaces generated go content from goframe v1 to v2.
func ReplaceGeneratedContentGFV2(folderPath string) (err error) {
	return gfile.ReplaceDirFunc(func(path, content string) string {
		if gstr.Contains(content, `"github.com/gogf/gf`) && !gstr.Contains(content, `"github.com/gogf/gf/v2`) {
			content = gstr.Replace(content, `"github.com/gogf/gf"`, `"github.com/gogf/gf/v2"`)
			content = gstr.Replace(content, `"github.com/gogf/gf/`, `"github.com/gogf/gf/v2/`)
			content = gstr.Replace(content, `"github.com/gogf/gf/v2/contrib/`, `"github.com/gogf/gf/contrib/`)
			return content
		}
		return content
	}, folderPath, "*.go", true)
}

// GetImportPath calculates and returns the golang import path for given `filePath`.
// Note that it needs a `go.mod` in current working directory or parent directories to detect the path.
func GetImportPath(filePath string) string {
	// If `filePath` does not exist, create it firstly to find the import path.
	var realPath = gfile.RealPath(filePath)
	if realPath == "" {
		_ = gfile.Mkdir(filePath)
		realPath = gfile.RealPath(filePath)
	}

	var (
		newDir     = gfile.Dir(realPath)
		oldDir     string
		suffix     string
		goModName  = "go.mod"
		goModPath  string
		importPath string
	)

	if gfile.IsDir(filePath) {
		suffix = gfile.Basename(filePath)
	}
	for {
		goModPath = gfile.Join(newDir, goModName)
		if gfile.Exists(goModPath) {
			match, _ := gregex.MatchString(`^module\s+(.+)\s*`, gfile.GetContents(goModPath))
			importPath = gstr.Trim(match[1]) + "/" + suffix
			importPath = gstr.Replace(importPath, `\`, `/`)
			importPath = gstr.TrimRight(importPath, `/`)
			return importPath
		}
		oldDir = newDir
		newDir = gfile.Dir(oldDir)
		if newDir == oldDir {
			return ""
		}
		suffix = gfile.Basename(oldDir) + "/" + suffix
	}
}

// GetModPath retrieves and returns the file path of go.mod for current project.
func GetModPath() string {
	var (
		oldDir    = gfile.Pwd()
		newDir    = gfile.Dir(oldDir)
		goModName = "go.mod"
		goModPath string
	)
	for {
		goModPath = gfile.Join(newDir, goModName)
		if gfile.Exists(goModPath) {
			return goModPath
		}
		oldDir = newDir
		newDir = gfile.Dir(oldDir)
		if newDir == oldDir {
			break
		}
	}
	return ""
}
