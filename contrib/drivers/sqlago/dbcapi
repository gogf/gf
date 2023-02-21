// vim:ts=4:sw=4:et

package sqlany

// dll bindings

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	API_VERSION_1     = 1
	API_VERSION_2     = 2
	SACAPI_ERROR_SIZE = 256

	libdbcapi_dll = "dbcapi.dll"
)

type dataType int32

const (
	// do not reorder
	A_INVALID_TYPE dataType = iota // invalid data type
	A_BINARY                       // Binary data: treated as is (no conversions performed)
	A_STRING                       // String data: character set conversion performed
	A_DOUBLE
	A_VAL64 // 64bit ints
	A_UVAL64
	A_VAL32 // 32bit ints
	A_UVAL32
	A_VAL16 // words
	A_UVAL16
	A_VAL8 // bytes
	A_UVAL8
)

type nativeType int32

const (
	DT_NOTYPE       = 0
	DT_DATE         = 384
	DT_TIME         = 388
	DT_TIMESTAMP    = 392
	DT_VARCHAR      = 448
	DT_FIXCHAR      = 452
	DT_LONGVARCHAR  = 456
	DT_STRING       = 460
	DT_DOUBLE       = 480
	DT_FLOAT        = 482
	DT_DECIMAL      = 484
	DT_INT          = 496
	DT_SMALLINT     = 500
	DT_BINARY       = 524
	DT_LONGBINARY   = 528
	DT_TINYINT      = 604
	DT_BIGINT       = 608
	DT_UNSINT       = 612
	DT_UNSSMALLINT  = 616
	DT_UNSBIGINT    = 620
	DT_BIT          = 624
	DT_LONGNVARCHAR = 640
)

// byte window used in *byte to slice conversion
const byteSliceWindow = 10 << 20

type dataValue struct {
	buffer     *byte
	buffersize uintptr
	length     *uintptr
	datatype   dataType
	isnull     *bool
}

// converts specified byte pointer to a proper slice object
func byteSlice(b *byte, size int) []byte {
	bs := make([]byte, size)
	s := bs[:]
	rawptr := uintptr(unsafe.Pointer(b))
	for size > byteSliceWindow {
		copy(s, (*[byteSliceWindow]byte)(unsafe.Pointer(rawptr))[:])
		s = s[byteSliceWindow:]
		size -= byteSliceWindow
		rawptr += byteSliceWindow
	}
	if size > 0 {
		copy(s, (*[byteSliceWindow]byte)(unsafe.Pointer(rawptr))[:size])
	}
	return bs
}

func (dv *dataValue) String() string {
	isnull := bool(*dv.isnull)
	s := fmt.Sprintf("type: %d, null: %t, length: %d, buffer size: %d, value: %s",
		dv.datatype, isnull, *dv.length, dv.buffersize,
		bytePtrToString(dv.buffer))
	return s
}

func (dv *dataValue) bufferValue() []byte {
	size := int(*dv.length)
	// [ap]: optimize by using a single cast for buffers upto 1mb in size
	// fall back to slower method if bigger
	if size < 1<<20 {
		b := make([]byte, int(*dv.length))
		copy(b, (*[1 << 20]byte)(unsafe.Pointer(dv.buffer))[:])
		return b
	}
	return byteSlice(dv.buffer, size)
}

func (dv *dataValue) isNull() bool {
	return bool(*dv.isnull)
}

// reference to resultset/statement/just character set?
func (dv *dataValue) Value() (v interface{}) {
	if dv.isNull() {
		// null
		v = nil
		return
	}
	switch dv.datatype {
	case A_BINARY:
		v = dv.bufferValue()
	case A_STRING:
		// currently, character set is configured as utf-8 (effective
		// for each connection, set as a connection option)
		// this will make the server provide text in unicode w/o having
		// to perform manual conversion
		v = byteSliceToString(dv.bufferValue())
	case A_DOUBLE:
		v = *(*float64)(unsafe.Pointer(dv.buffer))
	case A_VAL64:
		v = *(*int64)(unsafe.Pointer(dv.buffer))
	case A_UVAL64:
		v = *(*uint64)(unsafe.Pointer(dv.buffer))
	case A_VAL32:
		v = *(*int32)(unsafe.Pointer(dv.buffer))
	case A_UVAL32:
		v = *(*uint32)(unsafe.Pointer(dv.buffer))
	case A_VAL16:
		v = *(*int16)(unsafe.Pointer(dv.buffer))
	case A_UVAL16:
		v = *(*uint16)(unsafe.Pointer(dv.buffer))
	case A_VAL8:
		v = *(*int8)(unsafe.Pointer(dv.buffer))
	case A_UVAL8:
		v = *dv.buffer
	}
	return
}

