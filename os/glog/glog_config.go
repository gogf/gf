// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"io"
)

// SetConfig set configurations for the logger.
func SetConfig(config Config) error {
	return logger.SetConfig(config)
}

// SetConfigWithMap set configurations with map for the logger.
func SetConfigWithMap(m map[string]interface{}) error {
	return logger.SetConfigWithMap(m)
}

// SetPath sets the directory path for file logging.
func SetPath(path string) error {
	return logger.SetPath(path)
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
func GetPath() string {
	return logger.GetPath()
}

// SetFile sets the file name <pattern> for file logging.
// Datetime pattern can be used in <pattern>, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
func SetFile(pattern string) {
	logger.SetFile(pattern)
}

// SetLevel sets the default logging level.
func SetLevel(level int) {
	logger.SetLevel(level)
}

// GetLevel returns the default logging level value.
func GetLevel() int {
	return logger.GetLevel()
}

// SetWriter sets the customized logging <writer> for logging.
// The <writer> object should implements the io.Writer interface.
// Developer can use customized logging <writer> to redirect logging output to another service,
// eg: kafka, mysql, mongodb, etc.
func SetWriter(writer io.Writer) {
	logger.SetWriter(writer)
}

// GetWriter returns the customized writer object, which implements the io.Writer interface.
// It returns nil if no customized writer set.
func GetWriter() io.Writer {
	return logger.GetWriter()
}

// SetDebug enables/disables the debug level for default logger.
// The debug level is enabled in default.
func SetDebug(debug bool) {
	logger.SetDebug(debug)
}

// SetAsync enables/disables async logging output feature for default logger.
func SetAsync(enabled bool) {
	logger.SetAsync(enabled)
}

// SetStdoutPrint sets whether ouptput the logging contents to stdout, which is true in default.
func SetStdoutPrint(enabled bool) {
	logger.SetStdoutPrint(enabled)
}

// SetHeaderPrint sets whether output header of the logging contents, which is true in default.
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
}

// SetPrefix sets prefix string for every logging content.
// Prefix is part of header, which means if header output is shut, no prefix will be output.
func SetPrefix(prefix string) {
	logger.SetPrefix(prefix)
}

// SetFlags sets extra flags for logging output features.
func SetFlags(flags int) {
	logger.SetFlags(flags)
}

// GetFlags returns the flags of logger.
func GetFlags() int {
	return logger.GetFlags()
}

// SetCtxKeys sets the context keys for logger. The keys is used for retrieving values
// from context and printing them to logging content.
//
// Note that multiple calls of this function will overwrite the previous set context keys.
func SetCtxKeys(keys ...interface{}) {
	logger.SetCtxKeys(keys...)
}

// GetCtxKeys retrieves and returns the context keys for logging.
func GetCtxKeys() []interface{} {
	return logger.GetCtxKeys()
}

// PrintStack prints the caller stack,
// the optional parameter <skip> specify the skipped stack offset from the end point.
func PrintStack(skip ...int) {
	logger.PrintStack(skip...)
}

// GetStack returns the caller stack content,
// the optional parameter <skip> specify the skipped stack offset from the end point.
func GetStack(skip ...int) string {
	return logger.GetStack(skip...)
}

// SetStack enables/disables the stack feature in failure logging outputs.
func SetStack(enabled bool) {
	logger.SetStack(enabled)
}

// SetLevelStr sets the logging level by level string.
func SetLevelStr(levelStr string) error {
	return logger.SetLevelStr(levelStr)
}

// SetLevelPrefix sets the prefix string for specified level.
func SetLevelPrefix(level int, prefix string) {
	logger.SetLevelPrefix(level, prefix)
}

// SetLevelPrefixes sets the level to prefix string mapping for the logger.
func SetLevelPrefixes(prefixes map[int]string) {
	logger.SetLevelPrefixes(prefixes)
}

// GetLevelPrefix returns the prefix string for specified level.
func GetLevelPrefix(level int) string {
	return logger.GetLevelPrefix(level)
}
