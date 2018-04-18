// x2j_findPath - utility functions to retrieve path to node in dot-notation
// Copyright 2012-2018 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package x2j

import (
	"strings"

	"github.com/clbanning/mxj"
)

//----------------------------- find all paths to a key --------------------------------
// Want eventually to extract shortest path and call GetValuesAtKeyPath()
// This will get all the possible paths.  These can be scanned for len(path) and sequence.

// Get all paths through the doc (in dot-notation) that terminate with the specified tag.
// Results can be used with ValuesAtTagPath() and ValuesFromTagPath().
func PathsForTag(doc string, key string) ([]string, error) {
	m, err := mxj.NewMapXml([]byte(doc))
	if err != nil {
		return nil, err
	}

	ss := PathsForKey(m, key)
	return ss, nil
}

// Extract the shortest path from all possible paths - from PathsForTag().
// Paths are strings using dot-notation.
func PathForTagShortest(doc string, key string) (string, error) {
	m, err := mxj.NewMapXml([]byte(doc))
	if err != nil {
		return "", err
	}

	s := PathForKeyShortest(m, key)
	return s, nil
}

// Get all paths through the doc (in dot-notation) that terminate with the specified tag.
// Results can be used with ValuesAtTagPath() and ValuesFromTagPath().
func BytePathsForTag(doc []byte, key string) ([]string, error) {
	m, err := mxj.NewMapXml(doc)
	if err != nil {
		return nil, err
	}

	ss := PathsForKey(m, key)
	return ss, nil
}

// Extract the shortest path from all possible paths - from PathsForTag().
// Paths are strings using dot-notation.
func BytePathForTagShortest(doc []byte, key string) (string, error) {
	m, err := ByteDocToMap(doc)
	if err != nil {
		return "", err
	}

	s := PathForKeyShortest(m, key)
	return s, nil
}

// Get all paths through the map (in dot-notation) that terminate with the specified key.
// Results can be used with ValuesAtKeyPath() and ValuesFromKeyPath().
func PathsForKey(m map[string]interface{}, key string) []string {
	breadbasket := make(map[string]bool,0)
	breadcrumb := ""

	hasKeyPath(breadcrumb, m, key, &breadbasket)
	if len(breadbasket) == 0 {
		return nil
	}

	// unpack map keys to return
	res := make([]string,len(breadbasket))
	var i int
	for k,_ := range breadbasket {
		res[i] = k
		i++
	}

	return res
}

// Extract the shortest path from all possible paths - from PathsForKey().
// Paths are strings using dot-notation.
func PathForKeyShortest(m map[string]interface{}, key string) string {
	paths := PathsForKey(m,key)

	lp := len(paths)
	if lp == 0 {
		return ""
	}
	if lp == 1 {
		return paths[0]
	}

	shortest := paths[0]
	shortestLen := len(strings.Split(shortest,"."))

	for i := 1 ; i < len(paths) ; i++ {
		vlen := len(strings.Split(paths[i],"."))
		if vlen < shortestLen {
			shortest = paths[i]
			shortestLen = vlen
		}
	}

	return shortest
}

// hasKeyPath - if the map 'key' exists append it to KeyPath.path and increment KeyPath.depth
// This is really just a breadcrumber that saves all trails that hit the prescribed 'key'.
func hasKeyPath(crumb string, iv interface{}, key string, basket *map[string]bool) {
	switch iv.(type) {
	case map[string]interface{}:
		vv := iv.(map[string]interface{})
		if _, ok := vv[key]; ok {
			if crumb == "" {
				crumb = key
			} else {
				crumb += "." + key
			}
			// *basket = append(*basket, crumb)
			(*basket)[crumb] = true
		}
		// walk on down the path, key could occur again at deeper node
		for k, v := range vv {
			// create a new breadcrumb, add the one we're at to the crumb-trail
			var nbc string
			if crumb == "" {
				nbc = k
			} else {
				nbc = crumb + "." + k
			}
			hasKeyPath(nbc, v, key, basket)
		}
	case []interface{}:
		// crumb-trail doesn't change, pass it on
		for _, v := range iv.([]interface{}) {
			hasKeyPath(crumb, v, key, basket)
		}
	}
}

