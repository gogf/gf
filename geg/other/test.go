package main

import (
	"fmt"
)

func main() {
	array := make([]interface{}, 0, 10)
	array[8] = 1
	fmt.Println(array)
}