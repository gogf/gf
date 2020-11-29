// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"github.com/gogf/gf/internal/intlog"
	"io"

	"github.com/gogf/gf/os/gfile"
)

// Ctx is a chaining function,
// which sets the context for current logging.
func (l *Logger) Ctx(ctx context.Context, keys ...interface{}) *Logger {
	if ctx == nil {
		return l
	}
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.ctx = ctx
	if len(keys) > 0 {
		logger.SetCtxKeys(keys...)
	}
	return logger
}

// To is a chaining function,
// which redirects current logging content output to the specified <writer>.
func (l *Logger) To(writer io.Writer) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.SetWriter(writer)
	return logger
}

// Path is a chaining function,
// which sets the directory path to <path> for current logging content output.
//
// Note that the parameter <path> is a directory path, not a file path.
func (l *Logger) Path(path string) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	if path != "" {
		if err := logger.SetPath(path); err != nil {
			// panic(err)
			intlog.Error(err)
		}
	}
	return logger
}

// Cat is a chaining function,
// which sets the category to <category> for current logging content output.
// Param <category> can be hierarchical, eg: module/user.
func (l *Logger) Cat(category string) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	if logger.config.Path != "" {
		if err := logger.SetPath(gfile.Join(logger.config.Path, category)); err != nil {
			// panic(err)
			intlog.Error(err)
		}
	}
	return logger
}

// File is a chaining function,
// which sets file name <pattern> for the current logging content output.
func (l *Logger) File(file string) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.SetFile(file)
	return logger
}

// Level is a chaining function,
// which sets logging level for the current logging content output.
func (l *Logger) Level(level int) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.SetLevel(level)
	return logger
}

// LevelStr is a chaining function,
// which sets logging level for the current logging content output using level string.
func (l *Logger) LevelStr(levelStr string) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	if err := logger.SetLevelStr(levelStr); err != nil {
		// panic(err)
		intlog.Error(err)
	}
	return logger
}

// Skip is a chaining function,
// which sets stack skip for the current logging content output.
// It also affects the caller file path checks when line number printing enabled.
func (l *Logger) Skip(skip int) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.SetStackSkip(skip)
	return logger
}

// Stack is a chaining function,
// which sets stack options for the current logging content output .
func (l *Logger) Stack(enabled bool, skip ...int) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.SetStack(enabled)
	if len(skip) > 0 {
		logger.SetStackSkip(skip[0])
	}
	return logger
}

// StackWithFilter is a chaining function,
// which sets stack filter for the current logging content output .
func (l *Logger) StackWithFilter(filter string) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	logger.SetStack(true)
	logger.SetStackFilter(filter)
	return logger
}

// Stdout is a chaining function,
// which enables/disables stdout for the current logging content output.
// It's enabled in default.
func (l *Logger) Stdout(enabled ...bool) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	// stdout printing is enabled if <enabled> is not passed.
	if len(enabled) > 0 && !enabled[0] {
		logger.config.StdoutPrint = false
	} else {
		logger.config.StdoutPrint = true
	}
	return logger
}

// Header is a chaining function,
// which enables/disables log header for the current logging content output.
// It's enabled in default.
func (l *Logger) Header(enabled ...bool) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	// header is enabled if <enabled> is not passed.
	if len(enabled) > 0 && !enabled[0] {
		logger.SetHeaderPrint(false)
	} else {
		logger.SetHeaderPrint(true)
	}
	return logger
}

// Line is a chaining function,
// which enables/disables printing its caller file path along with its line number.
// The parameter <long> specified whether print the long absolute file path, eg: /a/b/c/d.go:23,
// or else short one: d.go:23.
func (l *Logger) Line(long ...bool) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	if len(long) > 0 && long[0] {
		logger.config.Flags |= F_FILE_LONG
	} else {
		logger.config.Flags |= F_FILE_SHORT
	}
	return logger
}

// Async is a chaining function,
// which enables/disables async logging output feature.
func (l *Logger) Async(enabled ...bool) *Logger {
	logger := (*Logger)(nil)
	if l.parent == nil {
		logger = l.Clone()
	} else {
		logger = l
	}
	// async feature is enabled if <enabled> is not passed.
	if len(enabled) > 0 && !enabled[0] {
		logger.SetAsync(false)
	} else {
		logger.SetAsync(true)
	}
	return logger
}
