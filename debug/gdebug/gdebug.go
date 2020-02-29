// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gdebug contains facilities for programs to debug themselves while they are running.
package gdebug

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gogf/gf/encoding/ghash"

	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/os/gfile"
)

const (
	gMAX_DEPTH  = 1000
	gFILTER_KEY = "/debug/gdebug/gdebug"
)

var (
	goRootForFilter  = runtime.GOROOT() // goRootForFilter is used for stack filtering purpose.
	binaryVersion    = ""               // The version of current running binary(uint64 hex).
	binaryVersionMd5 = ""               // The version of current running binary(MD5).
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

// BinVersion returns the version of current running binary.
// It uses ghash.BKDRHash+BASE36 algorithm to calculate the unique version of the binary.
func BinVersion() string {
	if binaryVersion == "" {
		binaryVersion = strconv.FormatInt(
			int64(ghash.BKDRHash(gfile.GetBytes(gfile.SelfPath()))),
			36,
		)
	}
	return binaryVersion
}

// BinVersionMd5 returns the version of current running binary.
// It uses MD5 algorithm to calculate the unique version of the binary.
func BinVersionMd5() string {
	if binaryVersionMd5 == "" {
		binaryVersionMd5, _ = gmd5.EncryptFile(gfile.SelfPath())
	}
	return binaryVersionMd5
}

// PrintStack prints to standard error the stack trace returned by runtime.Stack.
func PrintStack(skip ...int) {
	fmt.Print(Stack(skip...))
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack(skip ...int) string {
	return StackWithFilter("", skip...)
}

// StackWithFilter returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// The parameter <filter> is used to filter the path of the caller.
func StackWithFilter(filter string, skip ...int) string {
	return StackWithFilters([]string{filter}, skip...)
}

// StackWithFilters returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// The parameter <filters> is a slice of strings, which are used to filter the path of the caller.
func StackWithFilters(filters []string, skip ...int) string {
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	name := ""
	space := "  "
	index := 1
	buffer := bytes.NewBuffer(nil)
	filtered := false
	ok := true
	pc, file, line, start := callerFromIndex(filters)
	for i := start + number; i < gMAX_DEPTH; i++ {
		if i != start {
			pc, file, line, ok = runtime.Caller(i)
		}
		if ok {
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
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
			if strings.Contains(file, gFILTER_KEY) {
				continue
			}
			if fn := runtime.FuncForPC(pc); fn == nil {
				name = "unknown"
			} else {
				name = fn.Name()
			}
			if index > 9 {
				space = " "
			}
			buffer.WriteString(fmt.Sprintf("%d.%s%s\n    %s:%d\n", index, space, name, file, line))
			index++
		} else {
			break
		}
	}
	return buffer.String()
}

// CallerPath returns the function name and the absolute file path along with its line number of the caller.
func Caller(skip ...int) (function string, path string, line int) {
	return CallerWithFilter("", skip...)
}

// CallerPathWithFilter returns the function name and the absolute file path along with its line number of the caller.
//
// The parameter <filter> is used to filter the path of the caller.
func CallerWithFilter(filter string, skip ...int) (function string, path string, line int) {
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	ok := true
	pc, file, line, start := callerFromIndex([]string{filter})
	if start != -1 {
		for i := start + number; i < gMAX_DEPTH; i++ {
			if i != start {
				pc, file, line, ok = runtime.Caller(i)
			}
			if ok {
				if filter != "" && strings.Contains(file, filter) {
					continue
				}
				if strings.Contains(file, gFILTER_KEY) {
					continue
				}
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

// callerFromIndex returns the caller position and according information exclusive of the debug package.
func callerFromIndex(filters []string) (pc uintptr, file string, line int, index int) {
	var filtered, ok bool
	for index = 0; index < gMAX_DEPTH; index++ {
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
			if strings.Contains(file, gFILTER_KEY) {
				continue
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

// FuncPath returns the complete function path of given <f>.
func FuncPath(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// FuncName returns the function name of given <f>.
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
