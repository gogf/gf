package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleGetContents() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	// Output:
	// test contents
}

func ExampleGetBytes() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	content := gfile.GetBytes(tempFile)
	fmt.Println(string(content))

	// Output:
	// test contents
}

func ExamplePutContents() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	// Output:
	// test contents
}

func ExamplePutBytes() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutBytes(tempFile, []byte("test contents"))

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	// Output:
	// test contents
}

func ExamplePutContentsAppend() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	// write contents
	gfile.PutContentsAppend(tempFile, " append")

	// read contents
	content1 := gfile.GetContents(tempFile)
	fmt.Println(content1)

	// Output:
	// test contents
	// test contents append
}

func ExamplePutBytesAppend() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutBytes(tempFile, []byte("test contents"))

	// read contents
	content := gfile.GetContents(tempFile)
	fmt.Println(content)

	// write contents
	gfile.PutBytesAppend(tempFile, []byte(" append"))

	// read contents
	content1 := gfile.GetContents(tempFile)
	fmt.Println(content1)

	// Output:
	// test contents
	// test contents append
}

func ExampleGetNextCharOffsetByPath() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents index")

	// read contents
	index := gfile.GetNextCharOffsetByPath(tempFile, 'i', 0)
	fmt.Println(index)

	// Output:
	// 14
}

func ExampleGetBytesTilCharByPath() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents: hello")

	// read contents
	contents, index := gfile.GetBytesTilCharByPath(tempFile, ':', 0)
	fmt.Println(string(contents))
	fmt.Println(index)

	// Output:
	// test contents:
	// 13
}

func ExampleGetBytesByTwoOffsetsByPath() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	contents := gfile.GetBytesByTwoOffsetsByPath(tempFile, 0, 4)
	fmt.Println(string(contents))

	// Output:
	// test
}

func ExampleReadLines() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents\ntest contents")

	// read contents
	gfile.ReadLines(tempFile, func(text string) error {
		// Process each line
		fmt.Println(text)
		return nil
	})

	// Output:
	// test contents
	// test contents
}

func ExampleReadLinesBytes() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_content")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents\ntest contents")

	// read contents
	gfile.ReadLinesBytes(tempFile, func(bytes []byte) error {
		// Process each line
		fmt.Println(string(bytes))
		return nil
	})

	// Output:
	// test contents
	// test contents
}
