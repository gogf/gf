// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gcfg"
)

func Example_mapSliceChange() {
	intlog.SetEnabled(false)
	defer intlog.SetEnabled(true)
	// For testing/example only.
	content := `{"map":{"key":"value"}, "slice":[59,90]}`
	gcfg.SetContent(content)
	defer gcfg.RemoveContent()

	m := g.Cfg().GetMap("map")
	fmt.Println(m)

	// Change the key-value pair.
	m["key"] = "john"

	// It changes the underlying key-value pair.
	fmt.Println(g.Cfg().GetMap("map"))

	s := g.Cfg().GetArray("slice")
	fmt.Println(s)

	// Change the value of specified index.
	s[0] = 100

	// It changes the underlying slice.
	fmt.Println(g.Cfg().GetArray("slice"))

	// output:
	// map[key:value]
	// map[key:john]
	// [59 90]
	// [100 90]
}
