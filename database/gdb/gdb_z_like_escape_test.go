package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func ExampleEscapeLikeString() {
	// Escape special characters in user input for LIKE operations
	userInput := "test%_data\\with_special"
	escaped := gdb.EscapeLikeString(userInput)

	// The escaped string can now be safely used in LIKE patterns
	// Original: "test%_data\with_special"
	// Escaped:  "test\\%\\_data\\\\with\\_special"

	// Usage example with WhereLike
	// db.Model("table").WhereLike("column", "%"+escaped+"%")

	// Output: test\%\_data\\with\_special
	print(escaped)
}

func Test_EscapeLikeString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test escaping backslash
		t.Assert(gdb.EscapeLikeString("test\\data"), "test\\\\data")

		// Test escaping percent
		t.Assert(gdb.EscapeLikeString("test%data"), "test\\%data")

		// Test escaping underscore
		t.Assert(gdb.EscapeLikeString("test_data"), "test\\_data")

		// Test escaping all special characters
		t.Assert(gdb.EscapeLikeString("test\\%_data"), "test\\\\\\%\\_data")

		// Test empty string
		t.Assert(gdb.EscapeLikeString(""), "")

		// Test string with no special characters
		t.Assert(gdb.EscapeLikeString("normal"), "normal")
	})
}
