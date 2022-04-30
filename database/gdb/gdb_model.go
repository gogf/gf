// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Model is core struct implementing the DAO for ORM.
type Model struct {
	db              DB              // Underlying DB interface.
	tx              *TX             // Underlying TX interface.
	rawSql          string          // rawSql is the raw SQL string which marks a raw SQL based Model not a table based Model.
	schema          string          // Custom database schema.
	linkType        int             // Mark for operation on master or slave.
	tablesInit      string          // Table names when model initialization.
	tables          string          // Operation table names, which can be more than one table names and aliases, like: "user", "user u", "user u, user_detail ud".
	fields          string          // Operation fields, multiple fields joined using char ','.
	fieldsEx        string          // Excluded operation fields, multiple fields joined using char ','.
	withArray       []interface{}   // Arguments for With feature.
	withAll         bool            // Enable model association operations on all objects that have "with" tag in the struct.
	extraArgs       []interface{}   // Extra custom arguments for sql, which are prepended to the arguments before sql committed to underlying driver.
	whereBuilder    *WhereBuilder   // Condition builder for where operation.
	groupBy         string          // Used for "group by" statement.
	orderBy         string          // Used for "order by" statement.
	having          []interface{}   // Used for "having..." statement.
	start           int             // Used for "select ... start, limit ..." statement.
	limit           int             // Used for "select ... start, limit ..." statement.
	option          int             // Option for extra operation features.
	offset          int             // Offset statement for some databases grammar.
	data            interface{}     // Data for operation, which can be type of map/[]map/struct/*struct/string, etc.
	batch           int             // Batch number for batch Insert/Replace/Save operations.
	filter          bool            // Filter data and where key-value pairs according to the fields of the table.
	distinct        string          // Force the query to only return distinct results.
	lockInfo        string          // Lock for update or in shared lock.
	cacheEnabled    bool            // Enable sql result cache feature, which is mainly for indicating cache duration(especially 0) usage.
	cacheOption     CacheOption     // Cache option for query statement.
	hookHandler     HookHandler     // Hook functions for model hook feature.
	shardingHandler ShardingHandler // Custom sharding handler for sharding feature.
	unscoped        bool            // Disables soft deleting features when select/delete operations.
	safe            bool            // If true, it clones and returns a new model object whenever operation done; or else it changes the attribute of current model.
	onDuplicate     interface{}     // onDuplicate is used for ON "DUPLICATE KEY UPDATE" statement.
	onDuplicateEx   interface{}     // onDuplicateEx is used for excluding some columns ON "DUPLICATE KEY UPDATE" statement.
}

// ModelHandler is a function that handles given Model and returns a new Model that is custom modified.
type ModelHandler func(m *Model) *Model

// ChunkHandler is a function that is used in function Chunk, which handles given Result and error.
// It returns true if it wants to continue chunking, or else it returns false to stop chunking.
type ChunkHandler func(result Result, err error) bool

const (
	linkTypeMaster           = 1
	linkTypeSlave            = 2
	defaultFields            = "*"
	whereHolderOperatorWhere = 1
	whereHolderOperatorAnd   = 2
	whereHolderOperatorOr    = 3
	whereHolderTypeDefault   = "Default"
	whereHolderTypeNoArgs    = "NoArgs"
	whereHolderTypeIn        = "In"
)

// Model creates and returns a new ORM model from given schema.
// The parameter `tableNameQueryOrStruct` can be more than one table names, and also alias name, like:
// 1. Model names:
//    db.Model("user")
//    db.Model("user u")
//    db.Model("user, user_detail")
//    db.Model("user u, user_detail ud")
// 2. Model name with alias:
//    db.Model("user", "u")
// 3. Model name with sub-query:
//    db.Model("? AS a, ? AS b", subQuery1, subQuery2)
func (c *Core) Model(tableNameQueryOrStruct ...interface{}) *Model {
	var (
		ctx       = c.db.GetCtx()
		tableStr  string
		tableName string
		extraArgs []interface{}
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
		db:         c.db,
		schema:     c.schema,
		tablesInit: tableStr,
		tables:     tableStr,
		fields:     defaultFields,
		start:      -1,
		offset:     -1,
		filter:     true,
		extraArgs:  extraArgs,
	}
	m.whereBuilder = m.Builder()
	if defaultModelSafe {
		m.safe = true
	}
	return m
}

// Raw creates and returns a model based on a raw sql not a table.
// Example:
//     db.Raw("SELECT * FROM `user` WHERE `name` = ?", "john").Scan(&result)
func (c *Core) Raw(rawSql string, args ...interface{}) *Model {
	model := c.Model()
	model.rawSql = rawSql
	model.extraArgs = args
	return model
}

