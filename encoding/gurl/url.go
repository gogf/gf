// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gurl provides useful API for URL handling.
package gurl

import (
	"net/url"
	"strings"
)

// Encode escapes the string so it can be safely placed
// inside a URL query.
func Encode(str string) string {
	return url.QueryEscape(str)
}

// Decode does the inverse transformation of Encode,
// converting each 3-byte encoded substring of the form "%AB" into the
// hex-decoded byte 0xAB.
// It returns an error if any % is not followed by two hexadecimal
// digits.
func Decode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// RawEncode does encode the given string according
// URL-encode according to RFC 3986.
// See http://php.net/manual/en/function.rawurlencode.php.
func RawEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// RawDecode does decode the given string
// Decode URL-encoded strings.
// See http://php.net/manual/en/function.rawurldecode.php.
func RawDecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// BuildQuery Generate URL-encoded query string.
// See http://php.net/manual/en/function.http-build-query.php.
func BuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// ParseURL Parse a URL and return its components.
// -1: all; 1: scheme; 2: host; 4: port; 8: user; 16: pass; 32: path; 64: query; 128: fragment.
// See http://php.net/manual/en/function.parse-url.php.
func ParseURL(str string, component int) (map[string]string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if component == -1 {
		component = 1 | 2 | 4 | 8 | 16 | 32 | 64 | 128
	}
	var components = make(map[string]string)
	if (component & 1) == 1 {
		components["scheme"] = u.Scheme
	}
	if (component & 2) == 2 {
		components["host"] = u.Hostname()
	}
	if (component & 4) == 4 {
		components["port"] = u.Port()
	}
	if (component & 8) == 8 {
		components["user"] = u.User.Username()
	}
	if (component & 16) == 16 {
		components["pass"], _ = u.User.Password()
	}
	if (component & 32) == 32 {
		components["path"] = u.Path
	}
	if (component & 64) == 64 {
		components["query"] = u.RawQuery
	}
	if (component & 128) == 128 {
		components["fragment"] = u.Fragment
	}
	return components, nil
}
