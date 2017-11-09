package gfilespace

import (
    "g/core/types/gbtree"
)

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
    space.insertIntoSizeMap(block)

    // 对插入的数据进行合并检测
    space.checkMerge(block)
}

// 申请空间，返回文件地址及大小，返回成功后则在管理器中删除该空闲块
func (space *Space) GetBlock(size uint) (int, uint) {
    if size <= 0 {
        return -1, 0
    }
    space.mu.Lock()
    defer space.mu.Unlock()

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
    space.mu.RLock()
    defer space.mu.RUnlock()
    blocks := make([]Block, 0)
    space.blocks.Ascend(func(item gbtree.Item) bool {
        blocks = append(blocks, *(item.(*Block)))
        return true
    })
    return blocks
}

// 获得所有的碎片空间大小列表，按照size升序排序
func (space *Space) GetAllSizes() []uint {
    space.mu.RLock()
    defer space.mu.RUnlock()
    sizes := make([]uint, 0)
    space.sizetr.Ascend(func(item gbtree.Item) bool {
        sizes = append(sizes, uint(item.(gbtree.Int)))
        return true
    })
    return sizes
}

// 获取当前空间管理器中最大的空闲块大小
func (space *Space) GetMaxSize() uint {
    space.mu.RLock()
    defer space.mu.RUnlock()

    if item := space.sizetr.Max(); item != nil {
        return uint(item.(gbtree.Int))
    }
    return 0
}



