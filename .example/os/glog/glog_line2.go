package main

import (
	"github.com/gogf/gf/frame/g"
)

func PrintLog(content string) {
	g.Log().Skip(0).Line().Println("line number with skip:", content)
	g.Log().Line(true).Println("line number without skip:", content)
}

func main() {
	PrintLog("just test")
}
