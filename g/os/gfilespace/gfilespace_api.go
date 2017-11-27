package gfilespace

import (

    "../../encoding/gbinary"
    "../../container/gbtree"
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

// 删除指定索引位置的空间块
func (space *Space) RemoveBlock(index int) {
    space.mu.Lock()
    defer space.mu.Unlock()

    space.removeBlock(&Block{index, 0})
}

// 给定的空间块*整块*是否包含在管理器中
func (space *Space) Contains(index int, size int) bool {
    block := &Block{index, size}
    if r := space.blocks.Get(block); r != nil {
        if r.(*Block).size >= size {
            return true
        }
    } else {
        pblock := space.getPrevBlock(block)
        if pblock != nil && (pblock.index <= index && (pblock.index + pblock.size) >= (index + size)) {
            return true
        }
    }
    return false
}

// 获取索引最小的空间块
func (space *Space) GetMinBlock() *Block {
    space.mu.RLock()
    defer space.mu.RUnlock()
    var block *Block
    space.blocks.Ascend(func(item gbtree.Item) bool {
        block = item.(*Block)
        return true
    })
    return block
}

// 获取索引最大的空间块
func (space *Space) GetMaxBlock() *Block {
    space.mu.RLock()
    defer space.mu.RUnlock()
    var block *Block
    space.blocks.Descend(func(item gbtree.Item) bool {
        block = item.(*Block)
        return true
    })
    return block
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

// 计算总的空闲空间大小
func (space *Space) SumSize() int {
    space.mu.RLock()
    defer space.mu.RUnlock()
    size := 0
    space.blocks.Ascend(func(item gbtree.Item) bool {
        size += item.(*Block).size
        return true
    })
    return size
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




