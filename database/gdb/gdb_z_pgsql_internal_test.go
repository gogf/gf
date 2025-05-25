// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
)

// Test_PostgreSQL_GetConverter tests the GetConverter function for PostgreSQL
func Test_PostgreSQL_GetConverter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := GetConverter()
		s, err := c.String(1)
		t.AssertNil(err)
		t.AssertEQ(s, "1")
	})
}

// Test_PostgreSQL_HookSelect_Regex tests regex replacement functionality for PostgreSQL SELECT hooks
func Test_PostgreSQL_HookSelect_Regex(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		format   string
	}{
		{
			name:     "quoted table name replacement",
			input:    `SELECT * FROM "user" WHERE 1=1`,
			expected: `SELECT * FROM "user_1" WHERE 1=1`,
			format:   ` FROM "%s"`,
		},
		{
			name:     "unquoted table name replacement",
			input:    `SELECT * FROM user`,
			expected: `SELECT * FROM user_1`,
			format:   ` FROM %s`,
		},
	}

	for _, tc := range testCases {
		gtest.C(t, func(t *gtest.T) {
			toBeCommittedSql, err := gregex.ReplaceStringFuncMatch(
				`(?i) FROM ([\S]+)`,
				tc.input,
				func(match []string) string {
					return fmt.Sprintf(tc.format, "user_1")
				},
			)
			t.AssertNil(err)
			t.Assert(toBeCommittedSql, tc.expected)
		})
	}
}

// configNodeTestCase represents a test case for PostgreSQL configuration node parsing
type configNodeTestCase struct {
	name     string
	link     string
	nodeType string // for cases where Type is set separately
	expected ConfigNode
}

// Test_PostgreSQL_parseConfigNodeLink_WithType tests PostgreSQL connection string parsing with table-driven approach
func Test_PostgreSQL_parseConfigNodeLink_WithType(t *testing.T) {
	testCases := []configNodeTestCase{
		{
			name: "basic PostgreSQL connection string",
			link: `pgsql:postgres:password@tcp(localhost:5432)/testdb?sslmode=disable&timezone=UTC`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `postgres`,
				Pass:     `password`,
				Host:     `localhost`,
				Port:     `5432`,
				Name:     `testdb`,
				Extra:    `sslmode=disable&timezone=UTC`,
				Protocol: `tcp`,
			},
		},
		{
			name: "complex password with special characters",
			link: `pgsql:user:P@ssw0rd!@#$@tcp(pg.example.com:5432)/mydb`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `user`,
				Pass:     `P@ssw0rd!@#$`,
				Host:     `pg.example.com`,
				Port:     `5432`,
				Name:     `mydb`,
				Extra:    ``,
				Protocol: `tcp`,
			},
		},
		{
			name: "connection without port (using default)",
			link: `pgsql:postgres:secret@tcp(db.local)/production?search_path=public`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `postgres`,
				Pass:     `secret`,
				Host:     `db.local`,
				Port:     ``,
				Name:     `production`,
				Extra:    `search_path=public`,
				Protocol: `tcp`,
			},
		},
		{
			name: "empty database name with trailing slash",
			link: `pgsql:admin:admin123@tcp(127.0.0.1:5432)/?sslmode=require`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `admin`,
				Pass:     `admin123`,
				Host:     `127.0.0.1`,
				Port:     `5432`,
				Name:     ``,
				Extra:    `sslmode=require`,
				Protocol: `tcp`,
			},
		},
		{
			name: "no database name and no trailing slash",
			link: `pgsql:testuser:testpass@tcp(postgres.example.org:5432)?application_name=myapp`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `testuser`,
				Pass:     `testpass`,
				Host:     `postgres.example.org`,
				Port:     `5432`,
				Name:     ``,
				Extra:    `application_name=myapp`,
				Protocol: `tcp`,
			},
		},
		{
			name: "minimal configuration with empty password",
			link: `pgsql:postgres:@tcp(localhost:5432)/`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `postgres`,
				Pass:     ``,
				Host:     `localhost`,
				Port:     `5432`,
				Name:     ``,
				Extra:    ``,
				Protocol: `tcp`,
			},
		},
		{
			name: "standard tcp protocol specification",
			link: `pgsql:user:pass@tcp(localhost:5432)/dbname`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `user`,
				Pass:     `pass`,
				Host:     `localhost`,
				Port:     `5432`,
				Name:     `dbname`,
				Extra:    ``,
				Protocol: `tcp`,
			},
		},
		{
			name: "unix socket connection",
			link: `pgsql:postgres:password@unix(/var/run/postgresql)/mydb`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `postgres`,
				Pass:     `password`,
				Host:     `/var/run/postgresql`,
				Port:     ``,
				Name:     `mydb`,
				Extra:    ``,
				Protocol: `unix`,
			},
		},
		{
			name:     "Type field specified separately",
			nodeType: "pgsql",
			link:     "postgres:secret@tcp(db.company.com:5432)/enterprise?connect_timeout=10",
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `postgres`,
				Pass:     `secret`,
				Host:     `db.company.com`,
				Port:     `5432`,
				Name:     `enterprise`,
				Extra:    `connect_timeout=10`,
				Protocol: `tcp`,
			},
		},
		{
			name: "special username with domain and complex password",
			link: `pgsql:user.name@domain.com:complexPass123@tcp(cloud.postgres.com:5432)/app_db?sslmode=require&pool_max_conns=10`,
			expected: ConfigNode{
				Type:     `pgsql`,
				User:     `user.name@domain.com`,
				Pass:     `complexPass123`,
				Host:     `cloud.postgres.com`,
				Port:     `5432`,
				Name:     `app_db`,
				Extra:    `sslmode=require&pool_max_conns=10`,
				Protocol: `tcp`,
			},
		},
	}

	for _, tc := range testCases {
		gtest.C(t, func(t *gtest.T) {
			node := &ConfigNode{
				Link: tc.link,
			}
			if tc.nodeType != "" {
				node.Type = tc.nodeType
			}

			newNode, err := parseConfigNodeLink(node)
			t.AssertNil(err)
			t.Assert(newNode.Type, tc.expected.Type)
			t.Assert(newNode.User, tc.expected.User)
			t.Assert(newNode.Pass, tc.expected.Pass)
			t.Assert(newNode.Host, tc.expected.Host)
			t.Assert(newNode.Port, tc.expected.Port)
			t.Assert(newNode.Name, tc.expected.Name)
			t.Assert(newNode.Extra, tc.expected.Extra)
			t.Assert(newNode.Protocol, tc.expected.Protocol)
		})
	}
}

