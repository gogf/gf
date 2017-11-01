package main

import (
    "g/os/gfilespace"
    "fmt"
)



func main() {
    //t1 := gtime.Microsecond()
    space := gfilespace.New()
    //for i := 10; i > 0; i-- {
    //for i := 1; i <= 10; i++ {
    //    space.AddBlock(i, uint(i*10))
    //}

    space.AddBlock(319808, 64)
    space.AddBlock(319872, 64)
    space.AddBlock(319936, 64)

    //space.Empty()

    //fmt.Println(gtime.Microsecond() - t1)
    fmt.Println(space.GetAllBlocksByIndex())
    //fmt.Println(space.GetBlock(15))
    //fmt.Println(space.GetBlock(15))
}