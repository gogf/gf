package main

import (
	"fmt"
)

func main() {
	var i float64 = 0
	for index := 0; index < 10; index++ {
		i += 0.1
		fmt.Println(i)
	}
}