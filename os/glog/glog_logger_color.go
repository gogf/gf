// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import "github.com/fatih/color"

const (
	COLOR_BLACK = 30 + iota
	COLOR_RED
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_MAGENTA
	COLOR_CYAN
	COLOR_WHITE
)

// defaultLevelColor defines the default level and its mapping prefix string.
var defaultLevelColor = map[int]int{
	LEVEL_DEBU: COLOR_YELLOW,
	LEVEL_INFO: COLOR_GREEN,
	LEVEL_NOTI: COLOR_CYAN,
	LEVEL_WARN: COLOR_YELLOW,
	LEVEL_ERRO: COLOR_RED,
	LEVEL_CRIT: COLOR_RED,
	LEVEL_PANI: COLOR_RED,
	LEVEL_FATA: COLOR_RED,
}

// getColoredStr returns a string that is colored by given color.
func (l *Logger) getColoredStr(c int, s string) string {
	return color.New(color.Attribute(c)).Sprint(s)
}

func (l *Logger) getColorByLevel(level int) int {
	return defaultLevelColor[level]
}
