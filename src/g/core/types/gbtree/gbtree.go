// 这是一颗改进的B树
package gbtree

import (
    "g/core/types/gset"
    "fmt"
    "unsafe"
    "math"
)

// B树对象
type Tree struct {
    max   int   // 最大数据项数
    root *Node  // 根节点数据块
}

// B树节点
type Node struct {
    tree   *Tree   // 所属B树
    parent *Node   // 父级数据节点
    items  []*Item // 数据项链表头，最小值
}

// B树数据项
type Item struct {
    key    []byte  // 关键字
    node   *Node   // 所属节点
    childl *Node   // 左孩子节点(所有左边的数据都比自身小)
    childr *Node   // 右孩子节点(所有右边的数据都比自身大，不存在相等，这颗树过滤掉了相等情况)
    data   *Data   // 数据指针
}

// B树数据信息
type Data struct {
    value []byte   // 数据
    start int64    // 数据文件开始位置
    end   int64    // 数据文件结束位置
}

// 创建一棵树
func New(m int) *Tree {
    // 构建一棵m阶树
    tree := &Tree{
        max : m,
    }
    // 初始化根节点
    tree.root = &Node {
        tree   : tree,
        parent : nil,
        items  : make([]*Item, 0),
    }
    return tree
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
    if len(node.items) == node.tree.max {
        mid  := int(math.Ceil(float64(len(node.items))/2)) - 1
        item := node.items[mid]
        if node.parent != nil {
            // 新增分裂节点
            noden := &Node {
                tree   : node.tree,
                parent : node.parent,
                items  : node.items[mid + 1:],
            }
            // 当前节点分裂
            node.items  = node.items[0 : mid]
            item.node   = node.parent
            item.childl = node
            item.childr = noden
            // 替换分列节点中的item的node为新node
            for _, v := range noden.items {
                v.node = noden
                if v.childl != nil {
                    v.childl.parent = noden
                }
                if v.childr != nil {
                    v.childr.parent = noden
                }
            }
            //fmt.Printf("split insert %v, childl %v, childr %v\n", item.key, item.childl.node.items[0].key, item.childr.node.items[0].key)
            node.parent.insertWithItem(item)

        } else {
            // root节点满了，从node中的中间节点进行拆分，创建两个新分支，中间节点向上提为root节点
            root  := &Node {
                tree   : node.tree,
                parent : nil,
                items  : []*Item{ item },
            }
            // 新增分裂节点
            noden := &Node {
                tree   : node.tree,
                parent : root,
                items  : node.items[mid + 1:],
            }
            // 设置根节点
            node.tree.root = root
            // 当前节点分裂
            node.items     = node.items[0 : mid]
            // 原root节点降级为普通节点
            node.parent    = root
            // 重构item的上下节点链接关系
            item.childl    = node
            item.childr    = noden
            // 提升的item与节点的关联关系
            item.node      = root
            // 替换分列节点中的item的node为新node
            //fmt.Printf("new node: %v\n", noden.items[0].key)
            for _, v := range noden.items {
                v.node = noden
                if v.childl != nil {
                    //fmt.Printf("%v, update childl node: %v\n", v.key, v.childl.node.items[0].key)
                    v.childl.parent = noden
                }
                if v.childr != nil {
                    //fmt.Printf("%v, update childr node: %v\n", v.key, v.childr.node.items[0].key)
                    v.childr.parent = noden
                }
            }
        }
    }
}

// 节点合并检查
func (node *Node) checkMerge() {
    min := int(math.Ceil(float64(len(node.items))/2)) - 1
    if len(node.items) < min {
        // 不满足节点的最小数据要求
    }
}

// 从根节点开始遍历树，返回顺序的节点列表
func (tree *Tree) walk() []*Node {
    m    := gset.NewStringSet()
    list := make([]*Node, 0)
    list  = append(list, tree.root)
    temp := list
    for {
        temp2 := make([]*Node, 0)
        for _, v := range temp {
            node := v
            //fmt.Printf("scan node %d, %x: \n", node.items[0].key, unsafe.Pointer(node))
            for _, item := range node.items {
                if item.childl != nil {
                    key := fmt.Sprintf("%x", unsafe.Pointer(item.childl))
                    //fmt.Printf("%d childl is %d, node %s, items[0] is %d\n", item.key, item.childl.key, key, item.childl.node.items[0].key)
                    if !m.Contains(key) {
                        temp2 = append(temp2, item.childl)
                        m.Add(key)
                    }
                } else {
                    //fmt.Printf("%d childl is nil\n", item.key)
                }
                if item.childr != nil {
                    key := fmt.Sprintf("%x", unsafe.Pointer(item.childr))
                    //fmt.Printf("%d childr is %d, node %s, items[0] is %d\n", item.key, item.childr.key, key, item.childr.node.items[0].key)
                    if !m.Contains(key) {
                        temp2 = append(temp2, item.childr)
                        m.Add(key)
                    }
                } else {
                    //fmt.Printf("%d childr is nil\n", item.key)
                }
            }
            //fmt.Println()
        }
        if len(temp2) > 0 {
            // 插入一个nil表示分层
            list = append(list, nil)
            list = append(list, temp2...)
            temp = temp2
        } else {
            break
        }
    }
    return list
}

