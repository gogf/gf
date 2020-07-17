package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/gspath"
	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	sp := gspath.New()
	path := "/Users/john/Temp"
	rp, err := sp.Add(path)
	fmt.Println(err)
	fmt.Println(rp)
	fmt.Println(sp)

	gtime.SetInterval(5*time.Second, func() bool {
		g.Dump(sp.AllPaths())
		return true
	})

	select {}
}
