// Preserve list order with intermixed sub-elements.
// from: https://groups.google.com/forum/#!topic/golang-nuts/8KvlKsdh84k

package main

import (
	"fmt"
	"sort"

	"github.com/clbanning/mxj"
)

var data = `<node>
  <a>sadasd</a>
  <b>gfdfg</b>
  <c>bcbcbc</c>
  <a>hihihi</a>
  <b>jkjkjk</b>
  <a>lmlmlm</a>
</node>`

func main() {
	m, err := mxj.NewMapXmlSeq([]byte(data))
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	// Merge a b, and c members into a single list.
	// The value has "#text" key; "#seq" has the list sequence index.
	// Preserve both and the list key if you want to reencode.
	list := make([]*listval, 0)
	for k, v := range m["node"].(map[string]interface{}) {
		switch v.(type) {
		case map[string]interface{}:
			// This handles the lone 'c' element in the list.
			mem := v.(map[string]interface{})
			lval := &listval{k, mem["#text"].(string), mem["#seq"].(int)}
			list = append(list, lval)
		case []interface{}:
			// 'a' and 'b' were decoded as slices.
			for _, vv := range v.([]interface{}) {
				mem := vv.(map[string]interface{})
				lval := &listval{k, mem["#text"].(string), mem["#seq"].(int)}
				list = append(list, lval)
			}
		}
	}

	// Sort the list into orignal DOC sequence.
	sort.Sort(Lval(list))

	// Do some work with the list members - let's swap values.
	for i := 0; i < 3; i++ {
		list[i].val, list[5-i].val = list[5-i].val, list[i].val
	}

	// Rebuild map[string]interface{} value for "node".
	// Everything can be slice values - []interface{} - for encoding.
	a := make([]interface{}, 0)
	b := make([]interface{}, 0)
	c := make([]interface{}, 0)
	for _, v := range list {
		val := map[string]interface{}{"#text": v.val, "#seq": v.seq}
		switch v.list {
		case "a":
			a = append(a, interface{}(val))
		case "b":
			b = append(b, interface{}(val))
		case "c":
			c = append(c, interface{}(val))
		}
	}
	val := map[string]interface{}{"a": a, "b": b, "c": c}
	m["node"] = interface{}(val)

	x, err := m.XmlSeqIndent("", "  ")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(data)      // original
	fmt.Println(string(x)) // modified
}

// ======== sort interface implementation ========
type listval struct {
	list string
	val  string
	seq  int
}

type Lval []*listval

func (l Lval) Len() int {
	return len(l)
}

func (l Lval) Less(i, j int) bool {
	if l[i].seq < l[j].seq {
		return true
	}
	return false
}

func (l Lval) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
