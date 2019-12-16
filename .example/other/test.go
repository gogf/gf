package main

import (
	"github.com/gogf/gf/.example/frame/mvc/app/model/article"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Dump(article.FindAll(g.Slice{2, 3}))
}
