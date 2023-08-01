// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcharset implements character-set conversion functionality.
//
// Supported Character Set:
//
// Chinese : GBK/GB18030/GB2312/Big5
//
// Japanese: EUCJP/ISO2022JP/ShiftJIS
//
// Korean  : EUCKR
//
// Unicode : UTF-8/UTF-16/UTF-16BE/UTF-16LE
//
// Other   : macintosh/IBM*/Windows*/ISO-*
package gcharset

import (
	"bytes"
	"context"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

var (
	// Alias for charsets.
	charsetAlias = map[string]string{
		"HZGB2312": "HZ-GB-2312",
		"hzgb2312": "HZ-GB-2312",
		"GB2312":   "HZ-GB-2312",
		"gb2312":   "HZ-GB-2312",
	}
)

// Supported returns whether charset `charset` is supported.
func Supported(charset string) bool {
	return getEncoding(charset) != nil
}

// Convert converts `src` charset encoding from `srcCharset` to `dstCharset`,
// and returns the converted string.
// It returns `src` as `dst` if it fails converting.
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src
	// Converting `src` to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", gerror.Wrapf(err, `convert string "%s" to utf8 failed`, srcCharset)
			}
			src = string(tmp)
		} else {
			return dst, gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported srcCharset "%s"`, srcCharset)
		}
	}
	// Do the converting from UTF-8 to `dstCharset`.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", gerror.Wrapf(err, `convert string from utf8 to "%s" failed`, dstCharset)
			}
			dst = string(tmp)
		} else {
			return dst, gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported dstCharset "%s"`, dstCharset)
		}
	} else {
		dst = src
	}
	return dst, nil
}

// ToUTF8 converts `src` charset encoding from `srcCharset` to UTF-8 ,
// and returns the converted string.
func ToUTF8(srcCharset string, src string) (dst string, err error) {
	return Convert("UTF-8", srcCharset, src)
}

// UTF8To converts `src` charset encoding from UTF-8 to `dstCharset`,
// and returns the converted string.
func UTF8To(dstCharset string, src string) (dst string, err error) {
	return Convert(dstCharset, "UTF-8", src)
}

// getEncoding returns the encoding.Encoding interface object for `charset`.
// It returns nil if `charset` is not supported.
func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetAlias[charset]; ok {
		charset = c
	}
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	return enc
}
