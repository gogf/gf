// 文件空间管理， 可用于文件碎片空间维护及再利用，支持自动合并连续碎片空间

package gfilespace

import (
    "sync"
    "g/core/types/gbtree"
)

// 文件空间管理结构体
type Space struct {
    mu      sync.RWMutex           // 并发操作锁
    blocks  *gbtree.BTree          // 所有的空间块构建的B+树
    sizetr  *gbtree.BTree          // 空间块大小构建的B+树
    sizemap map[uint]*gbtree.BTree // 按照空间块大小构建的索引哈希表，便于检索，每个表项是一个B+树
}

// 文件空闲块
type Block struct {
    index int  // 文件偏移量
    size  uint // 区块大小(byte)
}

// 用于B+树的接口具体实现定义
func (block *Block) Less(item gbtree.Item) bool {
    if block.index < item.(*Block).index {
        return true
    }
    return false
}

// 创建一个空间管理器
func New() *Space {
    return &Space {
        blocks  : gbtree.New(10),
        sizetr  : gbtree.New(5),
        sizemap : make(map[uint]*gbtree.BTree),
    }
}

// 获取指定block的前一项block
func (space *Space) getPrevBlock(block *Block) *Block {
    var pblock *Block = nil
    space.blocks.DescendLessOrEqual(block, func(item gbtree.Item) bool {
        if item.(*Block).index != block.index {
            pblock = item.(*Block)
            return false
        }
        return true
    })
    return pblock
}

// 获取指定block的后一项block
func (space *Space) getNextBlock(block *Block) *Block {
    var nblock *Block = nil
    space.blocks.DescendGreaterThan(block, func(item gbtree.Item) bool {
        nblock = item.(*Block)
        return false
    })
    return nblock
}

// 获取指定block的前一项block size
func (space *Space) getPrevBlockSize(size uint) uint {
    psize := uint(0)
    space.sizetr.DescendLessOrEqual(gbtree.Int(size), func(item gbtree.Item) bool {
        if uint(item.(gbtree.Int)) != size {
            psize = uint(item.(gbtree.Int))
            return false
        }
        return true
    })
    return psize
}

// 获取指定block的后一项block size
func (space *Space) getNextBlockSize(size uint) uint {
    nsize := uint(0)
    space.sizetr.DescendGreaterThan(gbtree.Int(size), func(item gbtree.Item) bool {
        nsize = uint(item.(gbtree.Int))
        return false
    })
    return nsize
}

// 内部按照索引检查合并
func (space *Space) checkMerge(block *Block) {
    // 首先检查插入空间块的前一项往后是否可以合并，如果当前合并失败后，才会判断当前插入项和后续的空间块合并
    if b := space.checkMergeOfTwoBlock(space.getPrevBlock(block), block); b.index == block.index {
        // 其次检查插入空间块的当前项往后是否可以合并
        space.checkMergeOfTwoBlock(block, space.getNextBlock(block))
    }
}

// 连续检测两个空间块的合并，返回最后一个无法合并的空间块指针
func (space *Space) checkMergeOfTwoBlock(pblock, block *Block) *Block {
    if pblock == nil {
        return block
    }
    if block == nil {
        return pblock
    }
    for {
        if pblock.index + int(pblock.size) >= block.index {
            pblock.size = uint(block.index + int(block.size) - pblock.index)
            space.removeBlock(block)
            block = space.getNextBlock(pblock)
            if block == nil {
                return pblock
            }
        } else {
            break
        }
    }
    return block
}


// 删除一项
func (space *Space) removeBlock(block *Block) {
    space.blocks.Delete(block)
    if tree, ok := space.sizemap[block.size]; ok {
        tree.Delete(block)
    }
}

// 获得碎片偏移量
func (block *Block) Index() int {
    return block.index
}

// 获得碎片大小
func (block *Block) Size() uint {
    return block.size
}


