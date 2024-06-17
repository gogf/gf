// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
)

func ExampleConfig_GetWithEnv() {
	var (
		key = `ENV_TEST`
		ctx = gctx.New()
	)
	v, err := g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("env:%s\n", v)
	if err = genv.Set(key, "gf"); err != nil {
		panic(err)
	}
	v, err = g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("env:%s", v)

	// Output:
	// env:
	// env:gf
}

func ExampleConfig_GetWithCmd() {
	var (
		key = `cmd.test`
		ctx = gctx.New()
	)
	v, err := g.Cfg().GetWithCmd(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cmd:%s\n", v)
	// Re-Initialize custom command arguments.
	os.Args = append(os.Args, fmt.Sprintf(`--%s=yes`, key))
	gcmd.Init(os.Args...)
	// Retrieve the configuration and command option again.
	v, err = g.Cfg().GetWithCmd(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cmd:%s", v)

	// Output:
	// cmd:
	// cmd:yes
}

func Example_NewWithAdapter() {
	var (
		ctx          = gctx.New()
		content      = `{"a":"b", "c":1}`
		adapter, err = gcfg.NewAdapterContent(content)
	)
	if err != nil {
		panic(err)
	}
	config := gcfg.NewWithAdapter(adapter)
	fmt.Println(config.MustGet(ctx, "a"))
	fmt.Println(config.MustGet(ctx, "c"))

	// Output:
	// b
	// 1
}
