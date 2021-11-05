package gfile_test

import (
	"fmt"
	"regexp"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleReplaceFile() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_replace")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	gfile.ReplaceFile("test", "replace word", tempFile)

	content1 := gfile.GetContents(tempFile)
	fmt.Println(content1)

	// Output:
	// test contents
	// replace word contents
}

func ExampleReplaceFileFunc() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_replace")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "666 test contents 888 a1a2a3")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	// replace by yourself
	gfile.ReplaceFileFunc(func(path, content string) string {
		// Replace with regular match
		reg, _ := regexp.Compile(`\d{3}`)
		return reg.ReplaceAllString(content, "[num]")
	}, tempFile)

	content1 := gfile.GetContents(tempFile)
	fmt.Println(content1)

	// Output:
	// 666 test contents 888 a1a2a3
	// [num] test contents [num] a1a2a3
}

func ExampleReplaceDir() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_replace")
	tempFile := gfile.Join(tempDir, fileName)

	tempSubDir := gfile.Join(tempDir, "sub_dir")
	tempSubFile := gfile.Join(tempSubDir, fileName)

	gfile.Mkdir(tempSubDir)
	gfile.Create(tempFile)
	gfile.Create(tempSubFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")
	gfile.PutContents(tempSubFile, "test contents")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)
	contentSub := gfile.GetContents(tempSubFile)
	fmt.Println(contentSub)

	gfile.ReplaceDir("test", "replace word", tempDir, "123.txt", true)

	// read contents
	content1 := gfile.GetContents(tempFile)
	fmt.Println(content1)
	contentSub1 := gfile.GetContents(tempSubFile)
	fmt.Println(contentSub1)

	// Output:
	// test contents
	// test contents
	// replace word contents
	// replace word contents
}

func ExampleReplaceDirFunc() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_replace")
	tempFile := gfile.Join(tempDir, fileName)

	tempSubDir := gfile.Join(tempDir, "sub_dir")
	tempSubFile := gfile.Join(tempSubDir, fileName)

	gfile.Mkdir(tempSubDir)
	gfile.Create(tempFile)
	gfile.Create(tempSubFile)

	// write contents
	gfile.PutContents(tempFile, "666 test contents 888 a1a2a3")
	gfile.PutContents(tempSubFile, "666 test contents 888 a1a2a3")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)
	contentSub := gfile.GetContents(tempSubFile)
	fmt.Println(contentSub)

	gfile.ReplaceDirFunc(func(path, content string) string {
		// Replace with regular match
		reg, _ := regexp.Compile(`\d{3}`)
		return reg.ReplaceAllString(content, "[num]")
	}, tempDir, "123.txt", true)

	// read contents
	content1 := gfile.GetContents(tempFile)
	fmt.Println(content1)
	contentSub1 := gfile.GetContents(tempSubFile)
	fmt.Println(contentSub1)

	// Output:
	// 666 test contents 888 a1a2a3
	// 666 test contents 888 a1a2a3
	// [num] test contents [num] a1a2a3
	// [num] test contents [num] a1a2a3
}
