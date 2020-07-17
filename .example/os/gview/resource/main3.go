package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/gres"
	_ "github.com/jin502437344/gf/os/gres/testdata"
)

func main() {
	gres.Dump()

	v := g.View()
	s, err := v.Parse("index.html")
	fmt.Println(err)
	fmt.Println(s)
}
