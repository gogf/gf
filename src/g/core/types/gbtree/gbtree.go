package gbtree

import (
    "fmt"
    "math"
    "g/core/types/gset"
    "unsafe"
)

// B树对象
type Tree struct {
    max  uint   // 最大数据项数
    root *Node  // 根节点数据块
}

// B树数据节点
type Node struct {
    level  uint  // 层级数，主要用于调试
    tree   *Tree // 所属B树
    parent *Node // 父级数据节点
    item   *Item // 数据项链表头，最小值
    size   uint  // 节点数据大小
}

// B树数据项(链表)
type Item struct {
    key    []byte // 关键字
    node   *Node  // 所属节点
    left   *Item  // 左指针
    right  *Item  // 右指针
    childl *Item  // 左孩子
    childr *Item  // 右孩子
    data   *Data  // 数据指针
}

// B树数据信息
type Data struct {
    value []byte // 数据
    start int64  // 数据文件开始位置
    end   int64  // 数据文件结束位置
}

// 创建一棵树
func New(m uint) *Tree {
    return &Tree{
        max : m,
    }
}

// 两个[]byte进行比较，v1 > v2 = 1, v1 < v2 = -1, v1 == v2 = 0
func compareBytes (v1, v2 []byte) int8 {
    i1 := len(v1) - 1
    i2 := len(v2) - 1
    for i := 0; i <= i1; i++ {
        if v1[i] < v2[i] {
            return -1
        }
        if v1[i] > v2[i] {
            return 1
        }
        if v1[i] == v2[i] {
            if i == i1 && i < i2 {
                return -1
            }
            if i == i2 && i < i1 {
                return 1
            }
            if i == i1 && i == i2 {
                return 0
            }
        }
    }
    return 0
}

// 节点分裂检查
func (node *Node) checkSplit() {
    if node.size >= node.tree.max {
        index  := 0
        middle := int(math.Ceil(float64(node.size)/2)) - 1
        item   := node.item
        for item != nil {
            if index == middle {
                break;
            }
            item = item.right
            index++
        }
        node.level++
        if node.parent != nil {
            // 分裂节点
            noden := &Node {
                level  : node.level,
                tree   : node.tree,
                parent : node.parent,
                item   : item.right,
                size   : node.size - uint(middle) - 1,
            }
            node.size       = uint(middle)
            // 普通节点
            item.left.right = nil
            item.right.left = nil
            node.parent.insertItem(&Item {
                key    : item.key,
                node   : node.parent,
                data   : item.data,
                childr : item.right,
                childl : item.left,
            })
            node.size--
            // 替换分列节点中的item的node为新node
            item := noden.item
            for item != nil {
                item.node = noden
                item = item.right
            }
        } else {
            // root节点满了，从node中的中间节点进行拆分，创建两个新分支，中间节点向上提为root节点
            root  := &Node {
                level  : 0,
                tree   : node.tree,
                parent : nil,
                item   : item,
                size   : 1,
            }

            //fmt.Println(string(item.key))
            // 分裂节点
            noden := &Node {
                level  : node.level,
                tree   : node.tree,
                parent : root,
                item   : item.right,
                size   : node.size - uint(middle) - 1,
            }
            node.size      = uint(middle)
            node.tree.root = root
            // 原root节点降级为普通节点
            node.parent     = root
            // 解除item的左右item链接关系
            item.left.right = nil
            item.right.left = nil
            // 重构item的上下链接关系(注意和上面分裂的区别)
            item.childl     = node.item
            item.childr     = noden.item
            // 解除item的左右链接关系
            item.left  = nil
            item.right = nil
            // 替换分列节点中的item的node为新node
            item := noden.item
            for item != nil {
                item.node = noden
                item = item.right
            }
        }
    }
}

// 节点合并检查
func (node *Node) checkMerge() {

}

// 从上往下查找关键字的位置，返回最后一次搜索的数据项以及查找是否成功
func (node *Node) search(key []byte) (*Item, bool) {
    item := node.item
    for item != nil {
        r := compareBytes(key, item.key)
        if r > 0 {
            // 判断是否往右查找，或者往右下级
            if item.right == nil {
                if item.childr == nil {
                    return item, false
                } else {
                    item = item.childr
                }
            } else {
                item = item.right
            }
        } else if r < 0 {
            // 判断是否往左下级查找
            if item.childl == nil {
                return item, false
            } else {
                item = item.childl
            }
        } else {
            return item, true
        }
    }

    return nil, false
}

// 打印节点信息（测试）
func (tree *Tree) Print() {
    m    := gset.NewStringSet()
    list := make([]*Node, 0)
    list  = append(list, tree.root)
    for len(list) > 0 {
        fmt.Printf("level - %d: ", list[0].level)
        count := 0
        for _, v := range list {
            count++
            fmt.Printf("[ ")
            item := v.item
            for item != nil {
                if item.childl != nil {
                    key := fmt.Sprintf("%x", unsafe.Pointer(item.childl.node))
                    if !m.Contains(key) {
                        list  = append(list, item.childl.node)
                    }
                }
                if item.childr != nil {
                    key := fmt.Sprintf("%x", unsafe.Pointer(item.childr.node))
                    if !m.Contains(key) {
                        list  = append(list, item.childr.node)
                    }
                }
                fmt.Print(string(item.key), " ")
                item = item.right
            }
            fmt.Printf("] ")
        }
        if len(list) > 0 {
            list = list[count:]
        }
        fmt.Println()
    }

}

// 往节点中写入数据
func (node *Node) insertRoundItem(key, value []byte, item *Item) {
    node.insertItemRoundItem(&Item {
        key  : key,
        node : node,
        data : &Data {
            value: value,
        },
    }, item)
}

// 插入一个item
func (node *Node) insertItemRoundItem(itemn, item *Item) {
    if item != nil {
        // 插入数据
        if compareBytes(itemn.key, item.key) < 0 {
            itemn.left      = item.left
            itemn.right     = item
            if item.left != nil {
                item.left.right = itemn
            }
            item.left       = itemn
        } else {
            itemn.left      = item
            itemn.right     = item.right
            if item.right != nil {
                item.right.left = itemn
            }
            item.right      = itemn
        }
    } else {
        // 如果是第一条数据
        node.item = itemn
    }
    node.size++
    node.checkSplit()
}

// 插入一个item
func (node *Node) insertItem(itemn *Item) {
    item := node.item
    for item != nil {
        if compareBytes(itemn.key, item.key) > 0 {
            item = item.right
        } else {
            break;
        }
    }
    node.insertItemRoundItem(itemn, item)
}

// 插入到节点中，不做层级判断
func (node *Node) insert(key, value []byte) {
    node.insertItem(&Item {
        key  : key,
        node : node,
        data : &Data {
            value: value,
        },
    })
}

// 往树中写入数据
func (tree *Tree) Set(key, value []byte) {
    if tree.root == nil {
        tree.root = &Node {
            tree   : tree,
            parent : nil,
            item   : nil,
        }
    }
    item, ok := tree.root.search(key);
    if ok {
        item.data.value = value
    } else {
        var node *Node
        if item == nil {
            node = tree.root
        } else {
            node = item.node
        }
        node.insertRoundItem(key, value, item)
    }
}

// 从树中查找数据
func (tree *Tree) Get(key []byte) []byte {
    if tree.root != nil {
        if item, ok := tree.root.search(key); ok {
            return item.data.value
        }
    }
    return nil
}
