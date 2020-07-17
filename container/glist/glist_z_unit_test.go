// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package glist

import (
	"container/list"
	"github.com/jin502437344/gf/internal/json"
	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gconv"

	"testing"
)

func checkListLen(t *gtest.T, l *List, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers(t *gtest.T, l *List, es []*Element) {
	if !checkListLen(t, l, len(es)) {
		return
	}
	l.RLockFunc(func(list *list.List) {
		for i, e := 0, l.list.Front(); i < list.Len(); i, e = i+1, e.Next() {
			if e.Prev() != es[i].Prev() {
				t.Errorf("list[%d].Prev = %p, want %p", i, e.Prev(), es[i].Prev())
			}
			if e.Next() != es[i].Next() {
				t.Errorf("list[%d].Next = %p, want %p", i, e.Next(), es[i].Next())
			}
		}
	})
}

func TestVar(t *testing.T) {
	var l List
	l.PushFront(1)
	l.PushFront(2)
	if v := l.PopBack(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopBack(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopBack(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		//fmt.Println(v)
	}
	l.PushBack(1)
	l.PushBack(2)
	if v := l.PopFront(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopFront(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopFront(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		//fmt.Println(v)
	}
}

func TestBasic(t *testing.T) {
	l := New()
	l.PushFront(1)
	l.PushFront(2)
	if v := l.PopBack(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopBack(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopBack(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		//fmt.Println(v)
	}
	l.PushBack(1)
	l.PushBack(2)
	if v := l.PopFront(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopFront(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		//fmt.Println(v)
	}
	if v := l.PopFront(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		//fmt.Println(v)
	}
}

func TestList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		checkListPointers(t, l, []*Element{})

		// Single element list
		e := l.PushFront("a")
		checkListPointers(t, l, []*Element{e})
		l.MoveToFront(e)
		checkListPointers(t, l, []*Element{e})
		l.MoveToBack(e)
		checkListPointers(t, l, []*Element{e})
		l.Remove(e)
		checkListPointers(t, l, []*Element{})

		// Bigger list
		e2 := l.PushFront(2)
		e1 := l.PushFront(1)
		e3 := l.PushBack(3)
		e4 := l.PushBack("banana")
		checkListPointers(t, l, []*Element{e1, e2, e3, e4})

		l.Remove(e2)
		checkListPointers(t, l, []*Element{e1, e3, e4})

		l.MoveToFront(e3) // move from middle
		checkListPointers(t, l, []*Element{e3, e1, e4})

		l.MoveToFront(e1)
		l.MoveToBack(e3) // move from middle
		checkListPointers(t, l, []*Element{e1, e4, e3})

		l.MoveToFront(e3) // move from back
		checkListPointers(t, l, []*Element{e3, e1, e4})
		l.MoveToFront(e3) // should be no-op
		checkListPointers(t, l, []*Element{e3, e1, e4})

		l.MoveToBack(e3) // move from front
		checkListPointers(t, l, []*Element{e1, e4, e3})
		l.MoveToBack(e3) // should be no-op
		checkListPointers(t, l, []*Element{e1, e4, e3})

		e2 = l.InsertBefore(e1, 2) // insert before front
		checkListPointers(t, l, []*Element{e2, e1, e4, e3})
		l.Remove(e2)
		e2 = l.InsertBefore(e4, 2) // insert before middle
		checkListPointers(t, l, []*Element{e1, e2, e4, e3})
		l.Remove(e2)
		e2 = l.InsertBefore(e3, 2) // insert before back
		checkListPointers(t, l, []*Element{e1, e4, e2, e3})
		l.Remove(e2)

		e2 = l.InsertAfter(e1, 2) // insert after front
		checkListPointers(t, l, []*Element{e1, e2, e4, e3})
		l.Remove(e2)
		e2 = l.InsertAfter(e4, 2) // insert after middle
		checkListPointers(t, l, []*Element{e1, e4, e2, e3})
		l.Remove(e2)
		e2 = l.InsertAfter(e3, 2) // insert after back
		checkListPointers(t, l, []*Element{e1, e4, e3, e2})
		l.Remove(e2)

		// Check standard iteration.
		sum := 0
		for e := l.Front(); e != nil; e = e.Next() {
			if i, ok := e.Value.(int); ok {
				sum += i
			}
		}
		if sum != 4 {
			t.Errorf("sum over l = %d, want 4", sum)
		}

		// Clear all elements by iterating
		var next *Element
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			l.Remove(e)
		}
		checkListPointers(t, l, []*Element{})
	})
}

func checkList(t *gtest.T, l *List, es []interface{}) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {

		switch e.Value.(type) {
		case int:
			if le := e.Value.(int); le != es[i] {
				t.Errorf("elt[%d].Value() = %v, want %v", i, le, es[i])
			}
			// default string
		default:
			if le := e.Value.(string); le != es[i] {
				t.Errorf("elt[%v].Value() = %v, want %v", i, le, es[i])
			}
		}

		i++
	}

	//for e := l.Front(); e != nil; e = e.Next() {
	//	le := e.Value.(int)
	//	if le != es[i] {
	//		t.Errorf("elt[%d].Value() = %v, want %v", i, le, es[i])
	//	}
	//	i++
	//}
}

func TestExtending(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l1 := New()
		l2 := New()

		l1.PushBack(1)
		l1.PushBack(2)
		l1.PushBack(3)

		l2.PushBack(4)
		l2.PushBack(5)

		l3 := New()
		l3.PushBackList(l1)
		checkList(t, l3, []interface{}{1, 2, 3})
		l3.PushBackList(l2)
		checkList(t, l3, []interface{}{1, 2, 3, 4, 5})

		l3 = New()
		l3.PushFrontList(l2)
		checkList(t, l3, []interface{}{4, 5})
		l3.PushFrontList(l1)
		checkList(t, l3, []interface{}{1, 2, 3, 4, 5})

		checkList(t, l1, []interface{}{1, 2, 3})
		checkList(t, l2, []interface{}{4, 5})

		l3 = New()
		l3.PushBackList(l1)
		checkList(t, l3, []interface{}{1, 2, 3})
		l3.PushBackList(l3)
		checkList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

		l3 = New()
		l3.PushFrontList(l1)
		checkList(t, l3, []interface{}{1, 2, 3})
		l3.PushFrontList(l3)
		checkList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

		l3 = New()
		l1.PushBackList(l3)
		checkList(t, l1, []interface{}{1, 2, 3})
		l1.PushFrontList(l3)
		checkList(t, l1, []interface{}{1, 2, 3})
	})
}

func TestRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		e1 := l.PushBack(1)
		e2 := l.PushBack(2)
		checkListPointers(t, l, []*Element{e1, e2})
		//e := l.Front()
		//l.Remove(e)
		//checkListPointers(t, l, []*Element{e2})
		//l.Remove(e)
		//checkListPointers(t, l, []*Element{e2})
	})
}

