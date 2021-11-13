package gfile_test

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleGetContentsWithCache() {
	// init
	var (
		fileName = "gflie_example.txt"
		tempDir  = gfile.TempDir("gfile_example_cache")
		tempFile = gfile.Join(tempDir, fileName)
	)

	// write contents
	gfile.PutContents(tempFile, "goframe example content")

	// read contents
	fmt.Println(gfile.GetContentsWithCache(tempFile, time.Minute))

	// write new contents will clear its cache
	gfile.PutContents(tempFile, "new goframe example content")

	time.Sleep(time.Second * 1)

	// read contents
	fmt.Println(gfile.GetContentsWithCache(tempFile))

	// May Output:
	// goframe example content
	// new goframe example content
}
