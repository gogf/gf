package main

import (
	"fmt"
)

func main() {
	b1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	b2 := b1[:2]
	b1[0] = 9
	fmt.Println(b2)
}
