// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcode provides universal error code definition and common error codes implements.
package gcode

// Code is universal error code interface definition.
type Code interface {
	// Code returns the integer number of current error code.
	Code() int

	// Message returns the brief message for current error code.
	Message() string

	// Detail returns the detailed information of current error code,
	// which is mainly designed as an extension field for error code.
	Detail() interface{}
}

// ================================================================================================================
// Common error code definition.
// There are reserved internal error code by framework: code < 1000.
// ================================================================================================================

var (
	CodeNil                      = localCode{-1, "", nil}                            // No error code specified.
	CodeOK                       = localCode{0, "OK", nil}                           // It is OK.
	CodeInternalError            = localCode{50, "Internal Error", nil}              // An error occurred internally.
	CodeValidationFailed         = localCode{51, "Validation Failed", nil}           // Data validation failed.
	CodeDbOperationError         = localCode{52, "Database Operation Error", nil}    // Database operation error.
	CodeInvalidParameter         = localCode{53, "Invalid Parameter", nil}           // The given parameter for current operation is invalid.
	CodeMissingParameter         = localCode{54, "Missing Parameter", nil}           // Parameter for current operation is missing.
	CodeInvalidOperation         = localCode{55, "Invalid Operation", nil}           // The function cannot be used like this.
	CodeInvalidConfiguration     = localCode{56, "Invalid Configuration", nil}       // The configuration is invalid for current operation.
	CodeMissingConfiguration     = localCode{57, "Missing Configuration", nil}       // The configuration is missing for current operation.
	CodeNotImplemented           = localCode{58, "Not Implemented", nil}             // The operation is not implemented yet.
	CodeNotSupported             = localCode{59, "Not Supported", nil}               // The operation is not supported yet.
	CodeOperationFailed          = localCode{60, "Operation Failed", nil}            // I tried, but I cannot give you what you want.
	CodeNotAuthorized            = localCode{61, "Not Authorized", nil}              // Not Authorized.
	CodeSecurityReason           = localCode{62, "Security Reason", nil}             // Security Reason.
	CodeServerBusy               = localCode{63, "Server Is Busy", nil}              // Server is busy, please try again later.
	CodeUnknown                  = localCode{64, "Unknown Error", nil}               // Unknown error.
	CodeNotFound                 = localCode{65, "Not Found", nil}                   // Resource does not exist.
	CodeInvalidRequest           = localCode{66, "Invalid Request", nil}             // Invalid request.
	CodeBusinessValidationFailed = localCode{300, "Business Validation Failed", nil} // Business validation failed.
)

// New creates and returns an error code.
// Note that it returns an interface object of Code.
func New(code int, message string, detail interface{}) Code {
	return localCode{
		code:    code,
		message: message,
		detail:  detail,
	}
}
