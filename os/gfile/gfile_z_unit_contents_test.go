// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func createTestFile(filename, content string) error {
	TempDir := testpath()
	err := ioutil.WriteFile(TempDir+filename, []byte(content), 0666)
	return err
}

func delTestFiles(filenames string) {
	os.RemoveAll(testpath() + filenames)
}

func createDir(paths string) {
	TempDir := testpath()
	os.Mkdir(TempDir+paths, 0777)
}

func formatpaths(paths []string) []string {
	for k, v := range paths {
		paths[k] = filepath.ToSlash(v)
		paths[k] = strings.Replace(paths[k], "./", "/", 1)
	}

	return paths
}

func formatpath(paths string) string {
	paths = filepath.ToSlash(paths)
	paths = strings.Replace(paths, "./", "/", 1)
	return paths
}

func testpath() string {
	return gstr.TrimRight(os.TempDir(), "\\/")
}

func Test_GetContents(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		var (
			filepaths string = "/testfile_t1.txt"
		)
		createTestFile(filepaths, "my name is jroam")
		defer delTestFiles(filepaths)

		t.Assert(gfile.GetContents(testpath()+filepaths), "my name is jroam")
		t.Assert(gfile.GetContents(""), "")

	})
}

func Test_GetBinContents(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths1  = "/testfile_t1.txt"
			filepaths2  = testpath() + "/testfile_t1_no.txt"
			readcontent []byte
			str1        = "my name is jroam"
		)
		createTestFile(filepaths1, str1)
		defer delTestFiles(filepaths1)
		readcontent = gfile.GetBytes(testpath() + filepaths1)
		t.Assert(readcontent, []byte(str1))

		readcontent = gfile.GetBytes(filepaths2)
		t.Assert(string(readcontent), "")

		t.Assert(string(gfile.GetBytes(filepaths2)), "")

	})
}

func Test_Truncate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths1 = "/testfile_GetContentsyyui.txt"
			err        error
			files      *os.File
		)
		createTestFile(filepaths1, "abcdefghijkmln")
		defer delTestFiles(filepaths1)
		err = gfile.Truncate(testpath()+filepaths1, 10)
		t.AssertNil(err)

		files, err = os.Open(testpath() + filepaths1)
		t.AssertNil(err)
		defer files.Close()
		fileinfo, err2 := files.Stat()
		t.Assert(err2, nil)
		t.Assert(fileinfo.Size(), 10)

		err = gfile.Truncate("", 10)
		t.AssertNE(err, nil)

	})
}

func Test_PutContents(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)

		err = gfile.PutContents(testpath()+filepaths, "test!")
		t.AssertNil(err)

		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "test!")

		err = gfile.PutContents("", "test!")
		t.AssertNE(err, nil)

	})
}

func Test_PutContentsAppend(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)

		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)
		err = gfile.PutContentsAppend(testpath()+filepaths, "hello")
		t.AssertNil(err)

		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "ahello")

		err = gfile.PutContentsAppend("", "hello")
		t.AssertNE(err, nil)

	})

}

func Test_PutBinContents(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)

		err = gfile.PutBytes(testpath()+filepaths, []byte("test!!"))
		t.AssertNil(err)

		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "test!!")

		err = gfile.PutBytes("", []byte("test!!"))
		t.AssertNE(err, nil)

	})
}

func Test_PutBinContentsAppend(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "test!!")
		defer delTestFiles(filepaths)
		err = gfile.PutBytesAppend(testpath()+filepaths, []byte("word"))
		t.AssertNil(err)

		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "test!!word")

		err = gfile.PutBytesAppend("", []byte("word"))
		t.AssertNE(err, nil)

	})
}

func Test_GetBinContentsByTwoOffsetsByPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths   = "/testfile_GetContents.txt"
			readcontent []byte
		)

		createTestFile(filepaths, "abcdefghijk")
		defer delTestFiles(filepaths)
		readcontent = gfile.GetBytesByTwoOffsetsByPath(testpath()+filepaths, 2, 5)

		t.Assert(string(readcontent), "cde")

		readcontent = gfile.GetBytesByTwoOffsetsByPath("", 2, 5)
		t.Assert(len(readcontent), 0)

	})

}

func Test_GetNextCharOffsetByPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			filepaths  = "/testfile_GetContents.txt"
			localindex int64
		)
		createTestFile(filepaths, "abcdefghijk")
		defer delTestFiles(filepaths)
		localindex = gfile.GetNextCharOffsetByPath(testpath()+filepaths, 'd', 1)
		t.Assert(localindex, 3)

		localindex = gfile.GetNextCharOffsetByPath("", 'd', 1)
		t.Assert(localindex, -1)

	})
}

func Test_GetNextCharOffset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			localindex int64
		)
		reader := strings.NewReader("helloword")

		localindex = gfile.GetNextCharOffset(reader, 'w', 1)
		t.Assert(localindex, 5)

		localindex = gfile.GetNextCharOffset(reader, 'j', 1)
		t.Assert(localindex, -1)

	})
}

func Test_GetBinContentsByTwoOffsets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			reads []byte
		)
		reader := strings.NewReader("helloword")

		reads = gfile.GetBytesByTwoOffsets(reader, 1, 3)
		t.Assert(string(reads), "el")

		reads = gfile.GetBytesByTwoOffsets(reader, 10, 30)
		t.Assert(string(reads), "")

	})
}

func Test_GetBinContentsTilChar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			reads  []byte
			indexs int64
		)
		reader := strings.NewReader("helloword")

		reads, _ = gfile.GetBytesTilChar(reader, 'w', 2)
		t.Assert(string(reads), "llow")

		_, indexs = gfile.GetBytesTilChar(reader, 'w', 20)
		t.Assert(indexs, -1)

	})
}

func Test_GetBinContentsTilCharByPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			reads     []byte
			indexs    int64
			filepaths = "/testfile_GetContents.txt"
		)

		createTestFile(filepaths, "abcdefghijklmn")
		defer delTestFiles(filepaths)

		reads, _ = gfile.GetBytesTilCharByPath(testpath()+filepaths, 'c', 2)
		t.Assert(string(reads), "c")

		reads, _ = gfile.GetBytesTilCharByPath(testpath()+filepaths, 'y', 1)
		t.Assert(string(reads), "")

		_, indexs = gfile.GetBytesTilCharByPath(testpath()+filepaths, 'x', 1)
		t.Assert(indexs, -1)

	})
}

func Test_Home(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			reads string
			err   error
		)

		reads, err = gfile.Home("a", "b")
		t.AssertNil(err)
		t.AssertNE(reads, "")
	})
}

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
