package main

import (
	"github.com/gogf/gf/v2/util/gutil"
)

func Test(s *interface{}) {
	//debug.PrintStack()
	gutil.PrintStack()
}

func main() {
	Test(nil)
}
