package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gvalid"
)

func main() {
	g.I18n().SetLanguage("cn")
	err := gvalid.Check("", "required", nil)
	fmt.Println(err.String())
}
