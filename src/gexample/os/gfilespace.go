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
    //for i := 1; i <= 10; i++ {
    //    space.AddBlock(i, uint(i*10))
    //}
    space.AddBlock(640, 64)
    space.AddBlock(704, 64)
    space.AddBlock(768, 64)
    space.AddBlock(832, 64)
    space.AddBlock(896, 64)
    space.AddBlock(960, 64)
    space.AddBlock(1024, 64)
    space.AddBlock(1088, 64)
    space.AddBlock(1152, 64)
    space.AddBlock(1216, 64)

    space.Empty()

    fmt.Println(gtime.Microsecond() - t1)
    fmt.Println(space.GetAllBlocksByIndex())
    //fmt.Println(space.GetBlock(15))
    //fmt.Println(space.GetBlock(15))
}