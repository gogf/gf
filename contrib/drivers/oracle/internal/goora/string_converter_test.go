// This file verifies go-ora string converter regressions that affect the Oracle driver.

package goora_test

import (
	"testing"

	"github.com/sijms/go-ora/v2/converters"

	"github.com/gogf/gf/v2/test/gtest"
)

// goOraCharsetGBK is the Oracle charset ID for the go-ora GBK converter.
const goOraCharsetGBK = 0x354

// TestStringConverterDecodeGBKTrailingLeadByte verifies a truncated GBK sequence does not panic.
func TestStringConverterDecodeGBKTrailingLeadByte(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		converter := converters.NewStringConverter(goOraCharsetGBK)
		decoded := converter.Decode([]byte{0x81})

		t.Assert(decoded, string([]byte{0x81}))
	})
}
