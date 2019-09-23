package main

import (
	"fmt"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"

	"github.com/gogf/gf/debug/gdebug"
)

func main() {
	cdnUrl := "http://localhost"
	content := `
    <link rel="stylesheet" href="/plugin/amazeui-2.7.2/css/amazeui.min.css">
    <link rel="stylesheet" href="/plugin/markdown-css/github-markdown.min.js">
    <link rel="stylesheet" href="/plugin/prism/prism.css">
    <link rel="stylesheet" href="/resource/css/document/style.css">
    <link rel="icon" href="/resource/image/favicon.ico" type="image/x-icon">
`
	s, err := gregex.ReplaceStringFuncMatch(`(href|src)=['"](.+?)['"]`, content, func(match []string) string {
		link := match[2]
		if len(link) == 0 {
			return match[0]
		}
		if link[0:1] != "/" && link[0:1] != "#" {
			if len(link) > 10 && link[0:10] == "javascript" {
				return match[0]
			}
			if len(link) > 7 && link[0:7] == "mailto:" {
				return match[0]
			}
			if len(link) > 4 && link[0:4] == "http" {
				return match[0]
			}
			link = "/" + link
		}
		if link[0:1] == "/" {
			switch gfile.ExtName(link) {
			case "png", "jpg", "jpeg", "gif", "js", "css", "otf", "eot", "ttf", "woff", "woff2":
				return fmt.Sprintf(`%s="%s%s?%s"`, match[1], cdnUrl, link, gdebug.BinVersion())
			}
		}
		return match[0]
	})
	fmt.Println(err)
	fmt.Println(s)
}
