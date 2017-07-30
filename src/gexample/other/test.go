package main

import (
    "fmt"
    "g/util/gutil"
)

func main() {
    s := []byte{1,2,3,4,5}
    buffer1 := make([]byte, 5)
    buffer2 := make([]byte, 0)
    buffer2 = s[0:]
    fmt.Println(buffer1)
    fmt.Println(buffer2)
    fmt.Println(gutil.MergeSlice(buffer1, buffer2))



}