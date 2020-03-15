// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"os"
)

// MTime returns the modification time of file given by <path> in second.
func MTime(path string) int64 {
	s, e := os.Stat(path)
	if e != nil {
		return 0
	}
	return s.ModTime().Unix()
}

// MTimeMillisecond returns the modification time of file given by <path> in millisecond.
func MTimeMillisecond(path string) int64 {
	s, e := os.Stat(path)
	if e != nil {
		return 0
	}
	return s.ModTime().UnixNano() / 1000000
}
