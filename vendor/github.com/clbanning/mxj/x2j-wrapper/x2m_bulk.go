// Copyright 2012-2018 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

//	x2m_bulk.go: Process files with multiple XML messages.

package x2j

import (
	"bytes"
	"io"
	"os"
	"regexp"

	"github.com/clbanning/mxj"
)

// XmlMsgsFromFile()
//	'fname' is name of file
//	'phandler' is the map processing handler. Return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. Return of 'false' stops further processing and returns error.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or after processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func XmlMsgsFromFile(fname string, phandler func(map[string]interface{})(bool), ehandler func(error)(bool), recast ...bool) error {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	fi, fierr := os.Stat(fname)
	if fierr != nil {
		return fierr
	}
	fh, fherr := os.Open(fname)
	if fherr != nil {
		return fherr
	}
	defer fh.Close()
	buf := make([]byte,fi.Size())
	_, rerr  :=  fh.Read(buf)
	if rerr != nil {
		return rerr
	}
	doc := string(buf)

	// xml.Decoder doesn't properly handle whitespace in some doc
	// see songTextString.xml test case ... 
	reg,_ := regexp.Compile("[ \t\n\r]*<")
	doc = reg.ReplaceAllString(doc,"<")
	b := bytes.NewBufferString(doc)

	for {
		m, merr := mxj.NewMapXmlReader(b,r)
		if merr != nil && merr != io.EOF {
			if ok := ehandler(merr); !ok {
				// caused reader termination
				return merr
			 }
		}
		if m != nil {
			if ok := phandler(m); !ok {
				break
			}
		}
		if merr == io.EOF {
			break
		}
	}
	return nil
}

// XmlBufferToMap - process XML message from a bytes.Buffer
//	'b' is the buffer
//	Optional argument 'recast' coerces map values to float64 or bool where possible.
func XmlBufferToMap(b *bytes.Buffer,recast ...bool) (map[string]interface{},error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	return mxj.NewMapXmlReader(b, r)
}

// =============================  io.Reader version for stream processing  ======================

// XmlMsgsFromReader() - io.Reader version of XmlMsgsFromFile
//	'rdr' is an io.Reader for an XML message (stream)
//	'phandler' is the map processing handler. Return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. Return of 'false' stops further processing and returns error.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or after processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func XmlMsgsFromReader(rdr io.Reader, phandler func(map[string]interface{})(bool), ehandler func(error)(bool), recast ...bool) error {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	for {
		m, merr := ToMap(rdr,r)
		if merr != nil && merr != io.EOF {
			if ok := ehandler(merr); !ok {
				// caused reader termination
				return merr
			 }
		}
		if m != nil {
			if ok := phandler(m); !ok {
				break
			}
		}
		if merr == io.EOF {
			break
		}
	}
	return nil
}

