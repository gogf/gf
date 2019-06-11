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
package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/third/github.com/axgle/mahonia"
	"github.com/gogf/gf/third/golang.org/x/text/encoding"
	"github.com/gogf/gf/third/golang.org/x/text/encoding/japanese"
	"github.com/gogf/gf/third/golang.org/x/text/encoding/korean"
	"github.com/gogf/gf/third/golang.org/x/text/encoding/simplifiedchinese"
	"github.com/gogf/gf/third/golang.org/x/text/encoding/traditionalchinese"
	"github.com/gogf/gf/third/golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

var encodingMap *gmap.Map

func init() {
	encodingMap = gmap.New()
	encodingMap.Sets(
		map[interface{}]interface{}{
			"GBK":       simplifiedchinese.GBK,
			"GB18030":   simplifiedchinese.GB18030,
			"HZGB2312":  simplifiedchinese.HZGB2312,
			"GB2312":    simplifiedchinese.HZGB2312,
			"EUCJP":     japanese.EUCJP,
			"ISO2022JP": japanese.ISO2022JP,
			"SHIFTJIS":  japanese.ShiftJIS,
			"EUCKR":     korean.EUCKR,
			"BIG5":      traditionalchinese.Big5,
		})
}

func GetCharset(charset string) bool {
	c := strings.ToUpper(charset)
	if encodingMap.Contains(c) == false {
		if mahonia.GetCharset(c) == nil {
			return false
		}
	}
	return true
}

// 2个字符集之间的转换
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	srcCharsetUpper := strings.ToUpper(srcCharset)
	dstCharsetUpper := strings.ToUpper(dstCharset)

	if srcCharsetUpper == dstCharsetUpper {
		return src, nil
	}

	s := new(mahonia.Charset)
	d := new(mahonia.Charset)
	srctmp := src

	if srcCharset != "UTF-8" {
		enc := encodingMap.Get(srcCharset)
		if enc != nil {
			tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), enc.(encoding.Encoding).NewDecoder()))
			if err != nil {
				return "", fmt.Errorf("%s to utf8 failed. %v", srcCharset, err)
			}
			srctmp = string(tmp)
		} else {
			s = mahonia.GetCharset(srcCharsetUpper)
			if s == nil {
				return "", errors.New(fmt.Sprintf("not support charset:%s", srcCharset))
			}

			if s.Name != "UTF-8" {
				srctmp = s.NewDecoder().ConvertString(srctmp)
			}
		}
	}

	dst = srctmp

	if dstCharset != "UTF-8" {
		enc := encodingMap.Get(dstCharset)
		if enc != nil {
			tmp, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(srctmp)), enc.(encoding.Encoding).NewEncoder()))
			if err != nil {
				return "", fmt.Errorf("utf to %s failed. %v", dstCharset, err)
			}
			dst = string(tmp)
		} else {
			d = mahonia.GetCharset(dstCharsetUpper)
			if d == nil {
				return "", errors.New(fmt.Sprintf("not support charset:%s", dstCharset))
			}

			dst = srctmp
			if d.Name != "UTF-8" {
				dst = d.NewEncoder().ConvertString(dst)
			}
		}
	}
	return dst, nil
}
