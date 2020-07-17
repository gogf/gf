package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	formats := []string{
		"Y-m-d H:i:s.u",
		"D M d H:i:s T O Y",
		"\\T\\i\\m\\e \\i\\s: h:i:s a",
		"2006-01-02T15:04:05.000000000Z07:00",
	}
	t := gtime.Now()
	for _, f := range formats {
		fmt.Println(t.Format(f))
	}
}
