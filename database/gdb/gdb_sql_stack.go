// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "github.com/gogf/gf/v2/container/garray"

const (
	FirstSqlStackIndex = 0 // the first sql index
)

// SqlStack is a sql stack object.
type SqlStack struct {
	Stacks  *garray.StrArray
	MaxRows int
}

// NewSqlStack  creates and returns a SqlStack.
func NewSqlStack(maxRows ...int) *SqlStack {
	var maxRow int
	if len(maxRows) > 0 {
		maxRow = maxRows[0]
	} else {
		maxRow = DefaultMaxSqlStackRow
	}
	return &SqlStack{
		Stacks:  garray.NewStrArray(true),
		MaxRows: maxRow,
	}
}
func (s *SqlStack) SetMaxRows(rows int) {
	s.MaxRows = rows
}

// Append  add a sql and returns a garray.StrArray.
func (s *SqlStack) Append(sql string) *garray.StrArray {
	if s.MaxRows > 0 {
		if s.Stacks.Len() > s.MaxRows {
			s.Stacks.Clear()
		}
	}
	return s.Stacks.Append(sql)
}

// GetByIndex returns the value by the specified index.
// If the given `index` is out of range of the array, return empty string.
func (s *SqlStack) GetByIndex(idx int) string {
	if s.Stacks.IsEmpty() || idx < 0 || s.Stacks.Len() < idx {
		return ""
	}
	sql, found := s.Stacks.Get(idx)
	if !found {
		return ""
	}
	return sql
}

// Last returns the last sql .
func (s *SqlStack) Last() string {
	idx := s.Stacks.Len() - 1
	return s.GetByIndex(idx)
}

// First returns the first sql .
func (s *SqlStack) First() string {
	return s.GetByIndex(FirstSqlStackIndex)
}
func (s *SqlStack) All() *garray.StrArray {
	return s.Stacks
}
