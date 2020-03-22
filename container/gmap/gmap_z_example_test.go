package gmap_test

import (
	"fmt"

	"github.com/gogf/gf/container/gmap"
)

func Example_normalBasic() {
	m := gmap.New()

	//Add data
	m.Set("key1", "val1")

	//Print size
	fmt.Println(m.Size())
	//output 1

	add_map := make(map[interface{}]interface{})
	add_map["key2"] = "val2"
	add_map["key3"] = "val3"
	add_map[1] = 1

	fmt.Println(m.Values())

	//Batch add data
	m.Sets(add_map)

	//Gets the value of the corresponding key
	key3_val := m.Get("key3")
	fmt.Println(key3_val)

	//Get the value by key, or set it with given key-value if not exist.
	get_or_set_val := m.GetOrSet("key4", "val4")
	fmt.Println(get_or_set_val)

	// Set key-value if the key does not exist, then return true; or else return false.
	is_set := m.SetIfNotExist("key3", "val3")
	fmt.Println(is_set)

	//Remove key
	m.Remove("key2")
	fmt.Println(m.Keys())

	//Batch remove keys
	remove_keys := []interface{}{"key1", 1}
	m.Removes(remove_keys)
	fmt.Println(m.Keys())

	//Contains checks whether a key exists.
	is_contain := m.Contains("key3")
	fmt.Println(is_contain)

	//Flip exchanges key-value of the map, it will change key-value to value-key.
	m.Flip()
	fmt.Println(m.Map())

	// Clear deletes all data of the map,
	m.Clear()

	fmt.Println(m.Size())

}
func Example_normalMerge() {
	m1 := gmap.New()
	m2 := gmap.New()
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(m2)
	fmt.Println(m1.Map())
}
