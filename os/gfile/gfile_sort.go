// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"slices"
	"strings"
)

// fileSortFunc is the comparison function for files.
// It sorts the array in order of: directory -> file.
// If `path1` and `path2` are the same type, it then sorts them as strings.
func fileSortFunc(path1, path2 string) int {
	isDirPath1 := IsDir(path1)
	isDirPath2 := IsDir(path2)
	if isDirPath1 != isDirPath2 {
		if isDirPath1 {
			return -1
		}
		return 1
	}
	return strings.Compare(path1, path2)
}

// SortFiles sorts the `files` in order of: directory -> file.
// Note that the item of `files` should be absolute path.
func SortFiles(files []string) []string {
	sortedFiles := make([]string, len(files))
	copy(sortedFiles, files)
	slices.SortFunc(sortedFiles, fileSortFunc)
	return sortedFiles
}