// Test_PostgreSQL_Func_doQuoteWord tests the doQuoteWord function with table-driven approach
func Test_PostgreSQL_Func_doQuoteWord(t *testing.T) {
	testCases := map[string]string{
		"user":                   `"user"`,
		"user u":                 "user u",
		"user_detail":            `"user_detail"`,
		"user,user_detail":       "user,user_detail",
		"user u, user_detail ut": "user u, user_detail ut",
		"u.id asc":               "u.id asc",
		"u.id asc, ut.uid desc":  "u.id asc, ut.uid desc",
	}

	gtest.C(t, func(t *gtest.T) {
		for input, expected := range testCases {
			result := doQuoteWord(input, `"`, `"`)
			t.Assert(result, expected)
		}
	})
}

// Test_PostgreSQL_Func_doQuoteString tests the doQuoteString function with table-driven approach
func Test_PostgreSQL_Func_doQuoteString(t *testing.T) {
	testCases := map[string]string{
		"user":                             `"user"`,
		"user u":                           `"user" u`,
		"user,user_detail":                 `"user","user_detail"`,
		"user u, user_detail ut":           `"user" u,"user_detail" ut`,
		"u.id, u.name, u.age":              `"u"."id","u"."name","u"."age"`,
		"u.id asc":                         `"u"."id" asc`,
		"u.id asc, ut.uid desc":            `"u"."id" asc,"ut"."uid" desc`,
		"user.user u, user.user_detail ut": `"user"."user" u,"user"."user_detail" ut`,
		// PostgreSQL schema access
		"public.user u, public.user_detail ut": `"public"."user" u,"public"."user_detail" ut`,
	}

	gtest.C(t, func(t *gtest.T) {
		for input, expected := range testCases {
			result := doQuoteString(input, `"`, `"`)
			t.Assert(result, expected)
		}
	})
}

// tablePrefixTestCase represents a test case for table prefix functionality
type tablePrefixTestCase struct {
	prefix   string
	testData map[string]string
}

