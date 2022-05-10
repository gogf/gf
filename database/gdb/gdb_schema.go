// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// Schema is a schema object from which it can then create a Model.
type Schema struct {
	DB
}

// Schema creates and returns a schema.
func (c *Core) Schema(schema string) *Schema {
	// Do not change the schema of the original db,
	// it here creates a new db and changes its schema.
	db, err := NewByGroup(c.GetGroup())
	if err != nil {
		panic(err)
	}
	core := db.GetCore()
	// Different schema share some same objects.
	core.logger = c.logger
	core.cache = c.cache
	core.schema = schema
	return &Schema{
		DB: db,
	}
}
