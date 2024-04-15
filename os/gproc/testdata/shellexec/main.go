// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var content string
	var output string
	flag.StringVar(&content, "c", "", "写入内容")
	flag.StringVar(&output, "o", "", "写入路径")
	flag.Parse()
	fmt.Println(os.Args)
	fmt.Println(content)
	fmt.Println(output)
	if output != "" {
		file, err := os.Create(output)
		if err != nil {
			panic("create file fail: " + err.Error())
		}
		defer file.Close()
		file.WriteString(content)
	}
}
