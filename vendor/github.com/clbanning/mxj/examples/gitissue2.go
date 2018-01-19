// https://github.com/clbanning/mxj/issues/17

package main

import (
	"bytes"
	"fmt"
	"github.com/clbanning/mxj"
	"io"
	"io/ioutil"
)

func main() {
	b, err := ioutil.ReadFile("gitissue2.dat")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	r := bytes.NewReader(b)
	m := make(map[string]interface{})
	for {
		v, err := mxj.NewMapXmlSeqReader(r)
		// v, raw, err := mxj.NewMapXmlSeqReaderRaw(r)
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
		// fmt.Println(string(raw))
	}
	fmt.Printf("%v\n", m)
}
