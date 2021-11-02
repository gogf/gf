package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gvalid"
)

func main() {
	g.I18n().SetLanguage("cn")
	err := gvalid.Check(context.TODO(), "", "required", nil)
	fmt.Println(err.String())
}