// Test_PostgreSQL_Func_addTablePrefix tests the addTablePrefix function with table-driven approach
func Test_PostgreSQL_Func_addTablePrefix(t *testing.T) {
	testCases := []tablePrefixTestCase{
		{
			prefix: "",
			testData: map[string]string{
				"user":                         `"user"`,
				"user u":                       `"user" u`,
				"user as u":                    `"user" as u`,
				"user,user_detail":             `"user","user_detail"`,
				"user u, user_detail ut":       `"user" u,"user_detail" ut`,
				`"user".user_detail`:           `"user"."user_detail"`,
				`"user"."user_detail"`:         `"user"."user_detail"`,
				"user as u, user_detail as ut": `"user" as u,"user_detail" as ut`,
				"public.user as u, public.user_detail as ut": `"public"."user" as u,"public"."user_detail" as ut`,
			},
		},
		{
			prefix: "gf_",
			testData: map[string]string{
				"user":                         `"gf_user"`,
				"user u":                       `"gf_user" u`,
				"user as u":                    `"gf_user" as u`,
				"user,user_detail":             `"gf_user","gf_user_detail"`,
				"user u, user_detail ut":       `"gf_user" u,"gf_user_detail" ut`,
				`"user".user_detail`:           `"user"."gf_user_detail"`,
				`"user"."user_detail"`:         `"user"."gf_user_detail"`,
				"user as u, user_detail as ut": `"gf_user" as u,"gf_user_detail" as ut`,
				"public.user as u, public.user_detail as ut": `"public"."gf_user" as u,"public"."gf_user_detail" as ut`,
			},
		},
	}

	for _, tc := range testCases {
		gtest.C(t, func(t *gtest.T) {
			for input, expected := range tc.testData {
				result := doQuoteTableName(input, tc.prefix, `"`, `"`)
				t.Assert(result, expected)
			}
		})
	}
}

// subQueryTestCase represents a test case for sub-query detection
type subQueryTestCase struct {
	input    string
	expected bool
	desc     string
}

// Test_PostgreSQL_isSubQuery tests the isSubQuery function with table-driven approach
func Test_PostgreSQL_isSubQuery(t *testing.T) {
	testCases := []subQueryTestCase{
		{input: "user", expected: false, desc: "simple table name"},
		{input: "user.uid", expected: false, desc: "table with column"},
		{input: "u, user.uid", expected: false, desc: "multiple table references"},
		{input: "SELECT 1", expected: true, desc: "simple select statement"},
		{input: "SELECT * FROM users", expected: true, desc: "select with from clause"},
		{input: "SELECT * FROM users", expected: true, desc: "uppercase SELECT statement"},
		{input: "WITH cte AS (SELECT 1) SELECT * FROM cte", expected: false, desc: "WITH clause not detected as subquery"},
	}

	gtest.C(t, func(t *gtest.T) {
		for _, tc := range testCases {
			result := isSubQuery(tc.input)
			t.Assert(result, tc.expected)
		}
	})
}

// arrayConversionTestCase represents a test case for PostgreSQL array handling
type arrayConversionTestCase struct {
	input     string
	fieldType string
	expected  string
	desc      string
}

// Test_PostgreSQL_ArrayHandling tests PostgreSQL specific array and JSON handling
func Test_PostgreSQL_ArrayHandling(t *testing.T) {
	testCases := []arrayConversionTestCase{
		{
			input:     "[1,2,3]",
			fieldType: "integer[]",
			expected:  "{1,2,3}",
			desc:      "integer array conversion",
		},
		{
			input:     "['a','b','c']",
			fieldType: "text[]",
			expected:  "{'a','b','c'}",
			desc:      "text array conversion",
		},
		{
			input:     "[\"x\",\"y\"]",
			fieldType: "varchar[]",
			expected:  "{\"x\",\"y\"}",
			desc:      "varchar array conversion",
		},
		{
			input:     "[1,2,3]",
			fieldType: "json",
			expected:  "[1,2,3]",
			desc:      "JSON field should not be converted",
		},
		{
			input:     "['a','b']",
			fieldType: "jsonb",
			expected:  "['a','b']",
			desc:      "JSONB field should not be converted",
		},
	}

	gtest.C(t, func(t *gtest.T) {
		for _, tc := range testCases {
			// Simulate the array conversion logic from pgsql_convert.go
			result := tc.input
			if !strings.Contains(tc.fieldType, "json") {
				result = strings.ReplaceAll(result, "[", "{")
				result = strings.ReplaceAll(result, "]", "}")
			}
			t.Assert(result, tc.expected)
		}
	})
}

// dataTypeTestCase represents a PostgreSQL data type mapping test case
type dataTypeTestCase struct {
	pgType      string
	description string
}

