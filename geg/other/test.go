package main

import (
	"fmt"
	"math"

	"github.com/gogf/gf/g/encoding/gbinary"
)

func main() {
	v := math.MaxUint16
	//v := []byte{255, 127}
	//ve := gbinary.Encode(v)
	//ve1 := gbinary.BeEncodeByLength(len(ve), v)
	//fmt.Println(ve)
	//fmt.Println(ve1)

	//fmt.Println(gbinary.LeDecodeToInt(gbinary.LeEncode(v)))
	fmt.Println(gbinary.BeDecodeToInt(gbinary.BeEncode(v)))
}
