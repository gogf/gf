package gmd5

import (
    "crypto/md5"
    "fmt"
    "encoding/json"
    "reflect"
)

// 将任意类型的变量进行md5摘要
func Encode( v interface{}) string {
    h := md5.New()
    if "string" == reflect.TypeOf(v).String() {
        h.Write([]byte(v.(string)))
    } else {
        b, err := json.Marshal(v)
        if err != nil {
            return ""
        } else {
            h.Write(b)
        }
    }
    return fmt.Sprintf("%x", h.Sum(nil))
}