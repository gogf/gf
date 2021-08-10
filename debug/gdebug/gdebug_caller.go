// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"fmt"
	"github.com/gogf/gf/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

const (
	maxCallerDepth = 1000
	stackFilterKey = "/debug/gdebug/gdebug"
)

var (
	goRootForFilter  = runtime.GOROOT() // goRootForFilter is used for stack filtering purpose.
	binaryVersion    = ""               // The version of current running binary(uint64 hex).
	binaryVersionMd5 = ""               // The version of current running binary(MD5).
	selfPath         = ""               // Current running binary absolute path.
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
	// Initialize internal package variable: selfPath.
	selfPath, _ := exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

// Caller returns the function name and the absolute file path along with its line
// number of the caller.
func Caller(skip ...int) (function string, path string, line int) {
	return CallerWithFilter("", skip...)
}

// CallerWithFilter returns the function name and the absolute file path along with
// its line number of the caller.
//
// The parameter `filter` is used to filter the path of the caller.
func CallerWithFilter(filter string, skip ...int) (function string, path string, line int) {
	var (
		number = 0
		ok     = true
	)
	if len(skip) > 0 {
		number = skip[0]
	}
	pc, file, line, start := callerFromIndex([]string{filter})
	if start != -1 {
		for i := start + number; i < maxCallerDepth; i++ {
			if i != start {
				pc, file, line, ok = runtime.Caller(i)
			}
			if ok {
				function := ""
				if fn := runtime.FuncForPC(pc); fn == nil {
					function = "unknown"
				} else {
					function = fn.Name()
				}
				return function, file, line
			} else {
				break
			}
		}
	}
	return "", "", -1
}

// callerFromIndex returns the caller position and according information exclusive of the
// debug package.
//
// VERY NOTE THAT, the returned index value should be <index - 1> as the caller's start point.
func callerFromIndex(filters []string) (pc uintptr, file string, line int, index int) {
	var filtered, ok bool
	for index = 0; index < maxCallerDepth; index++ {
		if pc, file, line, ok = runtime.Caller(index); ok {
			filtered = false
			for _, filter := range filters {
				if filter != "" && strings.Contains(file, filter) {
					filtered = true
					break
				}
			}
			if filtered {
				continue
			}
			if !utils.IsDebugEnabled() {
				if strings.Contains(file, utils.StackFilterKeyForGoFrame) {
					continue
				}
			} else {
				if strings.Contains(file, stackFilterKey) {
					continue
				}
			}
			if index > 0 {
				index--
			}
			return
		}
	}
	return 0, "", -1, -1
}

// CallerPackage returns the package name of the caller.
func CallerPackage() string {
	function, _, _ := Caller()
	indexSplit := strings.LastIndexByte(function, '/')
	if indexSplit == -1 {
		return function[:strings.IndexByte(function, '.')]
	} else {
		leftPart := function[:indexSplit+1]
		rightPart := function[indexSplit+1:]
		indexDot := strings.IndexByte(function, '.')
		rightPart = rightPart[:indexDot-1]
		return leftPart + rightPart
	}
}

// CallerFunction returns the function name of the caller.
func CallerFunction() string {
	function, _, _ := Caller()
	function = function[strings.LastIndexByte(function, '/')+1:]
	function = function[strings.IndexByte(function, '.')+1:]
	return function
}

// CallerFilePath returns the file path of the caller.
func CallerFilePath() string {
	_, path, _ := Caller()
	return path
}

// CallerDirectory returns the directory of the caller.
func CallerDirectory() string {
	_, path, _ := Caller()
	return filepath.Dir(path)
}

// CallerFileLine returns the file path along with the line number of the caller.
func CallerFileLine() string {
	_, path, line := Caller()
	return fmt.Sprintf(`%s:%d`, path, line)
}

// CallerFileLineShort returns the file name along with the line number of the caller.
func CallerFileLineShort() string {
	_, path, line := Caller()
	return fmt.Sprintf(`%s:%d`, filepath.Base(path), line)
}

// FuncPath returns the complete function path of given `f`.
func FuncPath(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// FuncName returns the function name of given `f`.
func FuncName(f interface{}) string {
	path := FuncPath(f)
	if path == "" {
		return ""
	}
	index := strings.LastIndexByte(path, '/')
	if index < 0 {
		index = strings.LastIndexByte(path, '\\')
	}
	return path[index+1:]
}