// Raw sets current model as a raw sql model.
// Example:
//     db.Raw("SELECT * FROM `user` WHERE `name` = ?", "john").Scan(&result)
// See Core.Raw.
func (m *Model) Raw(rawSql string, args ...interface{}) *Model {
	model := m.db.Raw(rawSql, args...)
	model.db = m.db
	model.tx = m.tx
	return model
}

func (tx *TX) Raw(rawSql string, args ...interface{}) *Model {
	return tx.Model().Raw(rawSql, args...)
}

// With creates and returns an ORM model based on metadata of given object.
func (c *Core) With(objects ...interface{}) *Model {
	return c.db.Model().With(objects...)
}

// Model acts like Core.Model except it operates on transaction.
// See Core.Model.
func (tx *TX) Model(tableNameQueryOrStruct ...interface{}) *Model {
	model := tx.db.Model(tableNameQueryOrStruct...)
	model.db = tx.db
	model.tx = tx
	return model
}

// With acts like Core.With except it operates on transaction.
// See Core.With.
func (tx *TX) With(object interface{}) *Model {
	return tx.Model().With(object)
}

// Ctx sets the context for current operation.
func (m *Model) Ctx(ctx context.Context) *Model {
	if ctx == nil {
		return m
	}
	model := m.getModel()
	model.db = model.db.Ctx(ctx)
	if m.tx != nil {
		model.tx = model.tx.Ctx(ctx)
	}
	return model
}

// GetCtx returns the context for current Model.
// It returns `context.Background()` is there's no context previously set.
func (m *Model) GetCtx() context.Context {
	if m.tx != nil && m.tx.ctx != nil {
		return m.tx.ctx
	}
	return m.db.GetCtx()
}

// As sets an alias name for current table.
func (m *Model) As(as string) *Model {
	if m.tables != "" {
		model := m.getModel()
		split := " JOIN "
		if gstr.ContainsI(model.tables, split) {
			// For join table.
			array := gstr.Split(model.tables, split)
			array[len(array)-1], _ = gregex.ReplaceString(`(.+) ON`, fmt.Sprintf(`$1 AS %s ON`, as), array[len(array)-1])
			model.tables = gstr.Join(array, split)
		} else {
			// For base table.
			model.tables = gstr.TrimRight(model.tables) + " AS " + as
		}
		return model
	}
	return m
}

// DB sets/changes the db object for current operation.
func (m *Model) DB(db DB) *Model {
	model := m.getModel()
	model.db = db
	return model
}

// TX sets/changes the transaction for current operation.
func (m *Model) TX(tx *TX) *Model {
	model := m.getModel()
	model.db = tx.db
	model.tx = tx
	return model
}

// Schema sets the schema for current operation.
func (m *Model) Schema(schema string) *Model {
	model := m.getModel()
	model.schema = schema
	return model
}

// Clone creates and returns a new model which is a clone of current model.
// Note that it uses deep-copy for the clone.
func (m *Model) Clone() *Model {
	newModel := (*Model)(nil)
	if m.tx != nil {
		newModel = m.tx.Model(m.tablesInit)
	} else {
		newModel = m.db.Model(m.tablesInit)
	}
	// Basic attributes copy.
	*newModel = *m
	// WhereBuilder copy, note the attribute pointer.
	newModel.whereBuilder = m.whereBuilder.Clone()
	newModel.whereBuilder.model = newModel
	// Shallow copy slice attributes.
	if n := len(m.extraArgs); n > 0 {
		newModel.extraArgs = make([]interface{}, n)
		copy(newModel.extraArgs, m.extraArgs)
	}
	if n := len(m.withArray); n > 0 {
		newModel.withArray = make([]interface{}, n)
		copy(newModel.withArray, m.withArray)
	}
	return newModel
}

// Master marks the following operation on master node.
func (m *Model) Master() *Model {
	model := m.getModel()
	model.linkType = linkTypeMaster
	return model
}

// Slave marks the following operation on slave node.
// Note that it makes sense only if there's any slave node configured.
func (m *Model) Slave() *Model {
	model := m.getModel()
	model.linkType = linkTypeSlave
	return model
}

// Safe marks this model safe or unsafe. If safe is true, it clones and returns a new model object
// whenever the operation done, or else it changes the attribute of current model.
func (m *Model) Safe(safe ...bool) *Model {
	if len(safe) > 0 {
		m.safe = safe[0]
	} else {
		m.safe = true
	}
	return m
}

// Args sets custom arguments for model operation.
func (m *Model) Args(args ...interface{}) *Model {
	model := m.getModel()
	model.extraArgs = append(model.extraArgs, args)
	return model
}

// Handler calls each of `handlers` on current Model and returns a new Model.
// ModelHandler is a function that handles given Model and returns a new Model that is custom modified.
func (m *Model) Handler(handlers ...ModelHandler) *Model {
	model := m.getModel()
	for _, handler := range handlers {
		model = handler(model)
	}
	return model
}
