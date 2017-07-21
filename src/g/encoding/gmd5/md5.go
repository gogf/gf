package gmd5

import (
    "crypto/md5"
    "fmt"
)

// md5摘要
func Encode(s string) string {
    h := md5.New()
    h.Write([]byte(s))
    return fmt.Sprintf("%x", h.Sum(nil))
}
