// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// getProjectRootPath finds the project root path by looking for go.mod file
func getProjectRootPath(path string) string {
	// Start from the given path and look for go.mod file
	currentPath := path
	for {
		// Check if go.mod exists in current directory
		goModPath := filepath.Join(currentPath, "go.mod")
		if gfile.Exists(goModPath) {
			return currentPath
		}

		// Move to parent directory
		parentPath := filepath.Dir(currentPath)
		// If we've reached the root directory, stop
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	// If no go.mod found, return the original path
	return path
}

// TableMetadata represents the metadata of a table for incremental generation
type TableMetadata struct {
	TableName   string `json:"table_name"`
	Hash        string `json:"hash"`
	GeneratedAt int64  `json:"generated_at"`
}

// DaoGenMetadata represents the metadata for dao generation
type DaoGenMetadata struct {
	Tables  []TableMetadata `json:"tables"`
	Version string          `json:"version"`
}

const (
	// MetadataFileName is the name of the metadata file
	MetadataFileName = ".gf_daogen_meta"
	// CurrentVersion is the current version of the metadata format
	CurrentVersion = "1.0"
)

// calculateTableHash calculates the hash of table structure
func calculateTableHash(fieldMap map[string]*gdb.TableField) string {
	// Sort field names to ensure consistent hash calculation
	fieldNames := make([]string, 0, len(fieldMap))
	for fieldName := range fieldMap {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	// Build a string representation of the table structure
	var builder strings.Builder
	for _, fieldName := range fieldNames {
		field := fieldMap[fieldName]
		builder.WriteString(fmt.Sprintf("%s|%s|%v|%s|%s|%s|%s\n",
			field.Name, field.Type, field.Null, field.Key,
			field.Default, field.Extra, field.Comment))
	}

	// Calculate MD5 hash
	hash := md5.Sum([]byte(builder.String()))
	return hex.EncodeToString(hash[:])
}

// loadMetadata loads the metadata from file
func loadMetadata(path string) (*DaoGenMetadata, error) {
	// Always load metadata from project root path
	projectRoot := getProjectRootPath(path)
	metaFilePath := gfile.Join(projectRoot, MetadataFileName)
	if !gfile.Exists(metaFilePath) {
		return &DaoGenMetadata{
			Tables:  make([]TableMetadata, 0),
			Version: CurrentVersion,
		}, nil
	}

	content := gfile.GetContents(metaFilePath)
	var metadata DaoGenMetadata
	if err := json.Unmarshal([]byte(content), &metadata); err != nil {
		return nil, err
	}

	// Convert to current version if needed
	if metadata.Version != CurrentVersion {
		metadata.Version = CurrentVersion
		metadata.Tables = make([]TableMetadata, 0)
	}

	return &metadata, nil
}

// saveMetadata saves the metadata to file
func saveMetadata(path string, metadata *DaoGenMetadata) error {
	// Always save metadata to project root path
	projectRoot := getProjectRootPath(path)
	metaFilePath := gfile.Join(projectRoot, MetadataFileName)
	content, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return gfile.PutContents(metaFilePath, string(content))
}

// isTableChanged checks if the table structure has changed
func isTableChanged(ctx context.Context, db gdb.DB, tableName string, metadata *DaoGenMetadata) (bool, string, error) {
	// Get current table fields
	fieldMap, err := db.TableFields(ctx, tableName)
	if err != nil {
		return false, "", fmt.Errorf("failed to get table fields: %v", err)
	}

	// Calculate current hash
	currentHash := calculateTableHash(fieldMap)

	// Find existing metadata for this table
	for _, tableMeta := range metadata.Tables {
		if tableMeta.TableName == tableName {
			// Compare hashes
			return tableMeta.Hash != currentHash, currentHash, nil
		}
	}

	// Table not found in metadata, treat as changed
	return true, currentHash, nil
}

// updateTableMetadata updates the metadata for a table
func updateTableMetadata(metadata *DaoGenMetadata, tableName, hash string) {
	// Remove existing entry if exists
	tables := make([]TableMetadata, 0, len(metadata.Tables))
	for _, tableMeta := range metadata.Tables {
		if tableMeta.TableName != tableName {
			tables = append(tables, tableMeta)
		}
	}

	// Add new entry
	tables = append(tables, TableMetadata{
		TableName:   tableName,
		Hash:        hash,
		GeneratedAt: gtime.Timestamp(),
	})

	metadata.Tables = tables
}

// shouldGenerateTable determines if a table should be generated based on incremental logic
func shouldGenerateTable(ctx context.Context, db gdb.DB, tableName string, metadata *DaoGenMetadata, incremental bool) (bool, string, error) {
	// If incremental is disabled, always generate
	if !incremental {
		fieldMap, err := db.TableFields(ctx, tableName)
		if err != nil {
			return false, "", fmt.Errorf("failed to get table fields: %v", err)
		}
		return true, calculateTableHash(fieldMap), nil
	}

	// Check if table has changed
	changed, hash, err := isTableChanged(ctx, db, tableName, metadata)
	if err != nil {
		mlog.Printf("Failed to check if table %s has changed: %v, will generate anyway", tableName, err)
		fieldMap, err := db.TableFields(ctx, tableName)
		if err != nil {
			return false, "", fmt.Errorf("failed to get table fields: %v", err)
		}
		return true, calculateTableHash(fieldMap), nil
	}

	return changed, hash, nil
}