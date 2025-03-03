// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// Config is the configuration object for logger.
type Config struct {
	Handlers             []Handler      `json:"-"`                    // Logger handlers which implement feature similar as middleware.
	Writer               io.Writer      `json:"-"`                    // Customized io.Writer.
	Flags                int            `json:"flags"`                // Extra flags for logging output features.
	TimeFormat           string         `json:"timeFormat"`           // Logging time format
	Path                 string         `json:"path"`                 // Logging directory path.
	File                 string         `json:"file"`                 // Format pattern for logging file.
	Level                int            `json:"level"`                // Output level.
	Prefix               string         `json:"prefix"`               // Prefix string for every logging content.
	StSkip               int            `json:"stSkip"`               // Skipping count for stack.
	StStatus             int            `json:"stStatus"`             // Stack status(1: enabled - default; 0: disabled)
	StFilter             string         `json:"stFilter"`             // Stack string filter.
	CtxKeys              []interface{}  `json:"ctxKeys"`              // Context keys for logging, which is used for value retrieving from context.
	HeaderPrint          bool           `json:"header"`               // Print header or not(true in default).
	StdoutPrint          bool           `json:"stdout"`               // Output to stdout or not(true in default).
	LevelPrint           bool           `json:"levelPrint"`           // Print level format string or not(true in default).
	LevelPrefixes        map[int]string `json:"levelPrefixes"`        // Logging level to its prefix string mapping.
	RotateSize           int64          `json:"rotateSize"`           // Rotate the logging file if its size > 0 in bytes.
	RotateExpire         time.Duration  `json:"rotateExpire"`         // Rotate the logging file if its mtime exceeds this duration.
	RotateBackupLimit    int            `json:"rotateBackupLimit"`    // Max backup for rotated files, default is 0, means no backups.
	RotateBackupExpire   time.Duration  `json:"rotateBackupExpire"`   // Max expires for rotated files, which is 0 in default, means no expiration.
	RotateBackupCompress int            `json:"rotateBackupCompress"` // Compress level for rotated files using gzip algorithm. It's 0 in default, means no compression.
	RotateCheckInterval  time.Duration  `json:"rotateCheckInterval"`  // Asynchronously checks the backups and expiration at intervals. It's 1 hour in default.
	StdoutColorDisabled  bool           `json:"stdoutColorDisabled"`  // Logging level prefix with color to writer or not (false in default).
	WriterColorEnable    bool           `json:"writerColorEnable"`    // Logging level prefix with color to writer or not (false in default).
	internalConfig
}

type internalConfig struct {
	rotatedHandlerInitialized *gtype.Bool // Whether the rotation feature initialized.
}

// DefaultConfig returns the default configuration for logger.
func DefaultConfig() Config {
	c := Config{
		File:                defaultFileFormat,
		Flags:               F_TIME_STD,
		TimeFormat:          defaultTimeFormat,
		Level:               LEVEL_ALL,
		CtxKeys:             []interface{}{},
		StStatus:            1,
		HeaderPrint:         true,
		StdoutPrint:         true,
		LevelPrint:          true,
		LevelPrefixes:       make(map[int]string, len(defaultLevelPrefixes)),
		RotateCheckInterval: time.Hour,
		internalConfig: internalConfig{
			rotatedHandlerInitialized: gtype.NewBool(),
		},
	}
	for k, v := range defaultLevelPrefixes {
		c.LevelPrefixes[k] = v
	}
	if !defaultDebug {
		c.Level = c.Level & ^LEVEL_DEBU
	}
	return c
}

// GetConfig returns the configuration of current Logger.
func (l *Logger) GetConfig() Config {
	return l.config
}

// SetConfig set configurations for the logger.
func (l *Logger) SetConfig(config Config) error {
	l.config = config
	// Necessary validation.
	if config.Path != "" {
		if err := l.SetPath(config.Path); err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
			return err
		}
	}
	intlog.Printf(context.TODO(), "SetConfig: %+v", l.config)
	return nil
}

