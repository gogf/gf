package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleSortFiles() {
	files := []string{
		"/aaa/bbb/ccc.txt",
		"/aaa/bbb/",
		"/aaa/",
		"/aaa",
		"/aaa/ccc/ddd.txt",
		"/bbb",
		"/0123",
		"/ddd",
		"/ccc",
	}
	sortOut := gfile.SortFiles(files)
	fmt.Println(sortOut)

	// Output:
	// [/0123 /aaa /aaa/ /aaa/bbb/ /aaa/bbb/ccc.txt /aaa/ccc/ddd.txt /bbb /ccc /ddd]
}