// 打印节点信息（测试）
func (tree *Tree) Print() {
    list := tree.walk()
    for _, v := range list {
        if v == nil {
            fmt.Println()
            continue
        }
        fmt.Printf("[ ")
        for _, item := range v.items {
            fmt.Printf("%v ", item.key[0])
        }
        fmt.Printf("] ")
    }
}

// 往节点中写入数据
func (node *Node) insertRoundItem(key, value []byte, item *Item, index int, cmp int8) {
    itemn := &Item {
        key  : key,
        node : node,
        data : &Data {
            value: value,
        },
    }
    node.insertItemRoundItem(itemn, item, index, cmp)
}

// 插入一个自定义的item
func (node *Node) insertItemRoundItem (itemn *Item, item *Item, index int, cmp int8) {
    if item == nil {
        // 如果是第一条数据
        node.items = append(node.items, itemn)
    } else {
        // 插入数据
        i := index
        if cmp < 0 {
            if index > 0 {
                i -= 1
            }
        } else {
            i += 1
        }
        items     := node.items
        node.items = make([]*Item, 0)
        node.items = append(node.items, items[0:i]...)
        node.items = append(node.items, itemn)
        node.items = append(node.items, items[i: ]...)
    }
    node.checkSplit()
}

// 往节点插入一个带有关联关系item，
// 与insertRoundItem不同之处在于该方法支持自定义的item插入，该item一般是带关联关系，可以插入到任何节点中
func (node *Node) insertWithItem(itemn *Item) {
    item, index, cmp := node.search(itemn.key, false)
    //fmt.Printf("search %v, result %v\n", itemn.key, item.key)
    node.insertItemRoundItem(itemn, item, index, cmp)
}

// 删除数据项
func (node *Node) removeItem(index int) {
    items     := node.items
    node.items = make([]*Item, 0)
    node.items = append(node.items, items[0 : index]...)
    node.items = append(node.items, items[index + 1: ]...)
    node.checkMerge()
}

// 二分深度查找对应的数据项，返回匹配或者附近的item，以及该item在其node的索引位置，与key的比较结果
// deep参数用以控制是否需要进行深度查找，否则只在当前节点范围内水平查找
func (node *Node) search(key []byte, deep bool) (*Item, int, int8) {
    items := node.items
    for {
        min := 0
        max := len(items) - 1
        if min > max {
            break;
        }
        for {
            mid  := int((min + max) / 2)
            item := items[mid]
            cmp  := compareBytes(key, item.key)
            //fmt.Printf("%v VS %v: %d\n", key, item.key, cmp)
            if cmp < 0 {
                max = mid - 1
                // 深度查找
                if deep && min > max && item.childl != nil {
                    items = item.childl.items
                    break;
                }
            } else if cmp > 0 {
                min = mid + 1
                // 深度查找
                if deep && min > max && item.childr != nil {
                    items = item.childr.items
                    break;
                }
            } else {
                return item, mid, cmp
            }
            if min > max {
                return item, mid, cmp
            }
        }
    }
    return nil, -1, 0
}

// 往树中写入数据
func (tree *Tree) Set(key, value []byte) {
    item, index, cmp := tree.root.search(key, true)
    if index != -1 && cmp == 0 {
        item.data.value = value
    } else {
        if item == nil {
            tree.root.insertRoundItem(key, value, item, index, cmp)
        } else {
            item.node.insertRoundItem(key, value, item, index, cmp)
        }
    }
}

// 从树中查找数据
func (tree *Tree) Get(key []byte) []byte {
    if tree.root != nil {
        if item, index, cmp := tree.root.search(key, true); index != -1 && cmp == 0 {
            return item.data.value
        }
    }
    return nil
}

// 从树中删除数据
func (tree *Tree) Remove(key []byte) {
    if tree.root != nil {
        // 先进行查找，找到之后再进行删除
        if item, index, cmp := tree.root.search(key, true); index != -1 && cmp == 0 {
            item.node.removeItem(index)
        }
    }
}

