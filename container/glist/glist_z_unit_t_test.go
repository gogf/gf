// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glist

import (
	"container/list"
	"testing"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func checkTListLen(t *gtest.T, l *TList[any], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkTListPointers(t *gtest.T, l *TList[any], es []*TElement[any]) {
	if !checkTListLen(t, l, len(es)) {
		return
	}

	i := 0
	l.Iterator(func(e *TElement[any]) bool {
		if e.Prev() != es[i].Prev() {
			t.Errorf("list[%d].Prev = %p, want %p", i, e.Prev(), es[i].Prev())
			return false
		}
		if e.Next() != es[i].Next() {
			t.Errorf("list[%d].Next = %p, want %p", i, e.Next(), es[i].Next())
			return false
		}
		i++
		return true
	})
}

func TestTVar(t *testing.T) {
	var l TList[any]
	l.PushFront(1)
	l.PushFront(2)
	if v := l.PopBack(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopBack(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopBack(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
	l.PushBack(1)
	l.PushBack(2)
	if v := l.PopFront(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopFront(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopFront(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
}

func TestTBasic(t *testing.T) {
	l := NewT[any]()
	l.PushFront(1)
	l.PushFront(2)
	if v := l.PopBack(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopBack(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopBack(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
	l.PushBack(1)
	l.PushBack(2)
	if v := l.PopFront(); v != 1 {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopFront(); v != 2 {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v := l.PopFront(); v != nil {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
}

func TestTList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		checkTListPointers(t, l, []*TElement[any]{})

		// Single element list
		e := l.PushFront("a")
		checkTListPointers(t, l, []*TElement[any]{e})
		l.MoveToFront(e)
		checkTListPointers(t, l, []*TElement[any]{e})
		l.MoveToBack(e)
		checkTListPointers(t, l, []*TElement[any]{e})
		l.Remove(e)
		checkTListPointers(t, l, []*TElement[any]{})

		// Bigger list
		e2 := l.PushFront(2)
		e1 := l.PushFront(1)
		e3 := l.PushBack(3)
		e4 := l.PushBack("banana")
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e3, e4})

		l.Remove(e2)
		checkTListPointers(t, l, []*TElement[any]{e1, e3, e4})

		l.MoveToFront(e3) // move from middle
		checkTListPointers(t, l, []*TElement[any]{e3, e1, e4})

		l.MoveToFront(e1)
		l.MoveToBack(e3) // move from middle
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e3})

		l.MoveToFront(e3) // move from back
		checkTListPointers(t, l, []*TElement[any]{e3, e1, e4})
		l.MoveToFront(e3) // should be no-op
		checkTListPointers(t, l, []*TElement[any]{e3, e1, e4})

		l.MoveToBack(e3) // move from front
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e3})
		l.MoveToBack(e3) // should be no-op
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e3})

		e2 = l.InsertBefore(e1, 2) // insert before front
		checkTListPointers(t, l, []*TElement[any]{e2, e1, e4, e3})
		l.Remove(e2)
		e2 = l.InsertBefore(e4, 2) // insert before middle
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e4, e3})
		l.Remove(e2)
		e2 = l.InsertBefore(e3, 2) // insert before back
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e2, e3})
		l.Remove(e2)

		e2 = l.InsertAfter(e1, 2) // insert after front
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e4, e3})
		l.Remove(e2)
		e2 = l.InsertAfter(e4, 2) // insert after middle
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e2, e3})
		l.Remove(e2)
		e2 = l.InsertAfter(e3, 2) // insert after back
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e3, e2})
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
		var next *TElement[any]
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			l.Remove(e)
		}
		checkTListPointers(t, l, []*TElement[any]{})
	})
}

