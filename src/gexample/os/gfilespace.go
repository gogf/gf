package main

import (
    "g/os/gfilespace"
    "fmt"
    "g/os/gfile"
    "g/encoding/gbinary"
    "strings"
    "strconv"
)



func main() {

    //t1 := gtime.Microsecond()
    space := gfilespace.New()


    content := gfile.GetContents("/tmp/blockops")
    for _, v := range strings.Split(content, "\n") {
        ss := strings.Split(v, ",")
        if len(ss) != 3 {
            fmt.Println(v)
            continue
        }
        t, _     := strconv.Atoi(ss[0])
        index, _ := strconv.ParseInt(ss[1], 10, 64)
        size, _  := strconv.Atoi(ss[2])
        //fmt.Println(index, size)
        if t == 0 {
            //fmt.Println(index, size)
            space.AddBlock(int(index), uint(size))
        } else {
            //fmt.Println(size)
            space.GetBlock(uint(size))
        }
    }

    blocks  := space.GetAllBlocks()
    fmt.Println(blocks)
    buffer  := make([]byte, 0)
    for _, b := range blocks {
        buffer = append(buffer, gbinary.EncodeInt64(int64(b.Index()))...)
        buffer = append(buffer, gbinary.EncodeUint32(uint32(b.Size()))...)
    }
    gfile.PutBinContents("/tmp/blocks2", buffer)

    return
    //t1 := gtime.Microsecond()
    for i := 1; i <= 10; i++ {
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