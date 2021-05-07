package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	// 36*36^2+36*36+36
	var s string
	fmt.Println(strconv.ParseUint("zzz", 36, 3))
	fmt.Println(1 << 1)
	// MaxInt64
	s = strconv.FormatUint(math.MaxUint64, 16)
	fmt.Println(s, len(s))
	// PID
	s = strconv.FormatInt(1000000, 36)
	fmt.Println(s, len(s))
}
