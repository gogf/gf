package main

import (
    "g/os/gfilespace"
    "g/util/gtime"
    "fmt"
)



func main() {
    t1 := gtime.Microsecond()
    space := gfilespace.New()
    //for i := 10; i > 0; i-- {
    for i := 1; i <= 10; i++ {
        space.AddBlock(i, uint(i*10))
    }


    fmt.Println(gtime.Microsecond() - t1)
    fmt.Println(space.GetAllBlocksByIndex()[0])
    //fmt.Println(space.GetBlock(15))
    //fmt.Println(space.GetBlock(15))
}