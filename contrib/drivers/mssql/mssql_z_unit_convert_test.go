// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	mssqldriver "github.com/microsoft/go-mssqldb"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/drivers/mssql/v2"
)

// Test_CheckLocalTypeForField_UniqueIdentifier verifies that uniqueidentifier (any case)
// maps to LocalTypeUUID.
func Test_CheckLocalTypeForField_UniqueIdentifier(t *testing.T) {
	var (
		ctx    = context.Background()
		driver = mssql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		localType, err := driver.CheckLocalTypeForField(ctx, "uniqueidentifier", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeUUID)
	})

	// Case-insensitive match: SQL Server / driver layer may return upper-case names.
	gtest.C(t, func(t *gtest.T) {
		localType, err := driver.CheckLocalTypeForField(ctx, "UNIQUEIDENTIFIER", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeUUID)
	})
}

// Test_ConvertValueForLocal_UniqueIdentifier verifies the byte-order swap that turns
// SQL Server's wire-format UNIQUEIDENTIFIER bytes into a canonical uuid.UUID.
//
// The wire format puts the first 8 bytes in the little-endian COM/Win32 GUID layout
// and the remaining 8 bytes in big-endian RFC 4122 order.
func Test_ConvertValueForLocal_UniqueIdentifier(t *testing.T) {
	var (
		ctx    = context.Background()
		driver = mssql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Wire bytes for the UUID DA93D4F6-223F-42B2-A647-789371FFA693:
		//   first 4 bytes  little-endian:  F6 D4 93 DA  → DA93D4F6
		//   next 2 bytes   little-endian:  3F 22        → 223F
		//   next 2 bytes   little-endian:  B2 42        → 42B2
		//   last 8 bytes   big-endian:     A6 47 78 93 71 FF A6 93
		wireBytes := []byte{
			0xF6, 0xD4, 0x93, 0xDA,
			0x3F, 0x22,
			0xB2, 0x42,
			0xA6, 0x47, 0x78, 0x93, 0x71, 0xFF, 0xA6, 0x93,
		}
		want := uuid.MustParse("DA93D4F6-223F-42B2-A647-789371FFA693")

		got, err := driver.ConvertValueForLocal(ctx, "uniqueidentifier", wireBytes)
		t.AssertNil(err)
		t.Assert(got, want)
	})

	// go-mssqldb's UniqueIdentifier.Scan also accepts string form, so values
	// already converted server-side (e.g. via CAST AS NVARCHAR(36)) round-trip
	// correctly.
	gtest.C(t, func(t *gtest.T) {
		want := uuid.MustParse("DA93D4F6-223F-42B2-A647-789371FFA693")

		got, err := driver.ConvertValueForLocal(ctx, "uniqueidentifier",
			"DA93D4F6-223F-42B2-A647-789371FFA693")
		t.AssertNil(err)
		t.Assert(got, want)
	})

	// Sanity check: feeding the same wire bytes directly into go-mssqldb's
	// UniqueIdentifier scanner yields the same UUID, so this driver hook stays
	// in lock-step with the underlying go-mssqldb implementation.
	gtest.C(t, func(t *gtest.T) {
		var ms mssqldriver.UniqueIdentifier
		t.AssertNil(ms.Scan([]byte{
			0xF6, 0xD4, 0x93, 0xDA,
			0x3F, 0x22,
			0xB2, 0x42,
			0xA6, 0x47, 0x78, 0x93, 0x71, 0xFF, 0xA6, 0x93,
		}))
		t.Assert(uuid.UUID(ms).String(), "da93d4f6-223f-42b2-a647-789371ffa693")
	})

	// Invalid wire bytes (wrong length) propagate the underlying scan error.
	gtest.C(t, func(t *gtest.T) {
		_, err := driver.ConvertValueForLocal(ctx, "uniqueidentifier", []byte{0x01, 0x02})
		t.AssertNE(err, nil)
	})
}
