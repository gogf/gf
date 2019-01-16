// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gbase64 provides useful API for BASE64 encoding/decoding algorithms.
package gbase64

import (
    "encoding/base64"
)

// base64 encode
func Encode(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64 decode
func Decode(str string) (string, error) {
    s, e := base64.StdEncoding.DecodeString(str)
    return string(s), e
}