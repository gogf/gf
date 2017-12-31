
package gmmap

import (
    "bytes"
    "io/ioutil"
    "testing"
)

func TestMap(t *testing.T) {
    data, err := Map("mmap_test.go")
    if err != nil {
        t.Fatalf("Open: %v", err)
    }

    if exp, err := ioutil.ReadFile("mmap_test.go"); err != nil {
        t.Fatalf("ioutil.ReadFile: %v", err)
    } else if !bytes.Equal(data, exp) {
        t.Fatalf("got %q\nwant %q", string(data), string(exp))
    }
}