func checkTList(t *gtest.T, l *TList[any], es []any) {
	if !checkTListLen(t, l, len(es)) {
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

	// for e := l.Front(); e != nil; e = e.Next() {
	//	le := e.Value.(int)
	//	if le != es[i] {
	//		t.Errorf("elt[%d].Value() = %v, want %v", i, le, es[i])
	//	}
	//	i++
	// }
}

func TestTExtending(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l1 := NewT[any]()
		l2 := NewT[any]()

		l1.PushBack(1)
		l1.PushBack(2)
		l1.PushBack(3)

		l2.PushBack(4)
		l2.PushBack(5)

		l3 := NewT[any]()
		l3.PushBackList(l1)
		checkTList(t, l3, []any{1, 2, 3})
		l3.PushBackList(l2)
		checkTList(t, l3, []any{1, 2, 3, 4, 5})

		l3 = NewT[any]()
		l3.PushFrontList(l2)
		checkTList(t, l3, []any{4, 5})
		l3.PushFrontList(l1)
		checkTList(t, l3, []any{1, 2, 3, 4, 5})

		checkTList(t, l1, []any{1, 2, 3})
		checkTList(t, l2, []any{4, 5})

		l3 = NewT[any]()
		l3.PushBackList(l1)
		checkTList(t, l3, []any{1, 2, 3})
		l3.PushBackList(l3)
		checkTList(t, l3, []any{1, 2, 3, 1, 2, 3})

		l3 = NewT[any]()
		l3.PushFrontList(l1)
		checkTList(t, l3, []any{1, 2, 3})
		l3.PushFrontList(l3)
		checkTList(t, l3, []any{1, 2, 3, 1, 2, 3})

		l3 = NewT[any]()
		l1.PushBackList(l3)
		checkTList(t, l1, []any{1, 2, 3})
		l1.PushFrontList(l3)
		checkTList(t, l1, []any{1, 2, 3})
	})
}

func TestTRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		e1 := l.PushBack(1)
		e2 := l.PushBack(2)
		checkTListPointers(t, l, []*TElement[any]{e1, e2})
		// e := l.Front()
		// l.Remove(e)
		// checkTListPointers(t, l, []*TElement[any]{e2})
		// l.Remove(e)
		// checkTListPointers(t, l, []*TElement[any]{e2})
	})
}

func Test_T_Issue4103(t *testing.T) {
	l1 := NewT[any]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := NewT[any]()
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

func Test_T_Issue6349(t *testing.T) {
	l := NewT[any]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e)
	if e.Value != 1 {
		t.Errorf("e.value = %d, want 1", e.Value)
	}
	// if e.Next() != nil {
	//    t.Errorf("e.Next() != nil")
	// }
	// if e.Prev() != nil {
	//    t.Errorf("e.Prev() != nil")
	// }
}

func TestTMove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		e1 := l.PushBack(1)
		e2 := l.PushBack(2)
		e3 := l.PushBack(3)
		e4 := l.PushBack(4)

		l.MoveAfter(e3, e3)
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e3, e4})
		l.MoveBefore(e2, e2)
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e3, e4})

		l.MoveAfter(e3, e2)
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e3, e4})
		l.MoveBefore(e2, e3)
		checkTListPointers(t, l, []*TElement[any]{e1, e2, e3, e4})

		l.MoveBefore(e2, e4)
		checkTListPointers(t, l, []*TElement[any]{e1, e3, e2, e4})
		e2, e3 = e3, e2

		l.MoveBefore(e4, e1)
		checkTListPointers(t, l, []*TElement[any]{e4, e1, e2, e3})
		e1, e2, e3, e4 = e4, e1, e2, e3

		l.MoveAfter(e4, e1)
		checkTListPointers(t, l, []*TElement[any]{e1, e4, e2, e3})
		e2, e3, e4 = e4, e2, e3

		l.MoveAfter(e2, e3)
		checkTListPointers(t, l, []*TElement[any]{e1, e3, e2, e4})
		e2, e3 = e3, e2
	})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
