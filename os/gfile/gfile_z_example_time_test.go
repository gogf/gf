package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleMTime() {
	t := gfile.MTime(gfile.TempDir())
	fmt.Println(t)

	// May Output:
	// 2021-11-02 15:18:43.901141 +0800 CST
}

func ExampleMTimestamp() {
	t := gfile.MTimestamp(gfile.TempDir())
	fmt.Println(t)

	// May Output:
	// 1635838398
}

func ExampleMTimestampMilli() {
	t := gfile.MTimestampMilli(gfile.TempDir())
	fmt.Println(t)

	// May Output:
	// 1635838529330
}
