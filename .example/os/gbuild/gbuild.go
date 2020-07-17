package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/gbuild"
)

func main() {
	g.Dump(gbuild.Info())
	g.Dump(gbuild.Map())
}
