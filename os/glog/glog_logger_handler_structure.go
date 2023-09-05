// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/gogf/gf/v2/util/gconv"
)

type structuredBuffer struct {
	in     *HandlerInput
	buffer *bytes.Buffer
}

const (
	structureKeyTime       = "Time"
	structureKeyLevel      = "Level"
	structureKeyPrefix     = "Prefix"
	structureKeyContent    = "Content"
	structureKeyTraceId    = "TraceId"
	structureKeyCallerFunc = "CallerFunc"
	structureKeyCallerPath = "CallerPath"
	structureKeyCtxStr     = "CtxStr"
	structureKeyStack      = "Stack"
)

// Copied from encoding/json/tables.go.
//
// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

// HandlerStructure is a handler for output logging content as a structured string.
func HandlerStructure(ctx context.Context, in *HandlerInput) {
	s := newStructuredBuffer(in)
	in.Buffer.Write(s.Bytes())
	in.Buffer.Write([]byte("\n"))
	in.Next(ctx)
}

func newStructuredBuffer(in *HandlerInput) *structuredBuffer {
	return &structuredBuffer{
		in:     in,
		buffer: bytes.NewBuffer(nil),
	}
}

func (buf *structuredBuffer) Bytes() []byte {
	buf.addValue(structureKeyTime, buf.in.TimeFormat)
	if buf.in.TraceId != "" {
		buf.addValue(structureKeyTraceId, buf.in.TraceId)
	}
	if buf.in.CtxStr != "" {
		buf.addValue(structureKeyCtxStr, buf.in.CtxStr)
	}
	if buf.in.LevelFormat != "" {
		buf.addValue(structureKeyLevel, buf.in.LevelFormat)
	}
	if buf.in.CallerPath != "" {
		buf.addValue(structureKeyCallerPath, buf.in.CallerPath)
	}
	if buf.in.CallerFunc != "" {
		buf.addValue(structureKeyCallerFunc, buf.in.CallerFunc)
	}
	if buf.in.Prefix != "" {
		buf.addValue(structureKeyPrefix, buf.in.Prefix)
	}
	// If the values cannot be the pair, move the first one to content.
	values := buf.in.Values
	if len(values)%2 != 0 {
		if buf.in.Content != "" {
			buf.in.Content += " "
		}
		buf.in.Content += gconv.String(values[0])
		values = values[1:]
	}
	if buf.in.Content != "" {
		buf.addValue(structureKeyContent, buf.in.Content)
	}
	// Values pairs.
	for i := 0; i < len(values); i += 2 {
		buf.addValue(values[i], values[i+1])
	}
	if buf.in.Stack != "" {
		buf.addValue(structureKeyStack, buf.in.Stack)
	}
	contentBytes := buf.buffer.Bytes()
	buf.buffer.Reset()
	contentBytes = bytes.ReplaceAll(contentBytes, []byte{'\n'}, []byte{' '})
	return contentBytes
}

func (buf *structuredBuffer) addValue(k, v any) {
	var (
		ks = gconv.String(k)
		vs = gconv.String(v)
	)
	if buf.buffer.Len() > 0 {
		buf.buffer.WriteByte(' ')
	}
	buf.appendString(ks)
	buf.buffer.WriteByte('=')
	buf.appendString(vs)
}

func (buf *structuredBuffer) appendString(s string) {
	if buf.needsQuoting(s) {
		s = strconv.Quote(s)
	}
	buf.buffer.WriteString(s)
}

func (buf *structuredBuffer) needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need quoting in a
			// JSON string, as well as space and '='
			if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}
