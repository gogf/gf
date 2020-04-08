// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
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
	force         bool           // Force select/delete without soft operation features.
	safe          bool           // If true, it clones and returns a new model object whenever operation done; or else it changes the attribute of current model.
}

// whereHolder is the holder for where condition preparing.
type whereHolder struct {
	operator int           // Operator for this holder.
	where    interface{}   // Where parameter.
	args     []interface{} // Arguments for where parameter.
}

const (
	gLINK_TYPE_MASTER   = 1
	gLINK_TYPE_SLAVE    = 2
	gWHERE_HOLDER_WHERE = 1
	gWHERE_HOLDER_AND   = 2
	gWHERE_HOLDER_OR    = 3
	OPTION_OMITEMPTY    = 1 << iota
	OPTION_ALLOWEMPTY
)

// Table creates and returns a new ORM model from given schema.
// The parameter <tables> can be more than one table names, like :
// "user", "user u", "user, user_detail", "user u, user_detail ud"
func (c *Core) Table(table string) *Model {
	table = c.DB.QuotePrefixTableName(table)
	return &Model{
		db:         c.DB,
		tablesInit: table,
		tables:     table,
		fields:     "*",
		start:      -1,
		offset:     -1,
		option:     OPTION_ALLOWEMPTY,
	}
}

// Model is alias of Core.Table.
// See Core.Table.
func (c *Core) Model(table string) *Model {
	return c.DB.Table(table)
}

// Table acts like Core.Table except it operates on transaction.
// See Core.Table.
func (tx *TX) Table(table string) *Model {
	table = tx.db.QuotePrefixTableName(table)
	return &Model{
		db:         tx.db,
		tx:         tx,
		tablesInit: table,
		tables:     table,
		fields:     "*",
		start:      -1,
		offset:     -1,
		option:     OPTION_ALLOWEMPTY,
	}
}

// Model is alias of tx.Table.
// See tx.Table.
func (tx *TX) Model(table string) *Model {
	return tx.Table(table)
}

// As sets an alias name for current table.
func (m *Model) As(as string) *Model {
	if m.tables != "" {
		model := m.getModel()
		model.tables = gstr.TrimRight(model.tables) + " AS " + as
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
	model.linkType = gLINK_TYPE_MASTER
	return model
}

// Slave marks the following operation on slave node.
// Note that it makes sense only if there's any slave node configured.
func (m *Model) Slave() *Model {
	model := m.getModel()
	model.linkType = gLINK_TYPE_SLAVE
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
