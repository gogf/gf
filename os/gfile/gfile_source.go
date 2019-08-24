// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"os"
	"runtime"
	"strings"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

var (
	// goRootForFilter is used for stack filtering purpose.
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

// MainPkgPath returns absolute file path of package main,
// which contains the entrance function main.
//
// It's only available in develop environment.
//
// Note1: Only valid for source development environments,
// IE only valid for systems that generate this executable.
//
// Note2: When the method is called for the first time, if it is in an asynchronous goroutine,
// the method may not get the main package path.
func MainPkgPath() string {
	// Only for source development environments.
	if goRootForFilter == "" {
		return ""
	}
	path := mainPkgPath.Val()
	if path != "" {
		return path
	}
	for i := 1; i < 10000; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			// <file> is separated by '/'
			if gstr.Contains(file, "/github.com/gogf/gf/") &&
				!gstr.Contains(file, "/github.com/gogf/gf/.example/") {
				continue
			}
			if Ext(file) != ".go" {
				continue
			}
			// separator of <file> '/' will be converted to Separator.
			for path = Dir(file); len(path) > 1 && Exists(path) && path[len(path)-1] != os.PathSeparator; {
				files, _ := ScanDir(path, "*.go")
				for _, v := range files {
					if gregex.IsMatchString(`package\s+main`, GetContents(v)) {
						mainPkgPath.Set(path)
						return path
					}
				}
				path = Dir(path)
			}

		} else {
			break
		}
	}
	return ""
}