func TestIssue4103(t *testing.T) {
	l1 := New()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(e, 8)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e)
	if e.Value != 1 {
		t.Errorf("e.value = %d, want 1", e.Value)
	}
	//if e.Next() != nil {
	//    t.Errorf("e.Next() != nil")
	//}
	//if e.Prev() != nil {
	//    t.Errorf("e.Prev() != nil")
	//}
}

func TestMove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		e1 := l.PushBack(1)
		e2 := l.PushBack(2)
		e3 := l.PushBack(3)
		e4 := l.PushBack(4)

		l.MoveAfter(e3, e3)
		checkListPointers(t, l, []*Element{e1, e2, e3, e4})
		l.MoveBefore(e2, e2)
		checkListPointers(t, l, []*Element{e1, e2, e3, e4})

		l.MoveAfter(e3, e2)
		checkListPointers(t, l, []*Element{e1, e2, e3, e4})
		l.MoveBefore(e2, e3)
		checkListPointers(t, l, []*Element{e1, e2, e3, e4})

		l.MoveBefore(e2, e4)
		checkListPointers(t, l, []*Element{e1, e3, e2, e4})
		e2, e3 = e3, e2

		l.MoveBefore(e4, e1)
		checkListPointers(t, l, []*Element{e4, e1, e2, e3})
		e1, e2, e3, e4 = e4, e1, e2, e3

		l.MoveAfter(e4, e1)
		checkListPointers(t, l, []*Element{e1, e4, e2, e3})
		e2, e3, e4 = e4, e2, e3

		l.MoveAfter(e2, e3)
		checkListPointers(t, l, []*Element{e1, e3, e2, e4})
		e2, e3 = e3, e2
	})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