func TestTZeroList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var l1 = NewT[any]()
		l1.PushFront(1)
		checkTList(t, l1, []any{1})

		var l2 = NewT[any]()
		l2.PushBack(1)
		checkTList(t, l2, []any{1})

		var l3 = NewT[any]()
		l3.PushFrontList(l1)
		checkTList(t, l3, []any{1})

		var l4 = NewT[any]()
		l4.PushBackList(l2)
		checkTList(t, l4, []any{1})
	})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestTInsertBeforeUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.InsertBefore(new(TElement[any]), 1)
		checkTList(t, l, []any{1, 2, 3})
	})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestTInsertAfterUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.InsertAfter(new(TElement[any]), 1)
		checkTList(t, l, []any{1, 2, 3})
	})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestTMoveUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l1 := NewT[any]()
		e1 := l1.PushBack(1)

		l2 := NewT[any]()
		e2 := l2.PushBack(2)

		l1.MoveAfter(e1, e2)
		checkTList(t, l1, []any{1})
		checkTList(t, l2, []any{2})

		l1.MoveBefore(e1, e2)
		checkTList(t, l1, []any{1})
		checkTList(t, l2, []any{2})
	})
}

func TestTList_RemoveAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		l.PushBack(1)
		l.RemoveAll()
		checkTList(t, l, []any{})
		l.PushBack(2)
		checkTList(t, l, []any{2})
	})
}

func TestTList_PushFronts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2}
		l.PushFronts(a1)
		checkTList(t, l, []any{2, 1})
		a1 = []any{3, 4, 5}
		l.PushFronts(a1)
		checkTList(t, l, []any{5, 4, 3, 2, 1})
	})
}

func TestTList_PushBacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2}
		l.PushBacks(a1)
		checkTList(t, l, []any{1, 2})
		a1 = []any{3, 4, 5}
		l.PushBacks(a1)
		checkTList(t, l, []any{1, 2, 3, 4, 5})
	})
}

func TestTList_PopBacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		a2 := []any{"a", "c", "b", "e"}
		l.PushFronts(a1)
		i1 := l.PopBacks(2)
		t.Assert(i1, []any{1, 2})

		l.PushBacks(a2) // 4.3,a,c,b,e
		i1 = l.PopBacks(3)
		t.Assert(i1, []any{"e", "b", "c"})
	})
}

func TestTList_PopFronts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopFronts(2)
		t.Assert(i1, []any{4, 3})
		t.Assert(l.Len(), 2)
	})
}

func TestTList_PopBackAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopBackAll()
		t.Assert(i1, []any{1, 2, 3, 4})
		t.Assert(l.Len(), 0)
	})
}

func TestTList_PopFrontAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopFrontAll()
		t.Assert(i1, []any{4, 3, 2, 1})
		t.Assert(l.Len(), 0)
	})
}

func TestTList_FrontAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.FrontAll()
		t.Assert(i1, []any{4, 3, 2, 1})
		t.Assert(l.Len(), 4)
	})
}

func TestTList_BackAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.BackAll()
		t.Assert(i1, []any{1, 2, 3, 4})
		t.Assert(l.Len(), 4)
	})
}

func TestTList_FrontValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		l2 := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.FrontValue()
		t.Assert(gconv.Int(i1), 4)
		t.Assert(l.Len(), 4)

		i1 = l2.FrontValue()
		t.Assert(i1, nil)
	})
}

func TestTList_BackValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		l2 := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.BackValue()
		t.Assert(gconv.Int(i1), 1)
		t.Assert(l.Len(), 4)

		i1 = l2.FrontValue()
		t.Assert(i1, nil)
	})
}

func TestTList_Back(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		t.Assert(e1.Value, 1)
		t.Assert(l.Len(), 4)
	})
}

func TestTList_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		t.Assert(l.Size(), 4)
		l.PopFront()
		t.Assert(l.Size(), 3)
	})
}

func TestTList_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		l.Removes([]*TElement[any]{e1})
		t.Assert(l.Len(), 3)

		e2 := l.Back()
		l.Removes([]*TElement[any]{e2})
		t.Assert(l.Len(), 2)
		checkTList(t, l, []any{4, 3})
	})
}

