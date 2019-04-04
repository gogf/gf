package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gcfg"
)

func main() {
	fmt.Println(gcfg.Instance().GetString("viewpath"))
	fmt.Println(gcfg.Instance().GetString("database.default.0.host"))
}