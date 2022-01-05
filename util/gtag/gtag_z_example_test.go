// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtag_test

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
)

func ExampleSet() {
	type User struct {
		g.Meta `name:"User Struct" description:"{UserDescription}"`
	}
	gtag.Sets(g.MapStrStr{
		`UserDescription`: `This is a demo struct named "User Struct"`,
	})
	fmt.Println(gmeta.Get(User{}, `description`))

	// Output:
	// This is a demo struct named "User Struct"
}
