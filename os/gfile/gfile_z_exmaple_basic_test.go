// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleMkdir() {
	// init
	var (
		path = gfile.Temp("gfile_example_basic_dir")
	)

	// Creates directory
	gfile.Mkdir(path)

	// Check if directory exists
	fmt.Println(gfile.IsDir(path))

	// Output:
	// true
}

func ExampleCreate() {
	// init
	var (
		path     = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 50)
	)
	// Check whether the file exists
	isFile := gfile.IsFile(path)

	fmt.Println(isFile)

	// Creates file with given `path` recursively
	fileHandle, _ := gfile.Create(path)
	defer fileHandle.Close()

	// Write some content to file
	n, _ := fileHandle.WriteString("hello goframe")

	// Check whether the file exists
	isFile = gfile.IsFile(path)

	fmt.Println(isFile)

	// Reads len(b) bytes from the File
	fileHandle.ReadAt(dataByte, 0)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// false
	// true
	// hello goframe
}

func ExampleOpen() {
	// init
	var (
		path     = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 4096)
	)
	// Open file or directory with READONLY model
	file, _ := gfile.Open(path)
	defer file.Close()

	// Read data
	n, _ := file.Read(dataByte)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// hello goframe
}

func ExampleOpenFile() {
	// init
	var (
		path     = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 4096)
	)
	// Opens file/directory with custom `flag` and `perm`
	// Create if file does not exist,it is created in a readable and writable mode,prem 0777
	openFile, _ := gfile.OpenFile(path, os.O_CREATE|os.O_RDWR, gfile.DefaultPermCopy)
	defer openFile.Close()

	// Write some content to file
	writeLength, _ := openFile.WriteString("hello goframe test open file")

	fmt.Println(writeLength)

	// Read data
	n, _ := openFile.ReadAt(dataByte, 0)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// 28
	// hello goframe test open file
}

func ExampleOpenWithFlag() {
	// init
	var (
		path     = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 4096)
	)

	// Opens file/directory with custom `flag`
	// Create if file does not exist,it is created in a readable and writable mode with default `perm` is 0666
	openFile, _ := gfile.OpenWithFlag(path, os.O_CREATE|os.O_RDWR)
	defer openFile.Close()

	// Write some content to file
	writeLength, _ := openFile.WriteString("hello goframe test open file with flag")

	fmt.Println(writeLength)

	// Read data
	n, _ := openFile.ReadAt(dataByte, 0)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// 38
	// hello goframe test open file with flag
}

func ExampleJoin() {
	// init
	var (
		dirPath  = gfile.Temp("gfile_example_basic_dir")
		filePath = "file1"
	)

	// Joins string array paths with file separator of current system.
	joinString := gfile.Join(dirPath, filePath)

	fmt.Println(joinString)

	// May Output:
	// /tmp/gfile_example_basic_dir/file1
}

func ExampleExists() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)
	// Checks whether given `path` exist.
	joinString := gfile.Exists(path)

	fmt.Println(joinString)

	// Output:
	// true
}

func ExampleIsDir() {
	// init
	var (
		path     = gfile.Temp("gfile_example_basic_dir")
		filePath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)
	// Checks whether given `path` a directory.
	fmt.Println(gfile.IsDir(path))
	fmt.Println(gfile.IsDir(filePath))

	// Output:
	// true
	// false
}

func ExamplePwd() {
	// Get absolute path of current working directory.
	fmt.Println(gfile.Pwd())

	// May Output:
	// xxx/gf/os/gfile
}

func ExampleChdir() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)
	// Get current working directory
	fmt.Println(gfile.Pwd())

	// Changes the current working directory to the named directory.
	gfile.Chdir(path)

	// Get current working directory
	fmt.Println(gfile.Pwd())

	// May Output:
	// xxx/gf/os/gfile
	// /tmp/gfile_example_basic_dir/file1
}

func ExampleIsFile() {
	// init
	var (
		filePath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		dirPath  = gfile.Temp("gfile_example_basic_dir")
	)
	// Checks whether given `path` a file, which means it's not a directory.
	fmt.Println(gfile.IsFile(filePath))
	fmt.Println(gfile.IsFile(dirPath))

	// Output:
	// true
	// false
}

func ExampleStat() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)
	// Get a FileInfo describing the named file.
	stat, _ := gfile.Stat(path)

	fmt.Println(stat.Name())
	fmt.Println(stat.IsDir())
	fmt.Println(stat.Mode())
	fmt.Println(stat.ModTime())
	fmt.Println(stat.Size())
	fmt.Println(stat.Sys())

	// May Output:
	// file1
	// false
	// -rwxr-xr-x
	// 2021-12-02 11:01:27.261441694 +0800 CST
	// &{16777220 33261 1 8597857090 501 20 0 [0 0 0 0] {1638414088 192363490} {1638414087 261441694} {1638414087 261441694} {1638413480 485068275} 38 8 4096 0 0 0 [0 0]}
}

func ExampleMove() {
	// init
	var (
		srcPath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		dstPath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file2")
	)
	// Check is file
	fmt.Println(gfile.IsFile(dstPath))

	//  Moves `src` to `dst` path.
	// If `dst` already exists and is not a directory, it'll be replaced.
	gfile.Move(srcPath, dstPath)

	fmt.Println(gfile.IsFile(srcPath))
	fmt.Println(gfile.IsFile(dstPath))

	// Output:
	// false
	// false
	// true
}

