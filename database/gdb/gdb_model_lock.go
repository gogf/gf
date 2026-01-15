// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// Lock clause constants for different databases.
// These constants provide type-safe and IDE-friendly access to various lock syntaxes.
const (
	// Common lock clauses (supported by most databases)
	LockForUpdate           = "FOR UPDATE"
	LockForUpdateSkipLocked = "FOR UPDATE SKIP LOCKED"

	// MySQL lock clauses
	LockInShareMode     = "LOCK IN SHARE MODE" // MySQL legacy syntax
	LockForShare        = "FOR SHARE"          // MySQL 8.0+ and PostgreSQL
	LockForUpdateNowait = "FOR UPDATE NOWAIT"  // MySQL 8.0+ and Oracle

	// PostgreSQL specific lock clauses
	LockForNoKeyUpdate           = "FOR NO KEY UPDATE"
	LockForKeyShare              = "FOR KEY SHARE"
	LockForShareNowait           = "FOR SHARE NOWAIT"
	LockForShareSkipLocked       = "FOR SHARE SKIP LOCKED"
	LockForNoKeyUpdateNowait     = "FOR NO KEY UPDATE NOWAIT"
	LockForNoKeyUpdateSkipLocked = "FOR NO KEY UPDATE SKIP LOCKED"
	LockForKeyShareNowait        = "FOR KEY SHARE NOWAIT"
	LockForKeyShareSkipLocked    = "FOR KEY SHARE SKIP LOCKED"

	// Oracle specific lock clauses
	LockForUpdateWait5  = "FOR UPDATE WAIT 5"
	LockForUpdateWait10 = "FOR UPDATE WAIT 10"
	LockForUpdateWait30 = "FOR UPDATE WAIT 30"

	// SQL Server lock hints (use with WITH clause)
	LockWithUpdLock         = "WITH (UPDLOCK)"
	LockWithHoldLock        = "WITH (HOLDLOCK)"
	LockWithXLock           = "WITH (XLOCK)"
	LockWithTabLock         = "WITH (TABLOCK)"
	LockWithNoLock          = "WITH (NOLOCK)"
	LockWithUpdLockHoldLock = "WITH (UPDLOCK, HOLDLOCK)"
)

// Lock sets a custom lock clause for the current operation.
// This is a generic method that allows you to specify any lock syntax supported by your database.
// You can use predefined constants or custom strings.
//
// Database-specific lock syntax support:
//
// PostgreSQL (most comprehensive):
//   - "FOR UPDATE"                    - Exclusive lock, blocks all access
//   - "FOR NO KEY UPDATE"            - Weaker exclusive lock, doesn't block FOR KEY SHARE
//   - "FOR SHARE"                    - Shared lock, allows reads but blocks writes
//   - "FOR KEY SHARE"                - Weakest lock, only locks key values
//   - All above can be combined with:
//   - "NOWAIT"                     - Return immediately if lock cannot be acquired
//   - "SKIP LOCKED"               - Skip locked rows instead of waiting
//
// MySQL:
//   - "FOR UPDATE"                    - Exclusive lock (all versions)
//   - "LOCK IN SHARE MODE"           - Shared lock (legacy syntax)
//   - "FOR SHARE"                    - Shared lock (MySQL 8.0+)
//   - "FOR UPDATE NOWAIT"            - MySQL 8.0+ only
//   - "FOR UPDATE SKIP LOCKED"       - MySQL 8.0+ only
//
// Oracle:
//   - "FOR UPDATE"                    - Exclusive lock
//   - "FOR UPDATE NOWAIT"            - Exclusive lock, no wait
//   - "FOR UPDATE SKIP LOCKED"       - Exclusive lock, skip locked rows
//   - "FOR UPDATE WAIT n"            - Exclusive lock, wait n seconds
//   - "FOR UPDATE OF column_list"    - Lock specific columns
//
// SQL Server (uses WITH hints):
//   - "WITH (UPDLOCK)"               - Update lock
//   - "WITH (HOLDLOCK)"              - Hold lock until transaction end
//   - "WITH (XLOCK)"                 - Exclusive lock
//   - "WITH (TABLOCK)"               - Table lock
//   - "WITH (NOLOCK)"                - No lock (dirty read)
//   - "WITH (UPDLOCK, HOLDLOCK)"     - Combined update and hold lock
//
// SQLite:
//   - Limited locking support, database-level locks only
//   - No row-level lock syntax supported
//
// Usage examples:
//
//	db.Model("users").Lock("FOR UPDATE NOWAIT").Where("id", 1).One()
//	db.Model("users").Lock("FOR SHARE SKIP LOCKED").Where("status", "active").All()
//	db.Model("users").Lock("WITH (UPDLOCK)").Where("id", 1).One() // SQL Server
//	db.Model("users").Lock("FOR UPDATE OF name, email").Where("id", 1).One() // Oracle
//	db.Model("users").Lock("FOR UPDATE WAIT 15").Where("id", 1).One() // Oracle custom wait
//
// Or use predefined constants for better IDE support:
//
//	db.Model("users").Lock(gdb.LockForUpdateNowait).Where("id", 1).One()
//	db.Model("users").Lock(gdb.LockForShareSkipLocked).Where("status", "active").All()
func (m *Model) Lock(lockClause string) *Model {
	model := m.getModel()
	model.lockInfo = lockClause
	return model
}

// LockUpdate sets the lock for update for current operation.
// This is equivalent to Lock("FOR UPDATE").
func (m *Model) LockUpdate() *Model {
	model := m.getModel()
	model.lockInfo = LockForUpdate
	return model
}

// LockUpdateSkipLocked sets the lock for update for current operation.
// It skips the locked rows.
// This is equivalent to Lock("FOR UPDATE SKIP LOCKED").
// Note: Supported by PostgreSQL, Oracle, and MySQL 8.0+.
func (m *Model) LockUpdateSkipLocked() *Model {
	model := m.getModel()
	model.lockInfo = LockForUpdateSkipLocked
	return model
}

// LockShared sets the lock in share mode for current operation.
// This is equivalent to Lock("LOCK IN SHARE MODE") for MySQL or Lock("FOR SHARE") for PostgreSQL.
// Note: For maximum compatibility, this uses MySQL's legacy syntax.
func (m *Model) LockShared() *Model {
	model := m.getModel()
	model.lockInfo = LockInShareMode
	return model
}
