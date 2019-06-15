package gfpool_test

import (
	"fmt"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfpool"
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	testFile := start()

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

	stop()
}

func TestOpenErr(t *testing.T) {
	testFile := "errorPath"

	gtest.Case(t, func() {
		_, err := gfpool.Open(testFile, os.O_RDWR, 0666)
		gtest.AssertNE(err, nil)
	})
}

func TestOpenExipre(t *testing.T) {
	testFile := start()

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

	stop()
}

func start() string {
	testFile := os.TempDir() + string(os.PathSeparator) + "testGfpool.txt"
	fmt.Println(testFile)
	if gfile.Exists(testFile) {
		gfile.Remove(testFile)
	}
	content := "123"
	gfile.PutContents(testFile, content)
	return testFile
}

func stop() {
	testFile := os.TempDir() + string(os.PathSeparator) + "testGfpool.txt"
	if gfile.Exists(testFile) {
		gfile.Remove(testFile)
	}
}
