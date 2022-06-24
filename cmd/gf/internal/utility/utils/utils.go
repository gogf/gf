package utils

import (
	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/tools/imports"
)

// GoFmt formats the source file.
func GoFmt(path string) {
	if err := pretty(path, true); err != nil {
		mlog.Fatalf(`error format "%s" go files: %v`, path, err)
	}
}

// GoImports adds or removes import statements as necessary for the source file.
func GoImports(path string) {
	if err := pretty(path); err != nil {
		mlog.Fatalf(`error update "%s" go file imports: %v`, path, err)
	}
}

// IsFileDoNotEdit checks and returns whether file contains `do not edit` key.
func IsFileDoNotEdit(filePath string) bool {
	if !gfile.Exists(filePath) {
		return true
	}
	return gstr.Contains(gfile.GetContents(filePath), consts.DoNotEditKey)
}

// pretty format go file and adds or removes import statements as necessary.
func pretty(filePath string, formatOnly ...bool) error {
	var genOpt *imports.Options
	if len(formatOnly) > 0 {
		genOpt = &imports.Options{
			Comments:   true,
			TabIndent:  true,
			TabWidth:   8,
			FormatOnly: true,
		}
	}
	replaceFunc := func(path, content string) string {
		res, err := imports.Process(path, []byte(content), genOpt)
		if err != nil {
			mlog.Printf(`pretty go file "%s" failed: %v`, path, err)
			return content
		}
		return string(res)
	}
	if gfile.IsFile(filePath) {
		if gfile.ExtName(filePath) != "go" {
			return nil
		}
		return gfile.ReplaceFileFunc(replaceFunc, filePath)
	}
	return gfile.ReplaceDirFunc(replaceFunc, filePath, "*.go", true)
}