type dataInfo struct {
	datatype dataType
	isnull   sacapi_bool
	datasize uintptr
}

type dataDirection byte

const (
	// do not reorder
	DD_INVALID      dataDirection = iota // Invalid data direction
	DD_INPUT                             // Input only host vars
	DD_OUTPUT                            // Output only host vars
	DD_INPUT_OUTPUT                      // Host vars of both directions
)

func (dd dataDirection) String() string {
	switch dd {
	case DD_INPUT:
		return "input"
	case DD_OUTPUT:
		return "input"
	case DD_INPUT_OUTPUT:
		return "input_output"
	}
	return "unknown direction"
}

type bindParam struct {
	dir   dataDirection
	value dataValue
	name  *byte // name of the bind param (used by DescribeBindParam)
}

func (bp *bindParam) String() string {
	s := fmt.Sprintf("name: %s; value: %s; dir: %s", bytePtrToString(bp.name),
		bp.value.String(), bp.dir)
	return s
}

type columnInfo struct {
	name       *byte
	datatype   dataType
	nativetype nativeType
	precision  uint16
	scale      uint16
	maxsize    uintptr
	nullable   sacapi_bool
}

func (ci *columnInfo) Name() string {
	return bytePtrToString(ci.name)
}

func (ci *columnInfo) String() string {
	s := fmt.Sprintf("name: %s, type: %d, native type: %d, size: %d, nullable: %v",
		ci.Name(), ci.datatype, ci.nativetype, ci.maxsize, ci.nullable)
	return s
}

// WIN32
type sacapi_u32 uint32
type sacapi_i32 int32
type sacapi_bool int32

var (
	dll = syscall.MustLoadDLL(libdbcapi_dll)

	sqlany_affected_rows       = dll.MustFindProc("sqlany_affected_rows")
	sqlany_bind_param          = dll.MustFindProc("sqlany_bind_param")
	sqlany_cancel              = dll.MustFindProc("sqlany_cancel")
	sqlany_clear_error         = dll.MustFindProc("sqlany_clear_error")
	sqlany_client_version      = dll.MustFindProc("sqlany_client_version")
	sqlany_client_version_ex   = dll.MustFindProc("sqlany_client_version_ex")
	sqlany_commit              = dll.MustFindProc("sqlany_commit")
	sqlany_connect             = dll.MustFindProc("sqlany_connect")
	sqlany_describe_bind_param = dll.MustFindProc("sqlany_describe_bind_param")
	sqlany_disconnect          = dll.MustFindProc("sqlany_disconnect")
	sqlany_error               = dll.MustFindProc("sqlany_error")
	sqlany_execute             = dll.MustFindProc("sqlany_execute")
	sqlany_execute_direct      = dll.MustFindProc("sqlany_execute_direct")
	sqlany_execute_immediate   = dll.MustFindProc("sqlany_execute_immediate")
	sqlany_fetch_absolute      = dll.MustFindProc("sqlany_fetch_absolute")
	sqlany_fetch_next          = dll.MustFindProc("sqlany_fetch_next")
	sqlany_fini                = dll.MustFindProc("sqlany_fini")
	sqlany_fini_ex             = dll.MustFindProc("sqlany_fini_ex")
	sqlany_free_connection     = dll.MustFindProc("sqlany_free_connection")
	sqlany_free_stmt           = dll.MustFindProc("sqlany_free_stmt")
	sqlany_get_bind_param_info = dll.MustFindProc("sqlany_get_bind_param_info")
	sqlany_get_column          = dll.MustFindProc("sqlany_get_column")
	sqlany_get_column_info     = dll.MustFindProc("sqlany_get_column_info")
	sqlany_get_data            = dll.MustFindProc("sqlany_get_data")
	sqlany_get_data_info       = dll.MustFindProc("sqlany_get_data_info")
	sqlany_get_next_result     = dll.MustFindProc("sqlany_get_next_result")
	sqlany_init                = dll.MustFindProc("sqlany_init")
	sqlany_init_ex             = dll.MustFindProc("sqlany_init_ex")
	sqlany_make_connection     = dll.MustFindProc("sqlany_make_connection")
	sqlany_make_connection_ex  = dll.MustFindProc("sqlany_make_connection_ex")
	sqlany_new_connection      = dll.MustFindProc("sqlany_new_connection")
	sqlany_new_connection_ex   = dll.MustFindProc("sqlany_new_connection_ex")
	sqlany_num_cols            = dll.MustFindProc("sqlany_num_cols")
	sqlany_num_params          = dll.MustFindProc("sqlany_num_params")
	sqlany_num_rows            = dll.MustFindProc("sqlany_num_rows")
	sqlany_prepare             = dll.MustFindProc("sqlany_prepare")
	sqlany_reset               = dll.MustFindProc("sqlany_reset")
	sqlany_rollback            = dll.MustFindProc("sqlany_rollback")
	sqlany_send_param_data     = dll.MustFindProc("sqlany_send_param_data")
	sqlany_sqlstate            = dll.MustFindProc("sqlany_sqlstate")
)

