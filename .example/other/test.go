package main

import (
	"fmt"
	"github.com/gogf/gf/text/gregex"
)

func main() {
	data := "@var(.prefix)您收到的验证码为：@var(.code)，请在@var(.expire)内完成验证"
	result, err := gregex.ReplaceStringFuncMatch(`(@var\(\.\w+\))`, data, func(match []string) string {
		fmt.Println(match)
		return "#"
	})
	fmt.Println(err)
	fmt.Println(result)
}
