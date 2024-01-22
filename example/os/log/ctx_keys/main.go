// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	var ctx = context.Background()
	ctx = context.WithValue(ctx, "RequestId", "123456789")
	ctx = context.WithValue(ctx, "UserId", "10000")
	g.Log().Error(ctx, "runtime error")
}
