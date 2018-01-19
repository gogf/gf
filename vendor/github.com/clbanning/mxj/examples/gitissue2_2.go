// https://github.com/clbanning/mxj/issues/17

package main

import (
	"fmt"
	"github.com/clbanning/mxj"
	"io"
	"os"
)

func main() {
	fh, err := os.Open("gitissue2.dat")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	m := make(map[string]interface{})
	for {
		v, err := mxj.NewMapXmlSeqReader(fh)
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
