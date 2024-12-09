// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"hash/fnv"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// ShardingConfig defines the configuration for database/table sharding.
type ShardingConfig struct {
	// Table sharding configuration
	Table ShardingTableConfig
	// Schema sharding configuration
	Schema ShardingSchemaConfig
}

// ShardingSchemaConfig defines the configuration for database sharding.
type ShardingSchemaConfig struct {
	// Enable schema sharding
	Enable bool
	// Schema rule prefix, e.g., "db_"
	Prefix string
	// ShardingRule defines how to route data to different database nodes
	Rule ShardingRule
}

// ShardingTableConfig defines the configuration for table sharding
type ShardingTableConfig struct {
	// Enable table sharding
	Enable bool
	// Table rule prefix, e.g., "user_"
	Prefix string
	// ShardingRule defines how to route data to different tables
	Rule ShardingRule
}

// ShardingRule defines the interface for sharding rules
type ShardingRule interface {
	// SchemaName returns the target schema name based on sharding value.
	SchemaName(ctx context.Context, config ShardingSchemaConfig, value any) (string, error)
	// TableName returns the target table name based on sharding value.
	TableName(ctx context.Context, config ShardingTableConfig, value any) (string, error)
}

// DefaultShardingRule implements a simple modulo-based sharding rule
type DefaultShardingRule struct {
	// Number of schema count.
	SchemaCount int
	// Number of tables per schema.
	TableCount int
}

// Sharding creates a sharding model with given sharding configuration.
func (m *Model) Sharding(config ShardingConfig) *Model {
	model := m.getModel()
	model.shardingConfig = config
	return model
}

// ShardingValue sets the sharding value for routing
func (m *Model) ShardingValue(value any) *Model {
	model := m.getModel()
	model.shardingValue = value
	return model
}

// getActualSchema returns the actual schema based on sharding configuration.
// TODO it does not support schemas in different database config node.
func (m *Model) getActualSchema(ctx context.Context, defaultSchema string) (string, error) {
	if !m.shardingConfig.Schema.Enable {
		return defaultSchema, nil
	}
	if m.shardingValue == nil {
		return defaultSchema, gerror.NewCode(
			gcode.CodeInvalidParameter, "sharding value is required when sharding feature enabled",
		)
	}
	if m.shardingConfig.Schema.Rule == nil {
		return defaultSchema, gerror.NewCode(
			gcode.CodeInvalidParameter, "sharding rule is required when sharding feature enabled",
		)
	}
	return m.shardingConfig.Schema.Rule.SchemaName(ctx, m.shardingConfig.Schema, m.shardingValue)
}

// getActualTable returns the actual table name based on sharding configuration
func (m *Model) getActualTable(ctx context.Context, defaultTable string) (string, error) {
	if !m.shardingConfig.Table.Enable {
		return defaultTable, nil
	}
	if m.shardingValue == nil {
		return defaultTable, gerror.NewCode(
			gcode.CodeInvalidParameter, "sharding value is required when sharding feature enabled",
		)
	}
	if m.shardingConfig.Table.Rule == nil {
		return defaultTable, gerror.NewCode(
			gcode.CodeInvalidParameter, "sharding rule is required when sharding feature enabled",
		)
	}
	return m.shardingConfig.Table.Rule.TableName(ctx, m.shardingConfig.Table, m.shardingValue)
}

// SchemaName implements the default database sharding strategy
func (r *DefaultShardingRule) SchemaName(ctx context.Context, config ShardingSchemaConfig, value any) (string, error) {
	if r.SchemaCount == 0 {
		return "", gerror.NewCode(
			gcode.CodeInvalidParameter, "schema count should not be 0 using DefaultShardingRule when schema sharding enabled",
		)
	}
	hashValue, err := getHashValue(value)
	if err != nil {
		return "", err
	}
	nodeIndex := hashValue % uint64(r.SchemaCount)
	return fmt.Sprintf("%s%d", config.Prefix, nodeIndex), nil
}

// TableName implements the default table sharding strategy
func (r *DefaultShardingRule) TableName(ctx context.Context, config ShardingTableConfig, value any) (string, error) {
	if r.TableCount == 0 {
		return "", gerror.NewCode(
			gcode.CodeInvalidParameter, "table count should not be 0 using DefaultShardingRule when table sharding enabled",
		)
	}
	hashValue, err := getHashValue(value)
	if err != nil {
		return "", err
	}
	tableIndex := hashValue % uint64(r.TableCount)
	return fmt.Sprintf("%s%d", config.Prefix, tableIndex), nil
}

// getHashValue converts sharding value to uint64 hash
func getHashValue(value any) (uint64, error) {
	var rv = reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return gconv.Uint64(value), nil
	default:
		h := fnv.New64a()
		_, err := h.Write(gconv.Bytes(value))
		if err != nil {
			return 0, gerror.WrapCode(gcode.CodeInternalError, err)
		}
		return h.Sum64(), nil
	}
}
