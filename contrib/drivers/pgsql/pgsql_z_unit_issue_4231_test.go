package pgsql_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Driver_ConvertValueForField_Issue4231(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx    = context.Background()
			driver = db.(*pgsql.Driver)
		)
		// bytea
		v1, err := driver.ConvertValueForField(ctx, "bytea", []byte("123"))
		t.AssertNil(err)
		t.Assert(v1, []byte("123"))

		// jsonb
		v2, err := driver.ConvertValueForField(ctx, "jsonb", []string{"a", "b"})
		t.AssertNil(err)
		t.AssertNE(v2, "{a,b}")
		
		// _int4 (array)
		v3, err := driver.ConvertValueForField(ctx, "_int4", []int{1, 2})
		t.AssertNil(err)
		t.Assert(v3, "{1,2}")
	})
}
