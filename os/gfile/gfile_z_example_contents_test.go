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

func ExampleGetContents() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// It reads and returns the file content as string.
	// It returns empty string if it fails reading, for example, with permission or IO error.
	fmt.Println(gfile.GetContents(tempFile))

	// Output:
	// goframe example content
}

func ExampleGetBytes() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// It reads and returns the file content as []byte.
	// It returns nil if it fails reading, for example, with permission or IO error.
	fmt.Println(gfile.GetBytes(tempFile))

	// Output:
	// [103 111 102 114 97 109 101 32 101 120 97 109 112 108 101 32 99 111 110 116 101 110 116]
}

func ExamplePutContents() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// It creates and puts content string into specifies file path.
	// It automatically creates directory recursively if it does not exist.
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	fmt.Println(gfile.GetContents(tempFile))

	// Output:
	// goframe example content
}

func ExamplePutBytes() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutBytes(tempFile, []byte("goframe example content"))

	// read contents
	fmt.Println(gfile.GetContents(tempFile))

	// Output:
	// goframe example content
}

func ExamplePutContentsAppend() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	fmt.Println(gfile.GetContents(tempFile))

	// write contents
	gfile.PutContentsAppend(tempFile, " append content")

	// read contents
	fmt.Println(gfile.GetContents(tempFile))

	// Output:
	// goframe example content
	// goframe example content append content
}

func ExamplePutBytesAppend() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	fmt.Println(gfile.GetContents(tempFile))

	// write contents
	gfile.PutBytesAppend(tempFile, []byte(" append"))

	// read contents
	fmt.Println(gfile.GetContents(tempFile))

	// Output:
	// goframe example content
	// goframe example content append
}

func ExampleGetNextCharOffsetByPath() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	index := gfile.GetNextCharOffsetByPath(tempFile, 'f', 0)
	fmt.Println(index)

	// Output:
	// 2
}

func ExampleGetBytesTilCharByPath() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	fmt.Println(gfile.GetBytesTilCharByPath(tempFile, 'f', 0))

	// Output:
	// [103 111 102] 2
}

func ExampleGetBytesByTwoOffsetsByPath() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	fmt.Println(gfile.GetBytesByTwoOffsetsByPath(tempFile, 0, 7))

	// Output:
	// [103 111 102 114 97 109 101]
}

func ExampleReadLines() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "L1 goframe example content\nL2 goframe example content")

	// read contents
	gfile.ReadLines(tempFile, func(text string) error {
		// Process each line
		fmt.Println(text)
		return nil
	})

	// Output:
	// L1 goframe example content
	// L2 goframe example content
}

func ExampleReadLinesBytes() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_content")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "L1 goframe example content\nL2 goframe example content")

	// read contents
	gfile.ReadLinesBytes(tempFile, func(bytes []byte) error {
		// Process each line
		fmt.Println(bytes)
		return nil
	})

	// Output:
	// [76 49 32 103 111 102 114 97 109 101 32 101 120 97 109 112 108 101 32 99 111 110 116 101 110 116]
	// [76 50 32 103 111 102 114 97 109 101 32 101 120 97 109 112 108 101 32 99 111 110 116 101 110 116]
}