func ExampleRename() {
	// init
	var (
		srcPath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file2")
		dstPath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)
	// Check is file
	fmt.Println(gfile.IsFile(dstPath))

	//  renames (moves) `src` to `dst` path.
	// If `dst` already exists and is not a directory, it'll be replaced.
	gfile.Rename(srcPath, dstPath)

	fmt.Println(gfile.IsFile(srcPath))
	fmt.Println(gfile.IsFile(dstPath))

	// Output:
	// false
	// false
	// true
}

func ExampleDirNames() {
	// init
	var (
		path = gfile.Temp("gfile_example_basic_dir")
	)
	// Get sub-file names of given directory `path`.
	dirNames, _ := gfile.DirNames(path)

	fmt.Println(dirNames)

	// May Output:
	// [file1]
}

func ExampleGlob() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "*_example_basic_test.go"
	)
	// Get sub-file names of given directory `path`.
	// Only show file name
	matchNames, _ := gfile.Glob(path, true)

	fmt.Println(matchNames)

	// Show full path of the file
	matchNames, _ = gfile.Glob(path, false)

	fmt.Println(matchNames)

	// May Output:
	// [gfile_z_example_basic_test.go]
	// [xxx/gf/os/gfile/gfile_z_example_basic_test.go]
}

func ExampleIsReadable() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "testdata/readline/file.log"
	)

	// Checks whether given `path` is readable.
	fmt.Println(gfile.IsReadable(path))

	// Output:
	// true
}

func ExampleIsWritable() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "testdata/readline/"
		file = "file.log"
	)

	// Checks whether given `path` is writable.
	fmt.Println(gfile.IsWritable(path))
	fmt.Println(gfile.IsWritable(path + file))

	// Output:
	// true
	// true
}

func ExampleChmod() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)

	// Get a FileInfo describing the named file.
	stat, err := gfile.Stat(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Show original mode
	fmt.Println(stat.Mode())

	// Change file model
	gfile.Chmod(path, gfile.DefaultPermCopy)

	// Get a FileInfo describing the named file.
	stat, _ = gfile.Stat(path)
	// Show the modified mode
	fmt.Println(stat.Mode())

	// Output:
	// -rw-r--r--
	// -rwxrwxrwx
}

func ExampleAbs() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)

	// Get an absolute representation of path.
	fmt.Println(gfile.Abs(path))

	// May Output:
	// /tmp/gfile_example_basic_dir/file1
}

func ExampleRealPath() {
	// init
	var (
		realPath  = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
		worryPath = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "worryFile")
	)

	// fetch an absolute representation of path.
	fmt.Println(gfile.RealPath(realPath))
	fmt.Println(gfile.RealPath(worryPath))

	// May Output:
	// /tmp/gfile_example_basic_dir/file1
}

func ExampleSelfPath() {

	// Get absolute file path of current running process
	fmt.Println(gfile.SelfPath())

	// May Output:
	// xxx/___github_com_gogf_gf_v2_os_gfile__ExampleSelfPath
}

func ExampleSelfName() {

	// Get file name of current running process
	fmt.Println(gfile.SelfName())

	// May Output:
	// ___github_com_gogf_gf_v2_os_gfile__ExampleSelfName
}

func ExampleSelfDir() {

	// Get absolute directory path of current running process
	fmt.Println(gfile.SelfDir())

	// May Output:
	// /private/var/folders/p6/gc_9mm3j229c0mjrjp01gqn80000gn/T
}

func ExampleBasename() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "testdata/readline/file.log"
	)

	// Get the last element of path, which contains file extension.
	fmt.Println(gfile.Basename(path))

	// Output:
	// file.log
}

func ExampleName() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "testdata/readline/file.log"
	)

	// Get the last element of path without file extension.
	fmt.Println(gfile.Name(path))

	// Output:
	// file
}

func ExampleDir() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)

	// Get all but the last element of path, typically the path's directory.
	fmt.Println(gfile.Dir(path))

	// May Output:
	// /tmp/gfile_example_basic_dir
}

func ExampleIsEmpty() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)

	// Check whether the `path` is empty
	fmt.Println(gfile.IsEmpty(path))

	// Truncate file
	gfile.Truncate(path, 0)

	// Check whether the `path` is empty
	fmt.Println(gfile.IsEmpty(path))

	// Output:
	// false
	// true
}

func ExampleExt() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "testdata/readline/file.log"
	)

	// Get the file name extension used by path.
	fmt.Println(gfile.Ext(path))

	// Output:
	// .log
}

func ExampleExtName() {
	// init
	var (
		path = gfile.Pwd() + gfile.Separator + "testdata/readline/file.log"
	)

	// Get the file name extension used by path but the result does not contains symbol '.'.
	fmt.Println(gfile.ExtName(path))

	// Output:
	// log
}

func ExampleTempDir() {
	// init
	var (
		fileName = "gfile_example_basic_dir"
	)

	// fetch an absolute representation of path.
	path := gfile.Temp(fileName)

	fmt.Println(path)

	// May Output:
	// /tmp/gfile_example_basic_dir
}

func ExampleRemove() {
	// init
	var (
		path = gfile.Join(gfile.Temp("gfile_example_basic_dir"), "file1")
	)

	// Checks whether given `path` a file, which means it's not a directory.
	fmt.Println(gfile.IsFile(path))

	// deletes all file/directory with `path` parameter.
	gfile.Remove(path)

	// Check again
	fmt.Println(gfile.IsFile(path))

	// Output:
	// true
	// false
}
