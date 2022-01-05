// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"os"
	"time"
)

// MTime returns the modification time of file given by `path` in second.
func MTime(path string) time.Time {
	s, e := os.Stat(path)
	if e != nil {
		return time.Time{}
	}
	return s.ModTime()
}

// MTimestamp returns the modification time of file given by `path` in second.
func MTimestamp(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.Unix()
}

// MTimestampMilli returns the modification time of file given by `path` in millisecond.
func MTimestampMilli(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.UnixNano() / 1000000
}
