package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gregex"
)

func main() {
	s := `-abc`
	m, err := gregex.MatchString(`^\-{1,2}a={0,1}(.*)`, s)
	g.Dump(err)
	g.Dump(m)
}
