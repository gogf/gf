// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"/internal"
)

// internalTableUserDao is internal type for wrapping internal DAO implements.
type internalTableUserDao = *internal.TableUserDao

// tableUserDao is the data access object for table table_user.
// You can define custom methods on it to extend its functionality as you wish.
type tableUserDao struct {
	internalTableUserDao
}

var (
	// TableUser is globally public accessible object for table table_user operations.
	TableUser = tableUserDao{
		internal.NewTableUserDao(),
	}
)

// Fill with you ideas below.
