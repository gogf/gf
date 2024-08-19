// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gxml_test

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestEncodeing(t *testing.T) {
	type Address struct {
		City, State string
	}
	type Person struct {
		XMLName   xml.Name `xml:"person"`
		Id        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
		Slice   []interface{}
	}

	a := []interface{}{"string", true, 36.4}

	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42,
		Slice: []interface{}{"a", "b", true, 1, a}}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
}

func TestA(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// data := map[string]interface{}{
		// 	"Person": "<>&'\"AAA",
		// }

		// data := map[string]interface{}{
		// 	"Person": map[string]interface{}{
		// 		"Name":  "<>&'\"AAA",
		// 		"Email": "john.doe@example.com",
		// 		"Bio":   "I am a software developer & I love coding.",
		// 	},
		// }

		data := map[string]interface{}{
			// "map2": map[string]string{
			// 	"a": "b",
			// 	"c": "d",
			// },
			"Name": "<>&'\"AAA",
			// // "Email":  "john.doe@example.com",
			// // "Bio":    "I am a software developer & I love coding.",
			// // "Slices": []interface{}{"a", "b", true, 1},
			"map": map[string]interface{}{
				"Name":  "<>&'\"AAA",
				"Email": "john.doe@example.com",
				"Bio":   "I am a software developer & I love coding.",
				"map2": map[string]string{
					"a": "b",
					"c": "d",
				},
			},
		}
		str, err := gxml.Encode2(data)
		t.AssertNil(err)
		fmt.Println(string(str))
	})
}

func TestAi(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := map[string]interface{}{
			"Name":  "<>&'\"AAA",
			"Email": "john.doe@example.com",
			"Bio":   "I am a software developer & I love coding.",
			// "Slices": []interface{}{"a", "b", true, 1},
			// "map": map[string]interface{}{
			// 	"Name":  "<>&'\"AAA",
			// 	"Email": "john.doe@example.com",
			// 	"Bio":   "I am a software developer & I love coding.",
			// 	"map2": map[string]string{
			// 		"a": "b",
			// 		"c": "d",
			// 	},
			// },
		}
		str, err := gxml.Encode(data)
		t.AssertNil(err)
		fmt.Println(string(str))
	})
}

func TestB(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var x = `<doc><Name>&lt;&gt;&amp;&#39;&#34;AAA</Name><Email>john.doe@example.com</Email><Bio>I am a software developer &amp; I love coding.</Bio></doc>`
		m, err := gxml.Decode([]byte(x))
		t.AssertNil(err)
		fmt.Println(m)
	})
}
