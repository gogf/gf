package main

import (
	"github.com/gogf/gf/v2/frame/g"
)

func PrintLog(content string) {
	g.Log().Skip(0).Line().Print("line number with skip:", content)
	g.Log().Line(true).Print("line number without skip:", content)
}

func main() {
	PrintLog("just test")
}
