// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"time"

	"github.com/gogf/gf/v2/os/gcron"
)

func Example_cronAddSingleton() {
	array := garray.New(true)
	cron := gcron.New()
	fmt.Println(array.Len())
	cron.AddSingleton(ctx, "* * * * * *", func(ctx context.Context) {
		glog.Print(context.TODO(), "doing")
		if array.Len() < 1 {
			array.Append(1)

		} else {
			cron.Remove("cron1")
		}
	}, "cron1")
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len())
	// Output:
	// 0
	// 1
}

func Example_cronAddOnce() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddOnce(ctx, "@every 2s", func(ctx context.Context) {
		array.Append(1)
	})
	fmt.Println(cron.Size(), array.Len())
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(cron.Size(), array.Len())

	// Output:
	// 1 0
	// 0 1
}

func Example_cronAddTimes() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddTimes(ctx, "@every 2s", 2, func(ctx context.Context) {
		array.Append(1)
	})
	fmt.Println(cron.Size(), array.Len())
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(cron.Size(), array.Len())
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(cron.Size(), array.Len())

	// Output:
	// 1 0
	// 1 1
	// 0 2
}

func Example_cronEntries() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddTimes(ctx, "@every 1s", 2, func(ctx context.Context) {
		array.Append(1)
	}, "cron1")
	cron.AddOnce(ctx, "@every 1s", func(ctx context.Context) {
		array.Append(1)
	}, "cron2")
	entries := cron.Entries()
	for k, v := range entries {
		fmt.Println(k, v.Name, v.Time)

	}
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len())

	// May Output:
	// 0 cron2 2022-02-09 10:11:47.2421345 +0800 CST m=+0.159116501
	// 1 cron1 2022-02-09 10:11:47.2421345 +0800 CST m=+0.159116501
	// 3
}

func Example_cronSearch() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddTimes(ctx, "@every 1s", 2, func(ctx context.Context) {
		array.Append(1)
	}, "cron1")
	cron.AddOnce(ctx, "@every 1s", func(ctx context.Context) {
		array.Append(1)
	}, "cron2")
	search := cron.Search("cron2")

	g.Log().Print(ctx, search)

	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len())

	// Output:
	// 3
}

func Example_cronStop() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddTimes(ctx, "@every 2s", 1, func(ctx context.Context) {
		array.Append(1)
	}, "cron1")
	cron.AddOnce(ctx, "@every 2s", func(ctx context.Context) {

		array.Append(1)
	}, "cron2")
	fmt.Println(array.Len(), cron.Size())
	cron.Stop("cron2")
	fmt.Println(array.Len(), cron.Size())
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len(), cron.Size())
	// Output:
	// 0 2
	// 0 2
	// 1 1
}

func Example_cronRemove() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddTimes(ctx, "@every 2s", 1, func(ctx context.Context) {
		array.Append(1)
	}, "cron1")
	cron.AddOnce(ctx, "@every 2s", func(ctx context.Context) {

		array.Append(1)
	}, "cron2")
	fmt.Println(array.Len(), cron.Size())
	cron.Remove("cron2")
	fmt.Println(array.Len(), cron.Size())
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len(), cron.Size())
	// Output:
	// 0 2
	// 0 1
	// 1 0
}

func Example_cronStart() {
	var (
		ctx = gctx.New()
	)
	cron := gcron.New()
	array := garray.New(true)
	cron.AddOnce(ctx, "@every 2s", func(ctx context.Context) {

		array.Append(1)
	}, "cron2")
	cron.Stop("cron2")
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len(), cron.Size())
	cron.Start("cron2")
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(array.Len(), cron.Size())

	// Output:
	// 0 1
	// 1 0
}
