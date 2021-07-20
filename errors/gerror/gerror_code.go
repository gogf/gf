// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

// Reserved internal error code of framework: code < 1000.

const (
	CodeNil                  = -1 // No error code specified.
	CodeOk                   = 0  // It is OK without error.
	CodeInternalError        = 50 // An error occurred internally.
	CodeValidationFailed     = 51 // Data validation failed.
	CodeDbOperationError     = 52 // Database operation error.
	CodeInvalidParameter     = 53 // The given parameter for current operation is invalid.
	CodeMissingParameter     = 54 // Parameter for current operation is missing.
	CodeInvalidOperation     = 55 // The function cannot be used like this.
	CodeInvalidConfiguration = 56 // The configuration is invalid for current operation.
	CodeMissingConfiguration = 57 // The configuration is missing for current operation.
	CodeNotImplemented       = 58 // The operation is not implemented yet.
	CodeNotSupported         = 59 // The operation is not supported yet.
	CodeOperationFailed      = 60 // I tried, but I cannot give you what you want.
)
