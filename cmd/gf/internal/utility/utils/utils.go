package utils

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/tools/imports"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
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

// IsFileDoNotEdit checks and returns whether file contains `do not edit` key.
func IsFileDoNotEdit(filePath string) bool {
	if !gfile.Exists(filePath) {
		return true
	}
	return gstr.Contains(gfile.GetContents(filePath), consts.DoNotEditKey)
}
