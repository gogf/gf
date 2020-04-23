package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	b := []byte{3, 0, 0}
	fmt.Println(string(b))
	fmt.Println(hex.EncodeToString(b))
}
