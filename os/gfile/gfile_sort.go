// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"strings"

	"github.com/gogf/gf/container/garray"
)

// fileSortFunc is the comparison function for files.
// It sorts the array in order of: directory -> file.
// If <path1> and <path2> are the same type, it then sorts them as strings.
func fileSortFunc(path1, path2 string) int {
	isDirPath1 := IsDir(path1)
	isDirPath2 := IsDir(path2)
	if isDirPath1 && !isDirPath2 {
		return -1
	}
	if !isDirPath1 && isDirPath2 {
		return 1
	}
	if n := strings.Compare(path1, path2); n != 0 {
		return n
	} else {
		return -1
	}
}

// SortFiles sorts the <files> in order of: directory -> file.
// Note that the item of <files> should be absolute path.
func SortFiles(files []string) []string {
	array := garray.NewSortedStrArrayComparator(fileSortFunc)
	array.Add(files...)
	return array.Slice()
}
