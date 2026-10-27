// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains unit tests for the resource-agnostic table rendering
// and hint emission logic in the common subpackage. Tests use a tiny
// fake registry of fakeSpec values so they exercise the renderer without
// depending on any concrete resource subpackage.

package common

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderProducesTabularOutput(t *testing.T) {
	results := []Result{
		{Spec: fakeSpec{name: "claude-code", category: CategoryLink}, Status: StatusOK},
		{Spec: fakeSpec{name: "cursor", category: CategoryNative}, Status: StatusNative},
	}
	var buffer bytes.Buffer
	if err := Render(&buffer, results); err != nil {
		t.Fatalf("render: %v", err)
	}
	output := buffer.String()
	for _, header := range []string{"AGENT", "PROJECT PATH", "CATEGORY", "STATUS", "DETAIL"} {
		if !strings.Contains(output, header) {
			t.Fatalf("missing header %s; got %q", header, output)
		}
	}
	for _, name := range []string{"claude-code", "cursor"} {
		if !strings.Contains(output, name) {
			t.Fatalf("expected agent %s in output: %s", name, output)
		}
	}
}

func TestRenderSupportsExtraColumns(t *testing.T) {
	results := []Result{
		{Spec: fakeSpec{name: "claude-code", category: CategoryLink}, Status: StatusOK},
	}
	var buffer bytes.Buffer
	extra := ColumnSpec{Header: "SOURCE", Value: func(r Result) string { return r.Spec.SpecSourcePath() }}
	if err := Render(&buffer, results, extra); err != nil {
		t.Fatalf("render: %v", err)
	}
	output := buffer.String()
	if !strings.Contains(output, "SOURCE") {
		t.Fatalf("expected SOURCE header in output: %q", output)
	}
	// fakeSpec.SpecSourcePath() returns ".source"; check it shows up.
	if !strings.Contains(output, ".source") {
		t.Fatalf("expected per-row SOURCE value in output: %q", output)
	}
}

func TestEmitHintsReflectsStatuses(t *testing.T) {
	results := []Result{
		{Spec: fakeSpec{name: "claude-code"}, Status: StatusMismatch},
		{Spec: fakeSpec{name: "codebuddy"}, Status: StatusConflict},
		{Spec: fakeSpec{name: "openclaw"}, Status: StatusSkippedRootCollision},
		{Spec: fakeSpec{name: "demo"}, Status: StatusError, Detail: "boom"},
	}
	var buffer bytes.Buffer
	if err := EmitHints(&buffer, results); err != nil {
		t.Fatalf("emit: %v", err)
	}
	hint := buffer.String()
	for _, fragment := range []string{"force=1", "Resolve listed paths", "rootCollision", "DETAIL"} {
		if !strings.Contains(hint, fragment) {
			t.Fatalf("missing hint fragment %q in %q", fragment, hint)
		}
	}
}
