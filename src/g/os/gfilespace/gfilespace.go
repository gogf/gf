// 文件空间管理，包括内容空间分配、文件碎片维护及再利用
package gfilespace

import (
    "sync"
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

// 申请空间，返回文件地址及大小，返回成功后则在管理器中删除该空闲块
func (space *Space) GetBlock(size uint) (int, uint) {
    space.mu.RLock()
    defer space.mu.RUnlock()

    mid, cmp := space.searchBlockBySize(size)
    // 必须找到一块不比size小的区块
    if cmp != 0 {
        cmp = -1
        for i := mid; i < len(space.blocks); i++ {
            if space.blocks[i].size >= size {
                mid = i
                cmp = 1
                break;
            }
        }
    }
    // 找到符合要求的区块，返回前进行删除
    if cmp >= 0 {
        index := space.blocks[mid].index
        size  := space.blocks[mid].size
        space.blocks = space.removeBlock(space.blocks, mid)
        if indexpos := space.getIndexPositionByIndex(index); indexpos != -1 {
            space.indexes = space.removeBlock(space.indexes, indexpos)
        }
        return index, size
    }
    return -1, 0
}

// 删除一项
func (space *Space) removeBlock(slice []Block, index int) []Block {
    blocks  := make([]Block, 0)
    blocks   = append(blocks, space.indexes[ : index]...)
    if index + 1 <= len(slice) - 1 {
        blocks  = append(blocks, space.indexes[index + 1 : ]...)
    }
    return blocks
}

// 根据index和size准确查找blocks中对应的区块
func (space *Space) getBlockPositionByIndexAndSize(index int, size uint) int {
    mid, _ := space.searchBlockBySize(size - 1)
    for i := mid; i < len(space.blocks); i++ {
        if space.blocks[i].index == index {
            return i
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
            // 更新区块大小
            space.indexes[i].size = uint(space.indexes[next].index + int(space.indexes[next].size) - space.indexes[i].index)
            if blockpos := space.getBlockPositionByIndexAndSize(space.indexes[i].index, space.indexes[i].size); blockpos != -1 {
                space.blocks[blockpos].size = space.indexes[i].size
            }
            // 合并后删除next项，首先删除blocks中对应的区块
            if blockpos := space.getBlockPositionByIndexAndSize(space.indexes[next].index, space.indexes[next].size); blockpos != -1 {
                space.blocks = space.removeBlock(space.blocks, blockpos)
            }
            // 其次删除index对应区块
            space.indexes = space.removeBlock(space.indexes, next)
            // 递归处理
            space.checkIndexMergeFromIndex(i)
            break
        } else {
            break
        }
    }
}

// 添加空闲空间到管理器
func (space *Space) AddBlock(index int, size uint) {
    if size <= 0 {
        return
    }

    space.mu.Lock()
    defer space.mu.Unlock()

    block    := Block{index, size}
    // 首先按照索引搜索，插入到合适的位置
    mid, cmp := space.searchBlockByIndex(index)
    indexpos := mid
    if cmp == -1 {
        // 添加到前面
        indexpos = mid - 1
        if indexpos < 0 {
            indexpos = 0
        }
    } else {
        // 添加到后面
        indexpos = mid + 1
        if indexpos >= len(space.indexes) {
            indexpos = len(space.indexes)
        }
    }
    indexes := make([]Block, 0)
    indexes  = append(indexes, space.indexes[0 : indexpos]...)
    indexes  = append(indexes, block)
    indexes  = append(indexes, space.indexes[indexpos : ]...)
    space.indexes = indexes

    // 其次按照区块进行索引，插入到合适的位置
    mid, cmp  = space.searchBlockBySize(size)
    blockpos := mid
    if cmp == -1 {
        // 添加到前面
        blockpos = mid - 1
        if blockpos < 0 {
            blockpos = 0
        }
    } else {
        // 添加到后面
        blockpos = mid + 1
        if blockpos >= len(space.blocks) {
            blockpos = len(space.blocks)
        }
    }
    blocks := make([]Block, 0)
    blocks  = append(blocks, space.blocks[0 : blockpos]...)
    blocks  = append(blocks, block)
    blocks  = append(blocks, space.blocks[blockpos : ]...)
    space.blocks = blocks

    // 区块检查合并
    checkpos := indexpos - 1
    if checkpos < 0 {
        checkpos = 0
    }
    space.checkIndexMergeFromIndex(checkpos)

    //fmt.Println(space.indexes)
    //fmt.Println(space.blocks)
    //fmt.Println()
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


