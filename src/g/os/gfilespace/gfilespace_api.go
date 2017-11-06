package gfilespace


// 清空数据列表
func (space *Space) Empty() {
    space.blocks  = make([]Block, 0)
    space.indexes = make([]Block, 0)
}

// 添加空闲空间到管理器
func (space *Space) AddBlock(index int, size uint) {
    return
    if size <= 0 {
        return
    }

    space.mu.Lock()
    defer space.mu.Unlock()

    indexpos, cmp := space.searchBlockByIndex(index)
    if cmp == 0 && indexpos >= 0 {
        // 已经添加过就不能再添加
        return
    }

    // 插入到indexes和blocks
    block         := Block{index, size}
    space.indexes  = space.insertBlock(space.indexes, block, indexpos, cmp)
    blockpos, cmp := space.searchBlockBySize(size)
    space.blocks   = space.insertBlock(space.blocks, block, blockpos, cmp)

    // 对indexes进行区块检查合并，对插入位置的前后项检查
    checkpos := indexpos - 1
    if checkpos < 0 {
        checkpos = 0
    }

    space.checkIndexMergeFromIndex(checkpos)
    space.checkIndexMergeFromIndex(checkpos + 1)

    //fmt.Println(space.indexes)
    //fmt.Println(space.blocks)
    //fmt.Println()
}

// 申请空间，返回文件地址及大小，返回成功后则在管理器中删除该空闲块
func (space *Space) GetBlock(size uint) (int, uint) {
    return -1, 0
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
        ix := space.blocks[mid].index
        sz := space.blocks[mid].size
        // indexes和blocks必须同时删除
        if indexpos := space.getIndexPositionByIndex(ix); indexpos != -1 {
            space.indexes = space.removeBlock(space.indexes, indexpos)
            space.blocks  = space.removeBlock(space.blocks, mid)
        }
        return ix, sz
    }
    return -1, 0
}


// 获得所有的碎片空间，按照index升序排序
func (space *Space) GetAllBlocksByIndex() []Block {
    return space.indexes
}

// 获得所有的碎片空间，按照size升序排序
func (space *Space) GetAllBlocksBySize() []Block {
    return space.blocks
}



