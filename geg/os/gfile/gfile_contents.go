package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    path    := "/tmp/temp"
    content := `123
456
789`
    gfile.PutContents(path, content)
    f, err := gfile.Open(path)
    if err != nil {
        panic(err)
    }
    fmt.Println(gfile.GetBinContentsTilChar(f, '\n', 0))
}
