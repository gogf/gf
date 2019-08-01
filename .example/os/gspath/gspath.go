package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gspath"
	"github.com/gogf/gf/os/gtime"
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
