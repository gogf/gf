package main

import (
	"fmt"

	"github.com/jin502437344/gf/text/gstr"
)

func main() {
	fmt.Println(gstr.HideStr("热爱GF热爱生活", 20, "*"))
	fmt.Println(gstr.HideStr("热爱GF热爱生活", 50, "*"))
}
