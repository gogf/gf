package main

import (
    "gitee.com/johng/gf/g/container/gbtree"
    "fmt"
)

type Block struct {
    index int  // 文件偏移量
    size  uint // 区块大小(byte)
}

func (block *Block) Less(item gbtree.Item) bool {
    if block.index < item.(*Block).index {
        return true
    }
    return false
}

func main () {
    tr := gbtree.New(10)

    //t1 := gtime.Microsecond()
    for i := 0; i < 10; i++ {
        tr.ReplaceOrInsert(&Block{i, uint(i*10)})
    }
    //fmt.Println("create", gtime.Microsecond() - t1)

    //t2 := gtime.Microsecond()
    //b := &Block{9, 10}
    //fmt.Println(tr.Get(b))
    //fmt.Println(tr.Delete(b))
    //fmt.Println(tr.Get(b))
    //fmt.Println("get", gtime.Microsecond() - t2)

    //t3 := gtime.Microsecond()
    //var b Block
    tr.AscendGreaterOrEqual(&Block{2, 0}, func(item gbtree.Item) bool {
        fmt.Println(item)
        return true
    })
    //fmt.Println("asc fetch", gtime.Microsecond() - t3, b)

}
