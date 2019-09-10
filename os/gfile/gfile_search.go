// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gogf/gf/container/garray"
)

// Search searches file by name <name> in following paths with priority:
// prioritySearchPaths, Pwd()、SelfDir()、MainPkgPath().
// It returns the absolute file path of <name> if found, or en empty string if not found.
func Search(name string, prioritySearchPaths ...string) (realPath string, err error) {
	// Check if it's a absolute path.
	realPath = RealPath(name)
	if realPath != "" {
		return
	}
	// Search paths array.
	array := garray.NewStrArray()
	array.Append(prioritySearchPaths...)
	array.Append(Pwd(), SelfDir())
	if path := MainPkgPath(); path != "" {
		array.Append(path)
	}
	// Remove repeated items.
	array.Unique()
	// Do the searching.
	array.RLockFunc(func(array []string) {
		path := ""
		for _, v := range array {
			path = RealPath(v + Separator + name)
			if path != "" {
				realPath = path
				break
			}
		}
	})
	// If it fails searching, it returns formatted error.
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		buffer.WriteString(fmt.Sprintf("cannot find file/folder \"%s\" in following paths:", name))
		array.RLockFunc(func(array []string) {
			for k, v := range array {
				buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
			}
		})
		err = errors.New(buffer.String())
	}
	return
}
