// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package resolver defines APIs for name resolution in gRPC.
package resolver

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/resolver"
)

const (
	rawSvcKeyInSubConnInfo = `RawService`
)

func init() {
	// It uses default builder handling the DNS for grpc service requests.
	resolver.Register(&Builder{})
}

// SetRegistry sets the default Registry implements as your own implemented interface.
func SetRegistry(registry gsvc.Registry) {
	if registry == nil {
		panic(gerror.New(`invalid Registry value "nil" given`))
	}
	gsvc.SetRegistry(registry)
}
