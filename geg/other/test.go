package main

import "fmt"

func main() {
	x := uintptr(1)
	fmt.Println(x ^ 0)
	fmt.Println(2 ^ 0)
}
