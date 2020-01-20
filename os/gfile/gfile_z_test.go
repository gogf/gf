// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_IsDir(t *testing.T) {
	gtest.Case(t, func() {
		paths := "/testfile"
		createDir(paths)
		defer delTestFiles(paths)

		gtest.Assert(gfile.IsDir(testpath()+paths), true)
		gtest.Assert(gfile.IsDir("./testfile2"), false)
		gtest.Assert(gfile.IsDir("./testfile/tt.txt"), false)
		gtest.Assert(gfile.IsDir(""), false)
	})
}

func Test_IsEmpty(t *testing.T) {
	gtest.Case(t, func() {
		path := "/testdir_" + gconv.String(gtime.TimestampNano())
		createDir(path)
		defer delTestFiles(path)

		gtest.Assert(gfile.IsEmpty(testpath()+path), true)
		gtest.Assert(gfile.IsEmpty(testpath()+path+gfile.Separator+"test.txt"), true)
	})
	gtest.Case(t, func() {
		path := "/testfile_" + gconv.String(gtime.TimestampNano())
		createTestFile(path, "")
		defer delTestFiles(path)

		gtest.Assert(gfile.IsEmpty(testpath()+path), true)
	})
	gtest.Case(t, func() {
		path := "/testfile_" + gconv.String(gtime.TimestampNano())
		createTestFile(path, "1")
		defer delTestFiles(path)

		gtest.Assert(gfile.IsEmpty(testpath()+path), false)
	})
}

func Test_Create(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err       error
			filepaths []string
			fileobj   *os.File
		)
		filepaths = append(filepaths, "/testfile_cc1.txt")
		filepaths = append(filepaths, "/testfile_cc2.txt")
		for _, v := range filepaths {
			fileobj, err = gfile.Create(testpath() + v)
			defer delTestFiles(v)
			fileobj.Close()
			gtest.Assert(err, nil)
		}
	})
}

func Test_Open(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)

		file1 := "/testfile_nc1.txt"
		createTestFile(file1, "")
		defer delTestFiles(file1)

		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "./testfile/file1/c1.txt")
		flags = append(flags, false)

		for k, v := range files {
			fileobj, err = gfile.Open(testpath() + v)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func Test_OpenFile(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)

		files = append(files, "./testfile/file1/nc1.txt")
		flags = append(flags, false)

		f1 := "/testfile_tt.txt"
		createTestFile(f1, "")
		defer delTestFiles(f1)

		files = append(files, f1)
		flags = append(flags, true)

		for k, v := range files {
			fileobj, err = gfile.OpenFile(testpath()+v, os.O_RDWR, 0666)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func Test_OpenWithFlag(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)

		file1 := "/testfile_t1.txt"
		createTestFile(file1, "")
		defer delTestFiles(file1)
		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "/testfiless/dirfiles/t1_no.txt")
		flags = append(flags, false)

		for k, v := range files {
			fileobj, err = gfile.OpenWithFlag(testpath()+v, os.O_RDWR)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func Test_OpenWithFlagPerm(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)
		file1 := "/testfile_nc1.txt"
		createTestFile(file1, "")
		defer delTestFiles(file1)
		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "/testfileyy/tt.txt")
		flags = append(flags, false)

		for k, v := range files {
			fileobj, err = gfile.OpenWithFlagPerm(testpath()+v, os.O_RDWR, 666)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func Test_Exists(t *testing.T) {

	gtest.Case(t, func() {
		var (
			flag  bool
			files []string
			flags []bool
		)

		file1 := "/testfile_GetContents.txt"
		createTestFile(file1, "")
		defer delTestFiles(file1)

		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "./testfile/havefile1/tt_no.txt")
		flags = append(flags, false)

		for k, v := range files {
			flag = gfile.Exists(testpath() + v)
			if flags[k] {
				gtest.Assert(flag, true)
			} else {
				gtest.Assert(flag, false)
			}

		}

	})
}

func Test_Pwd(t *testing.T) {
	gtest.Case(t, func() {
		paths, err := os.Getwd()
		gtest.Assert(err, nil)
		gtest.Assert(gfile.Pwd(), paths)

	})
}

func Test_IsFile(t *testing.T) {
	gtest.Case(t, func() {
		var (
			flag  bool
			files []string
			flags []bool
		)

		file1 := "/testfile_tt.txt"
		createTestFile(file1, "")
		defer delTestFiles(file1)
		files = append(files, file1)
		flags = append(flags, true)

		dir1 := "/testfiless"
		createDir(dir1)
		defer delTestFiles(dir1)
		files = append(files, dir1)
		flags = append(flags, false)

		files = append(files, "./testfiledd/tt1.txt")
		flags = append(flags, false)

		for k, v := range files {
			flag = gfile.IsFile(testpath() + v)
			if flags[k] {
				gtest.Assert(flag, true)
			} else {
				gtest.Assert(flag, false)
			}

		}

	})
}

func Test_Info(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err    error
			paths  string = "/testfile_t1.txt"
			files  os.FileInfo
			files2 os.FileInfo
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)
		files, err = gfile.Info(testpath() + paths)
		gtest.Assert(err, nil)

		files2, err = os.Stat(testpath() + paths)
		gtest.Assert(err, nil)

		gtest.Assert(files, files2)

	})
}

