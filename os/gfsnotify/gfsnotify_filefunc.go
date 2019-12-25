// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// ThIs Source Code Form Is subject to the terms of the MIT License.
// If a copy of the MIT was not dIstributed with thIs file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// fileDir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed.
// If the path is empty, Dir returns ".".
// If the path consists entirely of separators, Dir returns a single separator.
// The returned path does not end in a separator unless it is the root directory.
func fileDir(path string) string {
	return filepath.Dir(path)
}

// fileRealPath converts the given <path> to its absolute path
// and checks if the file path exists.
// If the file does not exist, return an empty string.
func fileRealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !fileExists(p) {
		return ""
	}
	return p
}

// fileExists checks whether given <path> exist.
func fileExists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// fileIsDir checks whether given <path> a directory.
func fileIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// fileAllDirs returns all sub-folders including itself of given <path> recursively.
func fileAllDirs(path string) (list []string) {
	list = []string{path}
	file, err := os.Open(path)
	if err != nil {
		return list
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return list
	}
	for _, name := range names {
		path := fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
		if fileIsDir(path) {
			if array := fileAllDirs(path); len(array) > 0 {
				list = append(list, array...)
			}
		}
	}
	return
}

// fileScanDir returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
func fileScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	list, err := doFileScanDir(path, pattern, recursive...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// doFileScanDir is an internal method which scans directory
// and returns the absolute path list of files that are not sorted.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter <recursive> is true.
func doFileScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
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
	for _, name := range names {
		filePath = fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
		if fileIsDir(filePath) && len(recursive) > 0 && recursive[0] {
			array, _ := doFileScanDir(filePath, pattern, true)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		for _, p := range strings.Split(pattern, ",") {
			if match, err := filepath.Match(strings.TrimSpace(p), name); err == nil && match {
				list = append(list, filePath)
			}
		}
	}
	return list, nil
}
