// 这是一颗改进的B树
// @todo 未开发完成，暂时不能使用
package gbtree

import (
    "math"
)

// B树对象
type Tree struct {
    min   int   // 最小数据项数
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
        min : int(math.Ceil(float64(m)/2)) - 1,
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


