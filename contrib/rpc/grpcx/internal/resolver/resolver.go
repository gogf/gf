// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package resolver defines APIs for name resolution in gRPC.
package resolver

import (
	"google.golang.org/grpc/resolver"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
)

const (
	rawSvcKeyInSubConnInfo = `RawService`
)

func init() {
	// It registers default resolver here.
	// It uses default builder handling the name resolving for grpc service requests.
	// Use `grpc.WithResolver` to custom resolver for client.
	resolver.Register(NewBuilder(gsvc.GetRegistry()))
}

// SetRegistry sets the default Registry implements as your own implemented interface.
func SetRegistry(registry gsvc.Registry) {
	if registry == nil {
		panic(gerror.New(`invalid Registry value "nil" given`))
	}
	gsvc.SetRegistry(registry)
	resolver.Register(NewBuilder(registry))
}