func TestZeroList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var l1 = New()
		l1.PushFront(1)
		checkList(t, l1, []interface{}{1})

		var l2 = New()
		l2.PushBack(1)
		checkList(t, l2, []interface{}{1})

		var l3 = New()
		l3.PushFrontList(l1)
		checkList(t, l3, []interface{}{1})

		var l4 = New()
		l4.PushBackList(l2)
		checkList(t, l4, []interface{}{1})
	})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.InsertBefore(new(Element), 1)
		checkList(t, l, []interface{}{1, 2, 3})
	})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.InsertAfter(new(Element), 1)
		checkList(t, l, []interface{}{1, 2, 3})
	})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l1 := New()
		e1 := l1.PushBack(1)

		l2 := New()
		e2 := l2.PushBack(2)

		l1.MoveAfter(e1, e2)
		checkList(t, l1, []interface{}{1})
		checkList(t, l2, []interface{}{2})

		l1.MoveBefore(e1, e2)
		checkList(t, l1, []interface{}{1})
		checkList(t, l2, []interface{}{2})
	})
}

func TestList_RemoveAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		l.PushBack(1)
		l.RemoveAll()
		checkList(t, l, []interface{}{})
		l.PushBack(2)
		checkList(t, l, []interface{}{2})
	})
}

func TestList_PushFronts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2}
		l.PushFronts(a1)
		checkList(t, l, []interface{}{2, 1})
		a1 = []interface{}{3, 4, 5}
		l.PushFronts(a1)
		checkList(t, l, []interface{}{5, 4, 3, 2, 1})
	})
}

func TestList_PushBacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2}
		l.PushBacks(a1)
		checkList(t, l, []interface{}{1, 2})
		a1 = []interface{}{3, 4, 5}
		l.PushBacks(a1)
		checkList(t, l, []interface{}{1, 2, 3, 4, 5})
	})
}

func TestList_PopBacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		a2 := []interface{}{"a", "c", "b", "e"}
		l.PushFronts(a1)
		i1 := l.PopBacks(2)
		t.Assert(i1, []interface{}{1, 2})

		l.PushBacks(a2) //4.3,a,c,b,e
		i1 = l.PopBacks(3)
		t.Assert(i1, []interface{}{"e", "b", "c"})
	})
}

func TestList_PopFronts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopFronts(2)
		t.Assert(i1, []interface{}{4, 3})
		t.Assert(l.Len(), 2)
	})
}

func TestList_PopBackAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopBackAll()
		t.Assert(i1, []interface{}{1, 2, 3, 4})
		t.Assert(l.Len(), 0)
	})
}

func TestList_PopFrontAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopFrontAll()
		t.Assert(i1, []interface{}{4, 3, 2, 1})
		t.Assert(l.Len(), 0)
	})
}

func TestList_FrontAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.FrontAll()
		t.Assert(i1, []interface{}{4, 3, 2, 1})
		t.Assert(l.Len(), 4)
	})
}

func TestList_BackAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.BackAll()
		t.Assert(i1, []interface{}{1, 2, 3, 4})
		t.Assert(l.Len(), 4)
	})
}

func TestList_FrontValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		l2 := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.FrontValue()
		t.Assert(gconv.Int(i1), 4)
		t.Assert(l.Len(), 4)

		i1 = l2.FrontValue()
		t.Assert(i1, nil)
	})
}

