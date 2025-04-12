// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v3/text/gregex"
	"github.com/gogf/gf/v3/text/gstr"
	"github.com/gogf/gf/v3/util/gconv"
)

// Model is core struct implementing the DAO for ORM.
type Model struct {
	db             DB                // Underlying DB interface.
	tx             TX                // Underlying TX interface.
	rawSql         string            // rawSql is the raw SQL string which marks a raw SQL based Model not a table based Model.
	schema         string            // Custom database schema.
	linkType       int               // Mark for operation on master or slave.
	tablesInit     string            // Table names when model initialization.
	tables         string            // Operation table names, which can be more than one table names and aliases, like: "user", "user u", "user u, user_detail ud".
	fields         []any             // Operation fields, multiple fields joined using char ','.
	fieldsEx       []any             // Excluded operation fields, it here uses slice instead of string type for quick filtering.
	withArray      []any             // Arguments for With feature.
	withAll        bool              // Enable model association operations on all objects that have "with" tag in the struct.
	extraArgs      []any             // Extra custom arguments for sql, which are prepended to the arguments before sql committed to underlying driver.
	whereBuilder   *WhereBuilder     // Condition builder for where operation.
	groupBy        string            // Used for "group by" statement.
	orderBy        string            // Used for "order by" statement.
	having         []any             // Used for "having..." statement.
	start          int               // Used for "select ... start, limit ..." statement.
	limit          int               // Used for "select ... start, limit ..." statement.
	option         int               // Option for extra operation features.
	offset         int               // Offset statement for some databases grammar.
	partition      string            // Partition table partition name.
	data           any               // Data for operation, which can be type of map/[]map/struct/*struct/string, etc.
	batch          int               // Batch number for batch Insert/Replace/Save operations.
	filter         bool              // Filter data and where key-value pairs according to the fields of the table.
	distinct       string            // Force the query to only return distinct results.
	lockInfo       string            // Lock for update or in shared lock.
	cacheEnabled   bool              // Enable sql result cache feature, which is mainly for indicating cache duration(especially 0) usage.
	cacheOption    CacheOption       // Cache option for query statement.
	hookHandler    HookHandler       // Hook functions for model hook feature.
	unscoped       bool              // Disables soft deleting features when select/delete operations.
	onDuplicate    any               // onDuplicate is used for on Upsert clause.
	onDuplicateEx  any               // onDuplicateEx is used for excluding some columns on Upsert clause.
	onConflict     any               // onConflict is used for conflict keys on Upsert clause.
	tableAliasMap  map[string]string // Table alias to true table name, usually used in join statements.
	softTimeOption SoftTimeOption    // SoftTimeOption is the option to customize soft time feature for Model.
	shardingConfig ShardingConfig    // ShardingConfig for database/table sharding feature.
	shardingValue  any               // Sharding value for sharding feature.
	handlers       []ModelHandler
}

// ModelHandler is a function that handles given Model and returns a new Model that is custom modified.
type ModelHandler func(ctx context.Context, model *Model) *Model

// ChunkHandler is a function that is used in function Chunk, which handles given Result and error.
// It returns true if it wants to continue chunking, or else it returns false to stop chunking.
type ChunkHandler func(result Result, err error) bool

const (
	linkTypeMaster           = 1
	linkTypeSlave            = 2
	defaultField             = "*"
	whereHolderOperatorWhere = 1
	whereHolderOperatorAnd   = 2
	whereHolderOperatorOr    = 3
	whereHolderTypeDefault   = "Default"
	whereHolderTypeNoArgs    = "NoArgs"
	whereHolderTypeIn        = "In"
)

// Model creates and returns a new ORM model from given schema.
// The parameter `tableNameQueryOrStruct` can be more than one table names, and also alias name, like:
//  1. Model names:
//     db.Model("user")
//     db.Model("user u")
//     db.Model("user, user_detail")
//     db.Model("user u, user_detail ud")
//  2. Model name with alias:
//     db.Model("user", "u")
//  3. Model name with sub-query:
//     db.Model("? AS a, ? AS b", subQuery1, subQuery2)
func (c *Core) Model(tableNameQueryOrStruct ...any) *Model {
	newModel := &Model{}
	return newModel.Handler(func(ctx context.Context, model *Model) *Model {
		var (
			tableStr  string
			tableName string
			extraArgs []any
		)
		// Model creation with sub-query.
		if len(tableNameQueryOrStruct) > 1 {
			conditionStr := gconv.String(tableNameQueryOrStruct[0])
			if gstr.Contains(conditionStr, "?") {
				whereHolder := WhereHolder{
					Where: conditionStr,
					Args:  tableNameQueryOrStruct[1:],
				}
				tableStr, extraArgs = formatWhereHolder(ctx, c.db, formatWhereHolderInput{
					WhereHolder: whereHolder,
					OmitNil:     false,
					OmitEmpty:   false,
					Schema:      "",
					Table:       "",
				})
			}
		}
		// Normal model creation.
		if tableStr == "" {
			tableNames := make([]string, len(tableNameQueryOrStruct))
			for k, v := range tableNameQueryOrStruct {
				if s, ok := v.(string); ok {
					tableNames[k] = s
				} else if tableName = getTableNameFromOrmTag(v); tableName != "" {
					tableNames[k] = tableName
				}
			}
			if len(tableNames) > 1 {
				tableStr = fmt.Sprintf(
					`%s AS %s`, c.QuotePrefixTableName(tableNames[0]), c.QuoteWord(tableNames[1]),
				)
			} else if len(tableNames) == 1 {
				tableStr = c.QuotePrefixTableName(tableNames[0])
			}
		}
		m := &Model{
			db:            c.db,
			schema:        c.schema,
			tablesInit:    tableStr,
			tables:        tableStr,
			start:         -1,
			offset:        -1,
			filter:        true,
			extraArgs:     extraArgs,
			tableAliasMap: make(map[string]string),
		}
		m.whereBuilder = m.Builder()
		return m
	})

}

