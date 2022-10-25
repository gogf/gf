package clickhouse

import (
	"errors"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	errUnsupportedInsertIgnore = errors.New("unsupported method:InsertIgnore")
	errUnsupportedInsertGetId  = errors.New("unsupported method:InsertGetId")
	errUnsupportedReplace      = errors.New("unsupported method:Replace")
	errUnsupportedBegin        = errors.New("unsupported method:Begin")
	errUnsupportedTransaction  = errors.New("unsupported method:Transaction")
	//
	matchFieldPatternType = map[string]string{
		matchBigIntPattern:      gdb.LocalTypeBigInt,
		matchDatePattern:        gdb.LocalTypeDatetime,
		matchDecimalPattern:     gdb.LocalTypeDecimal,
		matchFixedStringPattern: gdb.LocalTypeString,
		matchArrayPattern:       gdb.LocalTypeArray,
	}
	// match base type map
	matchBaseTypeMap = map[string]string{
		"int8":    LocalTypeInt8,
		"int16":   LocalTypeInt16,
		"int32":   LocalTypeInt32,
		"int64":   gdb.LocalTypeInt64,
		"uint8":   LocalTypeUInt8,
		"uint16":  LocalTypeUInt16,
		"uint32":  LocalTypeUInt32,
		"uint64":  gdb.LocalTypeUint64,
		"float32": gdb.LocalTypeFloat32,
		"float64": gdb.LocalTypeFloat64,
		"string":  gdb.LocalTypeString,
		"ipv4":    LocalTypeUInt32,
		"ipv6":    gdb.LocalTypeString,
		"uuid":    gdb.LocalTypeUUID,
		"geo":     gdb.LocalTypeInterface,
		"bool":    LocalTypeUInt8,
		"json":    gdb.LocalTypeJson,
	}
)

const (
	updateFilterPattern                 = `(?i)UPDATE[\s]+?(\w+[\.]?\w+)[\s]+?SET`
	deleteFilterPattern                 = `(?i)DELETE[\s]+?FROM[\s]+?(\w+[\.]?\w+)`
	filterTypePattern                   = `(?i)^UPDATE|DELETE`
	replaceSchemaPattern                = `@(.+?)/([\w\.\-]+)+`
	needParsedSqlInCtx      gctx.StrKey = "NeedParsedSql"
	OrmTagForStruct                     = "orm"
	driverName                          = "clickhouse"
	matchBigIntPattern                  = "[u]?int(128|256)"
	matchDatePattern                    = "^date"
	matchDecimalPattern                 = "^decimal"
	matchFixedStringPattern             = "^fixedstring"
	matchArrayPattern                   = "^array"
)

const (
	LocalTypeInt8   = "int8"
	LocalTypeInt16  = "int16"
	LocalTypeInt32  = "int32"
	LocalTypeUInt8  = "uint8"
	LocalTypeUInt16 = "uint16"
	LocalTypeUInt32 = "uint32"
)
