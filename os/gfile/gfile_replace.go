// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"github.com/gogf/gf/text/gstr"
)

// Replace replaces content for files under <path>.
// The parameter <pattern> specifies the file pattern which matches to be replaced.
// It does replacement recursively if given parameter <recursive> is true.
func Replace(search, replace, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = PutContents(file, gstr.Replace(GetContents(file), search, replace)); err != nil {
			return err
		}
	}
	return err
}

// ReplaceFunc replaces content for files under <path> with callback function <f>.
// The parameter <pattern> specifies the file pattern which matches to be replaced.
// It does replacement recursively if given parameter <recursive> is true.
func ReplaceFunc(f func(path, content string) string, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	data := ""
	result := ""
	for _, file := range files {
		data = GetContents(file)
		result = f(file, data)
		if data != result {
			if err = PutContents(file, result); err != nil {
				return err
			}
		}
	}
	return err
}
