// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ShardingInput is input parameters for custom sharding handler.
type ShardingInput struct {
	Table  string           // Current operation table name.
	Schema string           // Current operation schema, usually empty string which means uses default schema from configuration.
	Data   map[string]Value // Accurate key-value pairs from SELECT/INSERT/UPDATE/DELETE statement.
}

// ShardingOutput is output parameters for custom sharding handler.
type ShardingOutput struct {
	Table  string // New table name for current operation. Use empty string for no changes of table name.
	Schema string // New schema name for current operation. Use empty string for using default schema from configuration.
}

type ShardingHandler func(ctx context.Context, in ShardingInput) (out *ShardingOutput, err error)

type callShardingHandlerInput struct {
	Table      string
	InsertData List
	UpdateData interface{}
	Condition  string
	Sql        string
}

func (m *Model) callShardingHandler(ctx context.Context, in callShardingHandlerInput) (out *ShardingOutput, err error) {
	if m.shardingHandler == nil {
		return &ShardingOutput{}, nil
	}
	return
}

func (m *Model) shardingDataFromInsertData(data List) (shardingData map[string]Value, err error) {
	if len(data) == 0 {
		return nil, nil
	}
	shardingData = make(map[string]Value)
	// If given batch data(in batch insert scenario), it uses the first data.
	for k, v := range data[0] {
		shardingData[k] = gvar.New(v)
	}
	return shardingData, nil
}

func (m *Model) shardingDataFromUpdateData(data interface{}) (shardingData map[string]Value, err error) {
	shardingData = make(map[string]Value)
	switch value := data.(type) {
	case map[string]interface{}:
		for k, v := range value {
			shardingData[k] = gvar.New(v)
		}
	case string:

	default:
		return nil, gerror.Newf(`unsupported data of type "%s" for sharding`, reflect.TypeOf(data))
	}
	return
}

func (m *Model) shardingDataFromSql(sql string, args []interface{}) (shardingData map[string]Value, err error) {
	return
}

func (m *Model) shardingDataFromCondition(condition string) (shardingData map[string]Value, err error) {
	return
}
