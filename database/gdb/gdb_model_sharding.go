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

// ShardingFunc is the custom function for records sharding by certain Model, which supports sharding on table or schema.
// It retrieves the original Table/Schema from ShardingInput, and returns the new Table/Schema by ShardingOutput.
// If the Table/Schema in ShardingOutput is empty string, it then ignores the returned value and uses the default
// Table/Schema to execute the sql statement.
type ShardingFunc func(ctx context.Context, in ShardingInput) (out *ShardingOutput, err error)

// Sharding sets custom sharding function for current model.
// More info please refer to ShardingFunc.
func (m *Model) Sharding(f ShardingFunc) *Model {
	model := m.getModel()
	model.shardingFunc = f
	return model
}
