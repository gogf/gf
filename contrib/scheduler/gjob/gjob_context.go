// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Lifecycle context helpers for job servers.

package gjob

import (
	"context"

	"github.com/gogf/gf/v2/os/gctx"
)

// normalizeCtx maps nil to gctx.GetInitCtx so job servers inherit framework defaults.
func normalizeCtx(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return gctx.GetInitCtx()
}
