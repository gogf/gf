package main

import (
	"fmt"
	"github.com/gogf/gf/util/guid"
)

func main() {
	for i := 0; i < 10; i++ {
		s := guid.New([]byte("123"))
		fmt.Println(s, len(s))
	}
	fmt.Println()
	for i := 0; i < 10; i++ {
		s := guid.New([]byte("123"), []byte("456"))
		fmt.Println(s, len(s))
	}
	fmt.Println()
	for i := 0; i < 10; i++ {
		s := guid.New([]byte("123"), []byte("456"), []byte("789"))
		fmt.Println(s, len(s))
	}
}
