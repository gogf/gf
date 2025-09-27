// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

// Context key constants for configuration operations.
const (
	// KeyFileName is the context key for file name
	KeyFileName = "fileName"
	// KeyFilePath is the context key for file path
	KeyFilePath = "filePath"
	// KeyFileType is the context key for file type
	KeyFileType = "fileType"
	// KeyOperation is the context key for operation type
	KeyOperation = "operation"
	// KeySetKey is the context key for set key
	KeySetKey = "setKey"
	// KeySetValue is the context key for set value
	KeySetValue = "setValue"
	// KeySetContent is the context key for set content
	KeySetContent = "setContent  "
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
