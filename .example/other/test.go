package main

import (
	"fmt"
	"github.com/gogf/gf/internal/utilbytes"
)

func main() {
	b := []byte{48, 49, 50, 51, 52, 53}
	fmt.Println(string(b))
	fmt.Println([]byte("\xff\xff"))
	fmt.Printf(utilbytes.Export(b))
}
