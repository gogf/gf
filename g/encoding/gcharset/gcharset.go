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
// 使用mahonia实现的字符集转换方法，支持的字符集包括常见的utf8/UTF-16/UTF-16LE/macintosh/big5/gbk/gb18030,支持的全量字符集可以参考mahonia包
package gcharset

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/third/github.com/axgle/mahonia"
	"github.com/gogf/gf/third/golang.org/x/text/encoding/simplifiedchinese"
	"github.com/gogf/gf/third/golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

// 2个字符集之间的转换
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {

	if strings.EqualFold(srcCharset, dstCharset) {
		return src, nil
	}

	s := new(mahonia.Charset)
	d := new(mahonia.Charset)
	srctmp := src

	switch {
	case strings.EqualFold("GBK", srcCharset):
		tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewDecoder()))
		if err != nil {
			return "", fmt.Errorf("gbk to utf8 failed. %v", err)
		}
		srctmp = string(tmp)
	case strings.EqualFold("GB18030", srcCharset):
		tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GB18030.NewDecoder()))
		if err != nil {
			return "", fmt.Errorf("GB18030 to utf8 failed. %v", err)
		}
		srctmp = string(tmp)
	case strings.EqualFold("GB2312", srcCharset) || strings.EqualFold("HZGB2312", srcCharset):
		tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.HZGB2312.NewDecoder()))
		if err != nil {
			return "", fmt.Errorf("GB2312 to utf8 failed. %v", err)
		}
		srctmp = string(tmp)
	case strings.EqualFold("UTF-8", srcCharset):
	default:
		s = mahonia.GetCharset(srcCharset)
		if s == nil {
			return "", errors.New(fmt.Sprintf("not support charset:%s", srcCharset))
		}

		if s.Name != "UTF-8" {
			srctmp = s.NewDecoder().ConvertString(srctmp)
		}
	}

	dst = srctmp

	switch {
	case strings.EqualFold("GBK", dstCharset):
		tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(srctmp)), simplifiedchinese.GBK.NewEncoder()))
		if err != nil {
			return "", fmt.Errorf("utf to gbk failed. %v", err)
		}
		dst = string(tmp)
	case strings.EqualFold("GB18030", dstCharset):
		tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(srctmp)), simplifiedchinese.GB18030.NewEncoder()))
		if err != nil {
			return "", fmt.Errorf("utf8 to gb18030 failed. %v", err)
		}
		dst = string(tmp)
	case strings.EqualFold("GB2312", dstCharset) || strings.EqualFold("HZGB2312", dstCharset):
		tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(srctmp)), simplifiedchinese.HZGB2312.NewEncoder()))
		if err != nil {
			return "", fmt.Errorf("utf8 to gb2312 failed. %v", err)
		}
		dst = string(tmp)
	case strings.EqualFold("UTF-8", dstCharset):
	default:
		d = mahonia.GetCharset(dstCharset)
		if d == nil {
			return "", errors.New(fmt.Sprintf("not support charset:%s", dstCharset))
		}

		dst = srctmp
		if d.Name != "UTF-8" {
			dst = d.NewEncoder().ConvertString(dst)
		}

	}

	return dst, nil
}

// 指定字符集转UTF8
func ToUTF8(charset string, src string) (dst string, err error) {
	return Convert("UTF-8", charset, src)
}

// UTF8转指定字符集
func UTF8To(charset string, src string) (dst string, err error) {
	return Convert(charset, "UTF-8", src)
}
