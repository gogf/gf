package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleCopy() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_copy")
	tempFile := gfile.Join(tempDir, fileName)

	dstFileName := "123copy.txt"
	dstTempDir := gfile.TempDir("gfile_example_copy_dst")
	dstTempFile := gfile.Join(tempDir, dstFileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// clear
	if gfile.Exists(dstTempFile) {
		gfile.Remove(dstTempFile)
	}
	if gfile.Exists(dstTempDir) {
		gfile.Remove(dstTempDir)
	}

	// write contents
	gfile.PutContents(tempFile, "test copy")

	// copy file
	gfile.Copy(tempFile, dstTempFile)

	// read contents
	content := gfile.GetContents(dstTempFile)
	fmt.Println(content)

	// copy dir
	gfile.Copy(tempDir, dstTempDir)

	fList, _ := gfile.ScanDir(dstTempDir, "*", false)
	for _, v := range fList {
		content := gfile.GetContents(v)
		fmt.Println(content)
	}

	// Output:
	// test copy
	// test copy
	// test copy
}