func Test_Move(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths     string = "/ovetest"
			filepaths string = "/testfile_ttn1.txt"
			topath    string = "/testfile_ttn2.txt"
		)
		createDir("/ovetest")
		createTestFile(paths+filepaths, "a")

		defer delTestFiles(paths)

		yfile := testpath() + paths + filepaths
		tofile := testpath() + paths + topath

		gtest.Assert(gfile.Move(yfile, tofile), nil)

		// 检查移动后的文件是否真实存在
		_, err := os.Stat(tofile)
		gtest.Assert(os.IsNotExist(err), false)

	})
}

func Test_Rename(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths  string = "/testfiles"
			ypath  string = "/testfilettm1.txt"
			topath string = "/testfilettm2.txt"
		)
		createDir(paths)
		createTestFile(paths+ypath, "a")
		defer delTestFiles(paths)

		ypath = testpath() + paths + ypath
		topath = testpath() + paths + topath

		gtest.Assert(gfile.Rename(ypath, topath), nil)
		gtest.Assert(gfile.IsFile(topath), true)

		gtest.AssertNE(gfile.Rename("", ""), nil)

	})

}

func Test_DirNames(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths    string = "/testdirs"
			err      error
			readlist []string
		)
		havelist := []string{
			"t1.txt",
			"t2.txt",
		}

		// 创建测试文件
		createDir(paths)
		for _, v := range havelist {
			createTestFile(paths+"/"+v, "")
		}
		defer delTestFiles(paths)

		readlist, err = gfile.DirNames(testpath() + paths)

		gtest.Assert(err, nil)
		gtest.AssertIN(readlist, havelist)

		_, err = gfile.DirNames("")
		gtest.AssertNE(err, nil)

	})
}

func Test_Glob(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths      string = "/testfiles/*.txt"
			dirpath    string = "/testfiles"
			err        error
			resultlist []string
		)

		havelist1 := []string{
			"t1.txt",
			"t2.txt",
		}

		havelist2 := []string{
			testpath() + "/testfiles/t1.txt",
			testpath() + "/testfiles/t2.txt",
		}

		//===============================构建测试文件
		createDir(dirpath)
		for _, v := range havelist1 {
			createTestFile(dirpath+"/"+v, "")
		}
		defer delTestFiles(dirpath)

		resultlist, err = gfile.Glob(testpath()+paths, true)
		gtest.Assert(err, nil)
		gtest.Assert(resultlist, havelist1)

		resultlist, err = gfile.Glob(testpath()+paths, false)

		gtest.Assert(err, nil)
		gtest.Assert(formatpaths(resultlist), formatpaths(havelist2))

		_, err = gfile.Glob("", true)
		gtest.Assert(err, nil)

	})
}

func Test_Remove(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths string = "/testfile_t1.txt"
		)
		createTestFile(paths, "")
		gtest.Assert(gfile.Remove(testpath()+paths), nil)

		gtest.Assert(gfile.Remove(""), nil)

		defer delTestFiles(paths)

	})
}

func Test_IsReadable(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_GetContents.txt"
			paths2 string = "./testfile_GetContents_no.txt"
		)

		createTestFile(paths1, "")
		defer delTestFiles(paths1)

		gtest.Assert(gfile.IsReadable(testpath()+paths1), true)
		gtest.Assert(gfile.IsReadable(paths2), false)

	})
}

func Test_IsWritable(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_GetContents.txt"
			paths2 string = "./testfile_GetContents_no.txt"
		)

		createTestFile(paths1, "")
		defer delTestFiles(paths1)
		gtest.Assert(gfile.IsWritable(testpath()+paths1), true)
		gtest.Assert(gfile.IsWritable(paths2), false)

	})
}

