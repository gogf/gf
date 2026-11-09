// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

import "github.com/gogf/gf/v2/os/gctx"

// Context key constants for configuration operations.
const (
	// ContextKeyFileName is the context key for file name
	ContextKeyFileName gctx.StrKey = "fileName"
	// ContextKeyFilePath is the context key for file path
	ContextKeyFilePath gctx.StrKey = "filePath"
	// ContextKeyFileType is the context key for file type
	ContextKeyFileType gctx.StrKey = "fileType"
	// ContextKeyOperation is the context key for operation type
	ContextKeyOperation gctx.StrKey = "operation"
	// ContextKeyKey is the context key for key
	ContextKeyKey gctx.StrKey = "key"
	// ContextKeyValue is the context key for value
	ContextKeyValue gctx.StrKey = "value"
	// ContextKeyContent is the context key for set content
	ContextKeyContent gctx.StrKey = "content"
)

// OperationType defines the type for configuration operation.
type OperationType string

// Operation constants for configuration operations.
const (
	// OperationSet represents set operation
	OperationSet OperationType = "set"
	// OperationWrite represents write operation
	OperationWrite OperationType = "write"
	// OperationRename represents rename operation
	OperationRename OperationType = "rename"
	// OperationRemove represents remove operation
	OperationRemove OperationType = "remove"
	// OperationCreate represents create operation
	OperationCreate OperationType = "create"
	// OperationChmod represents chmod operation
	OperationChmod OperationType = "chmod"
	// OperationClear represents clear operation
	OperationClear OperationType = "clear"
	// OperationUpdate represents update operation
	OperationUpdate OperationType = "update"
)
