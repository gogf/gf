package main

import (
	"fmt"
)

func Test() (int, int) {
	return 1, 1
}

func Assert(v1, v2, v3 interface{}) {
	fmt.Println(v1)
}

func F(v ...interface{}) []interface{} {
	return v
}

func main() {
	Assert(F(Test()), 2, 3)
}
