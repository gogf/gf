package gbtree

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

