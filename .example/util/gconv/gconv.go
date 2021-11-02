package main

import (
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	i := 123.456
	fmt.Printf("%10s %v\n", "Int:", gconv.Int(i))
	fmt.Printf("%10s %v\n", "Int8:", gconv.Int8(i))
	fmt.Printf("%10s %v\n", "Int16:", gconv.Int16(i))
	fmt.Printf("%10s %v\n", "Int32:", gconv.Int32(i))
	fmt.Printf("%10s %v\n", "Int64:", gconv.Int64(i))
	fmt.Printf("%10s %v\n", "Uint:", gconv.Uint(i))
	fmt.Printf("%10s %v\n", "Uint8:", gconv.Uint8(i))
	fmt.Printf("%10s %v\n", "Uint16:", gconv.Uint16(i))
	fmt.Printf("%10s %v\n", "Uint32:", gconv.Uint32(i))
	fmt.Printf("%10s %v\n", "Uint64:", gconv.Uint64(i))
	fmt.Printf("%10s %v\n", "Float32:", gconv.Float32(i))
	fmt.Printf("%10s %v\n", "Float64:", gconv.Float64(i))
	fmt.Printf("%10s %v\n", "Bool:", gconv.Bool(i))
	fmt.Printf("%10s %v\n", "String:", gconv.String(i))

	fmt.Printf("%10s %v\n", "Bytes:", gconv.Bytes(i))
	fmt.Printf("%10s %v\n", "Strings:", gconv.Strings(i))
	fmt.Printf("%10s %v\n", "Ints:", gconv.Ints(i))
	fmt.Printf("%10s %v\n", "Floats:", gconv.Floats(i))
	fmt.Printf("%10s %v\n", "Interfaces:", gconv.Interfaces(i))
}
