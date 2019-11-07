package main

import (
	"fmt"
	"regexp"
)

func main() {
	replaceCharReg, err := regexp.Compile(`[\-\.\_\s]+`)
	fmt.Println(err)
	fmt.Println(replaceCharReg.ReplaceAllString("s--s.s.a b", ""))
}
