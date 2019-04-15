package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//创建测试文件
func CreateTestFile(filename, content string) error {
	TempDir := Testpath()
	err := ioutil.WriteFile(TempDir+filename, []byte(content), 0666)
	return err
}

//测试完删除文件或目录
func DelTestFiles(filenames string) {
	os.RemoveAll(Testpath() + filenames)
}

//创建目录
func CreateDir(paths string) {
	TempDir := Testpath()
	os.Mkdir(TempDir+paths, 0777)
}

//统一格式化文件目录为"/"
func Formatpaths(paths []string) []string {
	for k, v := range paths {
		paths[k] = filepath.ToSlash(v)
		paths[k] = strings.Replace(paths[k], "./", "/", 1)
	}

	return paths
}

//统一格式化文件目录为"/"
func Formatpath(paths string) string {
	paths = filepath.ToSlash(paths)
	paths = strings.Replace(paths, "./", "/", 1)
	return paths
}

//指定返回要测试的目录
func Testpath() string {
	psths, err := filepath.Abs("./")
	if err != nil {
		return os.TempDir()
	}
	return strings.Replace(psths, "./", "/", 1)

}

func TestGetContents(t *testing.T) {
	gtest.Case(t, func() {

		var (
			filepaths string = "/testfile_t1.txt"
		)
		CreateTestFile(filepaths, "my name is jroam")
		defer DelTestFiles(filepaths)

		gtest.Assert(GetContents(Testpath()+filepaths), "my name is jroam")
		gtest.Assert(GetContents(""), "")

	})
}

func TestGetBinContents(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths1  string = "/testfile_t1.txt"                 //存在文件
			filepaths2  string = Testpath() + "/testfile_t1_no.txt" //不存大文件
			readcontent []byte
			str1        string = "my name is jroam"
		)
		CreateTestFile(filepaths1, str1)
		defer DelTestFiles(filepaths1)
		readcontent = GetBinContents(Testpath() + filepaths1)
		gtest.Assert(readcontent, []byte(str1))

		readcontent = GetBinContents(filepaths2)
		gtest.Assert(string(readcontent), "")

		//if readcontent!=nil{
		//	t.Error("文件应不存在")
		//}
		gtest.Assert(string(GetBinContents(filepaths2)), "")

	})
}

//截断文件为指定的大小
func TestTruncate(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths1 string = "/testfile_GetContents.txt" //存在文件
			err        error
		)
		CreateTestFile(filepaths1, "abcdefghijkmln")
		defer DelTestFiles(filepaths1)
		err = Truncate(Testpath()+filepaths1, 200)
		gtest.Assert(err, nil)

		err = Truncate("", 200)
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
		CreateTestFile(filepaths, "a")
		defer DelTestFiles(filepaths)

		err = PutContents(Testpath()+filepaths, "test!")
		gtest.Assert(err, nil)

		//==================判断是否真正写入
		readcontent, err = ioutil.ReadFile(Testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "test!")

		err = PutContents("", "test!")
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

		CreateTestFile(filepaths, "a")
		defer DelTestFiles(filepaths)
		err = PutContentsAppend(Testpath()+filepaths, "hello")
		gtest.Assert(err, nil)

		//==================判断是否真正写入
		readcontent, err = ioutil.ReadFile(Testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "ahello")

		err = PutContentsAppend("", "hello")
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
		CreateTestFile(filepaths, "a")
		defer DelTestFiles(filepaths)

		err = PutBinContents(Testpath()+filepaths, []byte("test!!"))
		gtest.Assert(err, nil)

		//==================判断是否真正写入
		readcontent, err = ioutil.ReadFile(Testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "test!!")

		err = PutBinContents("", []byte("test!!"))
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
		CreateTestFile(filepaths, "test!!")
		defer DelTestFiles(filepaths)
		err = PutBinContentsAppend(Testpath()+filepaths, []byte("word"))
		gtest.Assert(err, nil)

		//==================判断是否真正写入
		readcontent, err = ioutil.ReadFile(Testpath() + filepaths)
		gtest.Assert(err, nil)
		gtest.Assert(string(readcontent), "test!!word")

		err = PutBinContentsAppend("", []byte("word"))
		gtest.AssertNE(err, nil)

	})
}

func TestGetBinContentsByTwoOffsetsByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "/testfile_GetContents.txt" //文件内容: abcdefghijk
			readcontent []byte
		)

		CreateTestFile(filepaths, "abcdefghijk")
		defer DelTestFiles(filepaths)
		readcontent = GetBinContentsByTwoOffsetsByPath(Testpath()+filepaths, 2, 5)

		gtest.Assert(string(readcontent), "cde")

		readcontent = GetBinContentsByTwoOffsetsByPath("", 2, 5)
		gtest.Assert(len(readcontent), 0)

	})

}

func TestGetNextCharOffsetByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths  string = "/testfile_GetContents.txt" //文件内容: abcdefghijk
			localindex int64
		)
		CreateTestFile(filepaths, "abcdefghijk")
		defer DelTestFiles(filepaths)
		localindex = GetNextCharOffsetByPath(Testpath()+filepaths, 'd', 1)
		gtest.Assert(localindex, 3)

		localindex = GetNextCharOffsetByPath("", 'd', 1)
		gtest.Assert(localindex, -1)

	})
}

func TestGetNextCharOffset(t *testing.T) {
	gtest.Case(t, func() {
		var (
			localindex int64
		)
		reader := strings.NewReader("helloword")

		localindex = GetNextCharOffset(reader, 'w', 1)
		gtest.Assert(localindex, 5)

		localindex = GetNextCharOffset(reader, 'j', 1)
		gtest.Assert(localindex, -1)

	})
}

func TestGetBinContentsByTwoOffsets(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads []byte
		)
		reader := strings.NewReader("helloword")

		reads = GetBinContentsByTwoOffsets(reader, 1, 3)
		gtest.Assert(string(reads), "el")

		reads = GetBinContentsByTwoOffsets(reader, 10, 30)
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

		reads, _ = GetBinContentsTilChar(reader, 'w', 2)
		gtest.Assert(string(reads), "llow")

		_, indexs = GetBinContentsTilChar(reader, 'w', 20)
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

		CreateTestFile(filepaths, "abcdefghijklmn")
		defer DelTestFiles(filepaths)

		reads, _ = GetBinContentsTilCharByPath(Testpath()+filepaths, 'c', 2)
		gtest.Assert(string(reads), "c")

		reads, _ = GetBinContentsTilCharByPath(Testpath()+filepaths, 'y', 1)
		gtest.Assert(string(reads), "")

		_, indexs = GetBinContentsTilCharByPath(Testpath()+filepaths, 'x', 1)
		gtest.Assert(indexs, -1)

	})
}

func TestHome(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads string
			err   error
		)

		reads, err = Home()
		gtest.Assert(err, nil)
		gtest.AssertNE(reads, "")

	})
}
