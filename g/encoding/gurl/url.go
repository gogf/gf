// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gurl provides useful API for URL handling.
package gurl

import "net/url"

// url encode string, is + not %20
func Encode(str string) string {
    return url.QueryEscape(str)
}

// url decode string
func Decode(str string) (string, error) {
    return url.QueryUnescape(str)
}
