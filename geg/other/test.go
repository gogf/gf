package main

import (
	"fmt"
	"github.com/gogf/gf/g"
)

func main() {
	latestVersion := g.NewVar(nil, true)
	fmt.Println(latestVersion.IsNil())
}
