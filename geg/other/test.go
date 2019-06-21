package main

import (
	"fmt"
	"math"

	"github.com/gogf/gf/g/container/gtype"
)

func main() {
	v := gtype.NewInt32(math.MaxInt32)
	for i := 1; i < 100; i++ {
		fmt.Println(v.Add(int32(i)))
	}
}
