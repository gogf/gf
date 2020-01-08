// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	list, err := doScanDir(path, pattern, isRecursive, false)
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
	list, err := doScanDir(path, pattern, isRecursive, true)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// doScanDir is an internal method which scans directory
// and returns the absolute path list of files that are not sorted.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter <recursive> is true.
//
// It returns only files except folders if <onlyFile> is true.
func doScanDir(path string, pattern string, recursive bool, onlyFile bool) ([]string, error) {
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
	filePath := ""
	isDir := false
	patterns := strings.Split(pattern, ",")
	for i := 0; i < len(patterns); i++ {
		patterns[i] = strings.TrimSpace(patterns[i])
	}
	for _, name := range names {
		filePath = path + Separator + name
		isDir = IsDir(filePath)
		if isDir && recursive {
			array, _ := doScanDir(filePath, pattern, true, onlyFile)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// It returns only files.
		if isDir && onlyFile {
			continue
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
