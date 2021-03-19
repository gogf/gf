// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"github.com/gogf/gf/text/gstr"
	"os"
	"runtime"
	"strings"

	"github.com/gogf/gf/text/gregex"
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
	// It is only for source development environments.
	if goRootForFilter == "" {
		return ""
	}
	path := mainPkgPath.Val()
	if path != "" {
		return path
	}
	var lastFile string
	for i := 1; i < 10000; i++ {
		if pc, file, _, ok := runtime.Caller(i); ok {
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if Ext(file) != ".go" {
				continue
			}
			lastFile = file
			// Check if it is called in package initialization function,
			// in which it here cannot retrieve main package path,
			// it so just returns that can make next check.
			if fn := runtime.FuncForPC(pc); fn != nil {
				array := gstr.Split(fn.Name(), ".")
				if array[0] != "main" {
					continue
				}
			}
			if gregex.IsMatchString(`package\s+main`, GetContents(file)) {
				mainPkgPath.Set(Dir(file))
				return Dir(file)
			}
		} else {
			break
		}
	}
	// If it still cannot find the path of the package main,
	// it recursively searches the directory and its parents directory of the last go file.
	// It's usually necessary for uint testing cases of business project.
	if lastFile != "" {
		for path = Dir(lastFile); len(path) > 1 && Exists(path) && path[len(path)-1] != os.PathSeparator; {
			files, _ := ScanDir(path, "*.go")
			for _, v := range files {
				if gregex.IsMatchString(`package\s+main`, GetContents(v)) {
					mainPkgPath.Set(path)
					return path
				}
			}
			path = Dir(path)
		}
	}
	return ""
}
