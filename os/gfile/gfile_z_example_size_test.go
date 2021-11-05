package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleSize() {
	tempDir := gfile.TempDir("gfile_example")
	gfile.Mkdir(tempDir)
	size := gfile.Size(tempDir)
	fmt.Println(size)

	// Output:
	// 0
}

func ExampleSizeFormat() {
	tempDir := gfile.TempDir("gfile_example")
	gfile.Mkdir(tempDir)
	sizeStr := gfile.SizeFormat(tempDir)
	fmt.Println(sizeStr)

	// Output:
	// 0.00B
}

func ExampleReadableSize() {
	tempDir := gfile.TempDir("gfile_example")
	gfile.Mkdir(tempDir)
	sizeStr := gfile.ReadableSize(tempDir)
	fmt.Println(sizeStr)

	// Output:
	// 0.00B
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
