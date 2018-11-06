// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package glog

import (
    "gitee.com/johng/gf/g/os/gfile"
    "io"
)

// 链式操作，设置下一次写入日志内容的Writer
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

// 链式操作，设置下一次输出的日志分类(可以按照文件目录层级设置)，在当前logpath或者当前工作目录下创建category目录，
// 这是一个链式操作，可以设置多个分类，将会创建层级的日志分类目录。
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

// 日志文件格式
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

// 设置日志打印等级
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

// 设置文件调用回溯信息
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

// 是否允许在设置输出文件时同时也输出到终端
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

// 是否打印每行日志头信息(默认开启)
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