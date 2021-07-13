// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"github.com/fatih/color"
	"os"
	"time"
)

type Handler func(ctx context.Context, input *HandlerInput)

type HandlerInput struct {
	logger      *Logger
	index       int
	Ctx         context.Context
	Time        time.Time
	TimeFormat  string
	Level       int
	LevelFormat string
	CallerFunc  string
	CallerPath  string
	CtxStr      string
	Prefix      string
	Content     string
	IsAsync     bool
}

// defaultHandler is the default handler for logger.
func defaultHandler(ctx context.Context, input *HandlerInput) {
	input.logger.printToWriter(ctx, input)
}

func (i *HandlerInput) addStringToBuffer(buffer *bytes.Buffer, s string) {
	if buffer.Len() > 0 {
		buffer.WriteByte(' ')
	}
	buffer.WriteString(s)
}

func (i *HandlerInput) Buffer() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(i.TimeFormat)
	levelString := i.LevelFormat
	if i.logger.config.FileColorEnable {
		fg := i.getLevelFormatColor()
		levelString = color.New(fg).Sprintf(i.LevelFormat)
	}
	i.addStringToBuffer(buffer, levelString)
	msg := i.GetContent()
	i.addStringToBuffer(buffer, msg.String())
	return buffer
}

// Stdout print log to console
func (i *HandlerInput) Stdout() {
	_, _ = os.Stdout.Write([]byte(i.TimeFormat))
	fg := i.getLevelFormatColor()
	_, _ = color.New(fg).Print(" " + i.LevelFormat + " ")
	msg := i.GetContent()
	_, _ = os.Stdout.Write(msg.Bytes())
}

// GetContent returns the primary content.
func (i *HandlerInput) GetContent() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	if i.CallerFunc != "" {
		i.addStringToBuffer(buffer, i.CallerFunc)
	}
	if i.CallerPath != "" {
		i.addStringToBuffer(buffer, i.CallerPath)
	}
	if i.Prefix != "" {
		i.addStringToBuffer(buffer, i.Prefix)
	}
	if i.CtxStr != "" {
		i.addStringToBuffer(buffer, i.CtxStr)
	}
	if i.Content != "" {
		i.addStringToBuffer(buffer, i.Content)
	}
	i.addStringToBuffer(buffer, "\n")
	return buffer
}

// getLevelFormatColor returns the prefix string color.
func (i *HandlerInput) getLevelFormatColor() color.Attribute {
	fg := defaultLevelColor[i.Level]
	if i.logger.config.currentColor != 0 {
		fg = i.logger.config.currentColor
	}
	return fg
}

func (i *HandlerInput) String() string {
	return i.Buffer().String()
}

func (i *HandlerInput) Next() {
	if len(i.logger.config.Handlers)-1 > i.index {
		i.index++
		i.logger.config.Handlers[i.index](i.Ctx, i)
	}
}
