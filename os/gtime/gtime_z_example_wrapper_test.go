package gtime_test

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

func ExampleSetStringLayout() {
	gtime.SetStringLayout(time.RFC3339)
	fmt.Println(gtime.New("2025-01-21 14:08:05").String())

	// Output:
	// 2025-01-21T14:08:05+08:00
}

func ExampleSetStringISO8601() {
	gtime.SetStringISO8601()
	fmt.Println(gtime.New("2025-01-21 14:10:12").String())

	// Output:
	// 2025-01-21T14:10:12+08:00
}

func ExampleSetStringRFC822() {
	gtime.SetStringRFC822()
	fmt.Println(gtime.New("2025-01-21 14:11:32").String())

	// Output:
	// Tue, 21 Jan 25 14:11 CST
}
