package main

import (
	"github.com/gogf/gf/g/util/gutil"
)

func Test(s interface{}) {
	gutil.PrintStack()
}

func main() {
	Test(nil)
}
