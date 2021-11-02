package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gring"
)

func main() {
	r1 := gring.New(10)
	for i := 0; i < 5; i++ {
		r1.Set(i).Next()
	}
	fmt.Println("Len:", r1.Len())
	fmt.Println("Cap:", r1.Cap())
	fmt.Println(r1.SlicePrev())
	fmt.Println(r1.SliceNext())

	r2 := gring.New(10)
	for i := 0; i < 10; i++ {
		r2.Set(i).Next()
	}
	fmt.Println("Len:", r2.Len())
	fmt.Println("Cap:", r2.Cap())
	fmt.Println(r2.SlicePrev())
	fmt.Println(r2.SliceNext())

}
