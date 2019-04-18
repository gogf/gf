// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
    "github.com/gogf/gf/g/os/gfile"
    "io"
)

// To is a chaining function, 
// which redirects current logging content output to the specified <writer>.
func (l *Logger) To(writer io.Writer) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    logger.SetWriter(writer)
    return logger
}

// Path is a chaining function,
// which sets the directory path to <path> for current logging content output.
func (l *Logger) Path(path string) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    if path != "" {
        logger.SetPath(path)
    }
    return logger
}

// Cat is a chaining function, 
// which sets the category to <category> for current logging content output.
// Param <category> can be hierarchical, eg: module/user.
func (l *Logger) Cat(category string) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    path := l.path.Val()
    if path != "" {
        logger.SetPath(path + gfile.Separator + category)
    }
    return logger
}

// File is a chaining function, 
// which sets file name <pattern> for the current logging content output.
func (l *Logger) File(file string) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
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
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    logger.SetLevel(level)
    return logger
}

// Backtrace is a chaining function, 
// which sets backtrace options for the current logging content output .
func (l *Logger) Backtrace(enabled bool, skip...int) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    logger.SetBacktrace(enabled)
    if len(skip) > 0 {
        logger.SetBacktraceSkip(skip[0])
    }
    return logger
}

// StdPrint is a chaining function, 
// which enables/disables stdout for the current logging content output.
func (l *Logger) StdPrint(enabled bool) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    logger.SetStdPrint(enabled)
    return logger
}

// Header is a chaining function, 
// which enables/disables log header for the current logging content output.
func (l *Logger) Header(enabled bool) *Logger {
    logger := (*Logger)(nil)
    if l.pr == nil {
        logger = l.Clone()
    } else {
        logger = l
    }
    logger.printHeader.Set(enabled)
    return logger
}