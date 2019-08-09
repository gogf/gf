// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"os"
	"runtime"

	"github.com/gogf/gf/internal/debug"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

const (
	gPATH_FILTER_KEY = "/gf/os/gfile/gfile"
)

// SourcePath returns absolute file path of the current source file path.
//
// Note that it's only available in develop environment.
func SourcePath(skip ...int) string {
	_, path, _ := debug.CallerWithFilter(gPATH_FILTER_KEY, skip...)
	return path
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
	path := mainPkgPath.Val()
	if path != "" {
		if path == "-" {
			return ""
		}
		return path
	}
	for i := 1; i < 10000; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
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
	// If it fails finding the path, then mark it as "-",
	// which means it will never do this search again.
	mainPkgPath.Set("-")
	return ""
}
