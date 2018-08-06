package main

import (
	"fmt"
	"gitee.com/johng/gf/g/util/gregex"
)

func main() {
	name := "page"
	path := "/page/template/{page}.html"
	rule := fmt.Sprintf(`{%s}`, name, name)
	tpl, err := gregex.ReplaceString(rule, `{.page}`, path)
	fmt.Println(err)
	fmt.Println(tpl)
}
