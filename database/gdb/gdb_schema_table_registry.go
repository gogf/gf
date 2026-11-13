// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "sync"

// tableRegistryKey identifies a table within database group and schema.
type tableRegistryKey struct {
	group  string
	schema string
	table  string
}

// schemaKey identifies a (group, schema) pair for tracking loaded table lists.
type schemaKey struct {
	group  string
	schema string
}

// tableRegistry is the single source of truth for all schema metadata.
// Uses 3D addressing (group, schema, table) to eliminate schema-confusion bugs.
//
// Map value semantics:
//   - key absent: table not registered
//   - nil value:  table registered as existence marker
//   - non-nil:    table fields loaded from database
type tableRegistry struct {
	mu            sync.RWMutex
	data          map[tableRegistryKey]map[string]*TableField
	loadedSchemas map[schemaKey]struct{}
}

func newTableRegistry() *tableRegistry {
	return &tableRegistry{
		data:          make(map[tableRegistryKey]map[string]*TableField),
		loadedSchemas: make(map[schemaKey]struct{}),
	}
}

// Get retrieves the field map for the given table.
// Returns (fields, true) if the key exists (fields may still be nil for existence-only entries).
// Returns (nil, false) if the key is not present.
func (r *tableRegistry) Get(group, schema, table string) (map[string]*TableField, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fields, ok := r.data[tableRegistryKey{group, schema, table}]
	return fields, ok
}

// Set stores field data for the specified table, overwriting any existing entry.
// Passing nil fields registers the table as an existence marker without field data.
func (r *tableRegistry) Set(group, schema, table string, fields map[string]*TableField) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[tableRegistryKey{group, schema, table}] = fields
}

// SetIfNotExist marks the table as known without loading its fields.
// If the table is already registered (with or without fields), this is a no-op.
// It returns true if the table was newly registered, false if it already existed.
func (r *tableRegistry) SetIfNotExist(group, schema, table string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := tableRegistryKey{group, schema, table}
	if _, ok := r.data[key]; !ok {
		r.data[key] = nil
		return true
	}
	return false
}

// Sets registers multiple tables and marks schema as fully loaded.
func (r *tableRegistry) Sets(group, schema string, tables []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, table := range tables {
		key := tableRegistryKey{group, schema, table}
		if _, ok := r.data[key]; !ok {
			r.data[key] = nil
		}
	}
	r.loadedSchemas[schemaKey{group, schema}] = struct{}{}
}

// GetLoadedSchemaTables returns all registered table names for given group/schema
// and reports whether the full table list has been loaded from database.
func (r *tableRegistry) GetLoadedSchemaTables(group, schema string) (tables []string, loaded bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.loadedSchemas[schemaKey{group, schema}]; !ok {
		return nil, false
	}
	for key := range r.data {
		if key.group == group && key.schema == schema {
			tables = append(tables, key.table)
		}
	}
	return tables, true
}

// LockFunc executes callback function with write lock for batch operations.
func (r *tableRegistry) LockFunc(f func(data map[tableRegistryKey]map[string]*TableField)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	f(r.data)
}

// HasTable reports whether the specified table is registered (O(1)).
func (r *tableRegistry) HasTable(group, schema, table string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.data[tableRegistryKey{group, schema, table}]
	return ok
}

// Tables returns all registered table names for the given group and schema.
func (r *tableRegistry) Tables(group, schema string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var tables []string
	for key := range r.data {
		if key.group == group && key.schema == schema {
			tables = append(tables, key.table)
		}
	}
	return tables
}

// GetOrSet returns field data for specified table with cache support.
// Uses double-checked locking for concurrent safety.
func (r *tableRegistry) GetOrSet(
	group, schema, table string,
	loader func() (map[string]*TableField, error),
) (map[string]*TableField, error) {
	// Fast path: fields already loaded.
	r.mu.RLock()
	entry, ok := r.data[tableRegistryKey{group, schema, table}]
	r.mu.RUnlock()
	if ok && entry != nil {
		return entry, nil
	}

	// Slow path: acquire write lock and load.
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check under write lock.
	entry, ok = r.data[tableRegistryKey{group, schema, table}]
	if ok && entry != nil {
		return entry, nil
	}

	fields, err := loader()
	if err != nil {
		return nil, err
	}
	if fields == nil {
		fields = make(map[string]*TableField) // empty map marks "loaded, no fields found"
	}
	r.data[tableRegistryKey{group, schema, table}] = fields
	return fields, nil
}

// Delete removes the registry entry for the specified table.
func (r *tableRegistry) Delete(group, schema, table string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, tableRegistryKey{group, schema, table})
}

// ClearAll removes all entries and resets loaded schema markers.
func (r *tableRegistry) ClearAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = make(map[tableRegistryKey]map[string]*TableField)
	r.loadedSchemas = make(map[schemaKey]struct{})
}