// TODO(ap): using syscall.(*Proc).Call incurs a slight overhead of
// a dynamically created slice of arguments.
// Might refactor later to avoid the allocation by directly using
// scyscall.Syscall/syscall.Syscall6 instead.

func sqlaInit(name string) bool {
	ret, _, _ := sqlany_init.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
		uintptr(API_VERSION_1),
		0)
	return ret != 1
}

func sqlaFini() {
	sqlany_fini.Call()
}

type sqlaConn uintptr

func newConnection() sqlaConn {
	ret, _, _ := sqlany_new_connection.Call()
	return sqlaConn(ret)
}

func (conn sqlaConn) free() {
	sqlany_free_connection.Call(uintptr(conn))
}

func (conn sqlaConn) connect(opts string) (err error) {
	ret, _, _ := sqlany_connect.Call(uintptr(conn),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(opts))))
	if ret != 1 {
		code, msg := conn.queryError()
		err = &sqlaError{code: code, msg: msg}
		return
	}
	return nil
}

func (conn sqlaConn) disconnect() bool {
	ret, _, _ := sqlany_disconnect.Call(uintptr(conn))
	return ret == 1
}

type sqlaStmt uintptr

func (conn sqlaConn) prepare(query string) (_ sqlaStmt, err error) {
	ret, _, _ := sqlany_prepare.Call(uintptr(conn),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(query))))
	if ret == 0 {
		err = conn.newError()
		return sqlaStmt(0), err
	}
	return sqlaStmt(ret), nil
}

func (stmt sqlaStmt) free() {
	sqlany_free_stmt.Call(uintptr(stmt))
}

func (stmt sqlaStmt) execute() bool {
	ret, _, _ := sqlany_execute.Call(uintptr(stmt))
	return ret == 1
}

func (conn sqlaConn) executeDirect(query string) (_ sqlaStmt, err error) {
	ret, _, _ := sqlany_execute_direct.Call(uintptr(conn),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(query))))
	if ret == 0 {
		err = conn.newError()
		return sqlaStmt(0), err
	}
	return sqlaStmt(ret), nil
}

func (conn sqlaConn) executeImmediate(query string) (err error) {
	ret, _, _ := syscall.Syscall(
		sqlany_execute_immediate.Addr(),
		uintptr(2),
		uintptr(conn),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(query))),
		0)
	if ret == 0 {
		err = conn.newError()
		return
	}
	return
}

// Reset a statement to its prepared state condition
func (stmt sqlaStmt) reset() bool {
	ret, _, _ := sqlany_reset.Call(uintptr(stmt))
	return ret == 1
}

// returns number of columns in result set or -1 upon failure
func (stmt sqlaStmt) numCols() int {
	ret, _, _ := sqlany_num_cols.Call(uintptr(stmt))
	return int(ret)
}

// returns number of rows affected by execution of a previously prepared
// statement
// returns -1 upon failure
func (stmt sqlaStmt) affectedRows() int {
	ret, _, _ := sqlany_affected_rows.Call(uintptr(stmt))
	return int(ret)
}

// returns number of parameters expected for a prepared statement
// returns -1 if the statement is invalid
func (stmt sqlaStmt) numParams() int {
	ret, _, _ := sqlany_num_params.Call(uintptr(stmt))
	return int(ret)
}

