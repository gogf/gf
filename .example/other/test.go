package main

import (
	"bytes"
	"fmt"
)

func main() {
	//b := make([]byte, 10)
	r := bytes.NewBuffer(make([]byte, 1))

	n, err := r.Write([]byte("12345"))

	fmt.Println(n)
	fmt.Println(err)
	fmt.Println(r.String())
}
