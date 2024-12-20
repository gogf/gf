// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import (
	"bytes"
	"container/list"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/gogf/gf/v2/internal/consts"
)

// stackInfo manages stack info of certain error.
type stackInfo struct {
	Index   int        // Index is the index of current error in whole error stacks.
	Message string     // Error information string.
	Lines   *list.List // Lines contains all error stack lines of current error stack in sequence.
}

// stackLine manages each line info of stack.
type stackLine struct {
	Function string // Function name, which contains its full package path.
	FileLine string // FileLine is the source file name and its line number of Function.
}

// Stack returns the error stack information as string.
func (err *Error) Stack() string {
	if err == nil {
		return ""
	}

	id := 1

	count := 0
	stackInfos := make([]tempStackInfo, 0, 4)

	tsi := getErrorStackInfo(err, id)
	stackInfos = append(stackInfos, tsi)
	count += tsi.count

	temp := err.error
	for temp != nil {
		switch x := temp.(type) {
		case *Error:
			id++
			tsi = getErrorStackInfo(x, id)
			stackInfos = append(stackInfos, tsi)
			count += tsi.count + len("\n")
			temp = x.error
		default:
			// TODO sb.Write(x.Error())
			break
		}
	}

	var sb = getBytesBuffer()
	sb.Grow(count)
	defer putBytesBuffer(sb)

	for _, si := range stackInfos {
		sb.Write(si.funcLine)
		for k, b := range si.bufptrs {
			sb.Write(getStackTraceFuncIDHeader(k))
			sb.Write(*(*[]byte)(b))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Stack returns the error stack information as string.
func (err *Error) stackWithBuffer(buffer *bytes.Buffer, errError string) {
	if err == nil {
		return
	}
	var (
		id         = 1
		count      = 0
		stackInfos = make([]tempStackInfo, 0, 4)
	)
	tsi := getErrorStackInfo(err, id)
	stackInfos = append(stackInfos, tsi)
	count += tsi.count

	temp := err.error
	for temp != nil {
		switch x := temp.(type) {
		case *Error:
			id++
			tsi = getErrorStackInfo(x, id)
			stackInfos = append(stackInfos, tsi)
			count += tsi.count + len("\n")
			temp = x.error
		default:
			// TODO buffer.Write(x.Error())
			break
		}
	}

	if len(errError) > 0 {
		buffer.Grow(count + len(errError) + 1)
		buffer.WriteString(errError)
		buffer.WriteString("\n")
	} else {
		buffer.Grow(count)
	}

	for _, si := range stackInfos {
		buffer.Write(si.funcLine)
		for k, bytesPtr := range si.bufptrs {
			buffer.Write(getStackTraceFuncIDHeader(k))
			buffer.Write(*(*[]byte)(bytesPtr))
		}
		buffer.WriteByte('\n')
	}
}

type tempStackInfo struct {
	// buf [][]byte
	// bufptrs[0] = unsafe.Pointer(&buf[0])
	bufptrs []unsafe.Pointer
	// funcName
	// 		xx/yy/zz:10
	funcLine []byte
	// count+=len [for len range len(buf[i])]
	count int
}

func getErrorStackInfo(err *Error, id int) (tsi tempStackInfo) {
	var (
		count = getIntLength(id) + len(". ") + len(err.text)
		pcs   = err.stack
	)
	tsi.bufptrs = make([]unsafe.Pointer, len(pcs))

	for i, pc := range pcs {
		count += len("\n\t") + getIntLength(i+1) + len(").")

		bufptr := cacheStackTrace.getBytesPtr(pc)
		if bufptr != nil {
			count += len(*bufptr)
			tsi.bufptrs[i] = unsafe.Pointer(bufptr)
			continue
		}

		f, _ := runtime.CallersFrames(pcs[i : i+1]).Next()
		buf := []byte(f.Function + "\n\t\t" + f.File + ":" + strconv.Itoa(f.Line))
		bufptr = cacheStackTrace.addAndGetBytesPtr(pc, buf)
		count += len(buf)
		tsi.bufptrs[i] = unsafe.Pointer(bufptr)
	}
	tsi.funcLine = []byte(strconv.Itoa(id) + ". " + err.text)
	tsi.count = count
	return tsi
}

// 0-9     = 1
// 10-99   = 2
// 100-999 = 3
// ...
func getIntLength(n int) int {
	if n < 10 {
		return 1
	}
	if n < 100 {
		return 2
	}
	if n < 1000 {
		return 3
	}
	// 返回一个超大负数,strings.Builder.Grow() 就会panic
	return -10000000
}

// filterLinesOfStackInfos removes repeated lines, which exist in subsequent stacks, from top errors.
func filterLinesOfStackInfos(infos []*stackInfo) {
	var (
		ok      bool
		set     = make(map[string]struct{})
		info    *stackInfo
		line    *stackLine
		removes []*list.Element
	)
	for i := len(infos) - 1; i >= 0; i-- {
		info = infos[i]
		if info.Lines == nil {
			continue
		}
		for n, e := 0, info.Lines.Front(); n < info.Lines.Len(); n, e = n+1, e.Next() {
			line = e.Value.(*stackLine)
			if _, ok = set[line.FileLine]; ok {
				removes = append(removes, e)
			} else {
				set[line.FileLine] = struct{}{}
			}
		}
		if len(removes) > 0 {
			for _, e := range removes {
				info.Lines.Remove(e)
			}
		}
		removes = removes[:0]
	}
}

// formatStackInfos formats and returns error stack information as string.
func formatStackInfos(infos []*stackInfo) string {
	buffer := bytes.NewBuffer(nil)
	for i, info := range infos {
		buffer.WriteString(fmt.Sprintf("%d. %s\n", i+1, info.Message))
		if info.Lines != nil && info.Lines.Len() > 0 {
			formatStackLines(buffer, info.Lines)
		}
	}
	return buffer.String()
}

// formatStackLines formats and returns error stack lines as string.
func formatStackLines(buffer *bytes.Buffer, lines *list.List) string {
	var (
		line   *stackLine
		space  = "  "
		length = lines.Len()
	)
	for i, e := 0, lines.Front(); i < length; i, e = i+1, e.Next() {
		line = e.Value.(*stackLine)
		// Graceful indent.
		if i >= 9 {
			space = " "
		}
		buffer.WriteString(fmt.Sprintf(
			"   %d).%s%s\n        %s\n",
			i+1, space, line.Function, line.FileLine,
		))
	}
	return buffer.String()
}

// loopLinesOfStackInfo iterates the stack info lines and produces the stack line info.
func loopLinesOfStackInfo(st stack, info *stackInfo, isStackModeBrief bool) {
	if st == nil {
		return
	}
	for _, p := range st {
		if fn := runtime.FuncForPC(p - 1); fn != nil {
			file, line := fn.FileLine(p - 1)
			if isStackModeBrief {
				// filter whole GoFrame packages stack paths.
				if strings.Contains(file, consts.StackFilterKeyForGoFrame) {
					continue
				}
			} else {
				// package path stack filtering.
				if strings.Contains(file, stackFilterKeyLocal) {
					continue
				}
			}
			// Avoid stack string like "`autogenerated`"
			if strings.Contains(file, "<") {
				continue
			}
			// Ignore GO ROOT paths.
			if goRootForFilter != "" &&
				len(file) >= len(goRootForFilter) &&
				file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if info.Lines == nil {
				info.Lines = list.New()
			}
			info.Lines.PushBack(&stackLine{
				Function: fn.Name(),
				FileLine: fmt.Sprintf(`%s:%d`, file, line),
			})
		}
	}
}
