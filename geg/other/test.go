package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	fmt.Println(binary.BigEndian.Uint32([]byte{byte(1), byte(1), byte(1), byte(1), byte(1)}))
}