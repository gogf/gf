// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/text/gstr"
)

func Test_CamelCase(t *testing.T) {
	cases := [][]string{
		{"test_case", "TestCase"},
		{"test", "Test"},
		{"TestCase", "TestCase"},
		{" test  case ", "TestCase"},
		{"userLogin_log.bak", "UserLoginLogBak"},
		{"", ""},
		{"many_many_words", "ManyManyWords"},
		{"AnyKind of_string", "AnyKindOfString"},
		{"odd-fix", "OddFix"},
		{"numbers2And55with000", "Numbers2And55With000"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.CamelCase(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}

func Test_CamelLowerCase(t *testing.T) {
	cases := [][]string{
		{"foo-bar", "fooBar"},
		{"TestCase", "testCase"},
		{"", ""},
		{"AnyKind of_string", "anyKindOfString"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.CamelLowerCase(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}

func Test_SnakeCase(t *testing.T) {
	cases := [][]string{
		{"testCase", "test_case"},
		{"TestCase", "test_case"},
		{"Test Case", "test_case"},
		{" Test Case", "test_case"},
		{"Test Case ", "test_case"},
		{" Test Case ", "test_case"},
		{"test", "test"},
		{"test_case", "test_case"},
		{"Test", "test"},
		{"", ""},
		{"ManyManyWords", "many_many_words"},
		{"manyManyWords", "many_many_words"},
		{"AnyKind of_string", "any_kind_of_string"},
		{"numbers2and55with000", "numbers_2_and_55_with_000"},
		{"JSONData", "json_data"},
		{"userID", "user_id"},
		{"AAAbbb", "aa_abbb"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.SnakeCase(in)
		if result != out {
			t.Error("'" + in + "'('" + result + "' != '" + out + "')")
		}
	}
}

func Test_DelimitedCase(t *testing.T) {
	cases := [][]string{
		{"testCase", "test@case"},
		{"TestCase", "test@case"},
		{"Test Case", "test@case"},
		{" Test Case", "test@case"},
		{"Test Case ", "test@case"},
		{" Test Case ", "test@case"},
		{"test", "test"},
		{"test_case", "test@case"},
		{"Test", "test"},
		{"", ""},
		{"ManyManyWords", "many@many@words"},
		{"manyManyWords", "many@many@words"},
		{"AnyKind of_string", "any@kind@of@string"},
		{"numbers2and55with000", "numbers@2@and@55@with@000"},
		{"JSONData", "json@data"},
		{"userID", "user@id"},
		{"AAAbbb", "aa@abbb"},
		{"test-case", "test@case"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.DelimitedCase(in, '@')
		if result != out {
			t.Error("'" + in + "' ('" + result + "' != '" + out + "')")
		}
	}
}

func Test_SnakeScreamingCase(t *testing.T) {
	cases := [][]string{
		{"testCase", "TEST_CASE"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.SnakeScreamingCase(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}

func Test_KebabCase(t *testing.T) {
	cases := [][]string{
		{"testCase", "test-case"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.KebabCase(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}

func Test_KebabScreamingCase(t *testing.T) {
	cases := [][]string{
		{"testCase", "TEST-CASE"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.KebabScreamingCase(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}

func Test_DelimitedScreamingCase(t *testing.T) {
	cases := [][]string{
		{"testCase", "TEST.CASE"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.DelimitedScreamingCase(in, '.', true)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}

func TestSnakeFirstUpperCase(t *testing.T) {
	cases := [][]string{
		{"RGBCodeMd5", "rgb_code_md5"},
		{"testCase", "test_case"},
		{"Md5", "md5"},
		{"userID", "user_id"},
		{"RGB", "rgb"},
		{"RGBCode", "rgb_code"},
		{"_ID", "id"},
		{"User_ID", "user_id"},
		{"user_id", "user_id"},
		{"md5", "md5"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := gstr.SnakeFirstUpperCase(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}
