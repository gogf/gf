package main

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/os/gspath"
    "gitee.com/johng/gf/g/util/gutil"
)

func main() {
    gutil.Dump(gspath.Get("/Users/john/Temp/config").AllPaths())
=======
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    fmt.Println(gfile.RealPath("config"))
>>>>>>> 104613b056b4b2f85af786638d37aa50d7bc06e2
}