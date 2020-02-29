package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	content := `{{.name}} says "a{#hello}{#world}!"`
	result1, _ := g.View().ParseContent(content, g.Map{
		"name":         "john",
		"I18nLanguage": "zh-CN",
	})
	fmt.Println(result1)

	result2, _ := g.View().ParseContent(content, g.Map{
		"name":         "john",
		"I18nLanguage": "ja",
	})
	fmt.Println(result2)
}
