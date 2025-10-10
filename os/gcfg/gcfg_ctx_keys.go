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
	// KeyFileName is the context key for file name
	KeyFileName gctx.StrKey = "fileName"
	// KeyFilePath is the context key for file path
	KeyFilePath gctx.StrKey = "filePath"
	// KeyFileType is the context key for file type
	KeyFileType gctx.StrKey = "fileType"
	// KeyOperation is the context key for operation type
	KeyOperation gctx.StrKey = "operation"
	// KeyKey is the context key for key
	KeyKey gctx.StrKey = "key"
	// KeyValue is the context key for value
	KeyValue gctx.StrKey = "value"
	// KeyContent is the context key for set content
	KeyContent gctx.StrKey = "content"
)

// Operation constants for configuration operations.
const (
	// OperationSet represents set operation
	OperationSet = "set"
	// OperationWrite represents write operation
	OperationWrite = "write"
	// OperationRename represents rename operation
	OperationRename = "rename"
	// OperationRemove represents remove operation
	OperationRemove = "remove"
	// OperationCreate represents create operation
	OperationCreate = "create"
	// OperationChmod represents chmod operation
	OperationChmod = "chmod"
	// OperationClear represents clear operation
	OperationClear = "clear"
)
