package main

import (
	"fmt"
	"net/url"
)

func main() {
	parse1, _ := url.Parse("https://gf.cdn.johng.cn")
	parse2, _ := url.Parse("https://gf.cdn.johng.cn/cli/")
	fmt.Println(parse1.Host)
	fmt.Println(parse2.Host)
}
