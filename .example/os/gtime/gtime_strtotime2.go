package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	if t, err := gtime.StrToTimeFormat("Tue Oct 16 15:55:59 CST 2018", "D M d H:i:s T Y"); err == nil {
		fmt.Println(t.String())
	} else {
		panic(err)
	}
}
