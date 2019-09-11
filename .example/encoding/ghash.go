package main

import (
	"fmt"
	"strconv"

	"github.com/gogf/gf/encoding/ghash"
)

func main() {
	m := make(map[uint64]bool)
	for i := 0; i < 100000000; i++ {
		hash := ghash.BKDRHash64([]byte("key_" + strconv.Itoa(i)))
		if _, ok := m[hash]; ok {
			fmt.Printf("duplicated hash %d\n", hash)
		} else {
			m[hash] = true
		}
	}
}
