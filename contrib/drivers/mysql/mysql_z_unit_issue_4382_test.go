// Test case for MySQL key length issue #4382
package mysql_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4382_KeyLengthLimit tests the MySQL key length limitation issue
// This test reproduces the issue reported in #4382 where upgrading from GoFrame 2.6 to 2.9
// causes "Specified key was too long; max key length is 1000 bytes" error
func Test_Issue4382_KeyLengthLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// This test should not fail due to key length limitations
		// when using proper MySQL configuration
		table := createTable("test_key_length")
		defer dropTable(table)
		
		// Try to create a table with potentially long keys (using utf8mb4)
		// This scenario could trigger the key length issue
		longTableSQL := `
		CREATE TABLE test_long_keys (
			id INT PRIMARY KEY AUTO_INCREMENT,
			long_field_1 VARCHAR(255) CHARACTER SET utf8mb4,
			long_field_2 VARCHAR(255) CHARACTER SET utf8mb4,
			KEY idx_long_composite (long_field_1, long_field_2)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		_, err := db.Exec(ctx, "DROP TABLE IF EXISTS test_long_keys")
		t.AssertNil(err)
		
		// This should not fail with key length error in GoFrame 2.9
		// Our DoFilter enhancement should automatically add ROW_FORMAT=DYNAMIC
		_, err = db.Exec(ctx, longTableSQL)
		if err != nil {
			// If we get the specific key length error, this confirms the issue
			// With our fix, this should not happen
			t.Logf("Error creating table: %v", err)
		}
		t.AssertNil(err)
		
		// Clean up
		db.Exec(ctx, "DROP TABLE IF EXISTS test_long_keys")
	})
}