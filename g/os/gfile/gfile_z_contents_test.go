package gfile_test

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/test/gtest"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 创建测试文件
func createTestFile(filename, content string) error {
	TempDir := testpath()
	err := ioutil.WriteFile(TempDir+filename, []byte(content), 0666)
	return err
}

// 测试完删除文件或目录
func delTestFiles(filenames string) {
	os.RemoveAll(testpath() + filenames)
}

// 创建目录
func createDir(paths string) {
	TempDir := testpath()
	os.Mkdir(TempDir+paths, 0777)
}

// 统一格式化文件目录为"/"
func formatpaths(paths []string) []string {
	for k, v := range paths {
		paths[k] = filepath.ToSlash(v)
		paths[k] = strings.Replace(paths[k], "./", "/", 1)
	}

	return paths
}

// 统一格式化文件目录为"/"
func formatpath(paths string) string {
	paths = filepath.ToSlash(paths)
	paths = strings.Replace(paths, "./", "/", 1)
	return paths
}

// 指定返回要测试的目录
func testpath() string {
	return os.TempDir()
}

func TestGetContents(t *testing.T) {
	gtest.Case(t, func() {

		var (
			filepaths string = "/testfile_t1.txt"
		)
		createTestFile(filepaths, "my name is jroam")
		defer delTestFiles(filepaths)

		gtest.Assert(gfile.GetContents(testpath()+filepaths), "my name is jroam")
		gtest.Assert(gfile.GetContents(""), "")

	})
}

func TestGetBinContents(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths1  string = "/testfile_t1.txt"                 // 文件存在时
			filepaths2  string = testpath() + "/testfile_t1_no.txt" // 文件不存在时
			readcontent []byte
			str1        string = "my name is jroam"
		)
		createTestFile(filepaths1, str1)
		defer delTestFiles(filepaths1)
		readcontent = gfile.GetBinContents(testpath() + filepaths1)
		gtest.Assert(readcontent, []byte(str1))

		readcontent = gfile.GetBinContents(filepaths2)
		gtest.Assert(string(readcontent), "")

		gtest.Assert(string(gfile.GetBinContents(filepaths2)), "")

	})
}

// 截断文件为指定的大小
func TestTruncate(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths1 string = "/testfile_GetContentsyyui.txt" //文件存在时
			err        error
			files      *os.File
		)
		createTestFile(filepaths1, "abcdefghijkmln")
		defer delTestFiles(filepaths1)
		err = gfile.Truncate(testpath()+filepaths1, 10)
		gtest.Assert(err, nil)

		//=========================检查修改文后的大小，是否与期望一致
		files, err = os.Open(testpath() + filepaths1)
		defer files.Close()
		gtest.Assert(err, nil)
		fileinfo, err2 := files.Stat()
		gtest.Assert(err2, nil)
		gtest.Assert(fileinfo.Size(), 10)

		//====测试当为空时，是否报错
		err = gfile.Truncate("", 10)
		gtest.AssertNE(err, nil)

	})
}

func TestPutContents(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)

		err = gfile.PutContents(testpath()+filepaths, "test!")
		gtest.Assert(err, nil)

		//==================判断是否真正写入
		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "test!")

		err = gfile.PutContents("", "test!")
		gtest.AssertNE(err, nil)

	})
}

func TestPutContentsAppend(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)

		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)
		err = gfile.PutContentsAppend(testpath()+filepaths, "hello")
		gtest.Assert(err, nil)

		//==================判断是否真正写入
		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "ahello")

		err = gfile.PutContentsAppend("", "hello")
		gtest.AssertNE(err, nil)

	})

}

func TestPutBinContents(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)

		err = gfile.PutBinContents(testpath()+filepaths, []byte("test!!"))
		gtest.Assert(err, nil)

		// 判断是否真正写入
		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "test!!")

		err = gfile.PutBinContents("", []byte("test!!"))
		gtest.AssertNE(err, nil)

	})
}

func TestPutBinContentsAppend(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "/testfile_PutContents.txt" //原文件内容: yy
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "test!!")
		defer delTestFiles(filepaths)
		err = gfile.PutBinContentsAppend(testpath()+filepaths, []byte("word"))
		gtest.Assert(err, nil)

		// 判断是否真正写入
		readcontent, err = ioutil.ReadFile(testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "test!!word")

		err = gfile.PutBinContentsAppend("", []byte("word"))
		gtest.AssertNE(err, nil)

	})
}

func TestGetBinContentsByTwoOffsetsByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "/testfile_GetContents.txt" // 文件内容: abcdefghijk
			readcontent []byte
		)

		createTestFile(filepaths, "abcdefghijk")
		defer delTestFiles(filepaths)
		readcontent = gfile.GetBinContentsByTwoOffsetsByPath(testpath()+filepaths, 2, 5)

		gtest.Assert(string(readcontent), "cde")

		readcontent = gfile.GetBinContentsByTwoOffsetsByPath("", 2, 5)
		gtest.Assert(len(readcontent), 0)

	})

}

func TestGetNextCharOffsetByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths  string = "/testfile_GetContents.txt" // 文件内容: abcdefghijk
			localindex int64
		)
		createTestFile(filepaths, "abcdefghijk")
		defer delTestFiles(filepaths)
		localindex = gfile.GetNextCharOffsetByPath(testpath()+filepaths, 'd', 1)
		gtest.Assert(localindex, 3)

		localindex = gfile.GetNextCharOffsetByPath("", 'd', 1)
		gtest.Assert(localindex, -1)

	})
}

func TestGetNextCharOffset(t *testing.T) {
	gtest.Case(t, func() {
		var (
			localindex int64
		)
		reader := strings.NewReader("helloword")

		localindex = gfile.GetNextCharOffset(reader, 'w', 1)
		gtest.Assert(localindex, 5)

		localindex = gfile.GetNextCharOffset(reader, 'j', 1)
		gtest.Assert(localindex, -1)

	})
}

func TestGetBinContentsByTwoOffsets(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads []byte
		)
		reader := strings.NewReader("helloword")

		reads = gfile.GetBinContentsByTwoOffsets(reader, 1, 3)
		gtest.Assert(string(reads), "el")

		reads = gfile.GetBinContentsByTwoOffsets(reader, 10, 30)
		gtest.Assert(string(reads), "")

	})
}

func TestGetBinContentsTilChar(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads  []byte
			indexs int64
		)
		reader := strings.NewReader("helloword")

		reads, _ = gfile.GetBinContentsTilChar(reader, 'w', 2)
		gtest.Assert(string(reads), "llow")

		_, indexs = gfile.GetBinContentsTilChar(reader, 'w', 20)
		gtest.Assert(indexs, -1)

	})
}

func TestGetBinContentsTilCharByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads     []byte
			indexs    int64
			filepaths string = "/testfile_GetContents.txt"
		)

		createTestFile(filepaths, "abcdefghijklmn")
		defer delTestFiles(filepaths)

		reads, _ = gfile.GetBinContentsTilCharByPath(testpath()+filepaths, 'c', 2)
		gtest.Assert(string(reads), "c")

		reads, _ = gfile.GetBinContentsTilCharByPath(testpath()+filepaths, 'y', 1)
		gtest.Assert(string(reads), "")

		_, indexs = gfile.GetBinContentsTilCharByPath(testpath()+filepaths, 'x', 1)
		gtest.Assert(indexs, -1)

	})
}

func TestHome(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads string
			err   error
		)

		reads, err = gfile.Home()
		gtest.Assert(err, nil)
		gtest.AssertNE(reads, "")

	})
}
