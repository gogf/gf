// @todo 未开发完成，暂时不能使用
package gfilespace

import (
    "g/os/gfilepool"
    "g/os/gfile"
    "errors"
    "os"
    "sync"
    "fmt"
)

// 文件空间管理结构体
type Space struct {
    fp     *gfilepool.Pool
    lock   sync.RWMutex
    blocks []Block
}

// 文件空闲块
type Block struct {
    index int64
    size  uint32
}

// 创建一个空间管理器，基于给定的文件
func New(path string) (*Space, error) {
    if gfile.Exists(path) && !gfile.IsReadable(path) {
        return nil, errors.New("permission denied to file: " + path)
    }
    space := &Space{
        fp     : gfilepool.New(path, os.O_RDONLY|os.O_CREATE, 60),
        blocks : make([]Block, 0),
    }
    return space, nil
}

// 申请空间，返回文件地址及大小，返回成功后则在管理器中删除该空闲块
func (space *Space) GetBlock(size uint32) (int64, uint32) {
    space.lock.RLock()
    defer space.lock.RUnlock()

    mid, cmp := space.searchBlock(size)
    if cmp == -1 {
        for i := mid + 1; i < len(space.blocks); i++ {
            if space.blocks[i].size >= size {
                mid = i
                cmp = 1
                break;
            }
        }
    }
    if cmp >= 0 {
        index       := space.blocks[mid].index
        size        := space.blocks[mid].size
        blocks      := space.blocks
        space.blocks = blocks[0 : mid]
        space.blocks = append(space.blocks, blocks[mid : ]...)
        return index, size
    } else {
        if pf, err := space.fp.File(); err == nil {
            defer pf.Close()
            // 需要保证同一时间只有1个文件指针在写文件
            if pos, err := pf.File().Seek(0, 2); err != nil {
                return pos, 0
            }
        }

    }
    return -1, 0
}

// 添加空闲空间到管理器
func (space *Space) AddBlock(index int64, size uint32) {
    space.lock.Lock()
    defer space.lock.Unlock()
    fmt.Println(space.blocks)
    block    := Block{index, size}
    mid, cmp := space.searchBlock(size)
    fmt.Println(mid)
    fmt.Println(cmp)
    switch cmp {
        case 0:
        case 1:
        case -2:
            // 添加到mid后面
            length := mid + 1
            if length > len(space.blocks) {
                length = len(space.blocks)
            }
            blocks      := space.blocks
            space.blocks = blocks[0 : length]
            space.blocks = append(space.blocks, block)
            space.blocks = append(space.blocks, blocks[length : ]...)
        case -1:
            // 添加到前面
            blocks      := space.blocks
            space.blocks = blocks[0 : mid]
            space.blocks = append(space.blocks, block)
            space.blocks = append(space.blocks, blocks[mid : ]...)
    }
}

// 搜索空闲空间，返回空间与size相似(=, <, >)的空闲块索引地址
// 找不到则返回-2
func (space *Space) searchBlock(size uint32) (int, int) {
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
            if cmp == 0 || min > max {
                break
            }
        }
    }
    return mid, cmp
}

