// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"errors"
	"fmt"
	"strings"
)

// Note that the LEVEL_PANI and LEVEL_FATA levels are not used for logging output,
// but for prefix configurations.
const (
	LEVEL_ALL  = LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT
	LEVEL_DEV  = LEVEL_ALL
	LEVEL_PROD = LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT
	LEVEL_DEBU = 1 << iota // 8
	LEVEL_INFO             // 16
	LEVEL_NOTI             // 32
	LEVEL_WARN             // 64
	LEVEL_ERRO             // 128
	LEVEL_CRIT             // 256
	LEVEL_PANI             // 512
	LEVEL_FATA             // 1024
)

// defaultLevelPrefixes defines the default level and its mapping prefix string.
var defaultLevelPrefixes = map[int]string{
	LEVEL_DEBU: "DEBU",
	LEVEL_INFO: "INFO",
	LEVEL_NOTI: "NOTI",
	LEVEL_WARN: "WARN",
	LEVEL_ERRO: "ERRO",
	LEVEL_CRIT: "CRIT",
	LEVEL_PANI: "PANI",
	LEVEL_FATA: "FATA",
}

// levelStringMap defines level string name to its level mapping.
var levelStringMap = map[string]int{
	"ALL":      LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"DEV":      LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"DEVELOP":  LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"PROD":     LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"PRODUCT":  LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"DEBU":     LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"DEBUG":    LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"INFO":     LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"NOTI":     LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"NOTICE":   LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"WARN":     LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"WARNING":  LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
	"ERRO":     LEVEL_ERRO | LEVEL_CRIT,
	"ERROR":    LEVEL_ERRO | LEVEL_CRIT,
	"CRIT":     LEVEL_CRIT,
	"CRITICAL": LEVEL_CRIT,
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level int) {
	l.config.Level = level
}

// GetLevel returns the logging level value.
func (l *Logger) GetLevel() int {
	return l.config.Level
}

// SetLevelStr sets the logging level by level string.
func (l *Logger) SetLevelStr(levelStr string) error {
	if level, ok := levelStringMap[strings.ToUpper(levelStr)]; ok {
		l.config.Level = level
	} else {
		return errors.New(fmt.Sprintf(`invalid level string: %s`, levelStr))
	}
	return nil
}

// SetLevelPrefix sets the prefix string for specified level.
func (l *Logger) SetLevelPrefix(level int, prefix string) {
	l.config.LevelPrefixes[level] = prefix
}

// SetLevelPrefixes sets the level to prefix string mapping for the logger.
func (l *Logger) SetLevelPrefixes(prefixes map[int]string) {
	for k, v := range prefixes {
		l.config.LevelPrefixes[k] = v
	}
}

// GetLevelPrefix returns the prefix string for specified level.
func (l *Logger) GetLevelPrefix(level int) string {
	return l.config.LevelPrefixes[level]
}

// getLevelPrefixWithBrackets returns the prefix string with brackets for specified level.
func (l *Logger) getLevelPrefixWithBrackets(level int) string {
	if s, ok := l.config.LevelPrefixes[level]; ok {
		return "[" + s + "]"
	}
	return ""
}
