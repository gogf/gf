package gfilespace

import "g/core/types/gbtree"

// 添加空闲空间到管理器
func (space *Space) AddBlock(index int, size uint) {
    if size <= 0 {
        return
    }
    block := &Block{index, size}

    space.mu.Lock()
    defer space.mu.Unlock()

    // 插入进全局树
    space.blocks.ReplaceOrInsert(block)

    // 插入进入索引表
    tree, ok := space.sizemap[block.size]
    if !ok {
        tree                      = gbtree.New(10)
        space.sizemap[block.size] = tree
    }
    tree.ReplaceOrInsert(block)

    // 插入空间块大小记录表
    space.sizetr.ReplaceOrInsert(gbtree.Int(block.size))

    // 对插入的数据进行合并检测
    space.checkMerge(block)
}

// 申请空间，返回文件地址及大小，返回成功后则在管理器中删除该空闲块
func (space *Space) GetBlock(size uint) (int, uint) {
    space.mu.RLock()
    defer space.mu.RUnlock()

    for {
        if tree, ok := space.sizemap[size]; ok {
            if r := tree.Min(); r != nil {
                block := r.(*Block)
                space.removeBlock(block)
                return block.index, block.size
            }
        }
        size = space.getNextBlockSize(size)
        if size == 0 {
            break
        }
    }
    return -1, 0
}


// 获得所有的碎片空间，按照index升序排序
func (space *Space) GetAllBlocks() []Block {
    index  := 0
    blocks := make([]Block, space.blocks.Len())
    space.blocks.Ascend(func(item gbtree.Item) bool {
        blocks[index] = *(item.(*Block))
        index++
        return true
    })
    return blocks
}

// 获得所有的碎片空间大小列表，按照size升序排序
func (space *Space) GetAllSizes() []uint {
    index := 0
    sizes := make([]uint, space.sizetr.Len())
    space.sizetr.Ascend(func(item gbtree.Item) bool {
        sizes[index] = uint(item.(gbtree.Int))
        index++
        return true
    })
    return sizes
}


