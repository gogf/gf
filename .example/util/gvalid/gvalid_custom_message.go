package main

import (
	"fmt"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gvalid"
)

func main() {
	g.I18n().SetLanguage("cn")
	err := gvalid.Check("", "required", nil)
	fmt.Println(err.String())
}
