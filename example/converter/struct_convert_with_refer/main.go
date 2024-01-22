// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
)

type Src struct {
	A int
}

type Dst struct {
	B int
}

type SrcWrap struct {
	Value Src
}

type DstWrap struct {
	Value Dst
}

// SrcToDstConverter is custom converting function for custom type.
func SrcToDstConverter(src Src) (dst *Dst, err error) {
	return &Dst{B: src.A}, nil
}

func main() {
	// register custom converter function.
	err := gconv.RegisterConverter(SrcToDstConverter)
	if err != nil {
		panic(err)
	}

	// custom struct converting.
	var src = Src{A: 1}
	dst := gconv.ConvertWithRefer(src, Dst{})
	fmt.Println("src:", src)
	fmt.Println("dst:", dst)

	// custom struct attributes converting.
	var srcWrap = SrcWrap{Src{A: 1}}
	dstWrap := gconv.ConvertWithRefer(srcWrap, &DstWrap{})
	fmt.Println("srcWrap:", srcWrap)
	fmt.Println("dstWrap:", dstWrap)
}
