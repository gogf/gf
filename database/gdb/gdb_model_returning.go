// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
)

// Returning sets the RETURNING clause, supports method chaining.
// Used to specify fields to return after INSERT/UPDATE/DELETE operations.
func (m *Model) Returning(fields ...string) *Model {
	model := m.getModel()
	model.returningFields = fields
	model.returningAll = false
	model.returningExcept = nil
	return model
}

// ReturningAll returns all fields.
// Used to return all fields of the table.
func (m *Model) ReturningAll() *Model {
	model := m.getModel()
	model.returningAll = true
	model.returningFields = nil
	model.returningExcept = nil
	return model
}

// ReturningExcept returns all fields except the specified ones.
// Used to return all other fields except the specified fields.
func (m *Model) ReturningExcept(fields ...string) *Model {
	model := m.getModel()
	model.returningAll = true
	model.returningFields = nil
	model.returningExcept = fields
	return model
}

// hasReturning checks if RETURNING is set.
func (m *Model) hasReturning() bool {
	return len(m.returningFields) > 0 || m.returningAll
}

// buildReturningClause builds the RETURNING clause.
func (m *Model) buildReturningClause(ctx context.Context) (string, error) {
	if !m.hasReturning() {
		return "", nil
	}

	if m.returningAll {
		tableFields, err := m.db.TableFields(ctx, m.tables)
		if err != nil {
			return "", err
		}

		var fields []string
		for fieldName := range tableFields {
			if !gstr.InArray(m.returningExcept, fieldName) {
				fields = append(fields, fmt.Sprintf(`"%s"`, fieldName))
			}
		}
		if len(fields) == 0 {
			return "", nil
		}
		return " RETURNING " + strings.Join(fields, ", "), nil
	}
	if len(m.returningFields) == 0 {
		return "", nil
	}

	var quotedFields []string
	for _, field := range m.returningFields {
		quotedFields = append(quotedFields, fmt.Sprintf(`"%s"`, field))
	}
	return " RETURNING " + strings.Join(quotedFields, ", "), nil
}

// getReturningFields gets the list of RETURNING fields.
func (m *Model) getReturningFields(ctx context.Context) ([]string, error) {
	if !m.hasReturning() {
		return nil, nil
	}

	if m.returningAll {
		tableFields, err := m.db.TableFields(ctx, m.tables)
		if err != nil {
			return nil, err
		}

		var fields []string
		for fieldName := range tableFields {
			if !gstr.InArray(m.returningExcept, fieldName) {
				fields = append(fields, fieldName)
			}
		}
		return fields, nil
	}
	return m.returningFields, nil
}