// Raw creates and returns a model based on a raw sql not a table.
// Example:
//
//	db.Raw("SELECT * FROM `user` WHERE `name` = ?", "john").Scan(&result)
func (c *Core) Raw(rawSql string, args ...any) *Model {
	newModel := c.Model()
	return newModel.Handler(func(ctx context.Context, model *Model) *Model {
		model.rawSql = rawSql
		model.extraArgs = args
		return model
	})
}

// Raw sets current model as a raw sql model.
// Example:
//
//	db.Raw("SELECT * FROM `user` WHERE `name` = ?", "john").Scan(&result)
//
// See Core.Raw.
func (m *Model) Raw(rawSql string, args ...any) *Model {
	newModel := m.db.Raw(rawSql, args...)
	return newModel.Handler(func(ctx context.Context, model *Model) *Model {
		model.db = m.db
		model.tx = m.tx
		return model
	})
}

func (tx *TXCore) Raw(rawSql string, args ...any) *Model {
	return tx.Model().Raw(rawSql, args...)
}

// With creates and returns an ORM model based on metadata of given object.
func (c *Core) With(objects ...any) *Model {
	return c.db.Model().With(objects...)
}

// Partition sets Partition name.
// Example:
// dao.User.Partition（"p1","p2","p3").All()
func (m *Model) Partition(partitions ...string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.partition = gstr.Join(partitions, ",")
		return model
	})
}

// Model acts like Core.Model except it operates on transaction.
// See Core.Model.
func (tx *TXCore) Model(tableNameQueryOrStruct ...any) *Model {
	newModel := tx.db.Model(tableNameQueryOrStruct...)
	return newModel.Handler(func(ctx context.Context, model *Model) *Model {
		model.db = tx.db
		model.tx = tx
		return model
	})
}

// With acts like Core.With except it operates on transaction.
// See Core.With.
func (tx *TXCore) With(object any) *Model {
	return tx.Model().With(object)
}

// As sets an alias name for current table.
func (m *Model) As(as string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		if model.tables == "" {
			return model
		}
		split := " JOIN "
		if gstr.ContainsI(model.tables, split) {
			// For join table.
			array := gstr.Split(model.tables, split)
			array[len(array)-1], _ = gregex.ReplaceString(
				`(.+) ON`,
				fmt.Sprintf(`$1 AS %s ON`, as), array[len(array)-1],
			)
			model.tables = gstr.Join(array, split)
		} else {
			// For base table.
			model.tables = gstr.TrimRight(model.tables) + " AS " + as
		}
		return model
	})
}

// DB sets/changes the db object for current operation.
func (m *Model) DB(db DB) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.db = db
		return model
	})
}

// TX sets/changes the transaction for current operation.
func (m *Model) TX(tx TX) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.db = tx.GetDB()
		model.tx = tx
		return model
	})
}

// Schema sets the schema for current operation.
func (m *Model) Schema(schema string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.schema = schema
		return model
	})
}

// Clone creates and returns a new model which is a Clone of current model.
// Note that it uses deep-copy for the Clone.
func (m *Model) Clone() *Model {
	newModel := &Model{
		handlers: make([]ModelHandler, len(m.handlers)),
	}
	copy(newModel.handlers, m.handlers)
	return newModel
}

// Master marks the following operation on master node.
func (m *Model) Master() *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.linkType = linkTypeMaster
		return model
	})
}

// Slave marks the following operation on slave node.
// Note that it makes sense only if there's any slave node configured.
func (m *Model) Slave() *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.linkType = linkTypeSlave
		return model
	})
}

// Args sets custom arguments for model operation.
func (m *Model) Args(args ...any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.extraArgs = append(model.extraArgs, args)
		return model
	})
}

// Handler calls each of `handlers` on current Model and returns a new Model.
// ModelHandler is a function that handles given Model and returns a new Model that is custom modified.
func (m *Model) Handler(handlers ...ModelHandler) *Model {
	m.handlers = append(m.handlers, handlers...)
	return m
}
