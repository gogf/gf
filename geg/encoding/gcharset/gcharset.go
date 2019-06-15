package main

import (
	"fmt"
	"github.com/gogf/gf/g/encoding/gcharset"
)

func main() {
	src := "~{;(<dR;:x>F#,6@WCN^O`GW!#"
	srcCharset := "GB2312"
	dstCharset := "UTF-8"
	str, err := gcharset.Convert(dstCharset, srcCharset, src)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
