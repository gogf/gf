// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// extraTypeNames are the entries of localTypeMap that `oid.TypeName` does not list, so that
// Test_LocalTypeMapHasNoUnknownName can tell a deliberate addition from a stale one:
//
//   - `_decimal` is a SQL alias of `_numeric` that the driver never reports as such, and is
//     kept for callers passing a declared type name.
//   - the multirange types exist since PostgreSQL 14, while the oid list of the pinned
//     `lib/pq` predates them.
var extraTypeNames = map[string]struct{}{
	"_decimal":        {},
	"int4multirange":  {},
	"int8multirange":  {},
	"_int4multirange": {},
	"_int8multirange": {},
}

// localTypeMap maps every type name the underlying driver can report, as listed in
// `oid.TypeName` of `github.com/lib/pq`, to the local Go type of its value, plus the
// entries of extraTypeNames.
//
// The mapping is exhaustive by design and Test_LocalTypeCoverage fails when the driver
// gains a type name that is missing here. Matching the name exactly, rather than looking
// for keywords in it, is what keeps a type like `point` from being taken for an integer
// because its name happens to contain "int", which used to convert its value to 0.
var localTypeMap = map[string]gdb.LocalType{
	"_decimal":        gdb.LocalTypeFloat64Slice,
	"int4multirange":  gdb.LocalTypeString,
	"int8multirange":  gdb.LocalTypeString,
	"_int4multirange": gdb.LocalTypeStringSlice,
	"_int8multirange": gdb.LocalTypeStringSlice,

	"_abstime":         gdb.LocalTypeStringSlice, // was datetime
	"_aclitem":         gdb.LocalTypeString,
	"_bit":             gdb.LocalTypeString,
	"_bool":            gdb.LocalTypeBoolSlice,
	"_box":             gdb.LocalTypeString,
	"_bpchar":          gdb.LocalTypeStringSlice,
	"_bytea":           gdb.LocalTypeBytesSlice,
	"_char":            gdb.LocalTypeStringSlice,
	"_cid":             gdb.LocalTypeString,
	"_cidr":            gdb.LocalTypeString,
	"_circle":          gdb.LocalTypeString,
	"_cstring":         gdb.LocalTypeString,
	"_date":            gdb.LocalTypeStringSlice, // was datetime
	"_daterange":       gdb.LocalTypeStringSlice, // was datetime
	"_float4":          gdb.LocalTypeFloat32Slice,
	"_float8":          gdb.LocalTypeFloat64Slice,
	"_gtsvector":       gdb.LocalTypeString,
	"_inet":            gdb.LocalTypeString,
	"_int2":            gdb.LocalTypeInt32Slice,
	"_int2vector":      gdb.LocalTypeStringSlice, // was int
	"_int4":            gdb.LocalTypeInt32Slice,
	"_int4range":       gdb.LocalTypeStringSlice,
	"_int8":            gdb.LocalTypeInt64Slice,
	"_int8range":       gdb.LocalTypeStringSlice,
	"_interval":        gdb.LocalTypeStringSlice,
	"_json":            gdb.LocalTypeString,
	"_jsonb":           gdb.LocalTypeString,
	"_line":            gdb.LocalTypeString,
	"_lseg":            gdb.LocalTypeString,
	"_macaddr":         gdb.LocalTypeString,
	"_money":           gdb.LocalTypeFloat64Slice,
	"_name":            gdb.LocalTypeString,
	"_numeric":         gdb.LocalTypeFloat64Slice,
	"_numrange":        gdb.LocalTypeString,
	"_oid":             gdb.LocalTypeString,
	"_oidvector":       gdb.LocalTypeString,
	"_path":            gdb.LocalTypeString,
	"_pg_lsn":          gdb.LocalTypeString,
	"_point":           gdb.LocalTypeStringSlice,
	"_polygon":         gdb.LocalTypeString,
	"_record":          gdb.LocalTypeString,
	"_refcursor":       gdb.LocalTypeString,
	"_regclass":        gdb.LocalTypeString,
	"_regconfig":       gdb.LocalTypeString,
	"_regdictionary":   gdb.LocalTypeString,
	"_regnamespace":    gdb.LocalTypeString,
	"_regoper":         gdb.LocalTypeString,
	"_regoperator":     gdb.LocalTypeString,
	"_regproc":         gdb.LocalTypeString,
	"_regprocedure":    gdb.LocalTypeString,
	"_regrole":         gdb.LocalTypeString,
	"_regtype":         gdb.LocalTypeString,
	"_reltime":         gdb.LocalTypeStringSlice, // was datetime
	"_text":            gdb.LocalTypeStringSlice,
	"_tid":             gdb.LocalTypeString,
	"_time":            gdb.LocalTypeStringSlice, // was datetime
	"_timestamp":       gdb.LocalTypeStringSlice, // was datetime
	"_timestamptz":     gdb.LocalTypeStringSlice, // was datetime
	"_timetz":          gdb.LocalTypeStringSlice, // was datetime
	"_tinterval":       gdb.LocalTypeStringSlice, // was int
	"_tsquery":         gdb.LocalTypeString,
	"_tsrange":         gdb.LocalTypeString,
	"_tstzrange":       gdb.LocalTypeString,
	"_tsvector":        gdb.LocalTypeString,
	"_txid_snapshot":   gdb.LocalTypeString,
	"_uuid":            gdb.LocalTypeUUIDSlice,
	"_varbit":          gdb.LocalTypeString,
	"_varchar":         gdb.LocalTypeStringSlice,
	"_xid":             gdb.LocalTypeString,
	"_xml":             gdb.LocalTypeString,
	"abstime":          gdb.LocalTypeDatetime,
	"aclitem":          gdb.LocalTypeString,
	"any":              gdb.LocalTypeString,
	"anyarray":         gdb.LocalTypeString,
	"anyelement":       gdb.LocalTypeString,
	"anyenum":          gdb.LocalTypeString,
	"anynonarray":      gdb.LocalTypeString,
	"anyrange":         gdb.LocalTypeString,
	"bit":              gdb.LocalTypeInt64Bytes, // bit(1) is a boolean; see CheckLocalTypeForField
	"bool":             gdb.LocalTypeBool,
	"box":              gdb.LocalTypeString,
	"bpchar":           gdb.LocalTypeString,
	"bytea":            gdb.LocalTypeBytes,
	"char":             gdb.LocalTypeString,
	"cid":              gdb.LocalTypeString,
	"cidr":             gdb.LocalTypeString,
	"circle":           gdb.LocalTypeString,
	"cstring":          gdb.LocalTypeString,
	"date":             gdb.LocalTypeDate,
	"daterange":        gdb.LocalTypeString, // was datetime
	"event_trigger":    gdb.LocalTypeString,
	"fdw_handler":      gdb.LocalTypeString,
	"float4":           gdb.LocalTypeFloat64,
	"float8":           gdb.LocalTypeFloat64,
	"gtsvector":        gdb.LocalTypeString,
	"index_am_handler": gdb.LocalTypeString,
	"inet":             gdb.LocalTypeString,
	"int2":             gdb.LocalTypeInt,
	"int2vector":       gdb.LocalTypeString, // was int
	"int4":             gdb.LocalTypeInt,
	"int4range":        gdb.LocalTypeString,
	"int8":             gdb.LocalTypeInt64,
	"int8range":        gdb.LocalTypeString,
	"internal":         gdb.LocalTypeString, // was int
	"interval":         gdb.LocalTypeString,
	"json":             gdb.LocalTypeJson,
	"jsonb":            gdb.LocalTypeJsonb,
	"language_handler": gdb.LocalTypeString,
	"line":             gdb.LocalTypeString,
	"lseg":             gdb.LocalTypeString,
	"macaddr":          gdb.LocalTypeString,
	"money":            gdb.LocalTypeString,
	"name":             gdb.LocalTypeString,
	"numeric":          gdb.LocalTypeString,
	"numrange":         gdb.LocalTypeString,
	"oid":              gdb.LocalTypeString,
	"oidvector":        gdb.LocalTypeString,
	"opaque":           gdb.LocalTypeString,
	"path":             gdb.LocalTypeString,
	"pg_attribute":     gdb.LocalTypeString,
	"pg_auth_members":  gdb.LocalTypeString,
	"pg_authid":        gdb.LocalTypeString,
	"pg_class":         gdb.LocalTypeString,
	"pg_database":      gdb.LocalTypeString,
	"pg_ddl_command":   gdb.LocalTypeString,
	"pg_lsn":           gdb.LocalTypeString,
	"pg_node_tree":     gdb.LocalTypeString,
	"pg_proc":          gdb.LocalTypeString,
	"pg_shseclabel":    gdb.LocalTypeString,
	"pg_type":          gdb.LocalTypeString,
	"point":            gdb.LocalTypeString,
	"polygon":          gdb.LocalTypeString,
	"record":           gdb.LocalTypeString,
	"refcursor":        gdb.LocalTypeString,
	"regclass":         gdb.LocalTypeString,
	"regconfig":        gdb.LocalTypeString,
	"regdictionary":    gdb.LocalTypeString,
	"regnamespace":     gdb.LocalTypeString,
	"regoper":          gdb.LocalTypeString,
	"regoperator":      gdb.LocalTypeString,
	"regproc":          gdb.LocalTypeString,
	"regprocedure":     gdb.LocalTypeString,
	"regrole":          gdb.LocalTypeString,
	"regtype":          gdb.LocalTypeString,
	"reltime":          gdb.LocalTypeString, // was datetime
	"smgr":             gdb.LocalTypeString,
	"text":             gdb.LocalTypeString,
	"tid":              gdb.LocalTypeString,
	"time":             gdb.LocalTypeTime,
	"timestamp":        gdb.LocalTypeDatetime,
	"timestamptz":      gdb.LocalTypeDatetime,
	"timetz":           gdb.LocalTypeDatetime,
	"tinterval":        gdb.LocalTypeString, // was int
	"trigger":          gdb.LocalTypeString,
	"tsm_handler":      gdb.LocalTypeString,
	"tsquery":          gdb.LocalTypeString,
	"tsrange":          gdb.LocalTypeString,
	"tstzrange":        gdb.LocalTypeString,
	"tsvector":         gdb.LocalTypeString,
	"txid_snapshot":    gdb.LocalTypeString,
	"unknown":          gdb.LocalTypeString,
	"uuid":             gdb.LocalTypeUUID,
	"varbit":           gdb.LocalTypeString,
	"varchar":          gdb.LocalTypeString,
	"void":             gdb.LocalTypeString,
	"xid":              gdb.LocalTypeString,
	"xml":              gdb.LocalTypeString,
}