// SetConfigWithMap set configurations with map for the logger.
func (l *Logger) SetConfigWithMap(m map[string]interface{}) error {
	if len(m) == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "configuration cannot be empty")
	}
	// The m now is a shallow copy of m.
	// A little tricky, isn't it?
	m = gutil.MapCopy(m)
	// Change string configuration to int value for level.
	levelKey, levelValue := gutil.MapPossibleItemByKey(m, "Level")
	if levelValue != nil {
		if level, ok := levelStringMap[strings.ToUpper(gconv.String(levelValue))]; ok {
			m[levelKey] = level
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid level string: %v`, levelValue)
		}
	}
	// Change string configuration to int value for file rotation size.
	rotateSizeKey, rotateSizeValue := gutil.MapPossibleItemByKey(m, "RotateSize")
	if rotateSizeValue != nil {
		m[rotateSizeKey] = gfile.StrToSize(gconv.String(rotateSizeValue))
		if m[rotateSizeKey] == -1 {
			return gerror.NewCodef(gcode.CodeInvalidConfiguration, `invalid rotate size: %v`, rotateSizeValue)
		}
	}
	if err := gconv.Struct(m, &l.config); err != nil {
		return err
	}
	return l.SetConfig(l.config)
}

// SetDebug enables/disables the debug level for logger.
// The debug level is enabled in default.
func (l *Logger) SetDebug(debug bool) {
	if debug {
		l.config.Level = l.config.Level | LEVEL_DEBU
	} else {
		l.config.Level = l.config.Level & ^LEVEL_DEBU
	}
}

// SetAsync enables/disables async logging output feature.
func (l *Logger) SetAsync(enabled bool) {
	if enabled {
		l.config.Flags = l.config.Flags | F_ASYNC
	} else {
		l.config.Flags = l.config.Flags & ^F_ASYNC
	}
}

// SetFlags sets extra flags for logging output features.
func (l *Logger) SetFlags(flags int) {
	l.config.Flags = flags
}

// GetFlags returns the flags of logger.
func (l *Logger) GetFlags() int {
	return l.config.Flags
}

// SetStack enables/disables the stack feature in failure logging outputs.
func (l *Logger) SetStack(enabled bool) {
	if enabled {
		l.config.StStatus = 1
	} else {
		l.config.StStatus = 0
	}
}

// SetStackSkip sets the stack offset from the end point.
func (l *Logger) SetStackSkip(skip int) {
	l.config.StSkip = skip
}

// SetStackFilter sets the stack filter from the end point.
func (l *Logger) SetStackFilter(filter string) {
	l.config.StFilter = filter
}

// SetCtxKeys sets the context keys for logger. The keys is used for retrieving values
// from context and printing them to logging content.
//
// Note that multiple calls of this function will overwrite the previous set context keys.
func (l *Logger) SetCtxKeys(keys ...interface{}) {
	l.config.CtxKeys = keys
}

// AppendCtxKeys appends extra keys to logger.
// It ignores the key if it is already appended to the logger previously.
func (l *Logger) AppendCtxKeys(keys ...interface{}) {
	var isExist bool
	for _, key := range keys {
		isExist = false
		for _, ctxKey := range l.config.CtxKeys {
			if ctxKey == key {
				isExist = true
				break
			}
		}
		if !isExist {
			l.config.CtxKeys = append(l.config.CtxKeys, key)
		}
	}
}

// GetCtxKeys retrieves and returns the context keys for logging.
func (l *Logger) GetCtxKeys() []interface{} {
	return l.config.CtxKeys
}

// SetWriter sets the customized logging `writer` for logging.
// The `writer` object should implement the io.Writer interface.
// Developer can use customized logging `writer` to redirect logging output to another service,
// eg: kafka, mysql, mongodb, etc.
func (l *Logger) SetWriter(writer io.Writer) {
	l.config.Writer = writer
}

// GetWriter returns the customized writer object, which implements the io.Writer interface.
// It returns nil if no writer previously set.
func (l *Logger) GetWriter() io.Writer {
	return l.config.Writer
}

// SetPath sets the directory path for file logging.
func (l *Logger) SetPath(path string) error {
	if path == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "logging path is empty")
	}
	if !gfile.Exists(path) {
		if err := gfile.Mkdir(path); err != nil {
			return gerror.Wrapf(err, `Mkdir "%s" failed in PWD "%s"`, path, gfile.Pwd())
		}
	}
	l.config.Path = strings.TrimRight(path, gfile.Separator)
	return nil
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
func (l *Logger) GetPath() string {
	return l.config.Path
}

// SetFile sets the file name `pattern` for file logging.
// Datetime pattern can be used in `pattern`, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
func (l *Logger) SetFile(pattern string) {
	l.config.File = pattern
}

// SetTimeFormat sets the time format for the logging time.
func (l *Logger) SetTimeFormat(timeFormat string) {
	l.config.TimeFormat = timeFormat
}

// SetStdoutPrint sets whether output the logging contents to stdout, which is true in default.
func (l *Logger) SetStdoutPrint(enabled bool) {
	l.config.StdoutPrint = enabled
}

// SetHeaderPrint sets whether output header of the logging contents, which is true in default.
func (l *Logger) SetHeaderPrint(enabled bool) {
	l.config.HeaderPrint = enabled
}

// SetLevelPrint sets whether output level string of the logging contents, which is true in default.
func (l *Logger) SetLevelPrint(enabled bool) {
	l.config.LevelPrint = enabled
}

// SetPrefix sets prefix string for every logging content.
// Prefix is part of header, which means if header output is shut, no prefix will be output.
func (l *Logger) SetPrefix(prefix string) {
	l.config.Prefix = prefix
}

// SetHandlers sets the logging handlers for current logger.
func (l *Logger) SetHandlers(handlers ...Handler) {
	l.config.Handlers = handlers
}

// SetWriterColorEnable enables file/writer logging with color.
func (l *Logger) SetWriterColorEnable(enabled bool) {
	l.config.WriterColorEnable = enabled
}

// SetStdoutColorDisabled disables stdout logging with color.
func (l *Logger) SetStdoutColorDisabled(disabled bool) {
	l.config.StdoutColorDisabled = disabled
}
