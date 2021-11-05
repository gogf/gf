package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleSearch() {
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example")
	tempFile := gfile.Join(tempDir, fileName)
	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// search file
	realPath, _ := gfile.Search(fileName, tempDir)
	fmt.Println(gfile.Basename(realPath))

	// Output:
	// 123.txt
}