func Test_Chmod(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_GetContents.txt"
			paths2 string = "./testfile_GetContents_no.txt"
		)
		createTestFile(paths1, "")
		defer delTestFiles(paths1)

		gtest.Assert(gfile.Chmod(testpath()+paths1, 0777), nil)
		gtest.AssertNE(gfile.Chmod(paths2, 0777), nil)

	})
}

// 获取绝对目录地址
func Test_RealPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string = "/testfile_files"
			readlPath string

			tempstr string
		)

		createDir(paths1)
		defer delTestFiles(paths1)

		readlPath = gfile.RealPath("./")

		tempstr, _ = filepath.Abs("./")

		gtest.Assert(readlPath, tempstr)

		gtest.Assert(gfile.RealPath("./nodirs"), "")

	})
}

// 获取当前执行文件的目录
func Test_SelfPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string
			readlPath string
			tempstr   string
		)
		readlPath = gfile.SelfPath()
		readlPath = filepath.ToSlash(readlPath)

		tempstr, _ = filepath.Abs(os.Args[0])
		paths1 = filepath.ToSlash(tempstr)
		paths1 = strings.Replace(paths1, "./", "/", 1)

		gtest.Assert(readlPath, paths1)

	})
}

func Test_SelfDir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string
			readlPath string
			tempstr   string
		)
		readlPath = gfile.SelfDir()

		tempstr, _ = filepath.Abs(os.Args[0])
		paths1 = filepath.Dir(tempstr)

		gtest.Assert(readlPath, paths1)

	})
}

func Test_Basename(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string = "/testfilerr_GetContents.txt"
			readlPath string
		)

		createTestFile(paths1, "")
		defer delTestFiles(paths1)

		readlPath = gfile.Basename(testpath() + paths1)
		gtest.Assert(readlPath, "testfilerr_GetContents.txt")

	})
}

func Test_Dir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string = "/testfiless"
			readlPath string
		)
		createDir(paths1)
		defer delTestFiles(paths1)

		readlPath = gfile.Dir(testpath() + paths1)

		gtest.Assert(readlPath, testpath())

	})
}

func Test_Ext(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1   string = "/testfile_GetContents.txt"
			dirpath1        = "/testdirs"
		)
		createTestFile(paths1, "")
		defer delTestFiles(paths1)

		createDir(dirpath1)
		defer delTestFiles(dirpath1)

		gtest.Assert(gfile.Ext(testpath()+paths1), ".txt")
		gtest.Assert(gfile.Ext(testpath()+dirpath1), "")
	})

	gtest.Case(t, func() {
		gtest.Assert(gfile.Ext("/var/www/test.js"), ".js")
		gtest.Assert(gfile.Ext("/var/www/test.min.js"), ".js")
		gtest.Assert(gfile.Ext("/var/www/test.js?1"), ".js")
		gtest.Assert(gfile.Ext("/var/www/test.min.js?v1"), ".js")
	})
}

func Test_ExtName(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gfile.ExtName("/var/www/test.js"), "js")
		gtest.Assert(gfile.ExtName("/var/www/test.min.js"), "js")
		gtest.Assert(gfile.ExtName("/var/www/test.js?v=1"), "js")
		gtest.Assert(gfile.ExtName("/var/www/test.min.js?v=1"), "js")
	})
}

func Test_TempDir(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gfile.TempDir(), "/tmp")
	})
}

func Test_Mkdir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			tpath string = "/testfile/createdir"
			err   error
		)

		defer delTestFiles("/testfile")

		err = gfile.Mkdir(testpath() + tpath)
		gtest.Assert(err, nil)

		err = gfile.Mkdir("")
		gtest.AssertNE(err, nil)

		err = gfile.Mkdir(testpath() + tpath + "2/t1")
		gtest.Assert(err, nil)

	})
}

func Test_Stat(t *testing.T) {
	gtest.Case(t, func() {
		var (
			tpath1   = "/testfile_t1.txt"
			tpath2   = "./testfile_t1_no.txt"
			err      error
			fileiofo os.FileInfo
		)

		createTestFile(tpath1, "a")
		defer delTestFiles(tpath1)

		fileiofo, err = gfile.Stat(testpath() + tpath1)
		gtest.Assert(err, nil)

		gtest.Assert(fileiofo.Size(), 1)

		_, err = gfile.Stat(tpath2)
		gtest.AssertNE(err, nil)

	})
}

func Test_MainPkgPath(t *testing.T) {
	gtest.Case(t, func() {
		reads := gfile.MainPkgPath()
		gtest.Assert(reads, "")
	})
}
