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
