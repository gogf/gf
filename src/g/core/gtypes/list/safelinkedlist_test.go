package list

import (
	"testing"
)

func TestSafeLinkedListLimited(t *testing.T) {
	ll := NewSafeLinkedListLimited(2)
	val1 := "val1"
	val2 := "val1"
	val3 := "val1"

	ll.PushFront(val1)
	ll.PushFront(val2)

	if ll.Len() != 2 || ll.Front().Value.(string) != val2 {
		t.Error("TestSafeLinkedListLimited, Len|Front")
	}

	if ll.PushFront(val3) {
		t.Error("TestSafeLinkedListLimited, PushFront not Limited")
	}

	if ll.PopBack(1)[0].(string) != val1 {
		t.Error("TestSafeLinkedListLimited, PopBack")
	}
}
