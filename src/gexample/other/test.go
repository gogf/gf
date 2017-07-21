package main

import (
    "fmt"
    "g/encoding/gmd5"
    "time"
)

func FormatTs(ts int64) string {
    return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func FormatTsInt(ts int) string {
    return FormatTs(int64(ts))
}
func main() {
    fmt.Println(F)
}