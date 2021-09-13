package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/i18n/gi18n"
)

func main() {
	t := gi18n.New()
	t.SetLanguage("ja")
	err := t.SetPath("./i18n-dir")
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Translate(context.TODO(), `hello`))
	fmt.Println(t.Translate(context.TODO(), `{#hello}{#world}!`))
}
