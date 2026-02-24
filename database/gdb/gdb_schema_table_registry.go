// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "sync"

// tableRegistryKey identifies a table within a specific database group and schema.
// All three dimensions are required to avoid any cross-schema or cross-group confusion.
type tableRegistryKey struct {
	group  string
	schema string
	table  string
}

// tableRegistry is the single source of truth for all schema metadata.
// It replaces innerMemCache for table fields and table name lookups.
//
// Addressing is 3D: (group, schema, table), which eliminates schema-confusion bugs
// that existed when different schemas or database groups used the same cache key.
//
// Map value semantics:
//   - key absent: table not registered
//   - nil value:  table registered as an existence marker (fields not yet loaded)
//   - non-nil:    table fields have been loaded from the database
type tableRegistry struct {
	mu   sync.RWMutex
	data map[tableRegistryKey]map[string]*TableField
}

func newTableRegistry() *tableRegistry {
	return &tableRegistry{
		data: make(map[tableRegistryKey]map[string]*TableField),
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

// Sets marks multiple tables as known without loading their fields.
// This is more efficient than calling SetIfNotExist in a loop as it acquires the lock only once.
// If a table is already registered (with or without fields), it is skipped.
func (r *tableRegistry) Sets(group, schema string, tables []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, table := range tables {
		key := tableRegistryKey{group, schema, table}
		if _, ok := r.data[key]; !ok {
			r.data[key] = nil
		}
	}
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
// This allows batch operations on the registry with a single lock acquisition.
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

// GetOrSet returns field data for the specified table, invoking loader to populate
// the registry on a cache miss. Uses double-checked locking so that:
//   - concurrent reads on already-loaded tables never block each other (RLock),
//   - only one goroutine executes loader per table on cold start (Lock),
//   - unrelated tables' reads are not affected after the lock is released.
//
// A nil return from loader is stored as an empty map so that subsequent calls
// do not re-invoke loader (distinguishes "loaded with no fields" from "not loaded").
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

// ClearAll removes every entry from the registry.
func (r *tableRegistry) ClearAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = make(map[tableRegistryKey]map[string]*TableField)
}
