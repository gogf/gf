// =================================================================================
// This file is auto-generated by the GoFrame CLI tool. You may modify it as needed.
// =================================================================================

// I am not overwritten.

package dao

import (
	"/internal"
)

// internalUser2Dao is an internal type for wrapping the internal DAO implementation.
type internalUser2Dao = *internal.User2Dao

// user2Dao is the data access object for the table user2.
// You can define custom methods on it to extend its functionality as needed.
type user2Dao struct {
	internalUser2Dao
}

var (
	// User2 is a globally accessible object for table user2 operations.
	User2 = user2Dao{
		internal.NewUser2Dao(),
	}
)

// Add your custom methods and functionality below.
