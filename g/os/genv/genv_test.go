package genv_test

import (
	"github.com/gogf/gf/g/os/genv"
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"testing"
	"time"
)

func Test_Genv_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(os.Environ(), genv.All())
	})
}

func Test_Genv_Get(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertEQ("keke", genv.Get("LLL"+time.Now().String(), "keke"))
		gtest.AssertEQ("", genv.Get("LLL"+time.Now().String()))
	})
}

func Test_Genv_Set(t *testing.T) {
	gtest.Case(t, func() {
		err := genv.Set("LLL", "keke")
		gtest.Assert(err, nil)
	})
}

func Test_Genv_Remove(t *testing.T) {
	gtest.Case(t, func() {
		err := genv.Remove("LLL")
		gtest.Assert(err, nil)
	})
}
