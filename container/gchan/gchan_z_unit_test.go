package gchan_test

import (
	"errors"
	"testing"

	"github.com/gogf/gf/container/gchan"
	"github.com/gogf/gf/test/gtest"
)

func Test_Gchan(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ch := gchan.New(10)

		t.Assert(ch.Cap(), 10)
		t.Assert(ch.Push(1), nil)
		t.Assert(ch.Len(), 1)
		t.Assert(ch.Size(), 1)
		ch.Pop()
		t.Assert(ch.Len(), 0)
		t.Assert(ch.Size(), 0)
		ch.Close()
		t.Assert(ch.Push(1), errors.New("channel is closed"))

		ch = gchan.New(0)
		ch1 := gchan.New(0)
		go func() {
			var i = 0
			for {
				v := ch.Pop()
				if v == nil {
					ch1.Push(i)
					break
				}
				t.Assert(v, i)
				i++
			}
		}()

		for index := 0; index < 10; index++ {
			ch.Push(index)
		}
		ch.Close()
		t.Assert(ch1.Pop(), 10)
		ch1.Close()
	})
}
