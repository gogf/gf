package gfpool_test

import (
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gfpool"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/test/gtest"
)

// TestOpen test open file cache
func TestOpen(t *testing.T) {
	testFile := start("TestOpen.txt")

	gtest.Case(t, func() {
		f, err := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		gtest.AssertEQ(err, nil)
		f.Close()

		f2, err1 := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		gtest.AssertEQ(err1, nil)
		gtest.AssertEQ(f, f2)
		f2.Close()
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
		pool := gfpool.New(testFile, os.O_RDWR, 0666)
		_, err1 := pool.File()
		gtest.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		gtest.AssertNE(err1, nil)

		// append mode delete file error and create again
		testFile = start("TestOpenCreateErr.txt")
		pool = gfpool.New(testFile, os.O_CREATE, 0666)
		_, err1 = pool.File()
		gtest.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		gtest.AssertEQ(err1, nil)

		// append mode delete file error
		testFile = start("TestOpenAppendErr.txt")
		pool = gfpool.New(testFile, os.O_APPEND, 0666)
		_, err1 = pool.File()
		gtest.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		gtest.AssertNE(err1, nil)

		// trunc mode delete file error
		testFile = start("TestOpenTruncErr.txt")
		pool = gfpool.New(testFile, os.O_TRUNC, 0666)
		_, err1 = pool.File()
		gtest.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		gtest.AssertNE(err1, nil)
	})
}

// TestOpenExpire test open file cache expire
func TestOpenExpire(t *testing.T) {
	testFile := start("TestOpenExpire.txt")

	gtest.Case(t, func() {
		f, err := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100*time.Millisecond)
		gtest.AssertEQ(err, nil)
		f.Close()

		time.Sleep(150 * time.Millisecond)
		f2, err1 := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100*time.Millisecond)
		gtest.AssertEQ(err1, nil)
		//gtest.AssertNE(f, f2)
		f2.Close()
	})

	stop(testFile)
}

// TestNewPool test gfpool new function
func TestNewPool(t *testing.T) {
	testFile := start("TestNewPool.txt")

	gtest.Case(t, func() {
		f, err := gfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		gtest.AssertEQ(err, nil)
		f.Close()

		pool := gfpool.New(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f2, err1 := pool.File()
		// pool not equal
		gtest.AssertEQ(err1, nil)
		//gtest.AssertNE(f, f2)
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
