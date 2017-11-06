// 文件空间管理(不仅仅是碎片管理)，
// 可用于文件碎片维护及再利用，支持自动合并连续碎片空间

// @TODO 大数据量下数据合并存在问题
package gfilespace

import (
    "sync"
    "fmt"
    "os"
)

// 文件空间管理结构体
type Space struct {
    mu      sync.RWMutex    // 并发操作锁
    blocks  []Block         // 空间区块列表(按照区块大小排序)
    indexes []Block         // 空间区块列表(按照索引大小排序)
}

// 文件空闲块
type Block struct {
    index int  // 文件偏移量
    size  uint // 区块大小(byte)
}

// 创建一个空间管理器
func New() *Space {
    return &Space {
        blocks  : make([]Block, 0),
        indexes : make([]Block, 0),
    }
}

// 根据index和size准确查找blocks中对应的区块
func (space *Space) getBlockPositionByIndexAndSize(index int, size uint) int {
    mid, _ := space.searchBlockBySize(size)
    // 往后继续匹配index
    for i := mid; i < len(space.blocks); i++ {
        if space.blocks[i].index == index {
            return i
        }
        if space.blocks[i].size != size {
            break
        }
    }
    // 往前继续匹配index
    for i := mid - 1; i >= 0; i-- {
        if space.blocks[i].index == index {
            return i
        }
        if space.blocks[i].size != size {
            break
        }
    }
    return -1
}

// 根据index和size准确查找indexes中对应的区块
func (space *Space) getIndexPositionByIndex(index int) int {
    mid, cmp := space.searchBlockByIndex(index)
    if cmp == 0 {
        return mid
    }
    return -1
}

// 内部按照索引检查合并
func (space *Space) checkIndexMergeFromIndex(from int) {
    for i := from; i < len(space.indexes); i++ {
        next := from + 1
        if next >= len(space.indexes) {
            break
        }
        if space.indexes[i].index + int(space.indexes[i].size) >= space.indexes[next].index {
            // 更新区块大小，必须同时更新indexes和blocks
            if blockpos := space.getBlockPositionByIndexAndSize(space.indexes[i].index, space.indexes[i].size); blockpos != -1 {
                space.indexes[i].size = uint(space.indexes[next].index + int(space.indexes[next].size) - space.indexes[i].index)
                // 由于涉及到列表排序，因此需要删除后重新插入
                block         := Block{space.blocks[blockpos].index, space.indexes[i].size}
                space.blocks   = space.removeBlock(space.blocks, blockpos)
                blockpos, cmp := space.searchBlockBySize(block.size)
                space.blocks   = space.insertBlock(space.blocks, block, blockpos, cmp)
            } else {
                fmt.Printf("update block missing for (%d, %d)\n", space.indexes[i].index, space.indexes[i].size)
                os.Exit(1)
            }

            // 合并后删除next项(正常情况下是肯定存在的)
            if blockpos := space.getBlockPositionByIndexAndSize(space.indexes[next].index, space.indexes[next].size); blockpos != -1 {
                // 首先删除blocks中对应的区块
                space.blocks  = space.removeBlock(space.blocks, blockpos)
                // 其次删除index对应区块
                space.indexes = space.removeBlock(space.indexes, next)
            } else {
                fmt.Printf("remove block missing for (%d, %d)\n", space.indexes[next].index, space.indexes[next].size)
                os.Exit(1)
            }

            // 递归处理
            space.checkIndexMergeFromIndex(i)
            break
        } else {
            break
        }
    }
}

// 添加一项, cmp < 0往前插入，cmp >= 0往后插入
func (space *Space) insertBlock(slice []Block, block Block, index int, cmp int) []Block {
    pos := index
    if cmp == -1 {
        // 添加到前面
    } else {
        // 添加到后面
        pos = index + 1
        if pos >= len(slice) {
            pos = len(slice)
        }
    }
    rear  := append([]Block{}, slice[pos : ]...)
    slice  = append(slice[0 : pos], block)
    slice  = append(slice, rear...)
    return slice
}


// 删除一项
func (space *Space) removeBlock(slice []Block, index int) []Block {
    return append(slice[:index], slice[index + 1:]...)
}


// 搜索空闲空间，返回空间 匹配size或者无法匹配时其附近 的空闲块索引地址，并返回匹配结果
func (space *Space) searchBlockBySize(size uint) (int, int) {
    min := 0
    max := len(space.blocks) - 1
    mid := 0
    cmp := -2
    for {
        if cmp == 0 || min > max {
            break
        }
        for {
            mid   = int((min + max) / 2)
            item := space.blocks[mid]
            if size < item.size {
                max = mid - 1
                cmp = -1
            } else if size > item.size {
                min = mid + 1
                cmp = 1
            } else {
                cmp = 0
                break
            }
            //fmt.Println("min:", min, "max:", max)
            if cmp == 0 || min > max {
                break
            }
        }
    }
    //fmt.Println(space.blocks)
    //fmt.Println(mid, cmp)
    //fmt.Println()
    return mid, cmp
}

// 搜索索引位置
func (space *Space) searchBlockByIndex(index int) (int, int) {
    min := 0
    max := len(space.indexes) - 1
    mid := 0
    cmp := -2
    for {
        if cmp == 0 || min > max {
            break
        }
        for {
            mid   = int((min + max) / 2)
            item := space.indexes[mid]
            if index < item.index {
                max = mid - 1
                cmp = -1
            } else if index > item.index {
                min = mid + 1
                cmp = 1
            } else {
                cmp = 0
                break
            }
            //fmt.Println("min:", min, "max:", max)
            if cmp == 0 || min > max {
                break
            }
        }
    }
    //fmt.Println(space.blocks)
    //fmt.Println(mid, cmp)
    //fmt.Println()
    return mid, cmp
}

// 获得碎片偏移量
func (block *Block) Index() int {
    return block.index
}

// 获得碎片大小
func (block *Block) Size() uint {
    return block.size
}


