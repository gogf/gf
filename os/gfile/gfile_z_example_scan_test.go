// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleScanDir() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_scan_dir")
		tempFile = gfile.Join(tempDir, fileName)

		tempSubDir  = gfile.Join(tempDir, "sub_dir")
		tempSubFile = gfile.Join(tempSubDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")
	gfile.PutContents(tempSubFile, "goframe example content")

	// scans directory recursively
	list, _ := gfile.ScanDir(tempDir, "*", true)
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// gflie_example.txt
	// sub_dir
	// gflie_example.txt
}

func ExampleScanDirFile() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_scan_dir_file")
		tempFile = gfile.Join(tempDir, fileName)

		tempSubDir  = gfile.Join(tempDir, "sub_dir")
		tempSubFile = gfile.Join(tempSubDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")
	gfile.PutContents(tempSubFile, "goframe example content")

	// scans directory recursively exclusive of directories
	list, _ := gfile.ScanDirFile(tempDir, "*.txt", true)
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// gflie_example.txt
	// gflie_example.txt
}

func ExampleScanDirFunc() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_scan_dir_func")
		tempFile = gfile.Join(tempDir, fileName)

		tempSubDir  = gfile.Join(tempDir, "sub_dir")
		tempSubFile = gfile.Join(tempSubDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")
	gfile.PutContents(tempSubFile, "goframe example content")

	// scans directory recursively
	list, _ := gfile.ScanDirFunc(tempDir, "*", true, func(path string) string {
		// ignores some files
		if gfile.Basename(path) == "gflie_example.txt" {
			return ""
		}
		return path
	})
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// sub_dir
}

func ExampleScanDirFileFunc() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_scan_dir_file_func")
		tempFile = gfile.Join(tempDir, fileName)

		fileName1 = "gflie_example_ignores.txt"
		tempFile1 = gfile.Join(tempDir, fileName1)

		tempSubDir  = gfile.Join(tempDir, "sub_dir")
		tempSubFile = gfile.Join(tempSubDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")
	gfile.PutContents(tempFile1, "goframe example content")
	gfile.PutContents(tempSubFile, "goframe example content")

	// scans directory recursively exclusive of directories
	list, _ := gfile.ScanDirFileFunc(tempDir, "*.txt", true, func(path string) string {
		// ignores some files
		if gfile.Basename(path) == "gflie_example_ignores.txt" {
			return ""
		}
		return path
	})
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// gflie_example.txt
	// gflie_example.txt
}
