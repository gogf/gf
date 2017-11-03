package main

import (
    "g/os/gfilespace"
    "fmt"
    "g/util/grand"
    "g/util/gtime"
    "time"
)



func main() {
    //t1 := gtime.Microsecond()
    space := gfilespace.New()

    gtime.SetInterval(3*time.Second, func() bool {
        fmt.Println(len(space.GetAllBlocksByIndex()))
        fmt.Println(len(space.GetAllBlocksBySize()))
        return true
    })

    //for i := 10; i > 0; i-- {
    for i := 1; i <= 1000000; i++ {
        space.AddBlock(i*grand.Rand(0, 10000000), uint(i*10))
    }
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