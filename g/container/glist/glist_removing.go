// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://gitee.com/johng/gf.

package glist

// 删除数据项e, 并返回删除项的元素项
func (l *List) Remove(e *Element) *Element {
    if e.list != l {
        return nil
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.doRemove(e)
}

// 删除所有数据项
func (l *List) RemoveAll() {
    l.length.Set(0)
    l.mu.Lock()
    l.root.mu.Lock()
    l.root.prev = l.root
    l.root.next = l.root
    l.root.mu.Unlock()
    l.mu.Unlock()
}

func (l *List) doRemove(e *Element) *Element {
    e.mu.RLock()
    if e.prev.getNext() == e {
        e.prev.setNext(e.next)
        e.next.setPrev(e.prev)
    } else {
        e.mu.RUnlock()
        return e
    }
    e.mu.RUnlock()
    l.length.Add(-1)
    return e
}