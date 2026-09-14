// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements tabular result rendering and aggregate hint
// emission shared by every resource subpackage. Render supports optional
// ExtraColumns so resources whose AgentSpec has additional discriminating
// fields (e.g. prompts' SourcePath) can surface them between the AGENT
// and PROJECT PATH columns without forking the renderer.

package common

import (
	"fmt"
	"io"
)

// ColumnSpec describes one extra column to insert into the result table.
// Header is the column header text; Value resolves the cell text for a
// given Result. The renderer inserts every ExtraColumn after the AGENT
// column and before the PROJECT PATH column, in the order supplied.
type ColumnSpec struct {
	// Header is the upper-case column header rendered in the table.
	Header string
	// Value resolves the cell text for a given Result.
	Value func(Result) string
}

// Render writes the result table for a list of results to out. ExtraColumns
// is variadic so callers that do not need additional columns omit it
// entirely; callers that need columns supply them in left-to-right order.
//
// All columns are left-padded to the width of the longest cell content so
// the table aligns cleanly across rows.
func Render(out io.Writer, results []Result, extraColumns ...ColumnSpec) error {
	const (
		columnAgent    = "AGENT"
		columnPath     = "PROJECT PATH"
		columnCategory = "CATEGORY"
		columnStatus   = "STATUS"
		columnDetail   = "DETAIL"
	)
	var (
		maxAgent    = len(columnAgent)
		maxPath     = len(columnPath)
		maxCategory = len(columnCategory)
		maxStatus   = len(columnStatus)
		extraWidths = make([]int, len(extraColumns))
	)
	for index, column := range extraColumns {
		extraWidths[index] = len(column.Header)
	}
	for _, result := range results {
		if width := len(result.Spec.SpecName()); width > maxAgent {
			maxAgent = width
		}
		if width := len(result.Spec.SpecProjectPath()); width > maxPath {
			maxPath = width
		}
		if width := len(string(result.Spec.SpecCategory())); width > maxCategory {
			maxCategory = width
		}
		if width := len(string(result.Status)); width > maxStatus {
			maxStatus = width
		}
		for index, column := range extraColumns {
			if width := len(column.Value(result)); width > extraWidths[index] {
				extraWidths[index] = width
			}
		}
	}
	// Header line.
	if _, err := fmt.Fprintf(out, "%-*s  ", maxAgent, columnAgent); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for index, column := range extraColumns {
		if _, err := fmt.Fprintf(out, "%-*s  ", extraWidths[index], column.Header); err != nil {
			return fmt.Errorf("write header: %w", err)
		}
	}
	if _, err := fmt.Fprintf(out, "%-*s  %-*s  %-*s  %s\n",
		maxPath, columnPath,
		maxCategory, columnCategory,
		maxStatus, columnStatus,
		columnDetail,
	); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	// Data rows.
	for _, result := range results {
		if _, err := fmt.Fprintf(out, "%-*s  ", maxAgent, result.Spec.SpecName()); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
		for index, column := range extraColumns {
			if _, err := fmt.Fprintf(out, "%-*s  ", extraWidths[index], column.Value(result)); err != nil {
				return fmt.Errorf("write row: %w", err)
			}
		}
		if _, err := fmt.Fprintf(out, "%-*s  %-*s  %-*s  %s\n",
			maxPath, result.Spec.SpecProjectPath(),
			maxCategory, string(result.Spec.SpecCategory()),
			maxStatus, string(result.Status),
			result.Detail,
		); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	return nil
}

// EmitHints writes follow-up hints derived from a result list. The hint
// set is identical across every resource subpackage so the helper lives
// in common.
func EmitHints(out io.Writer, results []Result) error {
	var (
		hasMismatch      = false
		hasConflict      = false
		hasRootCollision = false
		hasError         = false
	)
	for _, result := range results {
		switch result.Status {
		case StatusMismatch:
			hasMismatch = true
		case StatusConflict:
			hasConflict = true
		case StatusSkippedRootCollision:
			hasRootCollision = true
		case StatusError:
			hasError = true
		}
	}
	if hasMismatch {
		if _, err := fmt.Fprintln(out, "Hint: rerun with force=1 to rebuild mismatched links."); err != nil {
			return fmt.Errorf("write hint: %w", err)
		}
	}
	if hasConflict {
		if _, err := fmt.Fprintln(out, "Hint: real directories or files block linking. Resolve listed paths manually."); err != nil {
			return fmt.Errorf("write hint: %w", err)
		}
	}
	if hasRootCollision {
		if _, err := fmt.Fprintln(out, "Hint: rootCollision agents (e.g. openclaw) require explicit agent=<name> force=1."); err != nil {
			return fmt.Errorf("write hint: %w", err)
		}
	}
	if hasError {
		if _, err := fmt.Fprintln(out, "Hint: see DETAIL column for the underlying error message."); err != nil {
			return fmt.Errorf("write hint: %w", err)
		}
	}
	return nil
}