func TestList_BackValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		l2 := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.BackValue()
		t.Assert(gconv.Int(i1), 1)
		t.Assert(l.Len(), 4)

		i1 = l2.FrontValue()
		t.Assert(i1, nil)
	})
}

func TestList_Back(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		t.Assert(e1.Value, 1)
		t.Assert(l.Len(), 4)
	})
}

func TestList_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		t.Assert(l.Size(), 4)
		l.PopFront()
		t.Assert(l.Size(), 3)
	})
}

func TestList_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		l.Removes([]*Element{e1})
		t.Assert(l.Len(), 3)

		e2 := l.Back()
		l.Removes([]*Element{e2})
		t.Assert(l.Len(), 2)
		checkList(t, l, []interface{}{4, 3})
	})
}

func TestList_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

		t.Assert(l.PopBack(), 9)
		t.Assert(l.PopBacks(2), []interface{}{8, 7})
		t.Assert(l.PopFront(), 1)
		t.Assert(l.PopFronts(2), []interface{}{2, 3})
	})
}

func TestList_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		l.Clear()
		t.Assert(l.Len(), 0)
	})
}

func TestList_IteratorAsc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 5, 6, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e *Element) bool {
			if gconv.Int(e1.Value) > 2 {
				return true
			}
			return false
		}
		checkList(t, l, []interface{}{4, 3, 6, 5, 2, 1})
		l.IteratorAsc(fun1)
		checkList(t, l, []interface{}{4, 3, 6, 5, 2, 1})
	})
}

func TestList_IteratorDesc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e *Element) bool {
			if gconv.Int(e1.Value) > 6 {
				return true
			}
			return false
		}
		l.IteratorDesc(fun1)
		t.Assert(l.Len(), 4)
		checkList(t, l, []interface{}{4, 3, 2, 1})
	})
}

func TestList_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		a1 := []interface{}{"a", "b", "c", "d", "e"}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e *Element) bool {
			if gconv.String(e1.Value) > "c" {
				return true
			}
			return false
		}
		checkList(t, l, []interface{}{"e", "d", "c", "b", "a"})
		l.Iterator(fun1)
		checkList(t, l, []interface{}{"e", "d", "c", "b", "a"})
	})
}

func TestList_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewFrom([]interface{}{1, 2, "a", `"b"`, `\c`})
		t.Assert(l.Join(","), `1,2,a,"b",\c`)
		t.Assert(l.Join("."), `1.2.a."b".\c`)
	})
}

func TestList_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewFrom([]interface{}{1, 2, "a", `"b"`, `\c`})
		t.Assert(l.String(), `[1,2,a,"b",\c]`)
	})
}

func TestList_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		a := []interface{}{"a", "b", "c"}
		l := New()
		l.PushBacks(a)
		b1, err1 := json.Marshal(l)
		b2, err2 := json.Marshal(a)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		a := []interface{}{"a", "b", "c"}
		l := New()
		b, err := json.Marshal(a)
		t.Assert(err, nil)

		err = json.Unmarshal(b, l)
		t.Assert(err, nil)
		t.Assert(l.FrontAll(), a)
	})
	gtest.C(t, func(t *gtest.T) {
		var l List
		a := []interface{}{"a", "b", "c"}
		b, err := json.Marshal(a)
		t.Assert(err, nil)

		err = json.Unmarshal(b, &l)
		t.Assert(err, nil)
		t.Assert(l.FrontAll(), a)
	})
}

func TestList_UnmarshalValue(t *testing.T) {
	type TList struct {
		Name string
		List *List
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var tlist *TList
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"list": []byte(`[1,2,3]`),
		}, &tlist)
		t.Assert(err, nil)
		t.Assert(tlist.Name, "john")
		t.Assert(tlist.List.FrontAll(), []interface{}{1, 2, 3})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var tlist *TList
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"list": []interface{}{1, 2, 3},
		}, &tlist)
		t.Assert(err, nil)
		t.Assert(tlist.Name, "john")
		t.Assert(tlist.List.FrontAll(), []interface{}{1, 2, 3})
	})
}
