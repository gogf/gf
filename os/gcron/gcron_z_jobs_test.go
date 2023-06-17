package gcron

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
	"time"
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

func TestAddParallelSingleton(t *testing.T) {
	t0 := func(ctx context.Context) {
		fmt.Println("tick 0")
		time.Sleep(500 * time.Millisecond)
	}
	t1 := func(ctx context.Context) {
		fmt.Println("tick 1")
	}
	_, _ = AddParallelSingleton(gctx.New(), "* * * * * *", t0, t1)
	time.Sleep(3 * time.Second)
}

func TestAddSerialSingleton(t *testing.T) {
	t0 := func(ctx context.Context) {
		fmt.Println("tick 0")
		time.Sleep(500 * time.Millisecond)
	}
	t1 := func(ctx context.Context) {
		fmt.Println("tick 1")
	}
	_, _ = AddSerialSingleton(gctx.New(), "* * * * * *", t0, t1)
	time.Sleep(3 * time.Second)
}

func TestAddSerialGroupSingleton(t *testing.T) {
	group1 := []JobFunc{
		func(ctx context.Context) {
			fmt.Println("tick 0")
			time.Sleep(500 * time.Millisecond)
		},
		func(ctx context.Context) {
			fmt.Println("tick 1")
		},
	}
	group2 := []JobFunc{
		func(ctx context.Context) {
			fmt.Println("tick 2")
		},
	}
	_, _ = AddSerialGroupSingleton(gctx.New(), "* * * * * *", group1, group2)
	time.Sleep(3 * time.Second)
}
