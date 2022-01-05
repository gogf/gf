// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"io"
)

// SetConfig set configurations for the defaultLogger.
func SetConfig(config Config) error {
	return defaultLogger.SetConfig(config)
}

// SetConfigWithMap set configurations with map for the defaultLogger.
func SetConfigWithMap(m map[string]interface{}) error {
	return defaultLogger.SetConfigWithMap(m)
}

// SetPath sets the directory path for file logging.
func SetPath(path string) error {
	return defaultLogger.SetPath(path)
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
func GetPath() string {
	return defaultLogger.GetPath()
}

// SetFile sets the file name `pattern` for file logging.
// Datetime pattern can be used in `pattern`, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
func SetFile(pattern string) {
	defaultLogger.SetFile(pattern)
}

// SetLevel sets the default logging level.
func SetLevel(level int) {
	defaultLogger.SetLevel(level)
}

// GetLevel returns the default logging level value.
func GetLevel() int {
	return defaultLogger.GetLevel()
}

// SetWriter sets the customized logging `writer` for logging.
// The `writer` object should implements the io.Writer interface.
// Developer can use customized logging `writer` to redirect logging output to another service,
// eg: kafka, mysql, mongodb, etc.
func SetWriter(writer io.Writer) {
	defaultLogger.SetWriter(writer)
}

// GetWriter returns the customized writer object, which implements the io.Writer interface.
// It returns nil if no customized writer set.
func GetWriter() io.Writer {
	return defaultLogger.GetWriter()
}

// SetDebug enables/disables the debug level for default defaultLogger.
// The debug level is enabled in default.
func SetDebug(debug bool) {
	defaultLogger.SetDebug(debug)
}

// SetAsync enables/disables async logging output feature for default defaultLogger.
func SetAsync(enabled bool) {
	defaultLogger.SetAsync(enabled)
}

// SetStdoutPrint sets whether ouptput the logging contents to stdout, which is true in default.
func SetStdoutPrint(enabled bool) {
	defaultLogger.SetStdoutPrint(enabled)
}

// SetHeaderPrint sets whether output header of the logging contents, which is true in default.
func SetHeaderPrint(enabled bool) {
	defaultLogger.SetHeaderPrint(enabled)
}

// SetPrefix sets prefix string for every logging content.
// Prefix is part of header, which means if header output is shut, no prefix will be output.
func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

// SetFlags sets extra flags for logging output features.
func SetFlags(flags int) {
	defaultLogger.SetFlags(flags)
}

// GetFlags returns the flags of defaultLogger.
func GetFlags() int {
	return defaultLogger.GetFlags()
}

// SetCtxKeys sets the context keys for defaultLogger. The keys is used for retrieving values
// from context and printing them to logging content.
//
// Note that multiple calls of this function will overwrite the previous set context keys.
func SetCtxKeys(keys ...interface{}) {
	defaultLogger.SetCtxKeys(keys...)
}

// GetCtxKeys retrieves and returns the context keys for logging.
func GetCtxKeys() []interface{} {
	return defaultLogger.GetCtxKeys()
}

// PrintStack prints the caller stack,
// the optional parameter `skip` specify the skipped stack offset from the end point.
func PrintStack(ctx context.Context, skip ...int) {
	defaultLogger.PrintStack(ctx, skip...)
}

// GetStack returns the caller stack content,
// the optional parameter `skip` specify the skipped stack offset from the end point.
func GetStack(skip ...int) string {
	return defaultLogger.GetStack(skip...)
}

// SetStack enables/disables the stack feature in failure logging outputs.
func SetStack(enabled bool) {
	defaultLogger.SetStack(enabled)
}

// SetLevelStr sets the logging level by level string.
func SetLevelStr(levelStr string) error {
	return defaultLogger.SetLevelStr(levelStr)
}

// SetLevelPrefix sets the prefix string for specified level.
func SetLevelPrefix(level int, prefix string) {
	defaultLogger.SetLevelPrefix(level, prefix)
}

// SetLevelPrefixes sets the level to prefix string mapping for the defaultLogger.
func SetLevelPrefixes(prefixes map[int]string) {
	defaultLogger.SetLevelPrefixes(prefixes)
}

// GetLevelPrefix returns the prefix string for specified level.
func GetLevelPrefix(level int) string {
	return defaultLogger.GetLevelPrefix(level)
}

// SetHandlers sets the logging handlers for default defaultLogger.
func SetHandlers(handlers ...Handler) {
	defaultLogger.SetHandlers(handlers...)
}

//SetWriterColorEnable sets the file logging with color
func SetWriterColorEnable(enabled bool) {
	defaultLogger.SetWriterColorEnable(enabled)
}
