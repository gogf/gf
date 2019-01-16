// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// @author wenzi1
// @date   20180604

// Package gcharset provides converting string to requested character encoding.
//
// 字符集转换方法,
// 使用mahonia实现的字符集转换方法，支持的字符集包括常见的utf8/UTF-16/UTF-16LE/macintosh/big5/gbk/gb18030,支持的全量字符集可以参考mahonia包
package gcharset

import (
	"gitee.com/johng/gf/third/github.com/axgle/mahonia"
	"errors"
	"fmt"
)


// 2个字符集之间的转换
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	s := mahonia.GetCharset(srcCharset)
	if 	s == nil {
		return "", errors.New(fmt.Sprintf("not support charset:%s", srcCharset))
	}

	d := mahonia.GetCharset(dstCharset)
	if d == nil {
		return "", errors.New(fmt.Sprintf("not support charset:%s", dstCharset))
	}
	
	srctmp := src
	if s.Name != "UTF-8" {
		srctmp = s.NewDecoder().ConvertString(srctmp)
	}
	
	dst = srctmp
	if d.Name != "UTF-8" {
		dst = d.NewEncoder().ConvertString(dst)
	}
	
	return dst, nil
}

// 指定字符集转UTF8
func ToUTF8(charset string, src string) (dst string, err error) {
	return Convert("UTF-8", charset, src)
}

// UTF8转指定字符集
func UTF8To(charset string, src string) (dst string, err error) {
	return  Convert(charset, "UTF-8", src)
}