func TestTList_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, 3, 4, 5, 6, 7, 8, 9})

		t.Assert(l.PopBack(), 9)
		t.Assert(l.PopBacks(2), []any{8, 7})
		t.Assert(l.PopFront(), 1)
		t.Assert(l.PopFronts(2), []any{2, 3})
	})
}

func TestTList_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		l.Clear()
		t.Assert(l.Len(), 0)
	})
}

func TestTList_IteratorAsc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 5, 6, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e *TElement[any]) bool {
			return gconv.Int(e1.Value) > 2
		}
		checkTList(t, l, []any{4, 3, 6, 5, 2, 1})
		l.IteratorAsc(fun1)
		checkTList(t, l, []any{4, 3, 6, 5, 2, 1})
	})
}

func TestTList_IteratorDesc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e *TElement[any]) bool {
			return gconv.Int(e1.Value) > 6
		}
		l.IteratorDesc(fun1)
		t.Assert(l.Len(), 4)
		checkTList(t, l, []any{4, 3, 2, 1})
	})
}

func TestTList_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		a1 := []any{"a", "b", "c", "d", "e"}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e *TElement[any]) bool {
			return gconv.String(e1.Value) > "c"
		}
		checkTList(t, l, []any{"e", "d", "c", "b", "a"})
		l.Iterator(fun1)
		checkTList(t, l, []any{"e", "d", "c", "b", "a"})
	})
}

func TestTList_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, "a", `"b"`, `\c`})
		t.Assert(l.Join(","), `1,2,a,"b",\c`)
		t.Assert(l.Join("."), `1.2.a."b".\c`)
	})
}

func TestTList_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, "a", `"b"`, `\c`})
		t.Assert(l.String(), `[1,2,a,"b",\c]`)
	})
}

func TestTList_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		a := []any{"a", "b", "c"}
		l := NewT[any]()
		l.PushBacks(a)
		b1, err1 := json.Marshal(l)
		b2, err2 := json.Marshal(a)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		a := []any{"a", "b", "c"}
		l := NewT[any]()
		b, err := json.Marshal(a)
		t.AssertNil(err)

		err = json.UnmarshalUseNumber(b, l)
		t.AssertNil(err)
		t.Assert(l.FrontAll(), a)
	})
	gtest.C(t, func(t *gtest.T) {
		var l TList[any]
		a := []any{"a", "b", "c"}
		b, err := json.Marshal(a)
		t.AssertNil(err)

		err = json.UnmarshalUseNumber(b, &l)
		t.AssertNil(err)
		t.Assert(l.FrontAll(), a)
	})
}

func TestTList_UnmarshalValue(t *testing.T) {
	type list struct {
		Name string
		List *TList[any]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var tlist *list
		err := gconv.Struct(map[string]any{
			"name": "john",
			"list": []byte(`[1,2,3]`),
		}, &tlist)
		t.AssertNil(err)
		t.Assert(tlist.Name, "john")
		t.Assert(tlist.List.FrontAll(), []any{1, 2, 3})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var tlist *list
		err := gconv.Struct(map[string]any{
			"name": "john",
			"list": []any{1, 2, 3},
		}, &tlist)
		t.AssertNil(err)
		t.Assert(tlist.Name, "john")
		t.Assert(tlist.List.FrontAll(), []any{1, 2, 3})
	})
}

func TestTList_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, "a", `"b"`, `\c`})
		copyList := l.DeepCopy()
		cl := copyList.(*TList[any])
		cl.PopBack()
		t.AssertNE(l.Size(), cl.Size())
	})
	// Nil pointer deep copy
	gtest.C(t, func(t *gtest.T) {
		var l *TList[any]
		copyList := l.DeepCopy()
		t.AssertNil(copyList)
	})
}

func TestTList_ToList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, 3, 4, 5})
		nl := l.ToList()
		t.Assert(nl.Len(), 5)

		// Verify elements
		i := 1
		for e := nl.Front(); e != nil; e = e.Next() {
			t.Assert(e.Value, i)
			i++
		}
	})
	// Empty list
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		nl := l.ToList()
		t.Assert(nl.Len(), 0)
	})
}

