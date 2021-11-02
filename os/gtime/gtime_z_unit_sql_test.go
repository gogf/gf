package gtime_test

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestTime_Scan(t1 *testing.T) {
	gtest.C(t1, func(t *gtest.T) {
		tt := gtime.Time{}
		//test string
		s := gtime.Now().String()
		t.Assert(tt.Scan(s), nil)
		t.Assert(tt.String(), s)
		//test nano
		n := gtime.TimestampNano()
		t.Assert(tt.Scan(n), nil)
		t.Assert(tt.TimestampNano(), n)
		//test nil
		none := (*gtime.Time)(nil)
		t.Assert(none.Scan(nil), nil)
		t.Assert(none, nil)
	})

}

func TestTime_Value(t1 *testing.T) {
	gtest.C(t1, func(t *gtest.T) {
		tt := gtime.Now()
		s, err := tt.Value()
		t.Assert(err, nil)
		t.Assert(s, tt.Time)
		//test nil
		none := (*gtime.Time)(nil)
		s, err = none.Value()
		t.Assert(err, nil)
		t.Assert(s, nil)

	})
}
