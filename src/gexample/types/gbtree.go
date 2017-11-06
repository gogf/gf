package main

import (
    "g/core/types/gbtree"
    "fmt"
    "g/util/gtime"
    "g/util/grand"
)

type Block struct {
    index int  // 文件偏移量
    size  uint // 区块大小(byte)
}

func (block Block) Less(item gbtree.Item) bool {
    if block.index < item.(Block).index {
        return true
    }
    return false
}

func main () {
    tr := gbtree.New(100)
    t1 := gtime.Microsecond()
    for i := 0; i < 1000000; i++ {
        tr.ReplaceOrInsert(gbtree.Item(Block{i*grand.Rand(0, 10000000), uint(i*10)}))
    }
    fmt.Println("create", gtime.Microsecond() - t1)

    t2 := gtime.Microsecond()
    tr.Get(gbtree.Item(Block{99999, 0}))
    fmt.Println("get", gtime.Microsecond() - t2)

    t3 := gtime.Microsecond()
    var b Block
    tr.AscendGreaterOrEqual(gbtree.Item(Block{99999, 0}), func(item gbtree.Item) bool {
        b = item.(Block)
        return false
    })
    fmt.Println("asc fetch", gtime.Microsecond() - t3, b)

    fmt.Println(tr.Get(gbtree.Item(Block{1, 0})))
}
