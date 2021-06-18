// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"time"

	"github.com/gogf/gf/text/gregex"

	"github.com/gogf/gf/text/gstr"
)

// Model is core struct implementing the DAO for ORM.
type Model struct {
	db            DB             // Underlying DB interface.
	tx            *TX            // Underlying TX interface.
	rawSql        string         // rawSql is the raw SQL string which marks a raw SQL based Model not a table based Model.
	schema        string         // Custom database schema.
	linkType      int            // Mark for operation on master or slave.
	tablesInit    string         // Table names when model initialization.
	tables        string         // Operation table names, which can be more than one table names and aliases, like: "user", "user u", "user u, user_detail ud".
	fields        string         // Operation fields, multiple fields joined using char ','.
	fieldsEx      string         // Excluded operation fields, multiple fields joined using char ','.
	withArray     []interface{}  // Arguments for With feature.
	withAll       bool           // Enable model association operations on all objects that have "with" tag in the struct.
	extraArgs     []interface{}  // Extra custom arguments for sql, which are prepended to the arguments before sql committed to underlying driver.
	whereHolder   []*whereHolder // Condition strings for where operation.
	groupBy       string         // Used for "group by" statement.
	orderBy       string         // Used for "order by" statement.
	having        []interface{}  // Used for "having..." statement.
	start         int            // Used for "select ... start, limit ..." statement.
	limit         int            // Used for "select ... start, limit ..." statement.
	option        int            // Option for extra operation features.
	offset        int            // Offset statement for some databases grammar.
	data          interface{}    // Data for operation, which can be type of map/[]map/struct/*struct/string, etc.
	batch         int            // Batch number for batch Insert/Replace/Save operations.
	filter        bool           // Filter data and where key-value pairs according to the fields of the table.
	distinct      string         // Force the query to only return distinct results.
	lockInfo      string         // Lock for update or in shared lock.
	cacheEnabled  bool           // Enable sql result cache feature.
	cacheDuration time.Duration  // Cache TTL duration.
	cacheName     string         // Cache name for custom operation.
	unscoped      bool           // Disables soft deleting features when select/delete operations.
	safe          bool           // If true, it clones and returns a new model object whenever operation done; or else it changes the attribute of current model.
	onDuplicate   interface{}    // onDuplicate is used for ON "DUPLICATE KEY UPDATE" statement.
	onDuplicateEx interface{}    // onDuplicateEx is used for excluding some columns ON "DUPLICATE KEY UPDATE" statement.
}

// whereHolder is the holder for where condition preparing.
type whereHolder struct {
	operator int           // Operator for this holder.
	where    interface{}   // Where parameter.
	args     []interface{} // Arguments for where parameter.
}

const (
	OptionOmitEmpty  = 1
	OptionAllowEmpty = 2
	linkTypeMaster   = 1
	linkTypeSlave    = 2
	whereHolderWhere = 1
	whereHolderAnd   = 2
	whereHolderOr    = 3
)

// Table is alias of Core.Model.
// See Core.Model.
// Deprecated, use Model instead.
func (c *Core) Table(tableNameQueryOrStruct ...interface{}) *Model {
	return c.db.Model(tableNameQueryOrStruct...)
}

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
		tableStr   string
		tableName  string
		extraArgs  []interface{}
		tableNames = make([]string, len(tableNameQueryOrStruct))
	)
	// Model creation with sub-query.
	if len(tableNameQueryOrStruct) > 1 {
		conditionStr := gconv.String(tableNameQueryOrStruct[0])
		if gstr.Contains(conditionStr, "?") {
			tableStr, extraArgs = formatWhere(
				c.db, conditionStr, tableNameQueryOrStruct[1:], false,
			)
		}
	}
	// Normal model creation.
	if tableStr == "" {
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
	return &Model{
		db:         c.db,
		tablesInit: tableStr,
		tables:     tableStr,
		fields:     "*",
		start:      -1,
		offset:     -1,
		option:     OptionAllowEmpty,
		filter:     true,
		extraArgs:  extraArgs,
	}
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

// With creates and returns an ORM model based on meta data of given object.
func (c *Core) With(objects ...interface{}) *Model {
	return c.db.Model().With(objects...)
}

// Table is alias of tx.Model.
// Deprecated, use Model instead.
func (tx *TX) Table(tableNameQueryOrStruct ...interface{}) *Model {
	return tx.Model(tableNameQueryOrStruct...)
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
		if gstr.Contains(model.tables, split) {
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
	*newModel = *m
	// Shallow copy slice attributes.
	if n := len(m.extraArgs); n > 0 {
		newModel.extraArgs = make([]interface{}, n)
		copy(newModel.extraArgs, m.extraArgs)
	}
	if n := len(m.whereHolder); n > 0 {
		newModel.whereHolder = make([]*whereHolder, n)
		copy(newModel.whereHolder, m.whereHolder)
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
