// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

//gcharset
//@author wenzi1
//@date   20180604

package gcharset

import (
	"github.com/axgle/mahonia"
	"errors"
	"fmt"
)

//获取支持的字符集
func GetCharset(name string) (*mahonia.Charset, error) {
	s := mahonia.GetCharset(name)
	if s == nil {
		return nil, errors.New(fmt.Sprintf("not support charset: %s", name))
	}
	return s, nil
}
//2个字符集之间的转换
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

//指定字符集转UTF8
func ToUTF8(charset string, src string) (dst string, err error) {
	return Convert("UTF-8", charset, src)
}

//UTF8转指定字符集
func UTF8To(charset string, src string) (dst string, err error) {
	return  Convert(charset, "UTF-8", src)
}