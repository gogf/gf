package main

import (
    "g/os/gfilespace"
    "fmt"
)



func main() {
    //t1 := gtime.Microsecond()
    space := gfilespace.New()

    //t1 := gtime.Microsecond()
    for i := 1; i <= 1000000; i++ {
        //space.AddBlock(i*grand.Rand(0, 10000000), uint(i*10))
        space.AddBlock(i, uint(i*10))
        //fmt.Println(space.GetAllBlocks())
    }
    //fmt.Println("create", gtime.Microsecond() - t1)

    //t2 := gtime.Microsecond()
    fmt.Println(space.GetBlock(10))
    fmt.Println(space.GetBlock(10))
    fmt.Println(space.GetBlock(10))
    fmt.Println(space.GetBlock(10))
    fmt.Println(space.GetBlock(10))
    fmt.Println(space.GetBlock(10))
    fmt.Println(space.GetBlock(10))
    //fmt.Println("get", gtime.Microsecond() - t2)
    //
    //fmt.Println(space.GetAllBlocks())
    //fmt.Println(space.GetAllSizes())



    //add block: 1792 192
    //[{0 192} {512 192} {768 384} {1408 960}]
    //add block: 320 192
    //[{0 192} {320 192} {512 192} {768 384} {1408 960}]

    //add mt block 1618432 64
    //[{1618432 64}]
    //[{1618432 64}]
    //add mt block 1618496 64
    //[{1618432 128}]
    //[{1618432 64}]
    //space.AddBlock(467264,    64)
    //space.AddBlock(467200,    128)




    //space.Empty()

    //fmt.Println(gtime.Microsecond() - t1)

    //fmt.Println(space.GetBlock(15))
    //fmt.Println(space.GetBlock(15))
}