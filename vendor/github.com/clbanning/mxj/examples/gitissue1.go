// https://github.com/clbanning/mxj/issues/17

package main

import (
	"bytes"
	"fmt"
	"github.com/clbanning/mxj"
	"io"
)

var data = []byte(`
<?xml version="1.0" encoding="utf-8"?>
<doc><elem>just something to demo</elem></doc>
`)

func main() {
	r := bytes.NewReader(data)
	m := make(map[string]interface{})
	var v map[string]interface{}
	var err error
	for {
		v, err = mxj.NewMapXmlSeqReader(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != mxj.NoRoot {
				// handle error
			}
		}
		for key, val := range v {
			m[key] = val
		}
	}
	fmt.Printf("%v\n", m)
}
