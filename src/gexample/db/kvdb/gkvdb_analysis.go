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
        if blocks[i].index + int(blocks[i].size) >= blocks[i+1].index {
            fmt.Println(blocks[i].index, "+", blocks[i].size, ">=", blocks[i+1].index)
            //break
        }
    }

    //fmt.Println(blocks)
}