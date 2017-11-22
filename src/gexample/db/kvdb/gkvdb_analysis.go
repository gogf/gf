package main

import (
    "fmt"
    "g/os/gfile"
    "g/encoding/gbinary"
)

type Block struct {
    index int
    size  uint
}

func main() {

    content := gfile.GetBinContents("/tmp/blocks")

    blocks  := make([]Block, 0)
    for i := 0; i < len(content); i += 12 {
        block := Block{
            int(gbinary.DecodeToInt64(content[i : i + 8])),
            uint(gbinary.DecodeToUint32(content[i + 8 : i + 12])),
        }
        blocks = append(blocks, block)
    }
    for i := 0; i < len(blocks); i++ {
        if i + 1 == len(blocks) {
            break
        }
        //fmt.Println(blocks[i].index, blocks[i].size)
        if blocks[i].index + int(blocks[i].size) >= blocks[i+1].index {
            fmt.Println(blocks[i].index, "+", blocks[i].size, ">=", blocks[i+1].index)
            break
        }
    }


    //fs      := gfilespace.New()
    //blocks  := make([]gfilespace.Block, 0)
    //for i := 0; i < len(content); i += 12 {
    //    fs.AddBlock(
    //        int(gbinary.DecodeToInt64(content[i : i + 8])),
    //        uint(gbinary.DecodeToUint32(content[i + 8 : i + 12])),
    //    )
    //
    //}
    //for _, v := range fs.GetAllBlocks() {
    //    blocks = append(blocks, v)
    //}
    //
    //for i := 0; i < len(blocks); i++ {
    //    if i + 1 == len(blocks) {
    //        break
    //    }
    //    fmt.Println(blocks[i].Index(), blocks[i].Size())
    //    if blocks[i].Index() + int(blocks[i].Size()) >= blocks[i+1].Index() {
    //        fmt.Println(blocks[i].Index(), "+", blocks[i].Size(), ">=", blocks[i+1].Index())
    //        break
    //    }
    //}

    //fmt.Println(blocks)
}