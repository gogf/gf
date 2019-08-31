package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"

	_ "github.com/gogf/gf/os/gres/testdata"
)

func main() {
	t := g.I18n()
	t.SetLanguage("ja")
	err := t.SetPath("/i18n-dir")
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Translate(`hello`))
	fmt.Println(t.Translate(`{{hello}}{{world}}!`))
}
