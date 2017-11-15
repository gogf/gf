package gfilespace

import (
    "g/core/types/gbtree"
    "g/encoding/gbinary"
)

// 添加空闲空间到管理器
func (space *Space) AddBlock(index int, size int) {
    if size <= 0 {
        return
    }
    space.mu.Lock()
    defer space.mu.Unlock()

    space.addBlock(index, size)
}

// 申请空间，返回文件地址及大小，返回成功后则在管理器中删除该空闲块
func (space *Space) GetBlock(size int) (int, int) {
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
func (space *Space) GetMaxSize() int {
    space.mu.RLock()
    defer space.mu.RUnlock()

    if item := space.sizetr.Max(); item != nil {
        return int(item.(gbtree.Int))
    }
    return 0
}

// 获取空间块的数量
func (space *Space) Len() int {
    space.mu.RLock()
    defer space.mu.RUnlock()

    return space.blocks.Len()
}

// 导出空间块数据
func (space *Space) Export() []byte {
    space.mu.RLock()
    defer space.mu.RUnlock()

    content := make([]byte, 0)
    space.blocks.Ascend(func(item gbtree.Item) bool {
        block   := item.(*Block)
        content  = append(content, gbinary.EncodeInt64(int64(block.Index()))...)
        content  = append(content, gbinary.EncodeInt32(int32(block.Size()))...)
        return true
    })

    return content
}

// 导入空间块数据
func (space *Space) Import(content []byte) {
    space.mu.Lock()
    defer space.mu.Unlock()

    for i := 0; i < len(content); i += 12 {
        space.addBlock(
            int(gbinary.DecodeToInt64(content[i : i + 8])),
            int(gbinary.DecodeToInt32(content[i + 8 : i + 12])),
        )
    }
}




