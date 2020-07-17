package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

func main() {
	s := "我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人我是中国人"
	tplContent := `
{{.str | strlimit 10  "..."}}
`
	content, err := g.View().ParseContent(tplContent, g.Map{
		"str": s,
	})
	fmt.Println(err)
	fmt.Println(content)
}
