// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// DBConfig represents database configuration
type DBConfig struct {
	Username string
	Password string
	Host     string
	DBName   string
}

// QueryAndScan executes a query and scans the results into a slice of struct pointers
// Parameters:
//   - query: SQL query string
//   - args: Query arguments
//   - dest: Pointer to slice of struct pointers where results will be stored
//
// Returns error if any occurs during the process
func QueryAndScan(config DBConfig, query string, args []interface{}, dest interface{}) error {
	// Validate input parameters
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to slice")
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get column names: %v", err)
	}

	// Get the type of slice elements
	sliceType := destValue.Elem().Type()
	elementType := sliceType.Elem()
	if elementType.Kind() == reflect.Ptr {
		elementType = elementType.Elem()
	}

	// Create a map of field names to struct fields
	fieldMap := make(map[string]int)
	for i := 0; i < elementType.NumField(); i++ {
		field := elementType.Field(i)
		// Check orm tag first, then json tag, then field name
		tagName := field.Tag.Get("orm")
		if tagName == "" {
			tagName = field.Tag.Get("json")
		}
		if tagName == "" {
			tagName = strings.ToLower(field.Name)
		}
		fieldMap[tagName] = i
	}

	// Prepare slice to store results
	sliceValue := destValue.Elem()

	// Scan rows
	for rows.Next() {
		// Create a new struct instance
		newElem := reflect.New(elementType)

		// Create scan destinations that point directly to struct fields
		scanDest := make([]interface{}, len(columns))
		for i, colName := range columns {
			if fieldIndex, ok := fieldMap[colName]; ok {
				field := newElem.Elem().Field(fieldIndex)
				if field.CanAddr() {
					scanDest[i] = field.Addr().Interface()
				} else {
					// For fields that can't be addressed, use a temporary variable
					var v interface{}
					scanDest[i] = &v
				}
			} else {
				// Column doesn't map to any field, use a placeholder
				var v interface{}
				scanDest[i] = &v
			}
		}

		// Scan the row directly into struct fields
		if err := rows.Scan(scanDest...); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		// Append the new element to the result slice
		if sliceType.Elem().Kind() == reflect.Ptr {
			sliceValue.Set(reflect.Append(sliceValue, newElem))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, newElem.Elem()))
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %v", err)
	}

	return nil
}

func Test_Issue4086_2(t *testing.T) {
	config := DBConfig{
		Username: "root",
		Password: "12345678",
		Host:     "127.0.0.1",
		DBName:   "test1",
	}
	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64    `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64  `json:"recommendIds" orm:"recommend_ids"`
			Photos       []string `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam

		err := QueryAndScan(config, "SELECT * FROM issue4086", nil, &proxyParamList)
		fmt.Println(err)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: []int64{584, 585},
				Photos:       nil,
			},
			{
				ProxyId:      2,
				RecommendIds: []int64{},
				Photos:       nil,
			},
		})
	})
}
