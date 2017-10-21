package gbtree

import (
    "fmt"
    "g/core/types/gset"
    "unsafe"
)

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
