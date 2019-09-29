package main

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	s := "3.4028235e+38"
	fmt.Println(gconv.String(gconv.Float64(s)))
}