func (stmt sqlaStmt) fetchNext() bool {
	ret, _, _ := sqlany_fetch_next.Call(uintptr(stmt))
	return ret == 1
}

func (stmt sqlaStmt) fetchAbsolute(rownum sacapi_i32) bool {
	ret, _, _ := sqlany_fetch_absolute.Call(uintptr(stmt),
		uintptr(rownum))
	return ret == 1
}

// index specified parameter index in [0..NumParams()-1]
// bindparam receives the bind parameter information
func (stmt sqlaStmt) describeBindParam(index sacapi_u32, bindparam *bindParam) bool {
	ret, _, _ := sqlany_describe_bind_param.Call(uintptr(stmt),
		uintptr(index),
		uintptr(unsafe.Pointer(bindparam)))
	return ret == 1
}

// index specified parameter index in [0..NumParams()-1]
// bindparam specifies the bind parameter data
func (stmt sqlaStmt) bindParam(index sacapi_u32, bindparam *bindParam) bool {
	ret, _, _ := sqlany_bind_param.Call(uintptr(stmt),
		uintptr(index),
		uintptr(unsafe.Pointer(bindparam)))
	return ret == 1
}

func (conn sqlaConn) commit() bool {
	ret, _, _ := sqlany_commit.Call(uintptr(conn))
	return ret == 1
}

func (conn sqlaConn) rollback() bool {
	ret, _, _ := sqlany_rollback.Call(uintptr(conn))
	return ret == 1
}

// Retrieve data for column `colindex` in `dataval`.
//
// For A_BINARY and A_STRING data types, dataval.buffer points to internal
// buffer associated with result set. Users should copy the data out of
// provided pointers into their own buffer as it changes each time a row
// is fetched.
// dataval.length indicates the number of valid characters dataval.buffer points
// to - do not rely on buffer being null-terminated.
// This function will fetch _all_ data from server, if you do not want to allocate
// memory for resources of considerable size, use GetData instead.
func (stmt sqlaStmt) getColumn(colindex uint, dataval *dataValue) bool {
	ret, _, _ := sqlany_get_column.Call(uintptr(stmt),
		uintptr(sacapi_u32(colindex)),
		uintptr(unsafe.Pointer(dataval)))
	return ret == 1
}

func (stmt sqlaStmt) getColumnInfo(colindex sacapi_u32, colinfo *columnInfo) bool {
	ret, _, _ := sqlany_get_column_info.Call(uintptr(stmt),
		uintptr(colindex),
		uintptr(unsafe.Pointer(colinfo)))
	return ret == 1
}

func (stmt sqlaStmt) getDataInfo(colindex sacapi_u32, datainfo *dataInfo) bool {
	ret, _, _ := sqlany_get_data_info.Call(uintptr(stmt),
		uintptr(colindex),
		uintptr(unsafe.Pointer(datainfo)))
	return ret == 1
}

func (stmt sqlaStmt) getData(colindex sacapi_u32, offset uintptr,
	buffer *uintptr, size uintptr) bool {
	ret, _, _ := sqlany_get_data.Call(uintptr(stmt),
		uintptr(colindex),
		offset,
		uintptr(unsafe.Pointer(buffer)),
		size)
	return ret == 1
}

// Moves to the next result set in multiple result sets return
func (stmt sqlaStmt) getNextResult() bool {
	ret, _, _ := sqlany_get_next_result.Call(uintptr(stmt))
	return ret == 1
}

func (conn sqlaConn) newError() (err error) {
	code, msg := conn.queryError()
	if code != 0 {
		return &sqlaError{code: code, msg: msg}
	}
	return nil
}

func (conn sqlaConn) queryError() (code sacapi_i32, err string) {
	buf := make([]byte, SACAPI_ERROR_SIZE)
	ret, _, _ := sqlany_error.Call(uintptr(conn),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)))
	code = sacapi_i32(ret)
	err = byteSliceToString(buf)
	return
}

func byteSliceToString(b []byte) string {
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func bytePtrToString(b *byte) string {
	return byteSliceToString((*[1024]byte)(unsafe.Pointer(b))[:])
}

// A generic error type signalled when any of the low-level functions
// fail
type sqlaError struct {
	code sacapi_i32
	msg  string
}

func (err *sqlaError) Error() string {
	return fmt.Sprintf("Error: %s, Code: %#v", err.msg, err.code)
}

func (err *sqlaError) Fatal() bool {
	return err.code < 0
}
