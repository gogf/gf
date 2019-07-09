// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gbase64 provides useful API for BASE64 encoding/decoding algorithm.
package gbase64

import (
	"encoding/base64"
)

// Encode encodes bytes with BASE64 algorithm.
func Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

// Decode decodes bytes with BASE64 algorithm.
func Decode(dst []byte) ([]byte, error) {
	src := make([]byte, base64.StdEncoding.DecodedLen(len(dst)))
	n, err := base64.StdEncoding.Decode(src, dst)
	return src[:n], err
}

// EncodeString encodes bytes with BASE64 algorithm.
func EncodeString(src []byte) string {
	return string(Encode(src))
}

// DecodeString decodes string with BASE64 algorithm.
func DecodeString(str string) ([]byte, error) {
	return Decode([]byte(str))
}
