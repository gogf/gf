package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
)

func main(){
    for i := 0; i < 100; i++ {
        gfile.Create(fmt.Sprintf(`/Users/john/Documents/test/%d`, i))
    }
}
