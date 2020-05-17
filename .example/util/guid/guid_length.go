package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	fmt.Println(strconv.FormatUint(4589634556, 36))
	fmt.Println(strconv.FormatUint(math.MaxUint64-2, 36))
	fmt.Println(strconv.FormatUint(math.MaxUint32-1, 36))
	fmt.Println(strconv.FormatUint(math.MaxUint32, 36))
	fmt.Println(strconv.FormatUint(2000000, 36))
}
