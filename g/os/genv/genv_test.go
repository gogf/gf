package genv_test

import (
	"github.com/gogf/gf/g/os/genv"
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"testing"
)

func Test_Genv_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(os.Environ(), genv.All())
	})
}

func Test_Genv_Get(t *testing.T) {
	gtest.Case(t, func() {
		key := "TEST_GET_ENV"
		err := os.Setenv(key, "TEST")
		gtest.Assert(err, nil)
		gtest.AssertEQ(genv.Get(key), "TEST")
	})
}

func Test_Genv_Set(t *testing.T) {
	gtest.Case(t, func() {
		key := "TEST_SET_ENV"
		err := genv.Set(key, "TEST")
		gtest.Assert(err, nil)
		gtest.AssertEQ(os.Getenv(key), "TEST")
	})
}

func Test_Genv_Remove(t *testing.T) {
	gtest.Case(t, func() {
		key := "TEST_REMOVE_ENV"
		err := os.Setenv(key, "TEST")
		gtest.Assert(err, nil)
		err = genv.Remove(key)
		gtest.Assert(err, nil)
		gtest.AssertEQ(os.Getenv(key), "")
	})
}
