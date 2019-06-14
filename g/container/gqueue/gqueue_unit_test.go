package gqueue_test
import (
	"testing"
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/test/gtest"
)


func TestQueue_Size(t *testing.T) {
	q1:=gqueue.New(3)
	q1.Push(1)
	q1.Push(2)
	q1.Push(3)
	gtest.Assert(q1.Size(),3)
}
