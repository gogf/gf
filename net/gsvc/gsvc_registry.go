// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Register registers `service` to default registry..
func Register(ctx context.Context, service Service) error {
	if defaultRegistry == nil {
		return gerror.NewCodef(gcode.CodeNotImplemented, `no Registry is registered`)
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return defaultRegistry.Register(ctx, service)
}

// Deregister removes `service` from default registry.
func Deregister(ctx context.Context, service Service) error {
	if defaultRegistry == nil {
		return gerror.NewCodef(gcode.CodeNotImplemented, `no Registry is registered`)
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return defaultRegistry.Deregister(ctx, service)
}
