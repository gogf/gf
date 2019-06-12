package gchan_test

import (
	"errors"
	"testing"

	"github.com/gogf/gf/g/container/gchan"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Gchan(t *testing.T) {
	gtest.Case(t, func() {
		ch := gchan.New(10)

		gtest.Assert(ch.Cap(), 10)
		gtest.Assert(ch.Push(1), nil)
		gtest.Assert(ch.Len(), 1)
		gtest.Assert(ch.Size(), 1)
		ch.Pop()
		gtest.Assert(ch.Len(), 0)
		gtest.Assert(ch.Size(), 0)
		ch.Close()
		gtest.Assert(ch.Push(1), errors.New("channel is closed"))

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
				gtest.Assert(v, i)
				i++
			}
		}()

		for index := 0; index < 10; index++ {
			ch.Push(index)
		}
		ch.Close()
		gtest.Assert(ch1.Pop(), 10)
		ch1.Close()
	})
}
