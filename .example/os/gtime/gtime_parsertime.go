package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	content := `

2018-11-01 14:30:39 -- 67 -- assignDoctor:request -- {"问诊ID":"1017467","操作类型":2,"操作ID":52339491,"医生ID":52339491,"是否主动接单":"是"}
2018-11-01 14:35:55 -- 73 -- throwIntoPool:request -- {"问诊Id":1017474,"当前Id":null,"当前角色":null}
`
	if t := gtime.ParseTimeFromContent(content); t != nil {
		fmt.Println(t.String())
		fmt.Println(t.UTC())
		fmt.Println(gtime.Now().UTC())
	} else {
		panic("cannot parse time from content")
	}

	//if t := gtime.ParseTimeFromContent(content, "d/M/Y:H:i:s +0800"); t != nil {
	//    fmt.Println(t.String())
	//} else {
	//    panic("cannot parse time from content")
	//}
}
