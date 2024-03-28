// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/util/gutil"
)

// Go creates a new asynchronous goroutine function with specified recover function.
//
// The parameter `recoverFunc` is called when any panic during executing of `goroutineFunc`.
// If `recoverFunc` is given nil, it ignores the panic from `goroutineFunc` and no panic will
// throw to parent goroutine.
//
// But, note that, if `recoverFunc` also throws panic, such panic will be thrown to parent goroutine.
func Go(
	ctx context.Context,
	goroutineFunc func(ctx context.Context),
	recoverFunc func(ctx context.Context, exception error),
) {
	gutil.Go(ctx, goroutineFunc, recoverFunc)
}

// NewVar returns a gvar.Var.
func NewVar(i interface{}, safe ...bool) *Var {
	return gvar.New(i, safe...)
}

// Wait is an alias of ghttp.Wait, which blocks until all the web servers shutdown.
// It's commonly used in multiple servers' situation.
func Wait() {
	ghttp.Wait()
}

// Listen is an alias of gproc.Listen, which handles the signals received and automatically
// calls registered signal handler functions.
// It blocks until shutdown signals received and all registered shutdown handlers done.
func Listen() {
	gproc.Listen()
}

// Dump dumps a variable to stdout with more manually readable.
func Dump(values ...interface{}) {
	gutil.Dump(values...)
}

// DumpTo writes variables `values` as a string in to `writer` with more manually readable
func DumpTo(writer io.Writer, value interface{}, option gutil.DumpOption) {
	gutil.DumpTo(writer, value, option)
}

// DumpWithType acts like Dump, but with type information.
// Also see Dump.
func DumpWithType(values ...interface{}) {
	gutil.DumpWithType(values...)
}

// DumpWithOption returns variables `values` as a string with more manually readable.
func DumpWithOption(value interface{}, option gutil.DumpOption) {
	gutil.DumpWithOption(value, option)
}

// DumpJson pretty dumps json content to stdout.
func DumpJson(value any) {
	gutil.DumpJson(value)
}

// Throw throws an exception, which can be caught by TryCatch function.
func Throw(exception interface{}) {
	gutil.Throw(exception)
}

// Try implements try... logistics using internal panic...recover.
// It returns error if any exception occurs, or else it returns nil.
func Try(ctx context.Context, try func(ctx context.Context)) (err error) {
	return gutil.Try(ctx, try)
}

// TryCatch implements try...catch... logistics using internal panic...recover.
// It automatically calls function `catch` if any exception occurs and passes the exception as an error.
//
// But, note that, if function `catch` also throws panic, the current goroutine will panic.
func TryCatch(ctx context.Context, try func(ctx context.Context), catch func(ctx context.Context, exception error)) {
	gutil.TryCatch(ctx, try, catch)
}

// IsNil checks whether given `value` is nil.
// Parameter `traceSource` is used for tracing to the source variable if given `value` is type
// of pointer that also points to a pointer. It returns nil if the source is nil when `traceSource`
// is true.
// Note that it might use reflect feature which affects performance a little.
func IsNil(value interface{}, traceSource ...bool) bool {
	return empty.IsNil(value, traceSource...)
}

// IsEmpty checks whether given `value` empty.
// It returns true if `value` is in: 0, nil, false, "", len(slice/map/chan) == 0.
// Or else it returns true.
//
// The parameter `traceSource` is used for tracing to the source variable if given `value` is type of pointer
// that also points to a pointer. It returns true if the source is empty when `traceSource` is true.
// Note that it might use reflect feature which affects performance a little.
func IsEmpty(value interface{}, traceSource ...bool) bool {
	return empty.IsEmpty(value, traceSource...)
}

// RequestFromCtx retrieves and returns the Request object from context.
func RequestFromCtx(ctx context.Context) *ghttp.Request {
	return ghttp.RequestFromCtx(ctx)
}
