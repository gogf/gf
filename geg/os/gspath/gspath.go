package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gspath"
)

func main() {
    sp   := gspath.New()
    path := "/Users/john/Workspace"
    rp, err := sp.Add(path)
    fmt.Println(err)
    fmt.Println(rp)
    fmt.Println(len(sp.AllPaths()))

    //gtime.SetInterval(5*time.Second, func() bool {
    //    g.Dump(sp.AllPaths())
    //    return true
    //})

    select {

    }
}
