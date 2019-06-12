// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// @author wenzi1
// @date   20180604

// Package gcharset provides converting string to requested character encoding.
//
// 字符集转换方法,
package gcharset

import (
	"github.com/gogf/gf/g/encoding/internal"
)

// 2个字符集之间的转换
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	return internal.Convert(dstCharset, srcCharset, src)
}

// 指定字符集转UTF8
func ToUTF8(charset string, src string) (dst string, err error) {
	return Convert("UTF-8", charset, src)
}

// UTF8转指定字符集
func UTF8To(charset string, src string) (dst string, err error) {
	return Convert(charset, "UTF-8", src)
}
