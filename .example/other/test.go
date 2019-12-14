package main

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/frame/g"
	"os"
)

func main() {
	g.Dump(os.Args)
	fmt.Println(gdebug.BuildInfo())
}
