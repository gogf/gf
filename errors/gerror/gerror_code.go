// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

// Reserved internal error code of framework: code < 1000.

const (
	// ===============================================================================
	// Common system codes.
	// ===============================================================================

	CodeNil                  = -1        // No error code specified.
	CodeOk                   = 0         // It is OK.
	CodeInternalError        = 50 + iota // An error occurred internally.
	CodeValidationFailed                 // Data validation failed.
	CodeDbOperationError                 // Database operation error.
	CodeInvalidParameter                 // The given parameter for current operation is invalid.
	CodeMissingParameter                 // Parameter for current operation is missing.
	CodeInvalidOperation                 // The function cannot be used like this.
	CodeInvalidConfiguration             // The configuration is invalid for current operation.
	CodeMissingConfiguration             // The configuration is missing for current operation.
	CodeNotImplemented                   // The operation is not implemented yet.
	CodeNotSupported                     // The operation is not supported yet.
	CodeOperationFailed                  // I tried, but I cannot give you what you want.
	CodeNotAuthorized                    // Not Authorized.
	CodeSecurityReason                   // Security Reason.
	CodeServerBusy                       // Server is busy, please try again later.
	CodeUnknown                          // Unknown error.
	CodeResourceNotExist                 // Resource does not exist.

	// ===============================================================================
	// Common business codes.
	// ===============================================================================

	CodeBusinessValidationFailed = 300 + iota // Business validation failed.
)
