// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"github.com/gogf/gf/text/gregex"
	"time"

	"github.com/gogf/gf/text/gstr"
)

// Model is the DAO for ORM.
type Model struct {
	db            DB             // Underlying DB interface.
	tx            *TX            // Underlying TX interface.
	schema        string         // Custom database schema.
	linkType      int            // Mark for operation on master or slave.
	tablesInit    string         // Table names when model initialization.
	tables        string         // Operation table names, which can be more than one table names and aliases, like: "user", "user u", "user u, user_detail ud".
	fields        string         // Operation fields, multiple fields joined using char ','.
	fieldsEx      string         // Excluded operation fields, multiple fields joined using char ','.
	extraArgs     []interface{}  // Extra custom arguments for sql.
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
	lockInfo      string         // Lock for update or in shared lock.
	cacheEnabled  bool           // Enable sql result cache feature.
	cacheDuration time.Duration  // Cache TTL duration.
	cacheName     string         // Cache name for custom operation.
	unscoped      bool           // Disables soft deleting features when select/delete operations.
	safe          bool           // If true, it clones and returns a new model object whenever operation done; or else it changes the attribute of current model.
}

// whereHolder is the holder for where condition preparing.
type whereHolder struct {
	operator int           // Operator for this holder.
	where    interface{}   // Where parameter.
	args     []interface{} // Arguments for where parameter.
}

const (
	OPTION_OMITEMPTY  = 1
	OPTION_ALLOWEMPTY = 2

	linkTypeMaster   = 1
	linkTypeSlave    = 2
	whereHolderWhere = 1
	whereHolderAnd   = 2
	whereHolderOr    = 3
)

// Table creates and returns a new ORM model from given schema.
// The parameter <table> can be more than one table names, and also alias name, like:
// 1. Table names:
//    Table("user")
//    Table("user u")
//    Table("user, user_detail")
//    Table("user u, user_detail ud")
// 2. Table name with alias: Table("user", "u")
func (c *Core) Table(table ...string) *Model {
	tables := ""
	if len(table) > 1 {
		tables = fmt.Sprintf(
			`%s AS %s`, c.DB.QuotePrefixTableName(table[0]), c.DB.QuoteWord(table[1]),
		)
	} else if len(table) == 1 {
		tables = c.DB.QuotePrefixTableName(table[0])
	} else {
		panic("table cannot be empty")
	}
	return &Model{
		db:         c.DB,
		tablesInit: tables,
		tables:     tables,
		fields:     "*",
		start:      -1,
		offset:     -1,
		option:     OPTION_ALLOWEMPTY,
	}
}

// Model is alias of Core.Table.
// See Core.Table.
func (c *Core) Model(table ...string) *Model {
	return c.DB.Table(table...)
}

// Table acts like Core.Table except it operates on transaction.
// See Core.Table.
func (tx *TX) Table(table ...string) *Model {
	model := tx.db.Table(table...)
	model.db = tx.db
	model.tx = tx
	return model
}

// Model is alias of tx.Table.
// See tx.Table.
func (tx *TX) Model(table ...string) *Model {
	return tx.Table(table...)
}

// Ctx sets the context for current operation.
func (m *Model) Ctx(ctx context.Context) *Model {
	if ctx == nil {
		return m
	}
	model := m.getModel()
	model.db = model.db.Ctx(ctx)
	return model
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
		newModel = m.tx.Table(m.tablesInit)
	} else {
		newModel = m.db.Table(m.tablesInit)
	}
	*newModel = *m
	// Deep copy slice attributes.
	if n := len(m.extraArgs); n > 0 {
		newModel.extraArgs = make([]interface{}, n)
		copy(newModel.extraArgs, m.extraArgs)
	}
	if n := len(m.whereHolder); n > 0 {
		newModel.whereHolder = make([]*whereHolder, n)
		copy(newModel.whereHolder, m.whereHolder)
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
