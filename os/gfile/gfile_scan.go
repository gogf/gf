// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gstr"
	"os"
	"path/filepath"
	"sort"
)

const (
	// Max recursive depth for directory scanning.
	gMAX_SCAN_DEPTH = 100000
)

// ScanDir returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
func ScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, nil)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFunc returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
//
// The pattern parameter <pattern> supports multiple file name patterns, using the ','
// symbol to separate multiple patterns.
//
// The parameter <recursive> specifies whether scanning the <path> recursively, which
// means it scans its sub-files and appends the files path to result array if the sub-file
// is also a folder. It is false in default.
//
// The parameter <handler> specifies the callback function handling each sub-file path of
// the <path> and its sub-folders. It ignores the sub-file path if <handler> returns an empty
// string, or else it appends the sub-file path to result slice.
func ScanDirFunc(path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	list, err := doScanDir(0, path, pattern, recursive, handler)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFile returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// Note that it returns only files, exclusive of directories.
func ScanDirFile(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, func(path string) string {
		if IsDir(path) {
			return ""
		}
		return path
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFileFunc returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
//
// The pattern parameter <pattern> supports multiple file name patterns, using the ','
// symbol to separate multiple patterns.
//
// The parameter <recursive> specifies whether scanning the <path> recursively, which
// means it scans its sub-files and appends the files path to result array if the sub-file
// is also a folder. It is false in default.
//
// The parameter <handler> specifies the callback function handling each sub-file path of
// the <path> and its sub-folders. It ignores the sub-file path if <handler> returns an empty
// string, or else it appends the sub-file path to result slice.
//
// Note that the parameter <path> for <handler> is not a directory but a file.
// It returns only files, exclusive of directories.
func ScanDirFileFunc(path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	list, err := doScanDir(0, path, pattern, recursive, func(path string) string {
		if IsDir(path) {
			return ""
		}
		return handler(path)
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// doScanDir is an internal method which scans directory and returns the absolute path
// list of files that are not sorted.
//
// The pattern parameter <pattern> supports multiple file name patterns, using the ','
// symbol to separate multiple patterns.
//
// The parameter <recursive> specifies whether scanning the <path> recursively, which
// means it scans its sub-files and appends the files path to result array if the sub-file
// is also a folder. It is false in default.
//
// The parameter <handler> specifies the callback function handling each sub-file path of
// the <path> and its sub-folders. It ignores the sub-file path if <handler> returns an empty
// string, or else it appends the sub-file path to result slice.
func doScanDir(depth int, path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	if depth >= gMAX_SCAN_DEPTH {
		return nil, gerror.Newf("directory scanning exceeds max recursive depth: %d", gMAX_SCAN_DEPTH)
	}
	list := ([]string)(nil)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	var (
		filePath = ""
		patterns = gstr.SplitAndTrim(pattern, ",")
	)
	for _, name := range names {
		filePath = path + Separator + name
		if IsDir(filePath) && recursive {
			array, _ := doScanDir(depth+1, filePath, pattern, true, handler)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// Handler filtering.
		if handler != nil {
			filePath = handler(filePath)
			if filePath == "" {
				continue
			}
		}
		// If it meets pattern, then add it to the result list.
		for _, p := range patterns {
			if match, err := filepath.Match(p, name); err == nil && match {
				filePath = Abs(filePath)
				if filePath != "" {
					list = append(list, filePath)
				}
			}
		}
	}
	return list, nil
}
