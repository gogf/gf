package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleScanDir() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example")
	tempFile := gfile.Join(tempDir, fileName)

	tempSubDir := gfile.Join(tempDir, "sub_dir")
	tempSubFile := gfile.Join(tempSubDir, fileName)

	gfile.Mkdir(tempSubDir)
	gfile.Create(tempFile)
	gfile.Create(tempSubFile)

	// scans directory recursively
	list, _ := gfile.ScanDir(tempDir, "123.txt,sub_dir", true)
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// 123.txt
	// sub_dir
	// 123.txt
}

func ExampleScanDirFile() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example")
	tempFile := gfile.Join(tempDir, fileName)

	tempSubDir := gfile.Join(tempDir, "sub_dir")
	tempSubFile := gfile.Join(tempSubDir, fileName)

	gfile.Mkdir(tempSubDir)
	gfile.Create(tempFile)
	gfile.Create(tempSubFile)

	// scans directory recursively exclusive of directories
	list, _ := gfile.ScanDirFile(tempDir, "123.txt,sub_dir", true)
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// 123.txt
	// 123.txt
}

func ExampleScanDirFunc() {
	// init
	fileName := "123.txt"
	fileName1 := "1234.txt"
	tempDir := gfile.TempDir("gfile_example_1")
	tempFile := gfile.Join(tempDir, fileName)

	tempSubDir := gfile.Join(tempDir, "sub_dir")
	tempSubFile := gfile.Join(tempSubDir, fileName1)

	gfile.Mkdir(tempSubDir)
	gfile.Create(tempFile)
	gfile.Create(tempSubFile)

	// scans directory recursively
	list, _ := gfile.ScanDirFunc(tempDir, "123.txt,1234.txt,sub_dir", true, func(path string) string {
		// ignores some files
		if gfile.Basename(path) == "1234.txt" {
			return ""
		}
		return path
	})
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// 123.txt
	// sub_dir
}

func ExampleScanDirFileFunc() {
	// init
	fileName := "123.txt"
	fileName1 := "1234.txt"
	tempDir := gfile.TempDir("gfile_example_1")
	tempFile := gfile.Join(tempDir, fileName)

	tempSubDir := gfile.Join(tempDir, "sub_dir")
	tempSubFile := gfile.Join(tempSubDir, fileName1)

	gfile.Mkdir(tempSubDir)
	gfile.Create(tempFile)
	gfile.Create(tempSubFile)

	// scans directory recursively exclusive of directories
	list, _ := gfile.ScanDirFileFunc(tempDir, "123.txt,1234.txt,sub_dir", true, func(path string) string {
		// ignores some files
		if gfile.Basename(path) == "1234.txt" {
			return ""
		}
		return path
	})
	for _, v := range list {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// 123.txt
}