func TestTList_AppendList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, 3})
		nl := list.New()
		nl.PushBack(4)
		nl.PushBack(5)

		l.AppendList(nl)
		t.Assert(l.Len(), 5)
		t.Assert(l.FrontAll(), []any{1, 2, 3, 4, 5})
	})
	// Append empty list
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, 3})
		nl := list.New()
		l.AppendList(nl)
		t.Assert(l.Len(), 3)
		t.Assert(l.FrontAll(), []any{1, 2, 3})
	})
	// Append with type mismatch (should skip)
	gtest.C(t, func(t *gtest.T) {
		l := NewT[int]()
		nl := list.New()
		nl.PushBack(1)
		nl.PushBack("string") // type mismatch
		nl.PushBack(2)

		l.AppendList(nl)
		t.Assert(l.Len(), 2)
		t.Assert(l.FrontAll(), []int{1, 2})
	})
}

func TestTList_AssignList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, 3})
		nl := list.New()
		nl.PushBack(4)
		nl.PushBack(5)
		nl.PushBack(6)

		skipped := l.AssignList(nl)
		t.Assert(skipped, 0)
		t.Assert(l.Len(), 3)
		t.Assert(l.FrontAll(), []any{4, 5, 6})
	})
	// Assign empty list
	gtest.C(t, func(t *gtest.T) {
		l := NewTFrom([]any{1, 2, 3})
		nl := list.New()

		skipped := l.AssignList(nl)
		t.Assert(skipped, 0)
		t.Assert(l.Len(), 0)
	})
	// Assign with type mismatch (should return skipped count)
	gtest.C(t, func(t *gtest.T) {
		l := NewT[int]()
		nl := list.New()
		nl.PushBack(1)
		nl.PushBack("string") // type mismatch
		nl.PushBack(2)
		nl.PushBack("another") // type mismatch

		skipped := l.AssignList(nl)
		t.Assert(skipped, 2)
		t.Assert(l.Len(), 2)
		t.Assert(l.FrontAll(), []int{1, 2})
	})
}

func TestTList_String_Nil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var l *TList[any]
		t.Assert(l.String(), "")
	})
}

func TestTList_UnmarshalJSON_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		err := l.UnmarshalJSON([]byte("invalid json"))
		t.AssertNE(err, nil)
	})
}

func TestTList_UnmarshalValue_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		err := l.UnmarshalValue(`[1,2,3]`)
		t.AssertNil(err)
		t.Assert(l.FrontAll(), []any{1, 2, 3})
	})
}

func TestTList_UnmarshalValue_Bytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		err := l.UnmarshalValue([]byte(`[1,2,3]`))
		t.AssertNil(err)
		t.Assert(l.FrontAll(), []any{1, 2, 3})
	})
}

func TestTList_DeepCopy_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		copyList := l.DeepCopy()
		cl := copyList.(*TList[any])
		t.Assert(cl.Len(), 0)
	})
}

func TestTList_AppendList_WithTypeMismatch(t *testing.T) {
	// Test appendList internal function through AppendList with mixed types
	gtest.C(t, func(t *gtest.T) {
		l := NewT[int]()
		nl := list.New()
		// Only add non-matching types
		nl.PushBack("string1")
		nl.PushBack("string2")

		l.AppendList(nl)
		t.Assert(l.Len(), 0)
	})
}

func TestTList_UnmarshalValue_Error(t *testing.T) {
	// Test UnmarshalValue with data through default case
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		// Pass a slice directly through default case
		_ = l.UnmarshalValue([]any{1, 2, 3})
		t.Assert(l.Len(), 3)
		t.Assert(l.FrontAll(), []any{1, 2, 3})
	})
	// Test UnmarshalValue error in string case
	gtest.C(t, func(t *gtest.T) {
		l := NewT[any]()
		err := l.UnmarshalValue("invalid json")
		t.AssertNE(err, nil)
	})
}
