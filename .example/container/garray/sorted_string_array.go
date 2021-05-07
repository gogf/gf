package main

import (
	"fmt"

	"github.com/gogf/gf/container/garray"
)

func main() {
	array := garray.NewSortedStrArray()
	array.Add("9")
	array.Add("8")
	array.Add("7")
	array.Add("6")
	array.Add("5")
	array.Add("4")
	array.Add("3")
	array.Add("2")
	array.Add("1")
	fmt.Println(array.Slice())
	// output:
	// [1 2 3 4 5 6 7 8 9]
}
