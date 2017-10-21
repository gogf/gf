package gbtree


import "math"

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
