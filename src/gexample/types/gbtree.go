package main

import (
    "g/core/types/gbtree"
    "fmt"
    "g/util/gtime"
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
    tr := gbtree.New(10)

    t1 := gtime.Microsecond()
    for i := 0; i < 10000000; i++ {
        tr.ReplaceOrInsert(gbtree.Item(Block{i, uint(i*10)}))
    }
    fmt.Println("create", gtime.Microsecond() - t1)

    t2 := gtime.Microsecond()
    tr.Get(gbtree.Item(Block{9999990, 0}))
    fmt.Println("get", gtime.Microsecond() - t2)

    t3 := gtime.Microsecond()
    var b Block
    tr.AscendGreaterOrEqual(gbtree.Item(Block{9999999, 0}), func(item gbtree.Item) bool {
        b = item.(Block)
        return false
    })
    fmt.Println("asc fetch", gtime.Microsecond() - t3, b)

}
