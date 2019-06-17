package gfpool_test

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfpool"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"testing"
	"time"
)

// TestOpen test open file cache
func TestOpen(t *testing.T) {
	testFile := start("TestOpen.txt")

	gtest.Case(t, func() {
		f, _ := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f.Close()

		f2, _ := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		gtest.AssertEQ(f, f2)
		f2.Close()

		// Deprecated test
		f3, _ := gfpool.OpenFile(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		gtest.AssertEQ(f, f3)
		f3.Close()

	})

	stop(testFile)
}

// TestOpenErr test open file error
func TestOpenErr(t *testing.T) {
	gtest.Case(t, func() {
		testErrFile := "errorPath"
		_, err := gfpool.Open(testErrFile, os.O_RDWR, 0666)
		gtest.AssertNE(err, nil)

		// delete file error
		testFile := start("TestOpenDeleteErr.txt")
		f, _ := gfpool.Open(testFile, os.O_RDWR, 0666)
		f.Close()
		stop(testFile)
		_, err = gfpool.Open(testFile, os.O_RDWR, 0666)
		gtest.AssertNE(err, nil)

		// append mode delete file error
		testFile = start("TestOpenCreateErr.txt")
		f, _ = gfpool.Open(testFile, os.O_CREATE, 0666)
		f.Close()
		stop(testFile)
		_, err = gfpool.Open(testFile, os.O_CREATE, 0666)
		gtest.AssertNE(err, nil)

		// append mode delete file error
		testFile = start("TestOpenAppendErr.txt")
		f, _ = gfpool.Open(testFile, os.O_APPEND, 0666)
		f.Close()
		stop(testFile)
		_, err = gfpool.Open(testFile, os.O_APPEND, 0666)
		gtest.AssertNE(err, nil)

		// trunc mode delete file error
		testFile = start("TestOpenTruncErr.txt")
		f, _ = gfpool.Open(testFile, os.O_TRUNC, 0666)
		f.Close()
		stop(testFile)
		_, err = gfpool.Open(testFile, os.O_TRUNC, 0666)
		gtest.AssertNE(err, nil)
	})
}

// TestOpenExpire test open file cache expire
func TestOpenExpire(t *testing.T) {
	testFile := start("TestOpenExpire.txt")

	gtest.Case(t, func() {
		f, _ := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100)
		f.Close()

		time.Sleep(150 * time.Millisecond)
		f2, _ := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100)
		gtest.AssertNE(f, f2)
		f2.Close()

		// Deprecated test
		f3, _ := gfpool.OpenFile(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100)
		gtest.AssertEQ(f2, f3)
		f3.Close()
	})

	stop(testFile)
}

// TestNewPool test gfpool new function
func TestNewPool(t *testing.T) {
	testFile := start("TestNewPool.txt")

	gtest.Case(t, func() {
		f, _ := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f.Close()

		pool := gfpool.New(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f2, _ := pool.File()
		// pool not equal
		gtest.AssertNE(f, f2)
		f2.Close()

		pool.Close()
	})

	stop(testFile)
}

// test before
func start(name string) string {
	testFile := os.TempDir() + string(os.PathSeparator) + name
	if gfile.Exists(testFile) {
		gfile.Remove(testFile)
	}
	content := "123"
	gfile.PutContents(testFile, content)
	return testFile
}

// test after
func stop(testFile string) {
	if gfile.Exists(testFile) {
		err := gfile.Remove(testFile)
		if err != nil {
			glog.Error(err)
		}
	}
}