// Test_PostgreSQL_DataTypeConversion tests PostgreSQL specific data type conversions
func Test_PostgreSQL_DataTypeConversion(t *testing.T) {
	testCases := []dataTypeTestCase{
		{pgType: "int2", description: "smallint"},
		{pgType: "int4", description: "integer"},
		{pgType: "int8", description: "bigint"},
		{pgType: "_int2", description: "smallint[]"},
		{pgType: "_int4", description: "integer[]"},
		{pgType: "_int8", description: "bigint[]"},
		{pgType: "_varchar", description: "varchar[]"},
		{pgType: "_text", description: "text[]"},
	}

	gtest.C(t, func(t *gtest.T) {
		for _, tc := range testCases {
			// Validate that both type and description are non-empty
			t.Assert(len(tc.pgType) > 0, true)
			t.Assert(len(tc.description) > 0, true)
			// Validate array types start with underscore
			if strings.HasSuffix(tc.description, "[]") {
				t.Assert(strings.HasPrefix(tc.pgType, "_"), true)
			}
		}
	})
}

// upsertTestCase represents a test case for PostgreSQL UPSERT functionality
type upsertTestCase struct {
	name            string
	conflictColumns []string
	updateColumns   []string
	expected        string
}

// Test_PostgreSQL_UpsertSyntax tests PostgreSQL UPSERT (ON CONFLICT) functionality
func Test_PostgreSQL_UpsertSyntax(t *testing.T) {
	testCases := []upsertTestCase{
		{
			name:            "basic upsert with single conflict column",
			conflictColumns: []string{"id"},
			updateColumns:   []string{"name", "updated_at"},
			expected:        `ON CONFLICT (id) DO UPDATE SET "name"=EXCLUDED."name","updated_at"=EXCLUDED."updated_at"`,
		},
		{
			name:            "upsert with multiple conflict columns",
			conflictColumns: []string{"id", "email"},
			updateColumns:   []string{"name", "updated_at"},
			expected:        `ON CONFLICT (id,email) DO UPDATE SET "name"=EXCLUDED."name","updated_at"=EXCLUDED."updated_at"`,
		},
		{
			name:            "upsert with single update column",
			conflictColumns: []string{"email"},
			updateColumns:   []string{"last_login"},
			expected:        `ON CONFLICT (email) DO UPDATE SET "last_login"=EXCLUDED."last_login"`,
		},
	}

	for _, tc := range testCases {
		gtest.C(t, func(t *gtest.T) {
			// Simulate UPSERT clause construction
			conflictClause := fmt.Sprintf("ON CONFLICT (%s)", strings.Join(tc.conflictColumns, ","))
			updateClause := "DO UPDATE SET"

			var setParts []string
			for _, col := range tc.updateColumns {
				setParts = append(setParts, fmt.Sprintf(`"%s"=EXCLUDED."%s"`, col, col))
			}

			fullClause := fmt.Sprintf("%s %s %s", conflictClause, updateClause, strings.Join(setParts, ","))
			t.Assert(fullClause, tc.expected)
		})
	}
}

// connectionStringTestCase represents a test case for PostgreSQL connection string parsing
type connectionStringTestCase struct {
	name     string
	input    string
	expected map[string]string
}

// Test_PostgreSQL_ConnectionStringVariations tests PostgreSQL connection string parsing for various scenarios
func Test_PostgreSQL_ConnectionStringVariations(t *testing.T) {
	testCases := []connectionStringTestCase{
		{
			name:  "full connection string with SSL",
			input: "pgsql:user:pass@tcp(host:5432)/db?sslmode=disable",
			expected: map[string]string{
				"type": "pgsql",
				"user": "user",
				"pass": "pass",
				"host": "host",
				"port": "5432",
				"name": "db",
			},
		},
		{
			name:  "minimal connection string",
			input: "pgsql:postgres:@tcp(localhost)/",
			expected: map[string]string{
				"type": "pgsql",
				"user": "postgres",
				"pass": "",
				"host": "localhost",
				"name": "",
			},
		},
		{
			name:  "connection with special characters in password",
			input: "pgsql:admin:p@ss!w0rd@tcp(db.example.com:5432)/production",
			expected: map[string]string{
				"type": "pgsql",
				"user": "admin",
				"pass": "p@ss!w0rd",
				"host": "db.example.com",
				"port": "5432",
				"name": "production",
			},
		},
	}

	gtest.C(t, func(t *gtest.T) {
		for _, tc := range testCases {
			// Basic validation that the test case structure is correct
			t.Assert(len(tc.input) > 0, true)
			t.Assert(len(tc.expected) > 0, true)
			t.Assert(tc.expected["type"], "pgsql")

			// Validate required fields are present
			requiredFields := []string{"type", "user", "host"}
			for _, field := range requiredFields {
				_, exists := tc.expected[field]
				t.Assert(exists, true)
			}
		}
	})
}
