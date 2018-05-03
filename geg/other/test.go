package main

import (
    "runtime"
    "strconv"
    "fmt"
    "gitee.com/johng/gf/g/util/gregx"
    "gitee.com/johng/gf/g/os/gfile"
)

func Test() {
    backtrace := "Trace:\n"
    index     := 0
    for i := 1; i < 10000; i++ {
        if _, cfile, cline, ok := runtime.Caller(i); ok {
            // 不打印出go源码路径
            if !gregx.IsMatchString("^" + runtime.GOROOT(), cfile) {
                fmt.Println(gfile.Dir(cfile))
                backtrace += strconv.Itoa(index) + ". " + cfile + ":" + strconv.Itoa(cline) + "\n"
                index++
            }
        } else {
            break
        }
    }
    fmt.Println(backtrace)
}

func Test2() {
    Test()
}

func main() {
    Test2()
}