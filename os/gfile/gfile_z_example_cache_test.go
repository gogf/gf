package gfile_test

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleGetContentsWithCache() {
	intlog.SetEnabled(false)
	// init
	fileName := "123.txt"
	tempDir := gfile.TempDir("gfile_example_cache")
	tempFile := gfile.Join(tempDir, fileName)

	gfile.Mkdir(tempDir)
	gfile.Create(tempFile)

	// write contents
	gfile.PutContents(tempFile, "test contents")

	// read contents
	content := gfile.GetContentsWithCache(tempFile, time.Minute)
	fmt.Println(content)

	time.Sleep(time.Second * 1)

	// read contents
	content1 := gfile.GetContentsWithCache(tempFile)
	fmt.Println(content1)

	// write new contents will clear its cache
	gfile.PutContents(tempFile, "new test contents")

	time.Sleep(time.Second * 1)

	// read contents
	content2 := gfile.GetContentsWithCache(tempFile)
	fmt.Println(content2)

	// Output:
	// test contents
	// test contents
	// new test contents
}
