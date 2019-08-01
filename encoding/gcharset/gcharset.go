// Copyright 2018-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package charset implements character-set conversion functionality.
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
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
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

// Supported returns whether charset <charset> is supported.
func Supported(charset string) bool {
	return getEncoding(charset) != nil
}

// Convert converts <src> charset encoding from <srcCharset> to <dstCharset>,
// and returns the converted string.
// It returns <src> as <dst> if it fails converting.
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src
	// Converting <src> to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", fmt.Errorf("%s to utf8 failed. %v", srcCharset, err)
			}
			src = string(tmp)
		} else {
			return dst, errors.New(fmt.Sprintf("unsupport srcCharset: %s", srcCharset))
		}
	}
	// Do the converting from UTF-8 to <dstCharset>.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", fmt.Errorf("utf to %s failed. %v", dstCharset, err)
			}
			dst = string(tmp)
		} else {
			return dst, errors.New(fmt.Sprintf("unsupport dstCharset: %s", dstCharset))
		}
	} else {
		dst = src
	}
	return dst, nil
}

// ToUTF8 converts <src> charset encoding from <srcCharset> to UTF-8 ,
// and returns the converted string.
func ToUTF8(srcCharset string, src string) (dst string, err error) {
	return Convert("UTF-8", srcCharset, src)
}

// UTF8To converts <src> charset encoding from UTF-8 to <dstCharset>,
// and returns the converted string.
func UTF8To(dstCharset string, src string) (dst string, err error) {
	return Convert(dstCharset, "UTF-8", src)
}

// getEncoding returns the encoding.Encoding interface object for <charset>.
// It returns nil if <charset> is not supported.
func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetAlias[charset]; ok {
		charset = c
	}
	if e, err := ianaindex.MIB.Encoding(charset); err == nil && e != nil {
		return e
	}
	return nil
}
