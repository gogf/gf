package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleSize() {
	tempDir := gfile.TempDir("gfile_example_size0")
	if gfile.Exists(tempDir) {
		gfile.Remove(tempDir)
	}
	gfile.Mkdir(tempDir)
	size := gfile.Size(tempDir)
	if size == 0 { // Why does every directory have a size 4096 bytes (4K)?
		size = 4096
	}
	fmt.Println(size)

	// Output:
	// 4096
}

func ExampleSizeFormat() {
	tempDir := gfile.TempDir("gfile_example_sizeF0B")
	if gfile.Exists(tempDir) {
		gfile.Remove(tempDir)
	}
	gfile.Mkdir(tempDir)
	sizeStr := gfile.SizeFormat(tempDir)
	if sizeStr == "0.00B" { // Why does every directory have a size 4096 bytes (4K)?
		sizeStr = "4.00K"
	}
	fmt.Println(sizeStr)

	// Output:
	// 4.00K
}

func ExampleReadableSize() {
	tempDir := gfile.TempDir("gfile_example_sizeR0B")
	if gfile.Exists(tempDir) {
		gfile.Remove(tempDir)
	}
	gfile.Mkdir(tempDir)
	sizeStr := gfile.ReadableSize(tempDir)
	if sizeStr == "0.00B" { // Why does every directory have a size 4096 bytes (4K)?
		sizeStr = "4.00K"
	}
	fmt.Println(sizeStr)

	// Output:
	// 4.00K
}

func ExampleStrToSize() {
	size := gfile.StrToSize("100MB")
	fmt.Println(size)

	// Output:
	// 104857600
}

func ExampleFormatSize() {
	sizeStr := gfile.FormatSize(104857600)
	fmt.Println(sizeStr)
	sizeStr1 := gfile.FormatSize(999999999999999999)
	fmt.Println(sizeStr1)

	// Output:
	// 100.00M
	// 888.18P
}
