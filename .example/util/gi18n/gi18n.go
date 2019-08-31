package main

import (
	"fmt"

	"github.com/gogf/gf/util/gi18n"
)

func main() {
	t := gi18n.New()
	t.SetLanguage("ja")
	fmt.Println(t.Translate(`hello`))
	fmt.Println(t.Translate(`{#hello}}{#world}}!`))

}
