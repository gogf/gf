// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package resolver

import (
	"google.golang.org/grpc/resolver"

	"github.com/gogf/gf/v2/net/gsvc"
)

// Manager for Builder creating.
type Manager struct{}

// New creates and returns a Builder.
func (m Manager) New(discovery gsvc.Discovery) resolver.Builder {
	return NewBuilder(discovery)
}

// Register sets the default Registry implements as your own implemented interface.
func (m Manager) Register(registry gsvc.Registry) {
	Register(registry)
}
