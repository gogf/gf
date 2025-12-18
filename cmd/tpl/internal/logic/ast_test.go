package logic

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
)

func TestASTReplacer_ReplaceInFile(t *testing.T) {
	ctx := context.Background()

	// Create a temp directory
	tempDir, err := os.MkdirTemp("", "ast-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test Go file
	testFile := filepath.Join(tempDir, "test.go")
	content := `package main

import (
	"fmt"
	"github.com/old/module/internal/service"
	"github.com/old/module/internal/model"
)

func main() {
	fmt.Println("hello")
}
`
	if err := gfile.PutContents(testFile, content); err != nil {
		t.Fatal(err)
	}

	// Create replacer and run
	replacer := NewASTReplacer("github.com/old/module", "new-project")
	if err := replacer.ReplaceInFile(ctx, testFile); err != nil {
		t.Fatal(err)
	}

	// Read result
	result := gfile.GetContents(testFile)

	// Verify replacements
	if !strings.Contains(result, `"new-project/internal/service"`) {
		t.Errorf("Expected import path to be replaced, got: %s", result)
	}
	if !strings.Contains(result, `"new-project/internal/model"`) {
		t.Errorf("Expected import path to be replaced, got: %s", result)
	}
	if strings.Contains(result, "github.com/old/module") {
		t.Errorf("Old import path should not exist, got: %s", result)
	}
	// Verify unrelated imports are not changed
	if !strings.Contains(result, `"fmt"`) {
		t.Errorf("Unrelated import should not be changed, got: %s", result)
	}
}

func TestASTReplacer_ReplaceInDir(t *testing.T) {
	ctx := context.Background()

	// Create a temp directory structure
	tempDir, err := os.MkdirTemp("", "ast-dir-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "internal", "logic")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	file1 := filepath.Join(tempDir, "main.go")
	file2 := filepath.Join(subDir, "logic.go")

	content1 := `package main

import "github.com/template/repo/internal/cmd"

func main() {}
`
	content2 := `package logic

import (
	"github.com/template/repo/internal/model"
	"github.com/template/repo/internal/dao"
)

func Init() {}
`
	if err := gfile.PutContents(file1, content1); err != nil {
		t.Fatal(err)
	}
	if err := gfile.PutContents(file2, content2); err != nil {
		t.Fatal(err)
	}

	// Run replacement
	replacer := NewASTReplacer("github.com/template/repo", "my-app")
	if err := replacer.ReplaceInDir(ctx, tempDir); err != nil {
		t.Fatal(err)
	}

	// Verify file1
	result1 := gfile.GetContents(file1)
	if !strings.Contains(result1, `"my-app/internal/cmd"`) {
		t.Errorf("File1: Expected import to be replaced, got: %s", result1)
	}

	// Verify file2
	result2 := gfile.GetContents(file2)
	if !strings.Contains(result2, `"my-app/internal/model"`) {
		t.Errorf("File2: Expected model import to be replaced, got: %s", result2)
	}
	if !strings.Contains(result2, `"my-app/internal/dao"`) {
		t.Errorf("File2: Expected dao import to be replaced, got: %s", result2)
	}
}

func TestASTReplacer_NoChangeForUnrelatedImports(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "ast-nochange-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.go")
	content := `package main

import (
	"fmt"
	"github.com/other/package"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {}
`
	if err := gfile.PutContents(testFile, content); err != nil {
		t.Fatal(err)
	}

	// Replace with a module that doesn't match any imports
	replacer := NewASTReplacer("github.com/nonexistent/module", "new-project")
	if err := replacer.ReplaceInFile(ctx, testFile); err != nil {
		t.Fatal(err)
	}

	// Content should remain unchanged
	result := gfile.GetContents(testFile)
	if !strings.Contains(result, `"github.com/other/package"`) {
		t.Errorf("Unrelated imports should not be changed")
	}
	if !strings.Contains(result, `"github.com/gogf/gf/v2/frame/g"`) {
		t.Errorf("GoFrame import should not be changed")
	}
}
