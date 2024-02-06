// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type MyTime = *gtime.Time

type Src struct {
	A MyTime
}

type Dst struct {
	B string
}

type SrcWrap struct {
	Value Src
}

type DstWrap struct {
	Value Dst
}

// SrcToDstConverter is custom converting function for custom type.
func SrcToDstConverter(src Src) (dst *Dst, err error) {
	return &Dst{B: src.A.Format("Y-m-d")}, nil
}

// SrcToDstConverter is custom converting function for custom type.
func main() {
	// register custom converter function.
	err := gconv.RegisterConverter(SrcToDstConverter)
	if err != nil {
		panic(err)
	}

	// custom struct converting.
	var (
		src = Src{A: gtime.Now()}
		dst *Dst
	)
	err = gconv.Scan(src, &dst)
	if err != nil {
		panic(err)
	}

	fmt.Println("src:", src)
	fmt.Println("dst:", dst)

	// custom struct attributes converting.
	var (
		srcWrap = SrcWrap{Src{A: gtime.Now()}}
		dstWrap *DstWrap
	)
	err = gconv.Scan(srcWrap, &dstWrap)
	if err != nil {
		panic(err)
	}

	fmt.Println("srcWrap:", srcWrap)
	fmt.Println("dstWrap:", dstWrap)
}
