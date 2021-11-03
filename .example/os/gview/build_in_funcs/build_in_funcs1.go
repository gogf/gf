package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	tplContent := `
{{"<div>测试</div>"|text}}
{{"<div>测试</div>"|html}}
{{"&lt;div&gt;测试&lt;/div&gt;"|htmldecode}}
{{"https://goframe.org"|url}}
{{"https%3A%2F%2Fgoframe.org"|urldecode}}
{{1540822968 | date "Y-m-d"}}
{{"1540822968" | date "Y-m-d H:i:s"}}
{{date "Y-m-d H:i:s"}}
{{compare "A" "B"}}
{{compare "1" "2"}}
{{compare 2 1}}
{{compare 1 1}}
{{"我是中国人" | substr 2 -1}}
{{"我是中国人" | substr 2  2}}
{{"我是中国人" | strlimit 2  "..."}}
{{"热爱GF热爱生活" | hidestr 20  "*"}}
{{"热爱GF热爱生活" | hidestr 50  "*"}}
{{"热爱GF热爱生活" | highlight "GF" "red"}}
{{"gf" | toupper}}
{{"GF" | tolower}}
{{"Go\nFrame" | nl2br}}
`
	content, err := g.View().ParseContent(tplContent, nil)
	fmt.Println(err)
	fmt.Println(string(content))
}
