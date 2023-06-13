package gcron

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestWaitAsync(t *testing.T) {
	f0 := func(ctx context.Context) {
		fmt.Println("f0")
	}
	f1 := func(ctx context.Context) {
		panic("f1")
	}
	f2 := func(ctx context.Context) {
		panic("f2")
	}
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	waitAsync(gctx.New(), f0, f1, f2)
}
