// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

// Reserved internal error code of framework: code < 1000.

const (
	CodeNil                      = -1  // No error code specified.
	CodeOk                       = 0   // It is OK.
	CodeInternalError            = 50  // An error occurred internally.
	CodeValidationFailed         = 51  // Data validation failed.
	CodeDbOperationError         = 52  // Database operation error.
	CodeInvalidParameter         = 53  // The given parameter for current operation is invalid.
	CodeMissingParameter         = 54  // Parameter for current operation is missing.
	CodeInvalidOperation         = 55  // The function cannot be used like this.
	CodeInvalidConfiguration     = 56  // The configuration is invalid for current operation.
	CodeMissingConfiguration     = 57  // The configuration is missing for current operation.
	CodeNotImplemented           = 58  // The operation is not implemented yet.
	CodeNotSupported             = 59  // The operation is not supported yet.
	CodeOperationFailed          = 60  // I tried, but I cannot give you what you want.
	CodeNotAuthorized            = 61  // Not Authorized.
	CodeSecurityReason           = 62  // Security Reason.
	CodeServerBusy               = 63  // Server is busy, please try again later.
	CodeUnknown                  = 64  // Unknown error.
	CodeResourceNotExist         = 65  // Resource does not exist.
	CodeBusinessValidationFailed = 300 // Business validation failed.
)

var (
	// codeMessageMap is the mapping from code to according string message.
	codeMessageMap = map[int]string{
		CodeNil:                      "",
		CodeOk:                       "OK",
		CodeInternalError:            "Internal Error",
		CodeValidationFailed:         "Validation Failed",
		CodeDbOperationError:         "Database Operation Error",
		CodeInvalidParameter:         "Invalid Parameter",
		CodeMissingParameter:         "Missing Parameter",
		CodeInvalidOperation:         "Invalid Operation",
		CodeInvalidConfiguration:     "Invalid Configuration",
		CodeMissingConfiguration:     "Missing Configuration",
		CodeNotImplemented:           "Not Implemented",
		CodeNotSupported:             "Not Supported",
		CodeOperationFailed:          "Operation Failed",
		CodeNotAuthorized:            "Not Authorized",
		CodeSecurityReason:           "Security Reason",
		CodeServerBusy:               "Server Is Busy",
		CodeUnknown:                  "Unknown Error",
		CodeResourceNotExist:         "Resource Not Exist",
		CodeBusinessValidationFailed: "Business Validation Failed",
	}
)

// RegisterCode registers custom error code to global error codes for gerror recognition.
func RegisterCode(code int, message string) {
	codeMessageMap[code] = message
}

// RegisterCodeMap registers custom error codes to global error codes for gerror recognition using map.
func RegisterCodeMap(codes map[int]string) {
	for k, v := range codes {
		codeMessageMap[k] = v
	}
}
