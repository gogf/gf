// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "context"

// ShardingInput holds the input parameters for sharding.
type ShardingInput struct {
	Table  string // The original table name.
	Schema string // The original schema name. Note that this might be empty according database configuration.
}

// ShardingOutput holds the output parameters for sharding.
type ShardingOutput struct {
	Table  string // The target table name.
	Schema string // The target schema name.
}

// ShardingFunc is the custom function for records sharding.
type ShardingFunc func(ctx context.Context, in ShardingInput) (out *ShardingOutput, err error)

// Sharding sets custom sharding function for current model.
func (m *Model) Sharding(f ShardingFunc) *Model {
	model := m.getModel()
	model.shardingFunc = f
	return model
}
