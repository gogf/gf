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