// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"testing"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_NotFound(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		teatFile := gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/readline/error.log"
		callback := func(line string) error {
			return nil
		}
		err := gfile.ReadLines(teatFile, callback)
		t.AssertNE(err, nil)
	})
}

func Test_ReadLines(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			expectList = []string{"a", "b", "c", "d", "e"}
			getList    = make([]string, 0)
			callback   = func(line string) error {
				getList = append(getList, line)
				return nil
			}
			teatFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/readline/file.log"
		)
		err := gfile.ReadLines(teatFile, callback)
		t.AssertEQ(getList, expectList)
		t.AssertEQ(err, nil)
	})
}

func Test_ReadLines_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			callback = func(line string) error {
				return gerror.New("custom error")
			}
			teatFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/readline/file.log"
		)
		err := gfile.ReadLines(teatFile, callback)
		t.AssertEQ(err.Error(), "custom error")
	})
}

func Test_ReadLinesBytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			expectList = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
			getList    = make([][]byte, 0)
			callback   = func(line []byte) error {
				getList = append(getList, line)
				return nil
			}
			teatFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/readline/file.log"
		)
		err := gfile.ReadLinesBytes(teatFile, callback)
		t.AssertEQ(getList, expectList)
		t.AssertEQ(err, nil)
	})
}

func Test_ReadLinesBytes_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			callback = func(line []byte) error {
				return gerror.New("custom error")
			}
			teatFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/readline/file.log"
		)
		err := gfile.ReadLinesBytes(teatFile, callback)
		t.AssertEQ(err.Error(), "custom error")
	})
}
